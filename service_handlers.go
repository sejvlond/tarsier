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
	for _, command := range h.commander.Commands {
		w.Write([]byte("<h1>"))
		w.Write([]byte(command.Name()))
		w.Write([]byte("</h1><p>"))
		w.Write([]byte(command.Description()))
		w.Write([]byte("</p>"))
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
	w.Write([]byte("<h1>MemStats</h1><pre>"))
	runtime.ReadMemStats(&h.stats)
	w.Write([]byte(spew.Sdump(h.stats)))
	w.Write([]byte("</pre>"))
}
