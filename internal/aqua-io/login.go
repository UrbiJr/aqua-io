package copy_io

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/UrbiJr/aqua-io/internal/user"
	"github.com/UrbiJr/aqua-io/internal/whop"
)

// loginDialog shows a login form dialog and eventually authenticates user provided licenseKey.
// Returns license authentication result data and true if login is persistent, false otherwise
func (app *Config) LoginDialog() {
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
	loginForm := dialog.NewForm("Login...", "Log In", "Cancel", items, func(b bool) {
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

		if !authResult.Success {
			if authResult.ErrorMessage == "" {
				authResult.ErrorMessage = "application error"
			}
			app.App.SendNotification(&fyne.Notification{
				Title:   "Login Failed",
				Content: authResult.ErrorMessage,
			})
			app.Quit()
		}

		// get logged user
		app.User = &user.User{
			Email:           authResult.Email,
			Discord:         authResult.Discord,
			Username:        "",
			LicenseKey:      authResult.LicenseKey,
			ExpiresAt:       authResult.ExpiresAt,
			PersistentLogin: isPersitent,
			Settings:        &user.Settings{},
			ProfileManager:  &user.ProfileManager{},
		}

	}, app.LoginWindow)

	loginForm.Resize(fyne.NewSize(500, 500))
	loginForm.Show()
}
