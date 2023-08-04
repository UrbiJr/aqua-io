package protection

import (
	"github.com/UrbiJr/aqua-io/backend/internal/utils"
	"log"
	"time"
)

func BlockNetworkSniffing() {
	if !utils.DebugEnabled {
		for range time.Tick(time.Second * 4) {
			running, err := utils.IsProcRunning("Charles.exe", "burp", "Burp", "wireshark", "Wireshark", "postman", "Postman", "mitm", "Mitm", "Proxyman", "proxyman")

			if err == nil && running {
				log.Fatal("Security Risk Detected. Closing Bot...")
			}
		}
	}
}
