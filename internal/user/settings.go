package user

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type Settings struct {
	LicenseKey string            `json:"license_key"`
	Webhook    string            `json:"webhook"`
	Delay      string            `json:"delay"`
	APIKeys    map[string]string `json:"api_keys"`
}

func init() {
	if _, err := os.Stat("settings.json"); errors.Is(err, os.ErrNotExist) {
		// file does not exist
		settings := Settings{
			APIKeys: map[string]string{},
		}
		settings.APIKeys["2captcha"] = ""
		settings.APIKeys["capmonster"] = ""
		settings.APIKeys["anticaptcha"] = ""
		file, _ := json.MarshalIndent(settings, "", " ")
		_ = ioutil.WriteFile("settings.json", file, 0644)
	}
}

// ReadSettings reads settings.json and returns read data as Settings
func ReadSettings() *Settings {
	jsonFile, err := os.Open("settings.json")
	if err != nil {
		return &Settings{}
	}
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var settings *Settings
	err = json.Unmarshal(byteValue, settings)
	if err != nil {
		return &Settings{}
	}

	return settings

}

// WriteSettings writes settings to settings.json
func WriteSettings(settings *Settings) {
	file, _ := json.MarshalIndent(settings, "", " ")
	_ = ioutil.WriteFile("settings.json", file, 0644)
}
