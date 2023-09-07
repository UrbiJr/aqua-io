package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	aqua_io "github.com/UrbiJr/aqua-io/backend/internal/aqua-io"
	"github.com/UrbiJr/aqua-io/backend/internal/resources"
	"github.com/UrbiJr/aqua-io/backend/internal/user"
	"github.com/UrbiJr/aqua-io/backend/internal/utils"
	"github.com/UrbiJr/aqua-io/backend/pkg/auth"
	"github.com/UrbiJr/aqua-io/backend/pkg/logger"
	"github.com/UrbiJr/aqua-io/backend/pkg/protection"

	"fyne.io/fyne/v2"
	fyne_app "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	_ "github.com/glebarez/go-sqlite"
)

func init() {
	rand.Seed(time.Now().Unix())

	debugArg := flag.Bool("debug", false, "enable debug mode") // go run ./backend/internal/main.go -debug
	flag.Parse()
	debug := *debugArg
	utils.SetDebug(debug)

	var appDataLogsDir string

	// get the "AppData" directory path for the current user on Windows.
	appDataDir, err := os.UserCacheDir()
	if err != nil {
		log.Println(err)
		return
	}
	// create subdirectory "Aqua IO" path inside "AppData/Local".
	// windows: C:\Users\<user>\AppData\Local\Roaming\Aqua IO\logs
	appDataLogsDir = filepath.Join(appDataDir, "Aqua IO", "logs")

	err = os.MkdirAll(appDataLogsDir, os.ModePerm)
	if err != nil {
		log.Println(err)
	}

	go protection.BlockNetworkSniffing()
}

// app entry point
func main() {

	//utils.Info("Booting up...")
	app := aqua_io.App
	app.UI = &aqua_io.UI{}

	// create needed directories if don't exist
	path := "downloads"
	// ignore the error
	_ = os.Mkdir(path, os.ModePerm)

	// create a fyne application
	fyneApp := fyne_app.NewWithID("trading.aqua-io.app")
	// set custom theme
	fyneApp.Settings().SetTheme(&resources.LightTheme{})
	app.App = fyneApp

	if utils.DebugEnabled {
		// add localhost proxy
		proxyStr := "http://127.0.0.1:8888"
		proxyURL, _ := url.Parse(proxyStr)
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
		app.Client = &http.Client{Transport: transport}
	} else {
		app.Client = &http.Client{}
	}

	// create our loggers
	app.Logger = new(logger.AppLogger)
	app.Logger.SetupLogger()
	app.Logger.Debug("debug logging enabled")

	// open a connection to the database
	sqlDB, err := app.ConnectSQL()
	if err != nil {
		log.Panic(err)
	}

	// create a database repository
	app.SetupDB(sqlDB)

	// get Whop config
	whopSettings := auth.InitWhop()
	app.Whop = whopSettings

	// create the login page
	app.LoginWindow = app.App.NewWindow("Aqua.io - Login")
	app.MakeLoginWindow()
	app.LoginWindow.Resize(fyne.NewSize(300, 300))
	app.LoginWindow.CenterOnScreen()
	app.LoginWindow.SetFixedSize(true)
	app.LoginWindow.SetIcon(resources.ResourceIconPng)
	app.LoginWindow.SetOnClosed(func() {
		app.Quit()
	})

	// create and size a fyne window
	win := fyneApp.NewWindow("Aqua.io")
	app.MainWindow = win
	os := runtime.GOOS
	switch os {
	case "windows":
		win.Resize(fyne.NewSize(1280, 720))
		win.CenterOnScreen()
		win.SetFixedSize(true)
		win.SetMaster()
	case "darwin":
		win.Resize(fyne.NewSize(1280, 720))
		win.CenterOnScreen()
		win.SetFixedSize(true)
		win.SetMaster()
	default:
		win.Resize(fyne.NewSize(415, 890))
		win.SetFixedSize(true)
		win.SetMaster()
	}

	win.SetMainMenu(app.MakeMenu())
	win.SetIcon(resources.ResourceIconPng)

	// retrieve user if stored locally
	dbUsers, err := app.DB.AllUsers()
	showLogin := false
	if err != nil {
		app.Logger.Error(err)
		showLogin = true
	} else if len(dbUsers) > 0 && dbUsers[0].ID > 0 && dbUsers[0].PersistentLogin {
		// and attempt automatic login if persistent was set
		authResult, err := app.Whop.ValidateLicense(dbUsers[0].LicenseKey)
		if err != nil {
			app.Logger.Error(err)
			showLogin = true
		} else {
			if !authResult.Success {
				app.App.SendNotification(fyne.NewNotification(
					"⚠️ Auto-Login Failed",
					fmt.Sprintf("Error: %s", authResult.ErrorMessage),
				))
				showLogin = true
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
					dbUsers[0].PersistentLogin)

				loggedUser.Theme = dbUsers[0].Theme
				// only check for dark since light is set as default
				if loggedUser.Theme == "dark" {
					fyneApp.Settings().SetTheme(&resources.DarkTheme{})
				}
				loggedUser.ID = dbUsers[0].ID
				loggedUser.ProfilePicturePath = dbUsers[0].ProfilePicturePath
				// update db with info fetched from whop
				err = app.DB.UpdateUser(loggedUser.ID, *loggedUser)
				if err != nil {
					app.Logger.Error(err)
				}
				app.User = loggedUser
			}
		}
	} else {
		showLogin = true
	}

	// otherwise show login window
	app.MakeTray()
	if showLogin {
		app.LoginWindow.Show()
	} else {
		app.SplashWindow = app.App.NewWindow("Aqua.io")
		appLogo := canvas.NewImageFromResource(resources.ResourceIconPng)
		appLogo.SetMinSize(fyne.NewSize(35, 35))
		appLogo.FillMode = canvas.ImageFillContain
		preloader := container.NewVBox(
			container.NewCenter(
				container.NewHBox(widget.NewRichTextFromMarkdown(`## Loading App...`), appLogo),
			),
			layout.NewSpacer(),
			widget.NewProgressBarInfinite(),
			layout.NewSpacer(),
		)
		app.SplashWindow.SetContent(preloader)
		app.SplashWindow.CenterOnScreen()
		app.SplashWindow.Show()

		go func() {
			app.MakeDesktopUI()
			app.MainWindow.Show()
			app.SplashWindow.Hide()
		}()
	}

	// show and run the application (blocking function)
	app.App.Run()
}
