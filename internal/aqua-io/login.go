package copy_io

import (
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/UrbiJr/aqua-io/internal/whop"
)

// loginDialog shows a login form dialog and eventually authenticates user provided licenseKey.
// Returns license authentication result data and true if login is persistent, false otherwise
func (app *Config) LoginDialog() (*whop.AuthResult, bool) {
	licenseKey := widget.NewEntry()
	licenseKey.Validator = validation.NewRegexp(`^(BETA|AQUA)-[0-9A-F]{6}-[0-9A-F]{8}-[0-9A-F]{6}`, "wrong license key format")

	isPersitent := false
	items := []*widget.FormItem{
		widget.NewFormItem("License Key", licenseKey),
		widget.NewFormItem("Remember me", widget.NewCheck("", func(checked bool) {
			isPersitent = checked
		})),
	}

	authResult := &whop.AuthResult{}
	dialog.ShowForm("Login...", "Log In", "Cancel", items, func(b bool) {
		if !b {
			return
		}

		result, err := app.Whop.ValidateLicense(licenseKey.Text)
		if err != nil {
			authResult.Success = false
			app.Logger.Error(err)
		} else {
			authResult = result
		}

	}, app.LoginWindow)

	return authResult, isPersitent
}
