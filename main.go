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

	"fyne.io/fyne/v2"
	fyne_app "fyne.io/fyne/v2/app"

	aqua_io "github.com/UrbiJr/aqua-io/internal/aqua-io"
	"github.com/UrbiJr/aqua-io/internal/captcha"
	"github.com/UrbiJr/aqua-io/internal/client"
	"github.com/UrbiJr/aqua-io/internal/resources"
	"github.com/UrbiJr/aqua-io/internal/user"
	"github.com/UrbiJr/aqua-io/internal/utils"
	"github.com/UrbiJr/aqua-io/internal/whop"

	tls_client "github.com/bogdanfinn/tls-client"

	_ "github.com/glebarez/go-sqlite"
)

func init() {
	rand.Seed(time.Now().Unix())

	debugArg := flag.Bool("debug", false, "enable debug mode") // go run ./main.go -debug
	flag.Parse()
	debug := *debugArg
	utils.SetDebug(debug)

	var appDataLogsDir string

	// Ottieni il percorso della cartella "AppData" per l'utente corrente su Windows.
	appDataDir, err := os.UserCacheDir()
	if err != nil {
		log.Println(err)
		return
	}
	// Crea il percorso della sottocartella "Aqua IO" all'interno di "AppData/Local".
	// windows: C:\Users\<user>\AppData\Local\Roaming\Aqua IO\logs
	appDataLogsDir = filepath.Join(appDataDir, "Aqua IO", "logs")

	err = os.MkdirAll(appDataLogsDir, os.ModePerm)
	if err != nil {
		log.Println(err)
	}

	go utils.BlockNetworkSniffing()
}

// app entry point
func main() {

	//utils.Info("Booting up...")

	var app aqua_io.Config

	// create needed directories if don't exist
	path := "downloads"
	// ignore the error
	_ = os.Mkdir(path, os.ModePerm)

	// create a fyne application
	fyneApp := fyne_app.NewWithID("trading.aqua-io.app")
	// set custom theme
	fyneApp.Settings().SetTheme(&resources.LightTheme{})
	app.App = fyneApp

	clientOptions := &client.ClientOptions{
		Timeout:          30,
		TlsClientProfile: tls_client.Chrome_110,
	}
	if utils.DebugEnabled {
		// enable charles proxy for tls client
		clientOptions.CharlesProxy = true
	}
	client, err := client.NewTLSClient(&captcha.SolverOptions{Provider: "2captcha"}, clientOptions)
	if err != nil {
		log.Panic(err)
	}
	app.TLSClient = client
	if utils.DebugEnabled {
		// enable charles proxy for http client
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
	app.Logger = new(utils.AppLogger)
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
	whopSettings := whop.InitWhop()
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
	dbUser, err := app.DB.GetAllUsers()
	showLogin := false
	if err != nil {
		app.Logger.Error(err)
		showLogin = true
	} else if dbUser != nil && dbUser.PersistentLogin {
		// and attempt automatic login if persistent was set
		authResult, err := app.Whop.ValidateLicense(dbUser.LicenseKey)
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
					authResult.ExpiresAt,
					dbUser.PersistentLogin)

				loggedUser.Theme = dbUser.Theme
				// only check for dark since light is set as default
				if loggedUser.Theme == "dark" {
					fyneApp.Settings().SetTheme(&resources.DarkTheme{})
				}
				loggedUser.ID = dbUser.ID
				loggedUser.ProfilePicturePath = dbUser.ProfilePicturePath
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
		app.MakeDesktopUI()
		app.MainWindow.Show()
	}

	// show and run the application (blocking function)
	app.App.Run()
}
