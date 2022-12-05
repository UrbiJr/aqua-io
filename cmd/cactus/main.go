package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/cactus-aio/go-cactus/internal/utils"
)

// app entry point
func main() {

	debugArg := flag.Bool("debug", false, "enable debug mode") // go run ./cmd/cactus/main.go -debug
	flag.Parse()
	debug := *debugArg

	utils.SetDebug(debug)
	go utils.BlockNetworkSniffing()

	log.Println("Booting up...")

	fmt.Printf("Press <ENTER> to exit...")
	fmt.Scanln() // wait for Enter Key
}
