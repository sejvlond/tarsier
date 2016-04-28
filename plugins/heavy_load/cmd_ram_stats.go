package heavy_load

import (
	"bytes"
	"fmt"

	"github.com/sejvlond/tarsier/human_size"
	tarsier_plugins "github.com/sejvlond/tarsier/plugins"
)

func NewRamStats(lgr tarsier_plugins.Logger, memory *Memory) *RamStats {
	return &RamStats{
		lgr:    lgr,
		memory: memory,
	}
}

type RamStats struct {
	lgr    tarsier_plugins.Logger
	memory *Memory
}

func (cmd *RamStats) Name() string {
	return NAME + "/ram_stats"
}

func (cmd *RamStats) Description() string {
	return `RamStats will print statistics of RAM usage`
}

func (cmd *RamStats) Execute(raw interface{}) (string, error) {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("Memory in use %v:\n",
		human_size.Format(cmd.memory.Size)))
	for i, reg := range cmd.memory.Regions {
		buf.WriteString(fmt.Sprintf("\t%v => %v\n", i,
			human_size.Format(len(reg))))
	}
	return buf.String(), nil
}
