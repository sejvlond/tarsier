package heavy_load

import (
	tarsier_plugins "github.com/sejvlond/tarsier/plugins"
)

const NAME = "heavy_load"

type HeavyLoad struct {
	lgr    tarsier_plugins.Logger
	memory *Memory
}

func (h *HeavyLoad) Init(lgr tarsier_plugins.Logger, config interface{}) error {
	h.lgr = lgr
	h.memory = NewMemory()
	h.lgr.Infof("Successfully initialized")
	return nil
}

func (h *HeavyLoad) Commands() []tarsier_plugins.Command {
	return []tarsier_plugins.Command{
		NewGobbleRam(h.lgr, h.memory),
		NewFreeRam(h.lgr, h.memory),
		NewRamStats(h.lgr, h.memory),
		NewSpinCPU(h.lgr),
	}
}

// Register plugin
func init() {
	tarsier_plugins.RegisterPlugin(NAME, func() interface{} {
		return new(HeavyLoad)
	})
}
