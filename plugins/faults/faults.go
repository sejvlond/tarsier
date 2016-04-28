package faults

import (
	tarsier_plugins "github.com/sejvlond/tarsier/plugins"
)

const NAME = "faults"

type Faults struct {
	lgr tarsier_plugins.Logger
}

func (f *Faults) Init(lgr tarsier_plugins.Logger, config interface{}) error {
	f.lgr = lgr
	f.lgr.Infof("Successfully initialized")
	return nil
}

func (f *Faults) Commands() []tarsier_plugins.Command {
	return []tarsier_plugins.Command{
		NewDeadlock(f.lgr),
		NewSegFault(f.lgr),
	}
}

// Register plugin
func init() {
	tarsier_plugins.RegisterPlugin(NAME, func() interface{} {
		return new(Faults)
	})
}
