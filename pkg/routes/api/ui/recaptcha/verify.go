package recaptcha

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

var recaptchaVerifyURL = "https://www.google.com/recaptcha/api/siteverify"

type Validator interface {
	Validate(code string) error
}

type recaptchaV2Validator struct {
	secret string
	url    string
}

func NewValidatorV2(secret string) Validator {
	return &recaptchaV2Validator{
		secret,
		recaptchaVerifyURL,
	}
}

func (s *recaptchaV2Validator) Validate(code string) error {
	if code == "" {
		return errors.New("Invalid recaptcha")
	}

	logrus.Info("Validating recaptcha code...")

	form := url.Values{}
	form.Add("secret", s.secret)
	form.Add("response", code)
	response, err := http.PostForm(s.url, form)
	if err != nil {
		return err
	}

	respBody, _ := ioutil.ReadAll(response.Body)
	respData := struct {
		Success bool `json:"success"`
	}{}
	if err := json.Unmarshal(respBody, &respData); err != nil {
		return err
	}

	if !respData.Success {
		logrus.Info("Captcha not valid")
		return errors.New("Server reports code invalid")
	}

	logrus.Info("Captcha validated")

	return nil
}
