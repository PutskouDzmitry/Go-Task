package telecom

import (
	"MommyCO/internal/config"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"regexp"
)

const PhoneCodeLen = 4

type telecom struct {
	login    string `yaml:"login"`
	password string `yaml:"password"`
	urlCall  string `yaml:"url-call"`
	urlSend  string `yaml:"url-send"`
}

func newTelecom(cfg *config.Config) *telecom {
	return &telecom{
		login:    cfg.Sms.Login,
		password: cfg.Sms.Password,
		urlCall:  cfg.Sms.URLCall,
		urlSend:  cfg.Sms.URLSend,
	}
}

type Sms struct {
	telecom *telecom
	log     *zap.Logger
}

func NewTelecomService(cfg *config.Config, log *zap.Logger) *Sms {
	return &Sms{
		telecom: newTelecom(cfg),
		log:     log,
	}
}

func (ts *Sms) Auth() (bool, error) { return true, nil }

func (ts *Sms) SendCode(phone string) (string, error) {
	phoneMatch, err := regexp.Match("^(7)\\d{10}$", []byte(phone))
	if err != nil {
		ts.log.Error("failed to match phone", zap.Error(err))
		return "", err
	}
	if !phoneMatch {
		ts.log.Error("incorrect phone")
		return "", errors.New("incorrect phone")
	}

	params := url.Values{}
	params.Add("user", ts.telecom.login)
	params.Add("pwd", ts.telecom.password)
	params.Add("name_deliver", "MamaCO")
	params.Add("sadr", "MamaCO")
	params.Add("dadr", phone)
	code := generateCode(PhoneCodeLen)
	params.Add("text", code)

	u := ts.telecom.urlSend + params.Encode()

	ts.log.Debug("send code", zap.String("u", u))

	response, err := http.Get(u)
	if err != nil {
		ts.log.Error("failed to send code", zap.Error(err))
		return "", err
	}

	ts.log.Debug("response status", zap.String("status", response.Status))

	return code, nil
}

func (ts *Sms) CallCode(phone string) (string, error) {
	phoneMatch, err := regexp.Match("^(7)\\d{10}$", []byte(phone))
	if err != nil {
		ts.log.Error("failed to match phone", zap.Error(err))
		return "", err
	}
	if !phoneMatch {
		ts.log.Error("incorrect phone")
		return "", errors.New("incorrect phone")
	}

	params := url.Values{}
	params.Add("login", ts.telecom.login)
	params.Add("pass", ts.telecom.password)
	params.Add("type", "flash")
	params.Add("code_gen", "true")
	params.Add("phone", phone)

	res, err := http.Get(ts.telecom.urlCall + params.Encode())
	if err != nil {
		ts.log.Error("failed to call code", zap.Error(err))
		return "", err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var resObj flashCallResponse
	err = json.Unmarshal(body, &resObj)
	if err != nil {
		ts.log.Error("failed to unmarshal response", zap.Error(err))
		return "", err
	}

	if resObj.Result != "Success" {
		ts.log.Error("failed to call code", zap.Error(err))
		return "", fmt.Errorf("wrong response status: %s", resObj.Message)
	}

	return resObj.Code, nil
}

type flashCallResponse struct {
	Result  string `json:"result,omitempty"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
