package main

import (
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/mike955/zrpc/examples/grpc/server"
)

var (
	version bool
	conf    string

	BuildTime       = ""
	GitCommitID     = ""
	GitCommitBranch = ""
	GoVersion       = runtime.Version()
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	flag.StringVar(&conf, "f", "", "-f <config>")
	flag.BoolVar(&version, "v", false, "-v")
	flag.Parse()
	if conf == "" {
		panic("not found config file, use: -f config.yaml")
	}
	server.InitConfig(conf)
}

func main() {
	if version == true {
		fmt.Println("BuildTime: ", BuildTime)
		fmt.Println("GitCommitID: ", GitCommitID)
		fmt.Println("GitCommitBranch: ", GitCommitBranch)
		fmt.Println("GoVersion: ", GoVersion)
		fmt.Println("GitCommitID: ", BuildTime)
		return
	}
	grpcServe := server.NewGRPCServer()
	if err := server.RunGRPCServer(grpcServe); err != nil {
		grpcServe.Logger.Errorf("grpc server run error:%s", err.Error())
	}
}
