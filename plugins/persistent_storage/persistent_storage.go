package persistent_storage

import (
	"fmt"

	tarsier_plugins "github.com/sejvlond/tarsier/plugins"
)

const NAME = "persistent_storage"

type Config struct {
	BaseDir string `yaml:"base_dir"`
}

type PersistentStorage struct {
	lgr       tarsier_plugins.Logger
	container *FDContainer
}

func (ps *PersistentStorage) ConfigStruct() interface{} {
	return new(Config)
}

func (ps *PersistentStorage) Init(lgr tarsier_plugins.Logger, config interface{}) (err error) {
	cfg, ok := config.(*Config)
	if !ok {
		return fmt.Errorf("Get invalid config '%v'", config)
	}
	ps.lgr = lgr
	ps.container, err = NewFDContainer(cfg.BaseDir)
	if err != nil {
		return err
	}
	ps.lgr.Infof("Successfully initialized")
	return nil
}

func (ps *PersistentStorage) Commands() []tarsier_plugins.Command {
	return []tarsier_plugins.Command{
		NewOpenFD(ps.lgr, ps.container),
		NewFDStats(ps.lgr, ps.container),
		NewCloseFD(ps.lgr, ps.container),
		NewWrite(ps.lgr, ps.container),
		NewRead(ps.lgr),
	}
}

func (ps *PersistentStorage) Close() error {
	ps.container.Finalize()
	return nil
}

// Register plugin
func init() {
	tarsier_plugins.RegisterPlugin(NAME, func() interface{} {
		return new(PersistentStorage)
	})
}
