package copy_io

import (
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/UrbiJr/aqua-io/internal/resources"
	"github.com/UrbiJr/aqua-io/internal/user"
	"github.com/UrbiJr/aqua-io/internal/whop"
)

// MakeLoginWindow creates content for the login window and eventually authenticates user provided licenseKey.
// If authenticated, hides the login window and show the main window. Otherwise, quits the app
func (app *Config) MakeLoginWindow() {
	licenseKey := widget.NewEntry()
	licenseKey.SetPlaceHolder("Your license key")
	licenseKey.Validator = validation.NewRegexp(`^(BETA|AQUA)-[0-9A-F]{6}-[0-9A-F]{8}-[0-9A-F]{6}`, "wrong license key format")
	errorText := canvas.NewText("", color.RGBA{R: 255, G: 50, B: 50, A: 255})
	isPersitent := false
	rememberMe := widget.NewCheck("Remember Me", func(checked bool) {
		isPersitent = checked
	})

	appLogo := canvas.NewImageFromResource(resources.ResourceIconPng)
	appLogo.SetMinSize(fyne.NewSize(25, 25))
	appLogo.FillMode = canvas.ImageFillContain
	authResult := &whop.AuthResult{}
	signInButton := widget.NewButtonWithIcon("Sign In", theme.LoginIcon(), func() {
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
			errorText.Text = fmt.Sprintf("Login Error: %s", authResult.ErrorMessage)
			errorText.Refresh()
		} else {
			// get logged user
			discordID := ""
			username := strings.Split(authResult.Email, "@")[0]
			profilePicture := ""
			if authResult.Discord != nil {
				discordInfo := authResult.Discord.(map[string]any)
				discordID = discordInfo["id"].(string)
				username = discordInfo["username"].(string)
				profilePicture = discordInfo["image_url"].(string)
			}

			loggedUser := user.NewUser(
				authResult.Email,
				discordID,
				username,
				profilePicture,
				authResult.LicenseKey,
				authResult.ManageMembershipURL,
				authResult.ExpiresAt,
				isPersitent)

			// save user to sqlite DB
			app.DB.DeleteAllUsers()
			loggedUser.Theme = "light"
			inserted, err := app.DB.InsertUser(*loggedUser)
			if err != nil {
				app.Logger.Error(err)
				app.Quit()
			}
			app.User = inserted
			app.MakeTray()
			app.MakeDesktopUI()
			app.MainWindow.Show()
			app.LoginWindow.Hide()
		}
	})
	cancelButton := widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
		app.Quit()
	})
	signInButton.Importance = widget.HighImportance
	cancelButton.Importance = widget.DangerImportance
	vBox := container.NewVBox(
		container.NewCenter(
			container.NewHBox(widget.NewRichTextFromMarkdown(`## User Login`), appLogo)),
		licenseKey,
		rememberMe,
		errorText,
		container.NewCenter(
			container.NewHBox(
				signInButton,
				cancelButton,
			)),
	)

	app.LoginWindow.SetContent(vBox)
}
