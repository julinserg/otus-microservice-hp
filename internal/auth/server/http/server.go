package auth_internalhttp

import (
	"context"
	"errors"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type SrvAuth interface {
	RequestTokenByCode(code string, chatId string) error
	GetToken(chatId string) (string, error)
	GetRequestAuthString() string
}

type Server struct {
	server   *http.Server
	logger   Logger
	endpoint string
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
	Warn(msg string)
}

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func NewServer(logger Logger, endpoint string, srvAuth SrvAuth) *Server {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    endpoint,
		Handler: loggingMiddleware(mux, logger),
	}
	initPrometheus()
	uh := csHandler{logger: logger, srvAuth: srvAuth}
	mux.HandleFunc("/api/v1/auth/health", hellowHandler)
	mux.HandleFunc("/api/v1/auth/auth", uh.authHandler)
	mux.HandleFunc("/api/v1/auth/token", uh.tokenHandler)
	mux.HandleFunc("/api/v1/auth/reqstring", uh.reqStringHandler)
	mux.Handle("/metrics", promhttp.Handler())
	return &Server{server, logger, endpoint}
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("http server started on " + s.endpoint)
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	s.logger.Info("http server stopped")
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
