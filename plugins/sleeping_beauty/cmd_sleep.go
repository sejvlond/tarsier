package sleeping_beauty

import (
	"fmt"
	"time"

	tarsier_plugins "github.com/sejvlond/tarsier/plugins"
)

func NewSleep(lgr tarsier_plugins.Logger) *Sleep {
	return &Sleep{
		lgr: lgr,
	}
}

type Sleep struct {
	lgr tarsier_plugins.Logger
}

type SleepData struct {
	Duration string
}

func (cmd *Sleep) Name() string {
	return NAME + "/sleep"
}

func (cmd *Sleep) Description() string {
	return fmt.Sprintf(`Sleep will send Sleeping Beauty to sleep until set time
	expires and her Prince wakes her<br>
Data:
<pre>
	duration: time to sleep (decimal with suffix (no space)  Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h")
</pre>`, NAME)
}

func (cmd *Sleep) DataStruct() interface{} {
	return &SleepData{}
}

func (cmd *Sleep) Execute(raw interface{}) (string, error) {
	data, ok := raw.(*SleepData)
	if !ok {
		return "", fmt.Errorf("Tarsier, those data are not mine!")
	}
	duration, err := time.ParseDuration(data.Duration)
	if err != nil {
		cmd.lgr.Errorf("Duration %v is invalid: %v", data.Duration, err)
		return err.Error(), err
	}
	time.Sleep(duration)
	return fmt.Sprintf("I slept for %v. Thanks :-*", duration), nil
}
