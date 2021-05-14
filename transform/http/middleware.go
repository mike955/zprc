package http

import "net/http"

func HealthCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func Recovery(next http.Handler) http.Handler {

}

func Limit(next http.Handler) http.Handler {

}

func Timeout(next http.Handler) http.Handler {

}

func Logger(next http.Handler) http.Handler {

}
