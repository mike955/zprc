package main

import (
	"fmt"

	"github.com/mike955/zrpc/server/tcp"
)

func main() {
	client, err := tcp.NewClient(tcp.Retry(1))
	if err != nil {
		fmt.Println("new error: ", err.Error())
	}
	data := []byte("hello")
	var res []byte
	for i := 0; i < 10000; i++ {
		res, err = client.Send(data)
		fmt.Println("send: ", string(data))
		fmt.Println("res: ", string(res))
	}
}
