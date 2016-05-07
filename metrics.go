package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	Commands *prometheus.CounterVec
}

func NewMetrics() (*Metrics, error) {
	prom := &Metrics{}
	prom.Commands = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "tarsier_incoming_commands",
		Help: "total number of commands labeled with named and success flag",
	}, []string{"name", "successfull"})

	err := prometheus.Register(prom.Commands)
	return prom, err
}
