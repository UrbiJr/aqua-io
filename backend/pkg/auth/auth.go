package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/UrbiJr/aqua-io/backend/internal/utils"
	"io"
	"net/http"
	"strings"
)

func AuthenticateUser(isDev bool) error {
	if !isDev {
		//protection.Init()

		//LookupInstance()

		//todo re-enable

		resp, err := validateLicense()
		if err != nil {
			return err
		}

		/*
			err = CheckForUpdate()
			if err != nil {
				return err
			}
		*/

		GlobalUserInformation.DiscordUsername = resp.Discord.Username
		GlobalUserInformation.DiscordUserImage = resp.Discord.ImageUrl

		//go enablePresence(response.Discord.Username)
	}

	return nil
}

func validateLicense() (Response, error) {
	var buf bytes.Buffer
	var response Response

	payload := map[string]interface{}{
		"key": "config.GlobalCfg.License",
		"metadata": map[string]string{
			"hwid": "",
		},
		"signature": "",
	}

	err := json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		return response, err
	}

	req := &http.Request{
		Method: http.MethodPost,
		URL:    validate,
		Body:   io.NopCloser(&buf),
		Header: http.Header{
			"Content-Type": {"application/json"},
			"Hash":         {""},
		},
	}

	authClient := &http.Client{}
	
	resp, err := authClient.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			return response, errors.New("ERROR Authenticate – User not connected to Internet")
		}
		return response, err
	}

	if resp.StatusCode != 200 {
		return response, fmt.Errorf("ERROR Authentication Failed – %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}

	err = resp.Body.Close()
	if err != nil {
		return response, err
	}

	return response, nil
}

type AuthResult struct {
	Success             bool
	Email               string
	LicenseKey          string
	ManageMembershipURL string
	Discord             any
	ExpiresAt           float64
	ErrorMessage        string
}

func (settings *Whop) ValidateLicense(licenseKey string) (*AuthResult, error) {
	hwid := utils.GetDeviceID()
	url := fmt.Sprintf("%smemberships/%s/validate_license", settings.APIBaseEndpoint, licenseKey)
	data := map[string]interface{}{
		"metadata": map[string]interface{}{
			"hwid": hwid,
		},
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+settings.AuthAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return nil, fmt.Errorf("HTTP status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var parsed map[string]any
	result := &AuthResult{}
	result.LicenseKey = licenseKey

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal(body, &parsed)

	for key, _ := range parsed {
		if key == "error" {
			result.Success = false
			errorObj := parsed["error"].(map[string]any)
			result.ErrorMessage = errorObj["message"].(string)
		}
		if key == "valid" {
			result.Success = parsed["valid"].(bool)
		}
		if key == "email" {
			result.Email = parsed["email"].(string)
		}
		if key == "discord" {
			result.Discord = parsed["discord"]
		}
		if key == "expires_at" {
			result.ExpiresAt = parsed["expires_at"].(float64)
		}
		if key == "manage_url" {
			result.ManageMembershipURL = parsed["manage_url"].(string)
		}
	}

	return result, nil
}

func (settings *Whop) ResetLicense(licenseKey string) error {
	url := fmt.Sprintf("%smemberships/%s", settings.APIBaseEndpoint, licenseKey)
	data := map[string]interface{}{
		"metadata": map[string]interface{}{},
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+settings.AuthAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return fmt.Errorf("HTTP status code %d", resp.StatusCode)
	}

	return nil
}
