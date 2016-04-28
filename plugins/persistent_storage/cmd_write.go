package persistent_storage

import (
	"fmt"
	"os"
	"sync"

	"github.com/sejvlond/tarsier/human_size"
	tarsier_plugins "github.com/sejvlond/tarsier/plugins"
)

func NewWrite(lgr tarsier_plugins.Logger, container *FDContainer) *Write {
	return &Write{
		lgr:       lgr,
		container: container,
	}
}

type Write struct {
	lgr       tarsier_plugins.Logger
	container *FDContainer
}

type WriteData struct {
	Amount string
	Files  uint
}

func (cmd *Write) Name() string {
	return NAME + "/write"
}

func (cmd *Write) Description() string {
	return `Write will write as many data as you say into randomly chosen files (count can be specify as well)<br>
Data:
<pre>
	amount: how much data to write (number in bytes, or string with GB, MB, kB, b)
	files: into how many files to write paralely
</pre>`
}

func (cmd *Write) DataStruct() interface{} {
	return &WriteData{}
}

func (cmd *Write) Execute(raw interface{}) (string, error) {
	data, ok := raw.(*WriteData)
	if !ok {
		return "", fmt.Errorf("Tarsier, those data are not mine!")
	}
	amount, ok := human_size.Parse(data.Amount)
	if !ok {
		return "Invalid amount", fmt.Errorf("Invalid amount")
	}
	files, err := cmd.container.Get(int(data.Files))
	if err != nil {
		return fmt.Sprintf("Error create file descriptors: %v", err), err
	}
	text := make([]byte, amount)
	var wg sync.WaitGroup
	var glErr error
	for _, file := range files {
		wg.Add(1)
		go func(file *os.File) {
			cmd.lgr.Infof("Writing to file '%v'", file.Name())
			_, err := file.Write(text)
			if err != nil {
				glErr = err
			}
			wg.Done()
			cmd.lgr.Infof("Everything was written")
		}(file)
	}
	wg.Wait()
	if glErr != nil {
		return fmt.Sprintf("Error writing data to file: %v", glErr), glErr
	}
	return fmt.Sprintf("%v of text was written to %v files", data.Amount,
		data.Files), nil
}
