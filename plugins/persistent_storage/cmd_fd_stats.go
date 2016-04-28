package persistent_storage

import (
	"bytes"
	"fmt"

	"github.com/sejvlond/tarsier/human_size"
	tarsier_plugins "github.com/sejvlond/tarsier/plugins"
)

func NewFDStats(lgr tarsier_plugins.Logger, container *FDContainer) *FDStats {
	return &FDStats{
		lgr:       lgr,
		container: container,
	}
}

type FDStats struct {
	lgr       tarsier_plugins.Logger
	container *FDContainer
}

func (cmd *FDStats) Name() string {
	return NAME + "/fd_stats"
}

func (cmd *FDStats) Description() string {
	return `FDStats will print statistics of file descriptor usage`
}

func (cmd *FDStats) Execute(raw interface{}) (string, error) {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("File descriptors in use %v:\n",
		cmd.container.Count()))
	for i, file := range cmd.container.Files {
		info, err := file.Stat()
		if err != nil {
			buf.WriteString(fmt.Sprintf("\t%v => %v\n", i, err))
			continue
		}
		buf.WriteString(fmt.Sprintf("\t%v => %v\n", i,
			human_size.Format(int(info.Size()))))
	}
	return buf.String(), nil
}
