package http

import "net/http"

type CodeMap struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type Response struct {
	CodeMap
	TotalCount int64       `json:"TotalCount,omitempty"`
	Data       interface{} `json:"Data,omitempty"`
}

type ResponseWriter struct {
	CodeMap map[string]CodeMap
	Action  string
	http.ResponseWriter
}
