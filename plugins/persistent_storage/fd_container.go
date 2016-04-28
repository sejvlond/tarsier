package persistent_storage

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type FDContainer struct {
	Files []*os.File
	mutex sync.Mutex
	dir   string
	cnt   int
}

func NewFDContainer(baseDir string) (*FDContainer, error) {
	c := new(FDContainer)
	var err error
	c.dir, err = ioutil.TempDir(baseDir, "tarsier_persistent_storage")
	if err != nil {
		return nil, err
	}
	c.init()
	return c, nil
}

func (c *FDContainer) Finalize() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.close(-1)
	os.RemoveAll(c.dir)
}

func (c *FDContainer) init() {
	c.Files = make([]*os.File, 0)
	c.cnt = 0
}

func (c *FDContainer) Count() int {
	return len(c.Files)
}

func (c *FDContainer) Open() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.open()
}

func (c *FDContainer) open() error {
	file, err := os.Create(filepath.Join(c.dir, strconv.Itoa(c.cnt)))
	if err != nil {
		return err
	}
	c.Files = append(c.Files, file)
	c.cnt++
	return nil
}

func (c *FDContainer) del(file *os.File) error {
	if err := file.Close(); err != nil {
		return err
	}
	if err := os.Remove(file.Name()); err != nil {
		return err
	}
	return nil
}

func (c *FDContainer) Close(index int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.close(index)
}

func (c *FDContainer) close(index int) error {
	if index == -1 {
		for _, file := range c.Files {
			if err := c.del(file); err != nil {
				return err
			}
		}
		c.init()
		return nil
	}
	if index < 0 || index >= len(c.Files) {
		return fmt.Errorf("Index out of bounds")
	}
	file := c.Files[index]
	if err := c.del(file); err != nil {
		return err
	}
	c.Files = append(c.Files[:index], c.Files[index+1:]...)
	return nil
}

// Get will return as many file descriptors as count specified
// if there is not enought, it will create new ones
func (c *FDContainer) Get(count int) ([]*os.File, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if count > c.Count() {
		missing := count - c.Count()
		for i := 0; i < missing; i++ {
			if err := c.open(); err != nil {
				return nil, err
			}
		}
	}
	files := make([]*os.File, count)
	used := make(map[int]bool)
	var idx int
	for i := 0; i < count; i++ {
		for { // generate random index from c.Files that was not used before
			idx = rand.Intn(c.Count())
			if _, exists := used[idx]; !exists {
				break
			}
		}
		used[idx] = true
		files[i] = c.Files[idx]
	}
	return files, nil
}

func init() {
	rand.Seed(time.Now().Unix())
}
