package capmonster

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"sync"
	"time"
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

type CapMonsterManager struct {
	Credential Credential
	mu         sync.Mutex
	loaded     bool
}

type Credential struct {
	APIKey string
}

var (
	baseURL       = &url.URL{Host: "api.capmonster.cloud", Scheme: "https", Path: "/"}
	checkInterval = 2 * time.Second
)

func (capManager *CapMonsterManager) GetRecaptchaV2CreateTask(websiteURL string, recaptchaKey string) (float64, error) {
	body := map[string]interface{}{
		"clientKey": capManager.Credential.APIKey,
		"task": map[string]interface{}{
			"type":       "NoCaptchaTaskProxyless",
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

func (capManager *CapMonsterManager) GetHCaptchaCreateTask(websiteURL string, recaptchaKey string) (float64, error) {
	body := map[string]interface{}{
		"clientKey": capManager.Credential.APIKey,
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

func (capManager *CapMonsterManager) getTaskResult(taskID float64) (map[string]interface{}, error) {
	// Mount the data to be sent
	body := map[string]interface{}{
		"clientKey": capManager.Credential.APIKey,
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
