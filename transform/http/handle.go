package http

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const (
	Version = "0.0.1"
)

func GenerateAwesomeData(r *http.Request, log *logrus.Entry) (logger *logrus.Entry, ctx context.Context) {
	var x_real_ip, traceId, path, params, method string
	if r.Header.Get("X-Real-IP") != "" {
		x_real_ip = r.Header.Get("X-Real-IP")
	}
	if r.Header.Get("traceId") != "" {
		traceId = r.Header.Get("traceId")
	}
	path = r.URL.String()
	method = r.Method

	md := metadata.New(map[string]string{
		"X-Real-IP": x_real_ip,
		"traceId":   traceId,
	})
	logger = log.WithFields(logrus.Fields{
		"app":       "app",
		"x_real_ip": x_real_ip,
		"traceId":   traceId,
		"path":      path,
		"method":    method,
		"md":        md,
		"params":    params,
	})
	ctx = context.Background()
	ctx = context.WithValue(ctx, "logger", logger)
	ctx = context.WithValue(ctx, "x_real_ip", x_real_ip)
	ctx = context.WithValue(ctx, "traceId", traceId)
	ctx = context.WithValue(ctx, "md", md)
	return
}

func DecodeRequest(r *http.Request, contentType string, v interface{}) (err error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	switch contentType {
	case "application/protobuf", "application/x-protobuf":
		if err = proto.Unmarshal(body, v.(proto.Message)); err != nil {
			return
		}
	case "application/json":
		if err = protojson.Unmarshal(body, v.(proto.Message)); err != nil {
			return
		}
	default:
	}
	return
}

func EncodeResponse(w http.ResponseWriter, contentType string, v proto.Message) (err error) {
	var buf []byte
	switch contentType {
	case "application/protobuf", "application/x-protobuf":
		w.Header().Set("Content-Type", contentType)
		buf, err = proto.Marshal(v)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(buf)
	default:
		w.Header().Set("Content-Type", "application/json")
		buf, err = protojson.Marshal(v)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(buf)
	}
	return
}
