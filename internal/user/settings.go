package user

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/cactus-aio/go-cactus/internal/sites"
)

type SiteSettings struct {
	Captcha string `json:"captcha"`
	Delay   int    `json:"delay"`
}

type Settings struct {
	LicenseKey    string                  `json:"license_key"`
	Webhook       string                  `json:"webhook"`
	Delay         int                     `json:"delay"`
	APIKeys       map[string]string       `json:"api_keys"`
	SitesSettings map[string]SiteSettings `json:"site_settings"`
}

// ReadSettings reads settings.json and returns read data as Settings
func ReadSettings() (*Settings, error) {

	jsonFile, err := os.Open("settings.json")
	if err != nil {
		// file does not exist
		settings := Settings{
			APIKeys: map[string]string{},
			SitesSettings: map[string]SiteSettings{
				sites.Nike.String():  {},
				sites.Schuh.String(): {},
			},
		}
		settings.APIKeys["2captcha"] = ""
		settings.APIKeys["capmonster"] = ""
		settings.APIKeys["anticaptcha"] = ""
		file, _ := json.MarshalIndent(settings, "", "   ")
		_ = ioutil.WriteFile("settings.json", file, 0644)
		return &settings, nil
	}

	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	settings := Settings{
		APIKeys:       map[string]string{},
		SitesSettings: map[string]SiteSettings{},
	}
	err = json.Unmarshal(byteValue, &settings)
	if err != nil {
		return &Settings{}, err
	}

	return &settings, nil

}

// WriteSettings writes settings to settings.json
func WriteSettings(settings *Settings) {
	file, _ := json.MarshalIndent(settings, "", "   ")
	_ = ioutil.WriteFile("settings.json", file, 0644)
}
