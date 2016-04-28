package persistent_storage

import (
	"fmt"

	tarsier_plugins "github.com/sejvlond/tarsier/plugins"
)

func NewOpenFD(lgr tarsier_plugins.Logger, container *FDContainer) *OpenFD {
	return &OpenFD{
		lgr:       lgr,
		container: container,
	}
}

type OpenFD struct {
	lgr       tarsier_plugins.Logger
	container *FDContainer
}

type OpenFDData struct {
	Amount uint
}

func (cmd *OpenFD) Name() string {
	return NAME + "/open_fd"
}

func (cmd *OpenFD) Description() string {
	return `OpenFD will open as many file descriptors as you say<br>
Data:
<pre>
	amount: how many FDs to open
</pre>`
}

func (cmd *OpenFD) DataStruct() interface{} {
	return &OpenFDData{}
}

func (cmd *OpenFD) Execute(raw interface{}) (string, error) {
	data, ok := raw.(*OpenFDData)
	if !ok {
		return "", fmt.Errorf("Tarsier, those data are not mine!")
	}
	var (
		cnt uint
		err error
	)
	cmd.lgr.Infof("Opening %v files", data.Amount)
	for cnt = 0; cnt < data.Amount; cnt++ {
		if err = cmd.container.Open(); err != nil {
			break
		}
	}
	if cnt < data.Amount {
		return fmt.Sprintf("Could not open more than %v file descriptors: %v",
			cnt, err), err
	}
	return fmt.Sprintf("%v file descriptors were opened", cnt), nil
}
