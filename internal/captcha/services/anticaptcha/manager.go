package anticaptcha

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/UrbiJr/anticaptcha"
)

var (
	// ErrInvalidCredentials is caused by an invalid API key or access token
	ErrInvalidCredentials = errors.New("credentials rejected by Capmonster")
	// ErrConnection is caused by general connection failure.
	ErrConnection = errors.New("failed to connect to Capmonster")
	// ErrUnknown is caused by some unknown client failure or connection failure.
	ErrUnknown = errors.New("unknown AutoSolve client error")
	// ErrTokenTimeout is caused by when AutoSolve fails to respond within 150 seconds.
	ErrTokenTimeout = errors.New("response wasn't received from Capmonster")
	// ErrTokenCancel may be the result of users canceling the task VIA. the companion app.
	ErrTokenCancel = errors.New("task reported as canceled by Capmonster")
	// ErrCaptchaUnsupported is used when methods for unsupported captcha types are called.
	// AutoSolve only supports reCaptcha, hCaptcha, and GeeTest
	ErrCaptchaUnsupported = errors.New("captcha type unsupported by Capmonster")
)

type AntiCaptchaManager struct {
	Credential Credential
	mu         sync.Mutex
	loaded     bool
}

type Credential struct {
	APIKey string
}

var (
	baseURL       = &url.URL{Host: "api.anti-captcha.com", Scheme: "https", Path: "/"}
	checkInterval = 2 * time.Second
)

func (antiManager *AntiCaptchaManager) GetGeetest(gtKey string, challengeKey string, apiServer string, pageURL string) (anticaptcha.GeeTestSolution, error) {
	c := &anticaptcha.Client{APIKey: antiManager.Credential.APIKey}

	solution, err := c.SendGeeTest(
		pageURL,
		gtKey,
		challengeKey,
		apiServer,
		time.Duration(150)*time.Second,
	)
	if err != nil {
		return anticaptcha.GeeTestSolution{}, err
	} else {
		return solution, nil
	}
}

func (antiManager *AntiCaptchaManager) GetRecaptchaV2(websiteURL string, recaptchaKey string) (string, error) {
	c := &anticaptcha.Client{APIKey: antiManager.Credential.APIKey}

	key, err := c.SendRecaptchaV2(
		websiteURL,
		recaptchaKey,
		time.Duration(150)*time.Second,
	)
	if err != nil {
		return "", err
	} else {
		return key, nil
	}
}
func (antiManager *AntiCaptchaManager) GetImageToText(imageString string) (string, error) {
	c := &anticaptcha.Client{APIKey: antiManager.Credential.APIKey}

	text, err := c.SendImage(
		imageString, // the image file encoded to base64
	)
	if err != nil {
		return "", err
	} else {
		return text, nil
	}

}

func (antiManager *AntiCaptchaManager) GetHCaptchaCreateTask(websiteURL string, recaptchaKey string) (float64, error) {
	// Mount the data to be sent
	body := map[string]interface{}{
		"clientKey": antiManager.Credential.APIKey,
		"task": map[string]interface{}{
			"type":       "HCaptchaTaskProxyless",
			"websiteURL": websiteURL,
			"websiteKey": recaptchaKey,
		},
	}

	b, err := json.Marshal(body)
	if err != nil {
		return 0, err
	}

	// Make the request
	u := baseURL.ResolveReference(&url.URL{Path: "/createTask"})
	resp, err := http.Post(u.String(), "application/json", bytes.NewBuffer(b))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Decode response
	responseBody := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&responseBody)
	// TODO treat api errors and handle them properly
	if _, ok := responseBody["taskId"]; ok {
		if taskId, ok := responseBody["taskId"].(float64); ok {
			return taskId, nil
		}

		return 0, errors.New("task number of irregular format")
	}

	return 0, errors.New("task number not found in server response")
}

func (antiManager *AntiCaptchaManager) getTaskResult(taskID float64) (map[string]interface{}, error) {
	// Mount the data to be sent
	body := map[string]interface{}{
		"clientKey": antiManager.Credential.APIKey,
		"taskId":    taskID,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Make the request
	u := baseURL.ResolveReference(&url.URL{Path: "/getTaskResult"})
	resp, err := http.Post(u.String(), "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode response
	responseBody := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&responseBody)
	return responseBody, nil
}
