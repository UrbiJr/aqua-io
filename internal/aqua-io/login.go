package copy_io

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
	"github.com/UrbiJr/aqua-io/internal/resources"
	"github.com/UrbiJr/aqua-io/internal/user"
	"github.com/UrbiJr/aqua-io/internal/whop"
)

// loginDialog shows a login form dialog and eventually authenticates user provided licenseKey.
// Returns license authentication result data and true if login is persistent, false otherwise
func (app *Config) ShowLogin() {
	licenseKey := widget.NewEntry()
	licenseKey.Validator = validation.NewRegexp(`^(BETA|AQUA)-[0-9A-F]{6}-[0-9A-F]{8}-[0-9A-F]{6}`, "wrong license key format")
	errorMsgLabel := widget.NewLabel("")
	isPersitent := false
	rememberMe := widget.NewCheck("Remember Me", func(checked bool) {
		isPersitent = checked
	})

	appLogo := canvas.NewImageFromResource(resources.ResourceIconPng)
	appLogo.SetMinSize(fyne.NewSize(25, 25))
	appLogo.FillMode = canvas.ImageFillContain
	authResult := &whop.AuthResult{}
	vBox := container.NewVBox(
		container.NewCenter(
			container.NewHBox(widget.NewRichTextFromMarkdown(`## User Login`), appLogo)),
		widget.NewLabel("License Key"),
		licenseKey,
		rememberMe,
		errorMsgLabel,
		container.NewCenter(
			container.NewHBox(
				widget.NewButton("Sign In", func() {
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
						errorMsgLabel.SetText(fmt.Sprintf("Login Error: %s", authResult.ErrorMessage))
					} else {
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
						app.LoginWindow.Close()
						app.MakeDesktopUI()
						app.MainWindow.Show()
					}
				}),
				widget.NewButton("Cancel", func() {
					app.Quit()
				}),
			)),
	)

	app.LoginWindow.SetContent(vBox)
}
