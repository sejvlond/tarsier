package dummy

import (
	"fmt"

	tarsier_plugins "github.com/sejvlond/tarsier/plugins"
)

const NAME = "dummy"

type Config struct {
	Test   int
	Struct struct {
		Arr []int
		Str string
	}
}

type Dummy struct {
	cfg *Config
	lgr tarsier_plugins.Logger
}

func (d *Dummy) Init(lgr tarsier_plugins.Logger, config interface{}) error {
	var ok bool
	d.cfg, ok = config.(*Config)
	if !ok {
		return fmt.Errorf("Get invalid config '%v'", config)
	}
	d.lgr = lgr
	d.lgr.Infof("Successfully initialized")
	return nil
}

func (d *Dummy) Commands() []tarsier_plugins.Command {
	return []tarsier_plugins.Command{
		NewSayHelloCmd(d.lgr),
	}
}

func (d *Dummy) ConfigStruct() interface{} {
	return new(Config)
}

// -- SayHelloCmd --

func NewSayHelloCmd(lgr tarsier_plugins.Logger) *SayHelloCmd {
	return &SayHelloCmd{lgr: lgr}
}

type SayHelloData struct {
	Amount int
}

type SayHelloCmd struct {
	lgr tarsier_plugins.Logger
}

func (cmd *SayHelloCmd) Name() string {
	return NAME + "/say_hello"
}

func (cmd *SayHelloCmd) Description() string {
	return `simple command that can say hello as many times, as user send in 
its data:<br><pre>amount: 20</pre>will produce message 
<pre>Hello 20 times</pre>`
}

func (cmd *SayHelloCmd) DataStruct() interface{} {
	return new(SayHelloData)
}

func (cmd *SayHelloCmd) Execute(raw interface{}) (string, error) {
	data, ok := raw.(*SayHelloData)
	if !ok {
		return "", fmt.Errorf("Tarsier, those data are not mine!")
	}
	// do something with those data
	cmd.lgr.Infof("Hello %v times", data.Amount)
	return fmt.Sprintf("Hello %v times", data.Amount), nil
}

// Register plugin
func init() {
	tarsier_plugins.RegisterPlugin(NAME, func() interface{} {
		return new(Dummy)
	})
}
