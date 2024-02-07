package cloud_storage_app

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	yadisk "github.com/nikitaksv/yandex-disk-sdk-go"
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

type SrvCloudStorage struct {
	logger         Logger
	uriAuthService string
	ctx            context.Context
	debugToken     string
}

func New(logger Logger, uriAuthService string, ctx context.Context, debugToken string) *SrvCloudStorage {
	return &SrvCloudStorage{logger, uriAuthService, ctx, debugToken}
}

func (s *SrvCloudStorage) DownloadAndSaveToStorage(fileEvent FileEvent) error {
	/*token, err := s.getToken(fileEvent.ChatId)
	if err != nil {
		return fmt.Errorf("Error receive token: " + err.Error())
	}*/
	/*file, err := s.downloadFile(fileEvent.URL)
	if err != nil {
		return fmt.Errorf("Error download file: " + err.Error())
	}*/
	yaDisk, err := yadisk.NewYaDisk(s.ctx, http.DefaultClient, &yadisk.Token{AccessToken: s.debugToken})
	if err != nil {
		return fmt.Errorf("Error connect to YaDisk: " + err.Error())
	}
	disk, err := yaDisk.GetDisk([]string{})
	if err != nil {
		return fmt.Errorf("Error get disk info: " + err.Error())
	}
	fmt.Println("TotalSpace", disk.TotalSpace)
	fmt.Println("UsedSpace", disk.UsedSpace)
	return nil
}

func (s *SrvCloudStorage) downloadFile(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("client: not create http request: %s\n", err)
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client: error making http request: %s\n", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, nil
}

func (s *SrvCloudStorage) getToken(chatId int64) (string, error) {
	req, err := http.NewRequest(http.MethodGet, s.uriAuthService+fmt.Sprintf("/token?chat_id=%d", chatId), nil)
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
