//go:build windows
// +build windows

package utils

import "C"
import (
	"bytes"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"syscall"

	"github.com/denisbrodbeck/machineid"
	"github.com/mitchellh/go-ps"
)

// IsProcRunning returns true if a process in the names list is running
func IsProcRunning(names ...string) (bool, error) {
	if len(names) == 0 {
		return false, nil
	}

	if runtime.GOOS == "windows" {
		cmd := exec.Command("tasklist.exe", "/fo", "csv", "/nh")
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		out, err := cmd.Output()
		if err != nil {
			return false, err
		}

		for _, name := range names {
			if bytes.Contains(out, []byte(name)) {
				return true, nil
			}
		}
	} else {
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
	}

	return false, nil
}

func GetDeviceID() string {
	id, err := machineid.ID()
	if err != nil {
		log.Fatal(err)
	}

	return id
}
