package server

import (
	"go-final/pkg/api"
	"log"
	"net/http"
	"time"
)

type Server struct {
	logger *log.Logger
	Srv    *http.Server
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf(
			"[%s] %s %s %s %s",
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			r.URL.RawQuery,
			time.Since(start),
		)
		next.ServeHTTP(w, r)
	})
}

func NewHttpServer(logger *log.Logger) *Server {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./web"))
	mux.Handle("/", fs)
	api.Init(mux)
	wrappedMux := loggingMiddleware(mux)

	server := Server{
		logger: logger,
		Srv: &http.Server{
			Addr:    ":7540",
			Handler: wrappedMux,
		},
	}
	return &server
}
