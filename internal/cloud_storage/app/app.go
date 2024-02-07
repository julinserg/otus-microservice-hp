package cloud_storage_app

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
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
	storageFolder  string
}

func New(logger Logger, uriAuthService string, ctx context.Context, debugToken string, storageFolder string) *SrvCloudStorage {
	return &SrvCloudStorage{logger, uriAuthService, ctx, debugToken, storageFolder}
}

func (s *SrvCloudStorage) DownloadAndSaveToStorage(fileEvent FileEvent) error {
	/*token, err := s.getToken(fileEvent.ChatId)
	if err != nil {
		return fmt.Errorf("Error receive token: " + err.Error())
	}*/
	yaDisk, err := yadisk.NewYaDisk(s.ctx, http.DefaultClient, &yadisk.Token{AccessToken: s.debugToken})
	if err != nil {
		return fmt.Errorf("Error NewYaDisk: " + err.Error())
	}
	//currentDate := time.Now().Format("DD-MM-YYYY")
	_, fileName := filepath.Split(fileEvent.URL)
	_, err = yaDisk.UploadExternalResource(s.storageFolder+"/"+fileName, fileEvent.URL, true, nil)
	if err != nil {
		return fmt.Errorf("Error UploadExternalResource: " + err.Error())
	}
	return nil
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
