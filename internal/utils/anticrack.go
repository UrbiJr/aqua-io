package utils

import (
	"log"
	"time"
)

func BlockNetworkSniffing() {
	if !DebugEnabled {
		for range time.Tick(time.Second * 4) {
			running, err := IsProcRunning("Charles.exe", "burp", "Burp", "wireshark", "Wireshark", "postman", "Postman", "mitm", "Mitm", "Proxyman", "proxyman")

			if err == nil && running {
				log.Fatal("Security Risk Detected. Closing Bot...")
			}
		}
	}
}
