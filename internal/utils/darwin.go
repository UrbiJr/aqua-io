//go:build darwin && cgo
// +build darwin,cgo

package utils

import "C"
import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/mitchellh/go-ps"
)

// darwinUtils implements Darwin-specific platform utilities.
type darwinUtils struct {
}

// IsProcRunning returns true if a process in the names list is running
func IsProcRunning(names ...string) (bool, error) {
	if len(names) == 0 {
		return false, nil
	}
	processList, err := ps.Processes()
	if err != nil {
		return false, nil
	}
	for x := range processList {
		for _, name := range names {
			if strings.Contains(processList[x].Executable(), name) {
				return true, nil
			}
		}
	}

	return false, nil
}

func GetDeviceID() string {
	machineUUIDStr := ""
	p, err := os.Open("ioreg -rd1 -c IOPlatformExpertDevice | grep -E '(UUID)'")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer p.Close()

	for {
		line := make([]byte, 1024)
		_, err := p.Read(line)
		if err != nil {
			break
		}
		machineUUIDStr += string(line)
	}

	matchObj := regexp.MustCompile("[A-Z,0-9]{8,8}-[A-Z,0-9]{4,4}-[A-Z,0-9]{4,4}-[A-Z,0-9]{4,4}-[A-Z,0-9]{12,12}")
	results := matchObj.FindAllString(machineUUIDStr, -1)

	return results[0]
}
