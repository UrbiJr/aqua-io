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
  "id": "mem_SRBmju7fPOspDv",
  "product": "prod_0f0r0QsOwrFcx",
  "user": "user_lH3nTj7WSHhVl",
  "plan": "plan_8XBGEmK6asLLV",
  "email": "imperatormmxl@hotmail.it",
  "status": "completed",
  "valid": true,
  "cancel_at_period_end": false,
  "payment_processor": "free",
  "license_key": "BETA-18DCE6-E146E733-5C8655W",
  "metadata": {
    "hwid": "New Value"
  },
  "quantity": 1,
  "wallet_address": null,
  "custom_fields_responses": {},
  "discord": {
    "id": "556976708503339008",
    "username": "elleb.j#4827",
    "image_url": "https://cdn.discordapp.com/avatars/556976708503339008/71d0df3ec33ee8172a260749a1d0b402"
  },
  "nft_tokens": [],
  "expires_at": 1699201445,
  "renewal_period_start": null,
  "renewal_period_end": null,
  "created_at": 1683649445,
  "manage_url": "https://whop.com/hub/mem_SRBmju7fPOspDv/",
  "affiliate_page_url": "https://whop.com/aqua-io/?a=lockianclothfb62",
  "checkout_session": null,
  "access_pass": "prod_0f0r0QsOwrFcx"
}
*/

type AuthResult struct {
	Success      bool
	Email        string
	LicenseKey   string
	Discord      any
	ExpiresAt    float64
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
			result.ExpiresAt = parsed["expires_at"].(float64)
		}
	}

	return result, nil
}
