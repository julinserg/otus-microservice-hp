package auth_internalhttp_private

import (
	"fmt"
	"net/http"
)

var htmlAuthError = `
<html>
	<head>
	</head>	
    <body>
		<h1>Authorization is Error:	%s</h1>
	</body>
</html>
`

func hellowHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("This is my auth service!"))
}

type csHandler struct {
	logger  Logger
	srvAuth SrvAuth
}

func (h *csHandler) returnError(w http.ResponseWriter, stringError string) {
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(fmt.Sprintf(htmlAuthError, stringError)))
}

func (h *csHandler) tokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.returnError(w, "Request is not GET type")
		return
	}
	chatId := r.URL.Query().Get("chat_id")
	token, err := h.srvAuth.GetToken(chatId)
	if err != nil {
		h.returnError(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}

func (h *csHandler) reqStringHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.returnError(w, "Request is not GET type")
		return
	}
	reqAuthStr := h.srvAuth.GetRequestAuthString()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(reqAuthStr))
}
