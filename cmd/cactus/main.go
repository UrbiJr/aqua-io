package main

import (
	"errors"
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/cactus-aio/go-cactus/internal/cactus"
	"github.com/cactus-aio/go-cactus/internal/user"
	"github.com/cactus-aio/go-cactus/internal/utils"
)

func init() {
	rand.Seed(time.Now().Unix())

	debugArg := flag.Bool("debug", false, "enable debug mode") // go run ./cmd/cactus/main.go -debug
	flag.Parse()
	debug := *debugArg
	utils.SetDebug(debug)

	// create missing directories
	path := "tmp/logs"
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Println(err)
	}

	if utils.DebugEnabled {
		utils.Info("Debug logging enabled")
	}

	go utils.BlockNetworkSniffing()
}

// app entry point
func main() {

	//utils.Info("Booting up...")

	cactus := cactus.NewCactus()
	cactus.DrawUI()

	settings, err := user.ReadSettings()
	if err != nil {
		cactus.ShowErrorAndExit(errors.New("settings file is corrupted or contains malformed JSON. You may rename it and start the app again to bypass this error"))
	}

	// logged in
	loggedUser := user.NewUser(
		"example@gmail.com",
		"",
		"cactus-user",
	)

	cactus.User = loggedUser
	cactus.User.Settings = settings
	cactus.User.Profiles, err = user.ReadProfiles()
	if err != nil {
		cactus.ShowErrorAndExit(errors.New("profiles file is corrupted or contains malformed JSON. You may rename it and start the app again to bypass this error"))
	}

	// fill "Profiles" view with user profiles list
	cactus.RefreshProfileView()

	cactus.Run()
}
