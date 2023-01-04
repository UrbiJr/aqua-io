package main

import (
	"errors"
	"flag"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/UrbiJr/go-cactus/internal/nyx"
	"github.com/UrbiJr/go-cactus/internal/user"
	"github.com/UrbiJr/go-cactus/internal/utils"
)

func init() {
	rand.Seed(time.Now().Unix())

	debugArg := flag.Bool("debug", false, "enable debug mode") // go run ./cmd/nyx/main.go -debug
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

	utils.NewLogger()

	utils.Debug("debug logging enabled")

	go utils.BlockNetworkSniffing()
}

// app entry point
func main() {

	//utils.Info("Booting up...")

	nyx := nyx.NewNyx()
	nyx.DrawUI()

	settings, err := user.ReadSettings()
	if err != nil {
		nyx.ShowErrorAndExit(errors.New("settings file is corrupted or contains malformed JSON. You may rename it or delete it and start the app again"))
	}

	// logged in
	loggedUser := user.NewUser(
		"example@gmail.com",
		"",
		"nyx-user",
	)

	nyx.User = loggedUser
	nyx.User.Settings = settings
	nyx.User.Profiles, err = user.ReadProfiles()
	if err != nil {
		nyx.ShowErrorAndExit(errors.New("profiles file is corrupted or contains malformed JSON. You may rename it or delete it and start the app again"))
	}

	nyx.User.ProxyProfiles, err = user.ReadProxies()
	if err != nil {
		nyx.ShowErrorAndExit(errors.New("proxies file is corrupted or contains malformed JSON. You may rename it or delete it and start the app again"))
	}

	// draw user profiles list to "Profiles" view
	nyx.RefreshProfilesView()

	// if user has at least 1 proxy profile
	if len(nyx.User.ProxyProfiles) > 0 {
		// draw the elements affected to "Proxies" view
		nyx.RefreshProxiesView(0)
	} else {
		// otherwise draw only the remaining elements
		nyx.RefreshProxiesView(-1)
	}
	/*
		greetings := []string{
			"how can Nyx assist you today? :-)",
			"how are you going to use Nyx today? :-)",
			"ready to have some fun with Nyx?",
			"Nyx is at your service :nyx-salute:",
			"it's been a while.",
			"time to make the success log go brrr.",
			"everyday is a perfect day to run Nyx!",
		}

		msg := fmt.Sprintf("Hello %s, %s", user.Username, greetings[rand.Intn(len(greetings))])
	*/
	nyx.Run()
}
