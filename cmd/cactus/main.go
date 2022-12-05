package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

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

	greetings := []string{
		"how can Cactus-AIO assist you today? :-)",
		"how are you going to use Cactus-AIO today? :-)",
		"ready to have some fun with Cactus-AIO?",
		"Cactus-AIO is at your service :cactus-salute:",
		"it's been a while.",
		"time to make the success log go brrr.",
		"everyday is a perfect day to run Cactus-AIO!",
	}

	log.Printf("Hello %s, %s", user.Username, greetings[rand.Intn(len(greetings))])

	fmt.Printf("Press <ENTER> to exit...")
	fmt.Scanln() // wait for Enter Key
}
