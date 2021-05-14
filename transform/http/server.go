package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

type ServerOption func(o *Server)

type Server struct {
	*http.Server
	Logger *logrus.Entry

	app                                   string
	prometheusEnableHandlingTimeHistogram bool
}

// Address with server address.
func Address(addr string) ServerOption {
	return func(s *Server) {
		s.Server.Addr = addr
	}
}

// Timeout with server timeout.
func ReadTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.Server.ReadTimeout = timeout
	}
}

func Logger(logger *logrus.Entry) ServerOption {
	return func(s *Server) {
		s.Logger = logger
	}
}

func (s *Server) healthCheck() {
}

func Prometheus(enableHandlingTimeHistogram bool) ServerOption {
	return func(s *Server) {
		s.prometheusEnableHandlingTimeHistogram = enableHandlingTimeHistogram
	}
}

func NewServer(app string, opts ...ServerOption) *Server {
	srv := &Server{
		app:    app,
		Server: &http.Server{},
		Logger: defaultLogger().WithFields(logrus.Fields{"app": app}),
	}
	for _, o := range opts {
		o(srv)
	}
	return srv
}

func (s *Server) Run() error {
	if s.prometheusEnableHandlingTimeHistogram {

	}
	go func() {
		s.handleHTTPServerSignals()
	}()
	s.Logger.Infof("http server listening on: %s", s.Server.Addr)
	return s.Server.ListenAndServe()
}

// Stop stop the http server.
func (s *Server) Stop() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s.Shutdown(ctx)
	s.Logger.Info("http server stopping")
	return nil
}

func (s *Server) SetHandler(handler http.Handler) {
	s.Handler = handler
}

func (s *Server) handleHTTPServerSignals() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt) // stop process

	s.Logger.Info("listen http quit signal ...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	select {
	case signal := <-signalCh:
		switch signal {
		case syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt:
			s.Logger.Infof("stopping http process on %s signal", fmt.Sprintf("%s", signal))
			if err := s.Shutdown(ctx); err != nil {
				s.Logger.Errorf(fmt.Sprintf("quit http process error|error:%s", err.Error()))
				os.Exit(1)
			}
			s.Logger.Infof(fmt.Sprintf("quit http process success"))
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}
}

func defaultLogger() (logger *logrus.Logger) {
	logger = logrus.New()
	logger.Out = os.Stdout
	logger.Formatter = &logrus.JSONFormatter{}
	return
}
