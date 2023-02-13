package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"github.com/UrbiJr/nyx/internal/nyx"
	"github.com/UrbiJr/nyx/internal/utils"
)

func init() {
	rand.Seed(time.Now().Unix())

	debugArg := flag.Bool("debug", false, "enable debug mode") // go run ./cmdnyx-aiomain.go -debug
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
	nyx.App = fyneApp
	//nyx.TLSClient = client.NewClient()

	// create our loggers
	nyx.Logger = new(utils.AppLogger)
	nyx.Logger.SetupLogger()
	nyx.Logger.Debug("debug logging enabled")

	// open a connection to the database

	// create a database repository

	// create the login page

	// create and size a fyne window
	win := fyneApp.NewWindow("Nyx AIO")
	nyx.MainWindow = win
	win.Resize(fyne.NewSize(900, 500))
	win.SetFixedSize(true)
	win.SetMaster()

	nyx.MakeUI()

	// show and run the application
	win.ShowAndRun()
}
