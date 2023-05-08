package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	fyne_app "fyne.io/fyne/v2/app"

	aqua_io "github.com/UrbiJr/aqua-io/internal/aqua-io"
	"github.com/UrbiJr/aqua-io/internal/captcha"
	"github.com/UrbiJr/aqua-io/internal/client"
	"github.com/UrbiJr/aqua-io/internal/resources"
	"github.com/UrbiJr/aqua-io/internal/user"
	"github.com/UrbiJr/aqua-io/internal/utils"

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
	fyneApp := fyne_app.NewWithID("io.aqua-trading.app")
	// set custom theme
	fyneApp.Settings().SetTheme(&resources.DarkTheme{})
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

	// create the login page

	// get logged user
	app.User = &user.User{
		Email:          "urbijr@app-robotics.eu",
		Username:       "urbijr",
		Settings:       &user.Settings{},
		ProfileManager: &user.ProfileManager{},
	}

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
		app.MakeDesktopUI()
	case "darwin":
		win.Resize(fyne.NewSize(1280, 720))
		win.CenterOnScreen()
		win.SetFixedSize(true)
		win.SetMaster()
		app.MakeDesktopUI()
	default:
		win.Resize(fyne.NewSize(415, 890))
		win.SetFixedSize(true)
		win.SetMaster()
		app.MakeMobileUI()
	}

	win.SetMainMenu(app.MakeMenu())

	win.SetIcon(resources.ResourceIconPng)

	// show and run the application (blocking function)
	win.ShowAndRun()
}
