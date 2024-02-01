package cloud_storage_internalhttp

import (
	"net/http"
)

func hellowHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("This is my cloud storage service!"))
}

type csHandler struct {
	logger       Logger
	clientSecret string
}

func (h *csHandler) authHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("auth is run!!!!!!!!!!! " + r.RequestURI)
	code := r.URL.Query().Get("code")
	h.logger.Info("code: " + code)

	/*req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}*/

	w.WriteHeader(http.StatusOK)
}
