package cloud_storage_debug_internalhttp

import (
	"context"
	"errors"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type SrvCloudStorage interface {
	CheckExistFile(name string) (bool, error)
	RemoveFile(name string) error
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

func NewServer(logger Logger, endpoint string, srvCS SrvCloudStorage) *Server {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    endpoint,
		Handler: loggingMiddleware(mux, logger),
	}
	initPrometheus()
	uh := csHandler{logger: logger, srvCS: srvCS}
	mux.HandleFunc("/api/v1/cs-debug/health", hellowHandler)
	mux.HandleFunc("/api/v1/cs-debug/exist", uh.existFileHandler)
	mux.HandleFunc("/api/v1/cs-debug/remove", uh.removeFileHandler)
	mux.Handle("/metrics", promhttp.Handler())
	return &Server{server, logger, endpoint}
}

func (s *Server) Start() error {
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
