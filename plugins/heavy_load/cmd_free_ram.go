package heavy_load

import (
	"fmt"

	"github.com/sejvlond/tarsier/human_size"
	tarsier_plugins "github.com/sejvlond/tarsier/plugins"
)

func NewFreeRam(lgr tarsier_plugins.Logger, memory *Memory) *FreeRam {
	return &FreeRam{
		lgr:    lgr,
		memory: memory,
	}
}

type FreeRam struct {
	lgr    tarsier_plugins.Logger
	memory *Memory
}

type FreeRamData struct {
	Index int
}

func (cmd *FreeRam) Name() string {
	return NAME + "/free_ram"
}

func (cmd *FreeRam) Description() string {
	return fmt.Sprintf(`FreeRam will free allocated RAM regions<br>
Data:
<pre>
	index: index to be freed - see %v/ram_stats for indices; -1 means free everything
</pre>`, NAME)
}

func (cmd *FreeRam) DataStruct() interface{} {
	return &FreeRamData{Index: -1}
}

func (cmd *FreeRam) Execute(raw interface{}) (string, error) {
	data, ok := raw.(*FreeRamData)
	if !ok {
		return "", fmt.Errorf("Tarsier, those data are not mine!")
	}
	cmd.lgr.Infof("Freeing memory")
	cnt := cmd.memory.Free(data.Index)
	if cnt < 0 {
		return "Error freeing memory", fmt.Errorf("Error freeing memory")
	}
	if cnt == 0 {
		return "Index out of range", fmt.Errorf("Index out of range")
	}
	return fmt.Sprintf("Region %v was freed (%v)", data.Index,
		human_size.Format(cnt)), nil
}
