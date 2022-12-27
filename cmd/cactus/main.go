package main

import (
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

	//utils.Log("Booting up...")

	user := user.NewUser(
		"LICENSE-KEY",
		"example@gmail.com",
		"",
		"cactus-user",
	)

	cactus := cactus.NewCactus()

	cactus.Run(user)
}
