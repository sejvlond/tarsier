package main

import (
	"fmt"
	"io"
	"os"
	"sync"
	"syscall"

	"github.com/sejvlond/kafkalog-logrus"
	"github.com/sejvlond/tarsier/plugins"

	"github.com/Sirupsen/logrus"
	"github.com/docopt/docopt-go"
	"gopkg.in/yaml.v2"
)

type LOGGER logrus.FieldLogger

type Tarsier struct {
	cfg          *Config
	lgr          LOGGER
	workerWG     sync.WaitGroup
	shutdownChan chan struct{}
	signalWG     sync.WaitGroup
	signalChan   chan os.Signal
	plugins      map[string]plugins.Plugin
	commander    *Commander
	consul       *Consul
}

func (t *Tarsier) LoadPlugins() error {
	t.plugins = make(map[string]plugins.Plugin)
	for name, factory := range plugins.AvailablePlugins {
		// create plugin
		plugin, ok := factory().(plugins.Plugin)
		if !ok {
			return fmt.Errorf("Unable to create '%v' plugin", name)
		}
		var config interface{}
		if pluginWithCustomCfg, ok := plugin.(plugins.HasConfigStruct); ok {
			data, err := yaml.Marshal(t.cfg.Plugins[name])
			if err != nil {
				return fmt.Errorf("Unable to marshal data '%v'", err)
			}
			config = pluginWithCustomCfg.ConfigStruct()
			err = yaml.Unmarshal(data, config)
			if err != nil {
				return fmt.Errorf(
					"Unable to unmarshal data to custom struct '%v'", err)
			}
		} else {
			config = t.cfg.Plugins[name]
		}
		// load its config section, or nil
		err := plugin.Init(t.lgr.WithFields(logrus.Fields{
			"name":   "PLUGIN",
			"plugin": name,
		}), config)
		if err != nil {
			return fmt.Errorf("Unable to initialize '%v' plugin: '%v'",
				name, err)
		}
		t.plugins[name] = plugin
		t.lgr.Infof("Loaded plugin '%v'", name)
	}
	return nil
}

func (t *Tarsier) ClosePlugins() {
	for name, plugin := range t.plugins {
		if pluginCloser, ok := plugin.(io.Closer); ok {
			if err := pluginCloser.Close(); err != nil {
				t.lgr.Errorf("Error while closing plugin %q: %q", name, err)
			}
		}
	}
}

func (t *Tarsier) InitCommander() (err error) {
	t.commander, err = NewCommander()
	if err != nil {
		return fmt.Errorf("Unable to initialize Commander '%v'", err)
	}
	for name, plugin := range t.plugins {
		pluginWithCmds, ok := plugin.(plugins.HasCommands)
		if !ok {
			continue
		}
		for _, command := range pluginWithCmds.Commands() {
			if err := t.commander.Register(command); err != nil {
				return fmt.Errorf("Unable to initialize commands from plugin "+
					"'%v': '%v'", name, err)
			}
		}
	}
	return nil
}

func (t *Tarsier) StartConsul() (err error) {
	t.consul, err = NewConsul(t.lgr.WithField("name", "CONSUL"),
		&t.cfg.Consul, &t.cfg.Server, t.shutdownChan, &t.workerWG)
	if err != nil {
		return fmt.Errorf("Unable to initialize Consul '%v'", err)
	}
	t.workerWG.Add(1)
	go t.consul.Run()
	return nil
}

func (t *Tarsier) StartSignalHandler() error {
	dispatch := map[os.Signal]SignalHandlerFunc{
		os.Interrupt:    t.ShutDown,
		os.Kill:         t.ShutDown,
		syscall.SIGTERM: t.ShutDown,
	}
	signal, err := NewSignalHandler(t.lgr.WithField("name", "SIGNAL HANDLER"),
		t.signalChan, &t.signalWG, dispatch)
	if err != nil {
		return err
	}
	t.signalWG.Add(1)
	go signal.Run()
	return nil
}

func (t *Tarsier) StartServer() error {
	server, err := NewServer(t.lgr.WithField("name", "SERVER"),
		&t.cfg.Server, t.commander, t.consul, t.ShutDown,
		t.shutdownChan, &t.workerWG)
	if err != nil {
		return err
	}
	t.workerWG.Add(1)
	go server.Run()
	return nil
}

func (t *Tarsier) Start() {
	t.shutdownChan = make(chan struct{})
	t.signalChan = make(chan os.Signal)
	if err := t.LoadPlugins(); err != nil {
		t.lgr.Errorf("Plugins initialization error: %q", err)
		goto shutdown
	}
	if err := t.InitCommander(); err != nil {
		t.lgr.Errorf("Commander initialization error: %q", err)
		goto shutdown
	}
	if err := t.StartSignalHandler(); err != nil {
		t.lgr.Errorf("Signal handler initialization error: %q", err)
		goto shutdown
	}
	if err := t.StartConsul(); err != nil {
		t.lgr.Errorf("Consul initialization error: %q", err)
		goto shutdown
	}
	if err := t.StartServer(); err != nil {
		t.lgr.Errorf("Server initialization error: %q", err)
		goto shutdown
	}
	goto stopping // validly here - skip shutdown
shutdown:
	t.ShutDown()
stopping:
	t.workerWG.Wait()
	return
}

func (t *Tarsier) ShutDown() {
	t.lgr.Infof("Shutdown initialized")
	close(t.shutdownChan)
}

func (t *Tarsier) Stop() {
	close(t.signalChan)
	t.signalWG.Wait()
	t.ClosePlugins()
	t.lgr.Infof("Bye")
}

func main() {
	lgr := logrus.Logger{
		Out:       os.Stderr,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}

	usage := `tarsier

Usage:
    tarsier -c <config_file>
    tarsier -h | --help

Options:
    -c --config         configuration file
    -h --help           Show this screen.`

	args, err := docopt.Parse(usage, nil, true, "", false)
	if err != nil {
		lgr.Fatalf("Error parsing arguments %q", err)
	}

	var cfg *Config
	if useCfg := args["--config"].(bool); useCfg {
		configFile, ok := args["<config_file>"].(string)
		if !ok {
			lgr.Fatalf("Error config argument %q", err)
		}
		cfg, err = NewConfig(configFile)
		if err != nil {
			lgr.Fatalf("Error loading config file %q", err)
		}
	} else {
		lgr.Fatalf("Config was not loaded")
	}

	kafkalog_hook, err := kafkalog_logrus.NewKafkalogHook(
		cfg.Logging.Name, cfg.Logging.Interval, cfg.Logging.Dir)
	if err != nil {
		lgr.Fatalf("Could not create kafkalog hook: '%v'", err)
	}
	lgr.Hooks.Add(kafkalog_hook)

	app := Tarsier{
		lgr: lgr.WithField("name", "MAIN"),
		cfg: cfg,
	}
	app.Start()
	app.Stop()
	return
}
