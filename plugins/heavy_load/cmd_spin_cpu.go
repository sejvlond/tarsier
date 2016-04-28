package heavy_load

import (
	"fmt"
	"sync"
	"time"

	tarsier_plugins "github.com/sejvlond/tarsier/plugins"
)

func NewSpinCPU(lgr tarsier_plugins.Logger) *SpinCPU {
	return &SpinCPU{
		lgr: lgr,
	}
}

type SpinCPU struct {
	lgr tarsier_plugins.Logger
}

type SpinCPUData struct {
	Duration string
	Cores    uint
}

func (cmd *SpinCPU) Name() string {
	return NAME + "/spin_cpu"
}

func (cmd *SpinCPU) Description() string {
	return `SpinCPU will spin CPU for as long as you say<br>
Data:
<pre>
	duration: for how long to spin (decimal with suffix (no space)  Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h").
	cores: how many cores to use
</pre>`
}

func (cmd *SpinCPU) DataStruct() interface{} {
	return &SpinCPUData{}
}

func (cmd *SpinCPU) Execute(raw interface{}) (string, error) {
	data, ok := raw.(*SpinCPUData)
	if !ok {
		return "", fmt.Errorf("Tarsier, those data are not mine!")
	}
	duration, err := time.ParseDuration(data.Duration)
	if err != nil {
		return err.Error(), err
	}
	ticker := time.NewTicker(duration)
	var wg sync.WaitGroup
	for i := uint(0); i < data.Cores; i++ {
		wg.Add(1)
		go func() {
			cmd.lgr.Infof("Let's calculate some nonsense")
			run := true
			foo := float64(1231548487332)
			for run {
				select {
				case <-ticker.C:
					run = false
					break
				default:
				}
				foo = foo * 1548643286513 / 1531845416
			}
			wg.Done()
			cmd.lgr.Infof("Enough of maths")
		}()
	}
	wg.Wait()
	ticker.Stop()
	return fmt.Sprintf("Ok, %v CPUs was spinned as hell for %v\n",
		data.Cores, duration), nil
}
