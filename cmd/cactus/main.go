package main

import (
	"flag"
	"log"
	"math/rand"
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
	go utils.BlockNetworkSniffing()
}

// app entry point
func main() {

	log.Println("Booting up...")

	user := user.NewUser(
		"LICENSE-KEY",
		"example@gmail.com",
		"",
		"cactus-user",
	)

	cactus := cactus.NewCactus()
	cactus.Run(user)
}
