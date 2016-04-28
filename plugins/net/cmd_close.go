package net

import (
	"fmt"

	tarsier_plugins "github.com/sejvlond/tarsier/plugins"
)

func NewClose(lgr tarsier_plugins.Logger, container *ConnContainer) *Close {
	return &Close{
		lgr:       lgr,
		container: container,
	}
}

type Close struct {
	lgr       tarsier_plugins.Logger
	container *ConnContainer
}

type CloseData struct {
	Index int
}

func (cmd *Close) Name() string {
	return NAME + "/close"
}

func (cmd *Close) Description() string {
	return fmt.Sprintf(`Close will close previously opened connection<br>
Data:
<pre>
	index: index to be closed - see %v/stats for indices; -1 means close everything
</pre>`, NAME)
}

func (cmd *Close) DataStruct() interface{} {
	return &CloseData{Index: -1}
}

func (cmd *Close) Execute(raw interface{}) (string, error) {
	data, ok := raw.(*CloseData)
	if !ok {
		return "", fmt.Errorf("Tarsier, those data are not mine!")
	}
	if err := cmd.container.Close(data.Index); err != nil {
		return fmt.Sprintf("Error while closing: %v", err), err
	}
	return fmt.Sprintf("Connection %v was closed", data.Index), nil
}
