package configs

import "time"

var GlobalConfig = &Global{}

type Global struct {
	Server server
}
type server struct {
	AppName  string        `yaml:"app_name"`
	GRPCAddr string        `yaml:"grpc_addr"`
	HttpAddr string        `yaml:"http_addr"`
	Timeout  time.Duration `yaml:"timeout"`
}
