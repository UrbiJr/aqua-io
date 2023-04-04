package user

type SiteSettings struct {
	Captcha string
	Delay   int
}

type Settings struct {
	LicenseKey    string
	Webhook       string
	Delay         int
	APIKeys       map[string]string
	SitesSettings map[string]SiteSettings
}

// ReadSettings reads settings from the DB and returns read data as Settings
func ReadSettings() (*Settings, error) {

	return &Settings{}, nil
}

// WriteSettings writes settings to the DB
func WriteSettings(settings *Settings) {

}
