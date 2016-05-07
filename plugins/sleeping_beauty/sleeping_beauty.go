package sleeping_beauty

import (
	tarsier_plugins "github.com/sejvlond/tarsier/plugins"
)

const NAME = "sleeping_beauty"

type SleepingBeauty struct {
	lgr tarsier_plugins.Logger
}

func (d *SleepingBeauty) Init(lgr tarsier_plugins.Logger, config interface{}) error {
	d.lgr = lgr
	d.lgr.Infof("Successfully initialized")
	return nil
}

func (d *SleepingBeauty) Commands() []tarsier_plugins.Command {
	return []tarsier_plugins.Command{
		NewSleep(d.lgr),
	}
}

// Register plugin
func init() {
	tarsier_plugins.RegisterPlugin(NAME, func() interface{} {
		return new(SleepingBeauty)
	})
}
