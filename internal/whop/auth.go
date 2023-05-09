package whop

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/UrbiJr/aqua-io/internal/utils"
)

/*
	{
		"id": "mem_oR1RrIIfrsCcut",
		"product": "prod_0f0r0QsOwrFcx",
		"user": "user_REycSOXwO8Ro5",
		"plan": "plan_8XBGEmK6asLLV",
		"email": "carino_piacevole.0m@icloud.com",
		"status": "completed",
		"valid": true,
		"cancel_at_period_end": false,
		"payment_processor": "free",
		"license_key": "BETA-540CD4-FF95D0AD-096D3AW",
		"metadata": {
		  "hwid": "New Value"
		},
		"quantity": 1,
		"wallet_address": null,
		"custom_fields_responses": {},
		"discord": null,
		"nft_tokens": [],
		"expires_at": 1699178770,
		"renewal_period_start": null,
		"renewal_period_end": null,
		"created_at": 1683626770,
		"manage_url": "https://whop.com/hub/mem_oR1RrIIfrsCcut/",
		"affiliate_page_url": "https://whop.com/aqua-io/?a=user21cd1193814",
		"checkout_session": null,
		"access_pass": "prod_0f0r0QsOwrFcx"
	  }
*/

type AuthResult struct {
	Success      bool
	Email        string
	LicenseKey   string
	Discord      any
	ExpiresAt    int64
	ErrorMessage string
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
			result.ExpiresAt = parsed["expires_at"].(int64)
		}
	}

	return result, nil
}
