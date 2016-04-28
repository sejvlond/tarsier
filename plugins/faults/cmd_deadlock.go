package faults

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	tarsier_plugins "github.com/sejvlond/tarsier/plugins"
)

func NewDeadlock(lgr tarsier_plugins.Logger) *Deadlock {
	return &Deadlock{
		lgr: lgr,
	}
}

type Deadlock struct {
	lgr tarsier_plugins.Logger
}

type DeadlockData struct {
	Amount   uint
	Duration string
}

func (cmd *Deadlock) Name() string {
	return NAME + "/deadlock"
}

func (cmd *Deadlock) Description() string {
	return `Deadlock will deadlock as many goroutines as you say<br>
Data:
<pre>
	amount: how many threads to lock
	duration: for how long to lock (decimal with suffix (no space)  Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h").
</pre>`
}

func (cmd *Deadlock) DataStruct() interface{} {
	return &DeadlockData{}
}

func (cmd *Deadlock) Execute(raw interface{}) (string, error) {
	data, ok := raw.(*DeadlockData)
	if !ok {
		return "", fmt.Errorf("Tarsier, those data are not mine!")
	}

	duration, err := time.ParseDuration(data.Duration)
	if err != nil {
		cmd.lgr.Errorf("Duration %v is invalid: %v", data.Duration, err)
		return err.Error(), err
	}
	var wg, lock sync.WaitGroup
	lock.Add(1)
	for i := uint(0); i < data.Amount; i++ {
		wg.Add(1)
		go func() {
			runtime.LockOSThread()
			cmd.lgr.Infof("Going to deadlock in 3.. 2.. 1..")
			lock.Wait()
			runtime.UnlockOSThread()
			wg.Done()
			cmd.lgr.Infof("Unlocked!")
		}()
	}
	time.Sleep(duration)
	cmd.lgr.Infof("Let's unlock goroutines")
	lock.Done()
	wg.Wait()

	return fmt.Sprintf("Ok %v threads were locked for %v", data.Amount,
		duration), nil
}
