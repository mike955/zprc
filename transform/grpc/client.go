package grpc

// import (
// 	"context"
// 	"os"
// 	"sync"
// 	"time"

// 	age_pb "github.com/mike955/zebra/api/age"
// 	bank_pb "github.com/mike955/zebra/api/bank"
// 	cellphone_pb "github.com/mike955/zebra/api/cellphone"
// 	email_pb "github.com/mike955/zebra/api/email"
// 	flake_pb "github.com/mike955/zebra/api/flake"
// 	"github.com/sirupsen/logrus"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/metadata"
// )

// var gRPCClientMap sync.Map

// func NewFlakeRpc(flakeAddr string) (flake_pb.FlakeServiceClient, error) {
// 	if _, ok := gRPCClientMap.Load("flake"); !ok {
// 		if os.Getenv("FlAKE_ADDR") != "" {
// 			flakeAddr = os.Getenv("FlAKE_ADDR")
// 		}
// 		conn, err := grpc.DialContext(context.Background(), flakeAddr, grpc.WithUnaryInterceptor(CllentUnaryInterceptor), grpc.WithStreamInterceptor(CllentStreamInterceptor), grpc.WithInsecure())
// 		if err != nil {
// 			return nil, err
// 		}
// 		gRPCClientMap.Store("flake", flake_pb.NewFlakeServiceClient(conn))
// 	}
// 	client, _ := gRPCClientMap.Load("flake")
// 	return client.(flake_pb.FlakeServiceClient), nil
// }

// func NewAgeRpc(ageAddr string) (age_pb.AgeServiceClient, error) {
// 	if _, ok := gRPCClientMap.Load("age"); !ok {
// 		if os.Getenv("AGE_ADDR") != "" {
// 			ageAddr = os.Getenv("AGE_ADDR")
// 		}
// 		conn, err := grpc.DialContext(context.Background(), ageAddr, grpc.WithUnaryInterceptor(CllentUnaryInterceptor), grpc.WithStreamInterceptor(CllentStreamInterceptor), grpc.WithInsecure())
// 		if err != nil {
// 			return nil, err
// 		}
// 		gRPCClientMap.Store("age", age_pb.NewAgeServiceClient(conn))
// 	}
// 	client, _ := gRPCClientMap.Load("age")
// 	return client.(age_pb.AgeServiceClient), nil
// }

// func NewEmailRpc(emailAddr string) (email_pb.EmailServiceClient, error) {
// 	if _, ok := gRPCClientMap.Load("email"); !ok {
// 		if os.Getenv("EMAIL_ADDR") != "" {
// 			emailAddr = os.Getenv("EMAIL_ADDR")
// 		}
// 		conn, err := grpc.DialContext(context.Background(), emailAddr, grpc.WithUnaryInterceptor(CllentUnaryInterceptor), grpc.WithStreamInterceptor(CllentStreamInterceptor), grpc.WithInsecure())
// 		if err != nil {
// 			return nil, err
// 		}
// 		gRPCClientMap.Store("email", email_pb.NewEmailServiceClient(conn))
// 	}
// 	client, _ := gRPCClientMap.Load("email")
// 	return client.(email_pb.EmailServiceClient), nil
// }

// func NewBankRpc(bankAddr string) (bank_pb.BankServiceClient, error) {
// 	if _, ok := gRPCClientMap.Load("bank"); !ok {
// 		if os.Getenv("BANK_ADDR") != "" {
// 			bankAddr = os.Getenv("BANK_ADDR")
// 		}
// 		conn, err := grpc.DialContext(context.Background(), bankAddr, grpc.WithUnaryInterceptor(CllentUnaryInterceptor), grpc.WithStreamInterceptor(CllentStreamInterceptor), grpc.WithInsecure())
// 		if err != nil {
// 			return nil, err
// 		}
// 		gRPCClientMap.Store("bank", bank_pb.NewBankServiceClient(conn))
// 	}
// 	client, _ := gRPCClientMap.Load("bank")
// 	return client.(bank_pb.BankServiceClient), nil
// }

// func NewCellphoneRpc(cellphoneAddr string) (cellphone_pb.CellphoneServiceClient, error) {
// 	if _, ok := gRPCClientMap.Load("cellphone"); !ok {
// 		if os.Getenv("CELLPHONE_ADDR") != "" {
// 			cellphoneAddr = os.Getenv("CELLPHONE_ADDR")
// 		}
// 		conn, err := grpc.DialContext(context.Background(), cellphoneAddr, grpc.WithUnaryInterceptor(CllentUnaryInterceptor), grpc.WithStreamInterceptor(CllentStreamInterceptor), grpc.WithInsecure())
// 		if err != nil {
// 			return nil, err
// 		}
// 		gRPCClientMap.Store("cellphone", cellphone_pb.NewCellphoneServiceClient(conn))
// 	}
// 	client, _ := gRPCClientMap.Load("cellphone")
// 	return client.(cellphone_pb.CellphoneServiceClient), nil
// }

// func CllentUnaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
// 	ctx = metadata.AppendToOutgoingContext(ctx, "traceId", ctx.Value("traceId").(string))
// 	ctx = metadata.AppendToOutgoingContext(ctx, "x_real_ip", ctx.Value("x_real_ip").(string))
// 	start := time.Now()
// 	err := invoker(ctx, method, req, reply, cc, opts...)
// 	end := time.Now()
// 	logger := ctx.Value("logger").(*logrus.Entry)
// 	logger.Infof("info: rpc call,method: %s start time: %s, end time: %s, err: %v", method, start.Format("Basic"), end.Format(time.RFC3339), err)
// 	return err
// }

// func CllentStreamInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
// 	ctx = metadata.AppendToOutgoingContext(ctx, "traceId", ctx.Value("traceId").(string))
// 	ctx = metadata.AppendToOutgoingContext(ctx, "x_real_ip", ctx.Value("x_real_ip").(string))
// 	logger := ctx.Value("logger").(*logrus.Entry)
// 	start := time.Now()
// 	s, err := streamer(ctx, desc, cc, method, opts...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return newWrappedStream(ctx, s, start, logger), nil
// }

// type wrappedStream struct {
// 	ctx context.Context
// 	grpc.ClientStream
// 	start  time.Time
// 	logger *logrus.Entry
// }

// func (w *wrappedStream) RecvMsg(m interface{}) error {
// 	w.logger.Infof("Receive a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
// 	return w.ClientStream.RecvMsg(m)
// }

// func (w *wrappedStream) SendMsg(m interface{}) error {
// 	w.logger.Infof("Send a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
// 	return w.ClientStream.SendMsg(m)
// }

// func newWrappedStream(ctx context.Context, s grpc.ClientStream, start time.Time, logger *logrus.Entry) grpc.ClientStream {
// 	return &wrappedStream{ctx, s, start, logger}
// }
