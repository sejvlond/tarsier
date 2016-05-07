package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type Buddies []*Buddy

type Buddy struct {
	ID      string `json:"ServiceID"`
	Address string `json:"ServiceAddress"`
	Port    uint   `json:"ServicePort"`
}

type Consul struct {
	id           string
	ip           string
	cfg          *ConsulConfig
	serverCfg    *ServerConfig
	lgr          LOGGER
	wg           *sync.WaitGroup
	shutdownChan chan struct{}
	ticker       *time.Ticker
	buddies      Buddies
	mutex        sync.RWMutex
}

func NewConsul(lgr LOGGER, cfg *ConsulConfig,
	serverCfg *ServerConfig,
	shutdownChan chan struct{}, wg *sync.WaitGroup) (*Consul, error) {

	duration, err := time.ParseDuration(cfg.RefreshInterval)
	if err != nil {
		return nil, err
	}
	ip, err := externalIP()
	if err != nil {
		return nil, err
	}
	id := fmt.Sprintf("%s_%s:%v", cfg.Service, ip, serverCfg.Public.Port)
	c := &Consul{
		ip:           ip,
		id:           id,
		lgr:          lgr,
		cfg:          cfg,
		wg:           wg,
		serverCfg:    serverCfg,
		shutdownChan: shutdownChan,
		ticker:       time.NewTicker(duration),
	}
	c.cfg.Deregister = fmt.Sprintf(c.cfg.Deregister, id)
	err = c.Register()
	return c, err
}

func (c *Consul) Register() error {
	type Data struct {
		ID      string `json:"ID"`
		Name    string `json:"Name"`
		Address string `json:"Address"`
		Port    uint   `json:"Port"`
		Check   struct {
			HTTP     string `json:"HTTP"`
			Interval string `json:"Interval"`
		} `json:"Check"`
	}
	address := "://" + c.ip
	if c.serverCfg.Public.UseHttps {
		address = "https" + address
	} else {
		address = "http" + address
	}
	data := &Data{
		ID:      c.id,
		Name:    c.cfg.Service,
		Address: address,
		Port:    c.serverCfg.Public.Port,
	}
	checkUrl := fmt.Sprintf("://%s:%v/health_check", c.ip,
		c.serverCfg.Service.Port)
	if c.serverCfg.Service.UseHttps {
		checkUrl = "https" + checkUrl
	} else {
		checkUrl = "http" + checkUrl
	}
	data.Check.HTTP = checkUrl
	data.Check.Interval = c.cfg.CheckInterval

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	buf := bytes.NewReader(b)
	req, err := http.NewRequest("PUT", c.cfg.Register, buf)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			c.lgr.Errorf("Error while closing response body: '%v'", err)
		}
	}()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Register returns %v != 200", resp.StatusCode)
	}
	c.lgr.Infof("successfully registered\n%+v", data)
	return nil
}

func (c *Consul) Refresh() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.buddies = make(Buddies, 0)
	// obtain fresh and health buddies
	resp, err := http.Get(c.cfg.List)
	if err != nil {
		c.lgr.Errorf("Error refreshing buddies list: get '%v'", err)
		return
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			c.lgr.Errorf("Error while closing response body: '%v'", err)
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &c.buddies)
	if err != nil {
		c.lgr.Errorf("Error refreshing buddies list: body unmarshal '%v'", err)
		return
	}
	c.lgr.Infof("Buddies are fresh and healthy")
	for _, buddy := range c.buddies {
		c.lgr.Debugf("Buddy: %+v", buddy)
	}
}

func (c *Consul) Buddies() Buddies {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.buddies
}

func (c *Consul) RandomBuddies(amount int) Buddies {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if amount >= len(c.buddies) {
		return c.buddies
	}
	buddies := make([]*Buddy, amount)
	used := make(map[int]bool)
	var idx int
	for i := 0; i < amount; i++ {
		for { // generate random index from c.buddies that was not used before
			idx = rand.Intn(len(c.buddies))
			if _, exists := used[idx]; !exists {
				break
			}
		}
		used[idx] = true
		buddies[i] = c.buddies[idx]
	}
	return buddies
}

func (c *Consul) Run() {
	c.lgr.Infof("starting")
	run := true
	for run {
		select {
		case <-c.ticker.C:
			c.Refresh()
			break
		case <-c.shutdownChan:
			run = false
			break
		}
	}
	c.ticker.Stop()
	c.wg.Done()
	c.lgr.Infof("stopped")
}

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}
