package cloud_storage_debug_internalhttp

import (
	"net/http"
)

func hellowHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("This is my cloud storage debug service!"))
}

type csHandler struct {
	logger Logger
	srvCS  SrvCloudStorage
}

func (h *csHandler) returnError(w http.ResponseWriter, stringError string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(stringError))
}

func (h *csHandler) existFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.returnError(w, "Request is not GET type")
		return
	}
	name := r.URL.Query().Get("name")
	isExist, err := h.srvCS.CheckExistFile(name)
	if err != nil {
		h.returnError(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	if isExist {
		w.Write([]byte("YES"))
	} else {
		w.Write([]byte("NO"))
	}
}

func (h *csHandler) removeFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.returnError(w, "Request is not POST type")
		return
	}
	name := r.URL.Query().Get("name")
	err := h.srvCS.RemoveFile(name)
	if err != nil {
		h.returnError(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}
