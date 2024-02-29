package auth_internalhttp_private

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

type ServerPrivate struct {
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

func NewServer(logger Logger, endpoint string, srvAuth SrvAuth) *ServerPrivate {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    endpoint,
		Handler: loggingMiddleware(mux, logger),
	}
	initPrometheus()
	uh := csHandler{logger: logger, srvAuth: srvAuth}
	mux.HandleFunc("/api/v1/auth-private/health", hellowHandler)
	mux.HandleFunc("/api/v1/auth-private/token", uh.tokenHandler)
	mux.HandleFunc("/api/v1/auth-private/reqstring", uh.reqStringHandler)
	mux.Handle("/metrics", promhttp.Handler())
	return &ServerPrivate{server, logger, endpoint}
}

func (s *ServerPrivate) Start(ctx context.Context) error {
	s.logger.Info("http server started on " + s.endpoint)
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	s.logger.Info("http server stopped")
	return nil
}

func (s *ServerPrivate) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
