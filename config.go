package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	Public struct {
		Port     uint   `yaml:"port"`
		UseHttps bool   `yaml:"use_https"`
		CertFile string `yaml:"cert_file"`
		KeyFile  string `yaml:"key_file"`
	} `yaml:"public"`
	Service struct {
		Port     uint   `yaml:"port"`
		UseHttps bool   `yaml:"use_https"`
		CertFile string `yaml:"cert_file"`
		KeyFile  string `yaml:"key_file"`
	} `yaml:"service"`
}

type LoggingConfig struct {
	Name     string `yaml:"name"`
	Interval uint   `yaml:"interval"`
	Dir      string `yaml:"dir"`
}

type ConsulConfig struct {
	Service         string `yaml:"service"`
	Url             string `yaml:"url"`
	Register        string `yaml:"register"`
	Deregister      string `yaml:"deregister"`
	List            string `yaml:"list"`
	CheckInterval   string `yaml:"check_interval"`
	RefreshInterval string `yaml:"refresh_interval"`
}

type Config struct {
	Server  ServerConfig           `yaml:"server"`
	Logging LoggingConfig          `yaml:"logging"`
	Plugins map[string]interface{} `yaml:"plugins"`
	Consul  ConsulConfig           `yaml:"consul"`
}

func NewConfig(filename string) (cfg *Config, err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	cfg = &Config{}
	if err = yaml.Unmarshal(data, cfg); err != nil {
		return
	}

	// validators --
	if cfg.Server.Public.Port <= 0 {
		err = fmt.Errorf("Invalid server port %q", cfg.Server.Public.Port)
		return
	}
	if cfg.Server.Service.Port <= 0 {
		err = fmt.Errorf("Invalid service port %q", cfg.Server.Service.Port)
		return
	}
	if cfg.Server.Service.Port == cfg.Server.Public.Port {
		err = fmt.Errorf("Server and Service port are the same")
		return
	}

	if cfg.Consul.Service == "" {
		err = fmt.Errorf("Consul service name could not be empty")
		return
	}
	if cfg.Consul.Url == "" {
		err = fmt.Errorf("Consul url could not be empty")
		return
	}
	if cfg.Consul.Register == "" {
		err = fmt.Errorf("Consul register url could not be empty")
		return
	}
	if cfg.Consul.Deregister == "" {
		err = fmt.Errorf("Consul deregister url could not be empty")
		return
	}
	if cfg.Consul.List == "" {
		err = fmt.Errorf("Consul list url could not be empty")
		return
	}
	if cfg.Consul.CheckInterval == "" {
		err = fmt.Errorf("Consul check interval could not be empty")
		return
	}
	if cfg.Consul.RefreshInterval == "" {
		err = fmt.Errorf("Consul refresh interval could not be empty")
		return
	}
	cfg.Consul.Register = cfg.Consul.Url + cfg.Consul.Register
	cfg.Consul.Deregister = cfg.Consul.Url + cfg.Consul.Deregister
	cfg.Consul.List = cfg.Consul.Url + cfg.Consul.List
	cfg.Consul.List = fmt.Sprintf(cfg.Consul.List, cfg.Consul.Service)

	return
}
