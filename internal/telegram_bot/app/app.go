package telegram_bot_app

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
	Warn(msg string)
}

type SrvBot struct {
	logger         Logger
	uriAuthService string
}

func New(logger Logger, uriAuthService string) *SrvBot {
	return &SrvBot{logger, uriAuthService}
}

func (s *SrvBot) GetAuthRequestString() (string, error) {
	req, err := http.NewRequest(http.MethodGet, s.uriAuthService+"/reqstring", nil)
	if err != nil {
		return "", fmt.Errorf("client: could not create request: %s\n", err)
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("client: error making http request: %s\n", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), nil
}
