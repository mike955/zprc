package server

import (
	"io/ioutil"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	pb "github.com/mike955/zrpc/examples/grpc/api/example"
	"github.com/mike955/zrpc/examples/grpc/configs"
	"github.com/mike955/zrpc/examples/grpc/service"
	"github.com/mike955/zrpc/server/grpc"
	"gopkg.in/yaml.v2"
)

func NewGRPCServer() (server *grpc.Server) {
	config := configs.GlobalConfig.Server
	var opts = []grpc.ServerOption{
		grpc.Address(config.GRPCAddr),
		grpc.Timeout(config.Timeout),
		grpc.GrpcUnaryServerInterceptor(grpc_prometheus.UnaryServerInterceptor),
		grpc.GrpcDefaultUnaryServerInterceptor(),

		grpc.Prometheus(true, configs.GlobalConfig.Server.HttpAddr),
		grpc.Reflection(),
		grpc.HealthCheck(),
	}

	server = grpc.NewServer(config.AppName, opts...)
	log := server.Logger.WithField(map[string]interface{}{"app": config.AppName})
	s := service.NewExampleService(log)
	pb.RegisterExampleServer(server, s)
	return
}

func RunGRPCServer(server *grpc.Server) (err error) {
	err = server.Start()
	if err != nil {
		server.Logger.Errorf("server start error: %s", err.Error())
	}
	return
}

func InitConfig(conf string) {
	confData, err := ioutil.ReadFile(conf)
	if err != nil {
		panic("read config file error: " + err.Error())
	}
	if err := yaml.Unmarshal(confData, configs.GlobalConfig); err != nil {
		panic("parse config file error: " + err.Error())
	}
	// dao.Init(configs.GlobalConfig.Mysql)
}
