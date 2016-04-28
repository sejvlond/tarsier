package persistent_storage

import (
	"fmt"

	tarsier_plugins "github.com/sejvlond/tarsier/plugins"
)

func NewCloseFD(lgr tarsier_plugins.Logger, container *FDContainer) *CloseFD {
	return &CloseFD{
		lgr:       lgr,
		container: container,
	}
}

type CloseFD struct {
	lgr       tarsier_plugins.Logger
	container *FDContainer
}

type CloseFDData struct {
	Index int
}

func (cmd *CloseFD) Name() string {
	return NAME + "/close_fd"
}

func (cmd *CloseFD) Description() string {
	return fmt.Sprintf(`CloseFD will close previously opened file descriptor<br>
Data:
<pre>
	index: index to be closed - see %v/fd_stats for indices; -1 means close everything
</pre>`, NAME)
}

func (cmd *CloseFD) DataStruct() interface{} {
	return &CloseFDData{Index: -1}
}

func (cmd *CloseFD) Execute(raw interface{}) (string, error) {
	data, ok := raw.(*CloseFDData)
	if !ok {
		return "", fmt.Errorf("Tarsier, those data are not mine!")
	}
	if err := cmd.container.Close(data.Index); err != nil {
		return fmt.Sprintf("Error while closing: %v", err), err
	}
	return fmt.Sprintf("File descriptor %v was closed", data.Index), nil
}
