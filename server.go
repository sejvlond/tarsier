package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"sync"
)

type Server struct {
	lgr          LOGGER
	wg           *sync.WaitGroup
	shutdownChan chan struct{}
	cfg          *ServerConfig
	mux          *http.ServeMux
	serviceMux   *http.ServeMux
	metrics      *Metrics
	CallShutdown func()
}

type Handler interface {
	http.Handler
	Init(lgr LOGGER, args ...interface{}) error
}

func NewServer(lgr LOGGER, cfg *ServerConfig,
	commander *Commander, consul *Consul, shutdownFunc func(),
	shutdownChan chan struct{}, wg *sync.WaitGroup) (*Server, error) {

	server := &Server{
		lgr:          lgr,
		wg:           wg,
		shutdownChan: shutdownChan,
		cfg:          cfg,
		CallShutdown: shutdownFunc,
		mux:          http.NewServeMux(),
		serviceMux:   http.NewServeMux(),
		metrics:      NewMetrics(),
	}
	if err := server.InitHandlers(commander, consul); err != nil {
		return nil, err
	}
	return server, nil
}

func registerHandler(err *error, lgr LOGGER, mux *http.ServeMux,
	query string, handler Handler, args ...interface{}) bool {

	name := reflect.Indirect(reflect.ValueOf(handler)).Type().Name()
	lgr = lgr.WithField("handler", name)
	e := handler.Init(lgr, args...)
	if e != nil {
		*err = fmt.Errorf("Error initializing %v handler: '%v'", name, e)
		return false
	}
	mux.Handle(query, handler)
	return true
}

func (s *Server) InitHandlers(commander *Commander, consul *Consul) error {
	err := new(error)
	publicHandler := func(
		query string,
		handler Handler, args ...interface{}) bool {
		return registerHandler(err, s.lgr, s.mux, query, handler, args...)
	}
	serviceHandler := func(
		query string,
		handler Handler, args ...interface{}) bool {
		return registerHandler(err, s.lgr, s.serviceMux, query, handler, args...)
	}

	_ = true &&
		publicHandler(
			"/exec",
			new(CommandHandler), commander, s.metrics, consul) &&

		serviceHandler(
			"/metrics",
			new(PrometheusHandler)) &&
		serviceHandler(
			"/commands",
			new(CommandServiceHandler), commander) &&
		serviceHandler(
			"/mem_stats",
			new(MemStatsHandler)) &&
		serviceHandler(
			"/health_check",
			new(HealthCheckHandler))

	return *err
}

func (s *Server) runPublic() net.Listener {
	s.lgr.Infof("Starting public interface :%v", s.cfg.Public.Port)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", s.cfg.Public.Port))
	if err != nil {
		s.lgr.Errorf("Could not create listener %q", err)
		return nil
	}
	// https
	if s.cfg.Public.UseHttps {
		config := new(tls.Config)
		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0], err = tls.LoadX509KeyPair(
			s.cfg.Public.CertFile, s.cfg.Public.KeyFile)
		if err != nil {
			s.lgr.Errorf("Could not create TLS listener %q", err)
			return nil
		}
		listener = tls.NewListener(listener, config)
	}
	s.wg.Add(1)
	go func() {
		s.lgr.Infof("Starting to serve on %v...", s.cfg.Public.Port)
		err = http.Serve(listener, s.mux)
		s.lgr.Infof("Serving stopped with %q", err)
		s.wg.Done()
	}()
	return listener
}

func (s *Server) runService() net.Listener {
	s.lgr.Infof("Starting service interface :%v", s.cfg.Service.Port)

	listener, err := net.Listen("tcp",
		fmt.Sprintf(":%v", s.cfg.Service.Port))
	if err != nil {
		s.lgr.Errorf("Could not create service listener %q", err)
		return nil
	}
	// https
	if s.cfg.Service.UseHttps {
		config := new(tls.Config)
		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0], err = tls.LoadX509KeyPair(
			s.cfg.Service.CertFile, s.cfg.Service.KeyFile)
		if err != nil {
			s.lgr.Errorf("Could not create TLS listener %q", err)
			return nil
		}
		listener = tls.NewListener(listener, config)
	}
	s.wg.Add(1)
	go func() {
		s.lgr.Infof("Starting service to serve on %v...", s.cfg.Service.Port)
		err = http.Serve(listener, s.serviceMux)
		s.lgr.Infof("Serving service stopped with %q", err)
		s.wg.Done()
	}()
	return listener
}

func (s *Server) Run() {
	s.lgr.Infof("starting")
	listener := s.runPublic()
	serviceListener := s.runService()

	if listener == nil || serviceListener == nil {
		s.CallShutdown()
	} else {
		<-s.shutdownChan
	}

	if listener != nil {
		listener.Close()
	}
	if serviceListener != nil {
		serviceListener.Close()
	}
	s.wg.Done()
	s.lgr.Infof("stopped")
}
