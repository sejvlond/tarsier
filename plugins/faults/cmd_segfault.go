package faults

import (
	"errors"
	"unsafe"

	tarsier_plugins "github.com/sejvlond/tarsier/plugins"
)

func NewSegFault(lgr tarsier_plugins.Logger) *SegFault {
	return &SegFault{
		lgr: lgr,
	}
}

type SegFault struct {
	lgr tarsier_plugins.Logger
}

func (cmd *SegFault) Name() string {
	return NAME + "/segfault"
}

func (cmd *SegFault) Description() string {
	return `SegFault will segfault`
}

func (cmd *SegFault) Execute(raw interface{}) (string, error) {
	cmd.lgr.Infof("Going to seg fault in 3.. 2.. 1..")
	x := 1
	x = *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&x)) + 50000000))
	return "Ok segfault probably goes wront", errors.New("Should seg fault...")
}
