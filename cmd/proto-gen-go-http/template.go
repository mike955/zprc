package main

import (
	"bytes"
	"html/template"
	"strings"
)

var httpTemplate = `
type {{.ServiceType}}HTTPServer interface {
{{range .MethodSets}}
	{{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error)
{{end}}
}

func New{{.ServiceType}}HTTPServer(srv {{.ServiceType}}Server, logger *logrus.Entry) http.Handler {
	r := mux.NewRouter()
	{{range .Methods}}
	r.HandleFunc("{{.Path}}", func(w http.ResponseWriter, r *http.Request) {
		var req {{.Request}}

		logger.Infof("receive http request")
		contentType, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
		err := http1.DecodeRequest(r, contentType, &req)
		start := time.Now()
		logger, ctx := http1.GenerateAwesomeData(r, logger)

		logger.Infof("receive http request")
		res, err := srv.{{.Name}}(ctx, &req)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		err = http1.EncodeResponse(w, contentType, res)
		logger = logger.WithFields(logrus.Fields{
			"cost": time.Now().Sub(start).Seconds(),
		})
		if err != nil {
			logger.Infof("http request failled | err: %s", err.Error())
		} else {
			logger.Infof("http request success")
		}
	}).Methods("{{.Method}}")
	{{end}}
	return r
}
`

type serviceDesc struct {
	ServiceType string
	ServiceName string
	Metadata    string
	Methods     []*methodDesc
	MethodSets  map[string]*methodDesc
}

type methodDesc struct {
	Name    string
	Request string
	Reply   string
	Path    string
	Method  string
}

func (s *serviceDesc) execute() string {
	s.MethodSets = make(map[string]*methodDesc)
	for _, m := range s.Methods {
		s.MethodSets[m.Name] = m
	}
	buf := new(bytes.Buffer)
	tmpl, err := template.New("http").Parse(strings.TrimSpace(httpTemplate))
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, s); err != nil {
		panic(err)
	}
	return string(buf.Bytes())
}
