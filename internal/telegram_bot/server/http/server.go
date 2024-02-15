package telegram_bot_imitation_internalhttp

import (
	"context"
	"errors"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type SrvBot interface {
	GetAuthRequestString() (string, error)
	SendFileEvent(url string, chatId int64, isDebugMode bool) error
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

func NewServer(logger Logger, endpoint string, srvBot SrvBot) *Server {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    endpoint,
		Handler: loggingMiddleware(mux, logger),
	}
	initPrometheus()
	uh := csHandler{logger: logger, srvBot: srvBot}
	mux.HandleFunc("/api/v1/bot-imitation/health", hellowHandler)
	mux.HandleFunc("/api/v1/bot-imitation/file", uh.fileHandler)
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
