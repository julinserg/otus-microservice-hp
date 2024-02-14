package cloud_storage_app

import (
	"bytes"
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
	ChatId      int64  `json:"chat_id"`
	URL         string `json:"url"`
	IsDebugMode bool   `json:"is_debug_mode"`
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
	file, err := s.downloadFile(fileEvent.URL)
	if err != nil {
		return fmt.Errorf("Error DownloadFile: " + err.Error())
	}
	token, err := s.getToken(fileEvent.ChatId, fileEvent.IsDebugMode)
	if err != nil {
		return fmt.Errorf("Error GetToken: " + err.Error())
	}
	yaDisk, err := yadisk.NewYaDisk(s.ctx, http.DefaultClient, &yadisk.Token{AccessToken: token})
	if err != nil {
		return fmt.Errorf("Error NewYaDisk: " + err.Error())
	}
	folder := s.getFolder()
	fileName := s.getFileName(fileEvent.URL)
	err = s.createFolder(&yaDisk, folder)
	if err != nil {
		return fmt.Errorf("Error CreateFolder: " + err.Error())
	}
	link, err := yaDisk.GetResourceUploadLink(folder+"/"+fileName, nil, true)
	if err != nil {
		return fmt.Errorf("Error GetResourceUploadLink: " + err.Error())
	}
	buffer := bytes.NewBuffer(file)
	_, err = yaDisk.PerformUpload(link, buffer)
	if err != nil {
		return fmt.Errorf("Error PerformUpload: " + err.Error())
	}
	return nil
}

func (s *SrvCloudStorage) CheckExistFile(name string) (bool, error) {
	token, err := s.getToken(0, true)
	if err != nil {
		return false, fmt.Errorf("Error GetToken: " + err.Error())
	}
	yaDisk, err := yadisk.NewYaDisk(s.ctx, http.DefaultClient, &yadisk.Token{AccessToken: token})
	if err != nil {
		return false, fmt.Errorf("Error NewYaDisk: " + err.Error())
	}
	res, err := yaDisk.GetResource(s.getFolder()+"/"+name, nil, 0, 0, false, "", "")
	if res != nil && len(res.ResourceID) != 0 {
		return true, nil
	}
	return false, nil
}

func (s *SrvCloudStorage) RemoveFile(name string) error {
	token, err := s.getToken(0, true)
	if err != nil {
		return fmt.Errorf("Error GetToken: " + err.Error())
	}
	yaDisk, err := yadisk.NewYaDisk(s.ctx, http.DefaultClient, &yadisk.Token{AccessToken: token})
	if err != nil {
		return fmt.Errorf("Error NewYaDisk: " + err.Error())
	}
	yaDisk.DeleteResource(s.getFolder()+"/"+name, nil, false, "", true)
	return nil
}

func (s *SrvCloudStorage) downloadFile(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("downloadFile: not create http request: %s\n", err)
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("downloadFile: error making http request: %s\n", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if len(body) == 0 {
		return nil, fmt.Errorf("downloadFile: empty file")
	}
	return body, nil
}

func (s *SrvCloudStorage) getFolder() string {
	currentDate := time.Now()
	folder := s.storageFolder + "/" + fmt.Sprintf("%02d-%02d-%d", currentDate.Day(), currentDate.Month(), currentDate.Year())
	return folder
}

func (s *SrvCloudStorage) getFileName(url string) string {
	_, fileName := filepath.Split(url)
	return fileName
}

func (s *SrvCloudStorage) createFolder(yaDisk *yadisk.YaDisk, folder string) error {
	var errorCF error
	for attempt := 0; attempt < 5; attempt++ {
		res, err := (*yaDisk).GetResource(folder, nil, 0, 0, false, "", "")
		if err != nil {
			s.logger.Error("Error GetResource: " + err.Error())
		}
		errorCF = err
		if res == nil {
			_, err = (*yaDisk).CreateResource(folder, nil)
			errorCF = err
			if err != nil {
				s.logger.Error("Error CreateResource: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		} else {
			return nil
		}
	}
	return fmt.Errorf("Error create folder on YDisk: %s", errorCF.Error())
}

func (s *SrvCloudStorage) getToken(chatId int64, isDebugMode bool) (string, error) {
	token := ""
	var err error
	if isDebugMode && len(s.debugToken) != 0 {
		token = s.debugToken
	} else {
		token, err = s.getTokenFromAuthService(chatId)
	}
	return token, err
}

func (s *SrvCloudStorage) getTokenFromAuthService(chatId int64) (string, error) {
	req, err := http.NewRequest(http.MethodGet, s.uriAuthService+fmt.Sprintf("/api/v1/auth/token?chat_id=%d", chatId), nil)
	if err != nil {
		return "", fmt.Errorf("getToken: not create http request: %s\n", err)
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("getToken: error making http request: %s\n", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), nil
}
