package service

import (
	"context"
	"errors"

	pb "github.com/mike955/zrpc/examples/grpc/api/example"
	"github.com/mike955/zrpc/examples/grpc/data"
	"github.com/mike955/zrpc/log"
)

type ExampleService struct {
	pb.UnimplementedExampleServer
	logger *log.Entry
	data   *data.ExampleData
}

func NewExampleService(logger *log.Entry) *ExampleService {
	return &ExampleService{
		logger: log.Helper(logger, map[string]interface{}{"service": "example"}),
		data:   data.NewExampleData(logger),
	}
}

func (s *ExampleService) Hello(ctx context.Context, request *pb.HelloRequest) (response *pb.HelloResponse, err error) {
	response = new(pb.HelloResponse)
	s.logger.Infof("func:Hello|request:%+v", request)
	res, err := s.data.Hello(request.Data)
	if err != nil {
		s.logger.Errorf("func:Hello|request:%+v|error:%s", request, err.Error())
		err = errors.New("request error")
		return
	}
	response.Data = res
	return
}
