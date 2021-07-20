package main

import (
	"fmt"

	"github.com/mike955/zrpc/log"
	"github.com/mike955/zrpc/server/tcp"
)

const (
	tcpSize = 4 // byte
)

func main() {
	appName := "tcpExample"
	logger := log.NewLogger()
	server, err := tcp.NewServer(
		appName,
		// tcp.HeadSize(tcpSize),
		tcp.Logger(logger),
		tcp.Handler(handler()),
	)
	if err != nil {
		panic("create tcp server error")
	}
	server.Serve()
}

func handler() tcp.TcpHandler {
	return func(req []byte) (res []byte) {
		fmt.Println("receive", string(req))
		res = []byte("response")
		return
	}
}
