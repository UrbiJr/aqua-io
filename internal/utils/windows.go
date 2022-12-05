//go:build windows
// +build windows

package utils

import (
	"bytes"
	"math/big"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"

	"github.com/mitchellh/go-ps"
)

// windowsUtils implements Windows-specific platform utilities.
type windowsUtils struct{}

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

// NewPlatformUtils returns a platform-independent interface to an
// instance of windowsUtils.
func NewPlatformUtils() PlatformUtils {
	return &windowsUtils{}
}

// GetSerialNumber returns the serial number of the Windows device.
func (u *windowsUtils) GetSerialNumber() (string, error) {
	return "", ErrInvalidSerialNumber
}

// GetFileID returns a 96 bit file ID as a string.
// https://docs.microsoft.com/en-us/windows/win32/api/fileapi/ns-fileapi-by_handle_file_information#remarks
func (u *windowsUtils) GetFileID(path string) (string, error) {
	info, err := u.getFileInfo(path)
	if err != nil {
		return "", err
	}

	// The file ID only requires 96 bits, though it's documented as being 128 bits by Microsoft.
	id := big.NewInt(int64(info.VolumeSerialNumber))
	id.Lsh(id, 64)
	id.Add(id, new(big.Int).SetUint64((uint64(info.FileIndexHigh)<<32)+uint64(info.FileIndexLow)))
	return id.String(), nil
}

// GetIno returns the 64 bit file index number as a string.
// This method may be removed in preference of a GetDeviceID method that internally hashes
// the result of GetFileID, GetSerialNumber, and GetMAC-48
// https://docs.microsoft.com/en-us/windows/win32/api/fileapi/ns-fileapi-by_handle_file_information#remarks
func (u *windowsUtils) GetIno(path string) (string, error) {
	info, err := u.getFileInfo(path)
	if err != nil {
		return "", err
	}

	ino := new(big.Int).SetUint64((uint64(info.FileIndexHigh) << 32) + uint64(info.FileIndexLow))
	return ino.String(), nil
}

// https://docs.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-getfileinformationbyhandle
// https://github.com/golang/go/blob/49dccf141f5e315739c5517b24572fff7cb13734/src/syscall/types_windows.go#L426-L437
func (u *windowsUtils) getFileInfo(path string) (*syscall.ByHandleFileInformation, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	info := &syscall.ByHandleFileInformation{}
	err = syscall.GetFileInformationByHandle(syscall.Handle(file.Fd()), info)
	if err != nil {
		return nil, &os.PathError{Op: "GetFileInformationByHandle", Path: path, Err: err}
	}

	return info, nil
}
