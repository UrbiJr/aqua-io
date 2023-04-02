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
	"fyne.io/fyne/v2/app"
	"github.com/UrbiJr/nyx/internal/captcha"
	"github.com/UrbiJr/nyx/internal/client"
	"github.com/UrbiJr/nyx/internal/nyx"
	"github.com/UrbiJr/nyx/internal/resources"
	"github.com/UrbiJr/nyx/internal/user"
	"github.com/UrbiJr/nyx/internal/utils"
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
	// Crea il percorso della sottocartella "NyxAIO" all'interno di "AppData/Local".
	// windows: C:\Users\<user>\AppData\Local\Roaming\NyxAIO\logs
	appDataLogsDir = filepath.Join(appDataDir, "NyxAIO", "logs")

	err = os.MkdirAll(appDataLogsDir, os.ModePerm)
	if err != nil {
		log.Println(err)
	}

	go utils.BlockNetworkSniffing()
}

// app entry point
func main() {

	//utils.Info("Booting up...")

	var nyx nyx.Config

	// create a fyne application
	fyneApp := app.NewWithID("com.nyx-aio.nyxapp.preferences")
	// set custom theme
	fyneApp.Settings().SetTheme(&resources.NyxDarkTheme{})
	nyx.App = fyneApp

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
	nyx.TLSClient = client
	if utils.DebugEnabled {
		// enable charles proxy for http client
		proxyStr := "http://127.0.0.1:8888"
		proxyURL, _ := url.Parse(proxyStr)
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
		nyx.Client = &http.Client{Transport: transport}
	} else {
		nyx.Client = &http.Client{}
	}

	// create our loggers
	nyx.Logger = new(utils.AppLogger)
	nyx.Logger.SetupLogger()
	nyx.Logger.Debug("debug logging enabled")

	// open a connection to the database
	sqlDB, err := nyx.ConnectSQL()
	if err != nil {
		log.Panic(err)
	}

	// create a database repository
	nyx.SetupDB(sqlDB)

	// create the login page

	// get logged user
	nyx.User = &user.User{
		Email:                "urbijr@nyx-robotics.eu",
		Username:             "urbijr",
		Settings:             &user.Settings{},
		CopiedTradersManager: &user.CopiedTradersManager{},
		ProfileManager:       &user.ProfileManager{},
	}

	// create and size a fyne window
	win := fyneApp.NewWindow("Nyx AIO")
	nyx.MainWindow = win
	os := runtime.GOOS
	switch os {
	case "windows":
		win.Resize(fyne.NewSize(1390, 848))
		win.CenterOnScreen()
		win.SetFixedSize(true)
		win.SetMaster()
		nyx.MakeDesktopUI()
	case "darwin":
		win.Resize(fyne.NewSize(1390, 848))
		win.CenterOnScreen()
		win.SetFixedSize(true)
		win.SetMaster()
		nyx.MakeDesktopUI()
	default:
		win.Resize(fyne.NewSize(415, 890))
		win.SetFixedSize(true)
		win.SetMaster()
		nyx.MakeMobileUI()
	}

	win.SetMainMenu(nyx.MakeMenu())

	// show and run the application (blocking function)
	win.ShowAndRun()
}
