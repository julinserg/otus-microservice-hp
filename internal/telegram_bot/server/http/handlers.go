package telegram_bot_imitation_internalhttp

import (
	"net/http"
	"strconv"
)

func hellowHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("This is my telegram bot imitation service!"))
}

type csHandler struct {
	logger Logger
	srvBot SrvBot
}

func (h *csHandler) returnError(w http.ResponseWriter, stringError string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(stringError))
}

func (h *csHandler) fileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.returnError(w, "Request is not POST type")
		return
	}
	url := r.URL.Query().Get("url")
	chatId := r.URL.Query().Get("chat_id")
	testMode := r.URL.Query().Get("test_mode")
	chatIdInt64, err := strconv.ParseInt(chatId, 10, 64)
	if err != nil {
		h.returnError(w, err.Error())
		return
	}
	err = h.srvBot.SendFileEvent(url, chatIdInt64, testMode)
	if err != nil {
		h.returnError(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}
