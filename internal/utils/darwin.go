//go:build darwin && cgo
// +build darwin,cgo

package utils

import "C"
import (
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
