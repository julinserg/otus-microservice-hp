package auth_internalhttp_public

import (
	"context"
	"errors"
	"net/http"
)

type SrvAuth interface {
	RequestTokenByCode(code string, chatId string) error
	GetToken(chatId string) (string, error)
	GetRequestAuthString() string
}

type ServerPublic struct {
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

func NewServer(logger Logger, endpoint string, srvAuth SrvAuth) *ServerPublic {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    endpoint,
		Handler: loggingMiddleware(mux, logger),
	}
	uh := csHandler{logger: logger, srvAuth: srvAuth}
	mux.HandleFunc("/api/v1/auth-public/auth", uh.authHandler)
	return &ServerPublic{server, logger, endpoint}
}

func (s *ServerPublic) Start(ctx context.Context) error {
	s.logger.Info("http server started on " + s.endpoint)
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	s.logger.Info("http server stopped")
	return nil
}

func (s *ServerPublic) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
