package auth_app

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
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

type SrvAuth struct {
	logger       Logger
	clientId     string
	clientSecret string
	tokenStorage map[string]string
}

type AuthTokenStruct struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

//const URI_AUTH_STR = "https://oauth.yandex.ru/authorize?response_type=code&client_id=%s&state=%d"

func New(logger Logger, clientId string, clientSecret string) *SrvAuth {
	return &SrvAuth{logger, clientId, clientSecret, make(map[string]string)}
}

func (sa *SrvAuth) RequestTokenByCode(code string, chatId string) error {
	if len(code) == 0 || len(chatId) == 0 {
		return fmt.Errorf("Error RequestTokenByCode: request params in null \n")
	}
	tokenRequestBody := "grant_type=authorization_code&code=" + code
	req, err := http.NewRequest(http.MethodPost, "https://oauth.yandex.ru/token", bytes.NewReader([]byte(tokenRequestBody)))
	if err != nil {
		return fmt.Errorf("client: could not create request: %s\n", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	sEnc := base64.StdEncoding.EncodeToString([]byte(sa.clientId + ":" + sa.clientSecret))
	req.Header.Set("Authorization", "Basic "+sEnc)

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("client: error making http request: %s\n", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("client: error read http body: %s\n", err)
	}
	tokenStr := &AuthTokenStruct{}
	err = json.Unmarshal(body, tokenStr)
	if err != nil {
		return fmt.Errorf("client: error parth http body: %s\n", err)
	}
	if len(tokenStr.AccessToken) == 0 {
		return fmt.Errorf("Error AccessToken: AccessToken in null \n")
	}
	sa.tokenStorage[chatId] = tokenStr.AccessToken
	return nil
}

func (sa *SrvAuth) GetToken(chatId string) (string, error) {
	if len(chatId) == 0 {
		return "", fmt.Errorf("Error GetToken: request params in null \n")
	}
	val, ok := sa.tokenStorage[chatId]
	if !ok {
		return "", fmt.Errorf("Error GetToken: no token for key %s \n", chatId)
	}
	return val, nil
}

func (sa *SrvAuth) GetRequestAuthString() string {
	return "https://oauth.yandex.ru/authorize?response_type=code&client_id=" + sa.clientId
}
