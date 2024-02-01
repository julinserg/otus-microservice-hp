package cloud_storage_internalhttp

import (
	"net/http"
)

func hellowHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("This is my cloud storage service!"))
}

type csHandler struct {
	logger Logger
}

func (h *csHandler) authHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("auth is run!!!!!!!!!!! " + r.RequestURI)
	code := r.URL.Query().Get("code")
	h.logger.Info("code: " + code)
}
