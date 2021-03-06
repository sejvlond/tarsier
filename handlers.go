package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"gopkg.in/yaml.v2"
)

type CommandHandler struct {
	lgr       LOGGER
	commander *Commander
	consul    *Consul
	metrics   *Metrics
	client    *http.Client
}

func (h *CommandHandler) Init(lgr LOGGER, args ...interface{}) error {
	h.lgr = lgr
	if len(args) != 3 {
		return fmt.Errorf("Invalid arguments")
	}
	var ok bool
	h.commander, ok = args[0].(*Commander)
	if !ok {
		return fmt.Errorf("Arg 0 is not *Commander")
	}
	h.metrics, ok = args[1].(*Metrics)
	if !ok {
		return fmt.Errorf("Arg 1 is not *Metrics")
	}
	h.consul, ok = args[2].(*Consul)
	if !ok {
		return fmt.Errorf("Arg 2 is not *Consul")
	}

	tr := &http.Transport{
		DisableKeepAlives: true,
	}
	h.client = &http.Client{Transport: tr}

	return nil
}

type CommandRequest struct {
	Command string
	Data    interface{}
	Wave    struct {
		Remains uint
		Buddies uint
	}
}

func (h *CommandHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(r.Body); err != nil {
		h.lgr.Errorf("Error reading from request body: '%v'", err)
		http.Error(w, "Invalid request body", http.StatusInternalServerError)
	}
	defer func() {
		err := r.Body.Close()
		if err != nil {
			h.lgr.Errorf("Error while closing request body: '%v'", err)
		}
	}()
	cmdReq := CommandRequest{}
	err := yaml.Unmarshal(buf.Bytes(), &cmdReq)
	if err != nil {
		h.lgr.Warnf("Cannot unmarshal request body '%v'", err)
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go h.wave(w, &cmdReq, &wg)
	wg.Add(1)
	go h.exec(w, &cmdReq, &wg)
	wg.Wait()
}

func (h *CommandHandler) wave(w http.ResponseWriter, cmdReq *CommandRequest,
	wg *sync.WaitGroup) {

	defer wg.Done()
	if cmdReq.Wave.Remains <= 0 {
		return
	}
	cmdReq.Wave.Remains--

	b, err := yaml.Marshal(cmdReq)
	if err != nil {
		h.lgr.Errorf("Error marshaling data for buddies: '%v'", err)
		return
	}

	bCnt := int(cmdReq.Wave.Buddies)
	for bCnt > 0 {
		buddies := h.consul.RandomBuddies(bCnt)
		if len(buddies) == 0 {
			h.lgr.Warnf("There are no buddies to send")
			return // no buddies? :(
		}
		bCnt -= len(buddies)
		for _, buddy := range buddies {
			// send wave in goroutine and do not wait for it
			wg.Add(1)
			go func() {
				defer wg.Done()
				buf := bytes.NewReader(b)
				resp, err := h.client.Post(
					fmt.Sprintf("%v:%v/exec", buddy.Address, buddy.Port),
					"application/x-yaml", buf)
				if err != nil {
					h.lgr.Errorf("Could not send wave to buddy %v:%v '%v'",
						buddy.Address, buddy.Port, err)
					return
				}
				data, _ := ioutil.ReadAll(resp.Body)
				err = resp.Body.Close()
				if err != nil {
					h.lgr.Errorf("Error while closing response body: '%v'", err)
				}
				h.lgr.Infof("Response from buddy %v:%v is '%s'",
					buddy.Address, buddy.Port, string(data))

				_, err = w.Write([]byte(fmt.Sprintf(
					"Buddy %v:%v says: %v", buddy.Address, buddy.Port,
					string(data))))
				if err != nil {
					h.lgr.Errorf("Writing command output from buddy %v:%v"+
						"failed '%v'", buddy.Address, buddy.Port, err)
				}
			}()
		}
	}
}

func (h *CommandHandler) exec(w http.ResponseWriter, cmdReq *CommandRequest,
	wg *sync.WaitGroup) {

	defer wg.Done()
	var out string
	var err error
	if out, err = h.commander.Execute(cmdReq.Command, cmdReq.Data); err != nil {
		if err == CommandNotRegistered {
			h.lgr.Warnf("No plugin was registered for command '%v'",
				cmdReq.Command)
			http.Error(w, "Unknown command", http.StatusNotImplemented)
		} else {
			h.lgr.Warnf("Executing command failed '%v'", err)
			http.Error(w, fmt.Sprintf("Command failed: %v", out),
				http.StatusInternalServerError)
		}
		h.metrics.Commands.WithLabelValues(cmdReq.Command, "false").Inc()
		return
	}
	h.metrics.Commands.WithLabelValues(cmdReq.Command, "true").Inc()
	_, err = w.Write([]byte(out + "\n"))
	if err != nil {
		h.lgr.Errorf("Writing command output '%v' failed '%v'", out, err)
		http.Error(w, fmt.Sprintf("Writing command output '%v' failed '%v'",
			out, err), http.StatusInternalServerError)
	}
	h.lgr.Infof("Command '%v' successfully executed", cmdReq.Command)
}
