package heavy_load

import (
	"fmt"

	"github.com/sejvlond/tarsier/human_size"
	tarsier_plugins "github.com/sejvlond/tarsier/plugins"
)

func NewGobbleRam(lgr tarsier_plugins.Logger, memory *Memory) *GobbleRam {
	return &GobbleRam{
		lgr:    lgr,
		memory: memory,
	}
}

type GobbleRam struct {
	lgr    tarsier_plugins.Logger
	memory *Memory
}

type GobbleRamData struct {
	Amount string
	Resize bool
	Ratio  float32
}

func (cmd *GobbleRam) Name() string {
	return NAME + "/gobble_ram"
}

func (cmd *GobbleRam) Description() string {
	return `GobbleRam will gobble as much RAM as you say<br>
Data:
<pre>
	amount: how much data to allocate (number in bytes, or string with GB, MB, kB, b)
	resize: true/false - whenever to try to resize the amount if error (default true)
	ration: resize ratio (default 0.5 - amount*ratio when resizing)
</pre>`
}

func (cmd *GobbleRam) DataStruct() interface{} {
	return &GobbleRamData{
		Resize: true,
		Ratio:  0.5,
	}
}

func (cmd *GobbleRam) Execute(raw interface{}) (string, error) {
	data, ok := raw.(*GobbleRamData)
	if !ok {
		return "", fmt.Errorf("Tarsier, those data are not mine!")
	}
	amount, ok := human_size.Parse(data.Amount)
	if !ok {
		return "Invalid amount", fmt.Errorf("Invalid amount")
	}
	cmd.lgr.Infof("Allocation memory")
	cnt := cmd.memory.Alloc(amount, data.Resize, data.Ratio)
	if cnt < 0 {
		return "Could not allocate any more memory",
			fmt.Errorf("Could not allocate any more memmory")
	}
	return fmt.Sprintf("%v RAM were gobbled", human_size.Format(cnt)), nil
}
