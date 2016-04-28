package plugins

var (
	AvailablePlugins = make(map[string]func() interface{})
)

// Adds a plugin to the set of usable plugins
func RegisterPlugin(name string, factory func() interface{}) {
	AvailablePlugins[name] = factory
}

type Command interface {
	Name() string
	Description() string
	// Execute command with its data (either raw yaml unmarshaled data, or
	// custom data struct, populated from request body)
	// returned string will be send to caller
	Execute(data interface{}) (string, error)
}

// Indicates a Command has a specific-to-itself data struct that should be
// passed in to plugin execute method
type HasDataStruct interface {
	// Returns a default-value-populated data structure into which
	// the request data will be deserialized
	DataStruct() interface{}
}

type Plugin interface {
	// Receives either raw yaml unmarshaled data, or custom config struct,
	// populated from the yaml config, and uses that data to initialize self
	Init(lgr Logger, config interface{}) error
}

type HasCommands interface {
	// Returns commands which can be executed by this plugin
	Commands() []Command
}

// Indicates a plugin has a specific-to-itself config struct that should be
// passed in to its Init method
type HasConfigStruct interface {
	// Returns a default-value-populated configuration structure into which
	// the yaml configuration will be deserialized
	ConfigStruct() interface{}
}

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})

	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warnln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
}
