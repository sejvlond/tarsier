package net

import (
	"bytes"
	"fmt"

	tarsier_plugins "github.com/sejvlond/tarsier/plugins"
)

func NewStats(lgr tarsier_plugins.Logger, container *ConnContainer) *Stats {
	return &Stats{
		lgr:       lgr,
		container: container,
	}
}

type Stats struct {
	lgr       tarsier_plugins.Logger
	container *ConnContainer
}

func (cmd *Stats) Name() string {
	return NAME + "/stats"
}

func (cmd *Stats) Description() string {
	return `Stats will print statistics of connections usage`
}

func (cmd *Stats) Execute(raw interface{}) (string, error) {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("Connections in use %v:\n",
		cmd.container.Count()))
	for i, conn := range cmd.container.Conns {
		buf.WriteString(fmt.Sprintf("\t%v => %v (%v)\n", i,
			conn.RemoteAddr().String(), conn.RemoteAddr().Network()))
	}
	return buf.String(), nil
}
