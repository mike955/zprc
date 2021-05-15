# zprc
zrpc is a go framework, it include cli tool， grpc and http framework.

## Contents

## Installation

1.install cli tool

```bash
go get -u github.com/mike955/zrpc/cmd/zrpc
```

2.install protobuf http plugin

```bash
go get -u github.com/mike955/zrpc/cmd/proto-gen-go-http
```

## Create a service

1.create a new service

```sh
zrpc new demo
```

2.enter service folder
```sh
cd demo
```

3.generate protobuf files
```sh
protoc --proto_path=./api/.  --go-grpc_out=./api/ --go-http_out=./api/ --go_out=./api/ demo.proto
```

## Update config

update global.yml

## Run service

```
go run cmd/demo/main.go -f global.yml
```