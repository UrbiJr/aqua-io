package main

import (
	"errors"
	"flag"
	"log"
	"math/rand"
	"os"
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

	// create missing directories
	path := "tmp/logs"
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Println(err)
	}

	utils.NewLogger()

	if utils.DebugEnabled {
		utils.Info("Debug logging enabled")
	}

	go utils.BlockNetworkSniffing()
}

// app entry point
func main() {

	//utils.Info("Booting up...")

	nyx := nyx.NewNyx()
	nyx.DrawUI()

	settings, err := user.ReadSettings()
	if err != nil {
		nyx.ShowErrorAndExit(errors.New("settings file is corrupted or contains malformed JSON. You may rename it and start the app again to bypass this error"))
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
		nyx.ShowErrorAndExit(errors.New("profiles file is corrupted or contains malformed JSON. You may rename it and start the app again to bypass this error"))
	}

	// fill "Profiles" view with user profiles list
	nyx.RefreshProfileView()

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
