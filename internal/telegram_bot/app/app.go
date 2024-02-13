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

type FileEvent struct {
	ChatId int64  `json:"chat_id"`
	URL    string `json:"url"`
}

type BotMQ interface {
	PublishFileEvent(event FileEvent) error
}

type SrvBot struct {
	logger         Logger
	uriAuthService string
	botMQ          BotMQ
}

func New(logger Logger, uriAuthService string, botMQ BotMQ) *SrvBot {
	return &SrvBot{logger, uriAuthService, botMQ}
}

func (s *SrvBot) GetAuthRequestString() (string, error) {
	req, err := http.NewRequest(http.MethodGet, s.uriAuthService+"/api/v1/auth/reqstring", nil)
	if err != nil {
		return "", fmt.Errorf("client: not create http request: %s\n", err)
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

func (s *SrvBot) SendFileEvent(url string, chatId int64) error {
	err := s.botMQ.PublishFileEvent(FileEvent{chatId, url})
	return err
}
