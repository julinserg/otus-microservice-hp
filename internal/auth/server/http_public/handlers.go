package auth_internalhttp_public

import (
	"fmt"
	"net/http"
)

var htmlAuthOk = `
<html>
	<head>
	</head>	
    <body>
		<h1>Authorization is OK! Welcome to CloudStorageYSBot.</h1>	
	</body>
</html>
`

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

func (h *csHandler) authHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.returnError(w, "Request is not GET type")
		return
	}
	code := r.URL.Query().Get("code")
	chatId := r.URL.Query().Get("state")
	err := h.srvAuth.RequestTokenByCode(code, chatId)
	if err != nil {
		h.returnError(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlAuthOk))
}
