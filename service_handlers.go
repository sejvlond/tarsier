package main

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/davecgh/go-spew/spew"
	"github.com/prometheus/client_golang/prometheus"
)

// ----------------------------------------------------------------------------

type HealthCheckHandler struct {
	lgr LOGGER
}

func (h *HealthCheckHandler) Init(lgr LOGGER, args ...interface{}) error {
	h.lgr = lgr
	if len(args) != 0 {
		return fmt.Errorf("Invalid arguments")
	}
	return nil
}

func (h *HealthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

// ----------------------------------------------------------------------------

type PrometheusHandler struct {
	lgr LOGGER
	ph  http.Handler
}

func (h *PrometheusHandler) Init(lgr LOGGER, args ...interface{}) error {
	h.lgr = lgr
	if len(args) != 0 {
		return fmt.Errorf("Invalid arguments")
	}
	h.ph = prometheus.Handler()
	return nil
}

func (h *PrometheusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.ph.ServeHTTP(w, r)
}

// ----------------------------------------------------------------------------

type CommandServiceHandler struct {
	lgr       LOGGER
	commander *Commander
}

func (h *CommandServiceHandler) Init(lgr LOGGER, args ...interface{}) error {
	h.lgr = lgr
	if len(args) != 1 {
		return fmt.Errorf("Invalid arguments")
	}
	var ok bool
	h.commander, ok = args[0].(*Commander)
	if !ok {
		return fmt.Errorf("Arg 0 is not *Commander")
	}
	return nil
}

func (h *CommandServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	var err error
	for _, command := range h.commander.Commands {
		if _, err = w.Write([]byte("<h1>")); err != nil {
			break
		}
		if _, err = w.Write([]byte(command.Name())); err != nil {
			break
		}
		if _, err = w.Write([]byte("</h1><p>")); err != nil {
			break
		}
		if _, err = w.Write([]byte(command.Description())); err != nil {
			break
		}
		if _, err = w.Write([]byte("</p>")); err != nil {
			break
		}
	}
	if err != nil {
		h.lgr.Errorf("Error writing to response writer: '%v'", err)
		http.Error(w, "Error writing to response writer",
			http.StatusInternalServerError)
		return
	}
}

// ----------------------------------------------------------------------------

type MemStatsHandler struct {
	lgr   LOGGER
	stats runtime.MemStats
}

func (h *MemStatsHandler) Init(lgr LOGGER, args ...interface{}) error {
	h.lgr = lgr
	if len(args) != 0 {
		return fmt.Errorf("Invalid arguments")
	}
	return nil
}

func (h *MemStatsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	var err error
	for {
		if _, err = w.Write([]byte("<h1>MemStats</h1><pre>")); err != nil {
			break
		}
		runtime.ReadMemStats(&h.stats)
		if _, err = w.Write([]byte(spew.Sdump(h.stats))); err != nil {
			break
		}
		if _, err = w.Write([]byte("</pre>")); err != nil {
			break
		}
		break
	}
	if err != nil {
		h.lgr.Errorf("Error writing to response writer: '%v'", err)
		http.Error(w, "Error writing to response writer",
			http.StatusInternalServerError)
		return
	}
}
