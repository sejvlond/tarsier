package heavy_load

import (
	"sync"

	"github.com/edsrzf/mmap-go"
)

type Memory struct {
	Size    int
	Regions []mmap.MMap
	mutex   sync.Mutex
}

func NewMemory() *Memory {
	m := new(Memory)
	m.init()
	return m
}

func (m *Memory) init() {
	m.Regions = make([]mmap.MMap, 0)
	m.Size = 0
}

func (m *Memory) Alloc(amount int, resize bool, ratio float32) int {
	reg, err := mmap.MapRegion(nil, amount, mmap.RDWR, mmap.ANON, 0)
	if err != nil {
		if resize {
			amount = int(float32(amount) * ratio)
			if amount == 0 {
				return -1
			}
			return m.Alloc(amount, resize, ratio)
		}
		return -1
	}
	// use the mmap
	for i := range reg {
		reg[i] = 0
	}
	m.mutex.Lock()
	m.Regions = append(m.Regions, reg)
	m.Size += len(reg)
	m.mutex.Unlock()
	return len(reg)
}

func (m *Memory) Free(index int) int {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if index == -1 {
		size := m.Size
		for _, region := range m.Regions {
			if err := region.Unmap(); err != nil {
				return -1
			}
		}
		m.init()
		return size
	}
	if index < 0 || index >= len(m.Regions) {
		return 0
	}
	region := m.Regions[index]
	size := len(region)
	if err := region.Unmap(); err != nil {
		return -1
	}
	m.Regions = append(m.Regions[:index], m.Regions[index+1:]...)
	m.Size -= size
	return size
}
