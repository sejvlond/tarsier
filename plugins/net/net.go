package net

import (
	tarsier_plugins "github.com/sejvlond/tarsier/plugins"
)

const NAME = "net"

type Net struct {
	lgr         tarsier_plugins.Logger
	connections *ConnContainer
}

func (n *Net) Init(lgr tarsier_plugins.Logger, config interface{}) error {
	n.lgr = lgr
	n.connections = NewConnContainer()
	n.lgr.Infof("Successfully initialized")
	return nil
}

func (n *Net) Commands() []tarsier_plugins.Command {
	return []tarsier_plugins.Command{
		NewDial(n.lgr, n.connections),
		NewClose(n.lgr, n.connections),
		NewStats(n.lgr, n.connections),
	}
}

// Register plugin
func init() {
	tarsier_plugins.RegisterPlugin(NAME, func() interface{} {
		return new(Net)
	})
}
