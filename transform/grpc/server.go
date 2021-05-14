package grpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type ServerOption func(o *Server)

type Server struct {
	*grpc.Server
	app      string
	network  string
	address  string
	timeout  time.Duration
	grpcOpts []grpc.ServerOption

	Logger *logrus.Entry

	prometheusEnableHandlingTimeHistogram bool
	prometheusAddr                        string
	reflectionStatus                      bool
	healthCheckStatus                     bool
}

// Network with server network.
func Network(network string) ServerOption {
	return func(s *Server) {
		s.network = network
	}
}

// Address with server address.
func Address(addr string) ServerOption {
	return func(s *Server) {
		s.address = addr
	}
}

// Timeout with server timeout.
func Timeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
	}
}

func Logger(logger *logrus.Entry) ServerOption {
	return func(s *Server) {
		s.Logger = logger
	}
}

func Prometheus(enableHandlingTimeHistogram bool, prometheusAddr string) ServerOption {
	return func(s *Server) {
		s.prometheusEnableHandlingTimeHistogram = enableHandlingTimeHistogram
		s.prometheusAddr = prometheusAddr
	}
}

func Reflection() ServerOption {
	return func(s *Server) {
		s.reflectionStatus = true
	}
}

func HealthCheck() ServerOption {
	return func(s *Server) {
		s.healthCheckStatus = true
	}
}

func GrpcOpts(opts ...grpc.ServerOption) ServerOption {
	return func(s *Server) {
		// s.grpcOpts = append(s.grpcOpts, opts...)
		s.grpcOpts = opts
	}
}

func GrpcKeepAlive(kp keepalive.ServerParameters) ServerOption {
	return func(s *Server) {
		s.grpcOpts = append(s.grpcOpts, grpc.KeepaliveParams(kp))
	}
}

func GrpcUnaryServerInterceptor(interceptors ...grpc.UnaryServerInterceptor) ServerOption {
	return func(s *Server) {
		s.grpcOpts = append(s.grpcOpts, grpc.ChainUnaryInterceptor(interceptors...))
	}
}

func GrpcStreamServerInterceptor(interceptors ...grpc.StreamServerInterceptor) ServerOption {
	return func(s *Server) {
		s.grpcOpts = append(s.grpcOpts, grpc.ChainStreamInterceptor(interceptors...))
	}
}

func GrpcDefaultUnaryServerInterceptor() ServerOption {
	return func(s *Server) {
		s.grpcOpts = append(s.grpcOpts, defaultGrpcOpt(s))
	}
}

func NewServer(app string, opts ...ServerOption) *Server {
	srv := &Server{
		app:      app,
		network:  "tcp",
		address:  ":5080",
		timeout:  time.Second,
		Logger:   defaultLogger().WithFields(logrus.Fields{"app": app}),
		grpcOpts: []grpc.ServerOption{},
	}
	for _, o := range opts {
		o(srv)
	}
	srv.Server = grpc.NewServer(srv.grpcOpts...)
	srv.prometheus()
	srv.reflection()
	srv.healthCheck()
	return srv
}

func (s *Server) Start() error {
	lis, err := net.Listen(s.network, s.address)
	if err != nil {
		return err
	}
	if s.prometheusAddr != "" {
		go func() {
			s.Logger.Infof("http listening on %s", s.prometheusAddr)
			if err := http.ListenAndServe(s.prometheusAddr, promhttp.Handler()); err != nil {
				panic("prometheus start error: " + err.Error())
			}
		}()
	}
	go func() {
		s.handleGRPCServerSignals()
	}()
	s.Logger.Infof("grpc server listening on: %s", lis.Addr().String())
	return s.Server.Serve(lis)
}

// Stop stop the gRPC server.
func (s *Server) Stop() error {
	s.Server.GracefulStop()
	s.Logger.Info("grpc server stopping")
	return nil
}

func (s *Server) prometheus() {
	if s.prometheusEnableHandlingTimeHistogram {
		grpc_prometheus.EnableHandlingTimeHistogram()
	}
	if s.prometheusAddr != "" {
		grpc_prometheus.Register(s.Server)
	}
}

func (s *Server) reflection() {
	if s.reflectionStatus {
		reflection.Register(s.Server)
	}
}

func (s *Server) healthCheck() {
	if s.healthCheckStatus {
		grpc_health_v1.RegisterHealthServer(s.Server, health.NewServer())
	}
}

func (s *Server) handleGRPCServerSignals() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt) // stop process

	s.Logger.Info("listen grpc quit signal ...")
	select {
	case signal := <-signalCh:
		switch signal {
		case syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt:
			s.Logger.Infof("stopping grpc process on %s signal", fmt.Sprintf("%s", signal))
			if err := s.Stop(); err != nil {
				s.Logger.Errorf(fmt.Sprintf("quit grpc process error|error:%s", err.Error()))
				os.Exit(1)
			}
			s.Logger.Infof(fmt.Sprintf("quit grpc  process"))
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

func defaultGrpcOpt(s *Server) (opt grpc.ServerOption) {
	return grpc.ChainUnaryInterceptor(
		recoveryInterceptor(s.Logger),
		timeoutInterceptor(s.Logger),
		logInterceptor(s),
	)
}

func recoveryInterceptor(logger *logrus.Entry) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if rerr := recover(); rerr != nil {
				buf := make([]byte, 64<<10)
				n := runtime.Stack(buf, false)
				buf = buf[:n]
				logger.Errorf("recovery: %v: %+v\n%s\n", rerr, req, buf)
				// add err handle
			}
		}()
		return handler(ctx, req)
	}
}

func timeoutInterceptor(logger *logrus.Entry) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Second*60)
		defer cancel()
		return handler(ctx, req)
	}
}

func logInterceptor(s *Server) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()
		var x_real_ip, traceId, path, params, method string
		var md metadata.MD
		var ok bool

		md, ok = metadata.FromIncomingContext(ctx)
		if ok {
			if len(md.Get("X-Real-IP")) > 0 {
				x_real_ip = md.Get("X-Real-IP")[0]
			}
			if len(md.Get("traceId")) > 0 {
				traceId = md.Get("traceId")[0]
			} else {
				traceId = "no-id"
			}
		}
		path = info.FullMethod
		params = req.(fmt.Stringer).String()
		method = "POST"
		logger := s.Logger.WithFields(logrus.Fields{
			"app":       s.app,
			"x_real_ip": x_real_ip,
			"traceId":   traceId,
			"path":      path,
			"method":    method,
			"md":        md,
			"params":    params,
		})
		logger.Infof("receive grpc request")
		ctx = context.WithValue(ctx, "logger", logger)
		ctx = context.WithValue(ctx, "x_real_ip", x_real_ip)
		ctx = context.WithValue(ctx, "traceId", traceId)
		ctx = context.WithValue(ctx, "md", md)
		resp, err = handler(ctx, req)
		logger = logger.WithFields(logrus.Fields{
			"cost": time.Now().Sub(start).Seconds(),
		})
		if err != nil {
			logger.Infof("grpc request failled | err: %s", err.Error())
		} else {
			logger.Infof("grpc request success")
		}
		return
	}
}
