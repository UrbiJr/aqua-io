//go:build darwin && cgo
// +build darwin,cgo

package utils

// #cgo LDFLAGS: -framework CoreFoundation -framework IOKit
// #include <CoreFoundation/CoreFoundation.h>
// #include <IOKit/IOKitLib.h>
//
// const char* SerialNumber()
// {
//   CFStringRef serialNumber = NULL;
//   io_service_t platformExpert = IOServiceGetMatchingService(kIOMasterPortDefault, IOServiceMatching("IOPlatformExpertDevice"));
//   if (platformExpert) {
//     CFTypeRef serialNumberAsCFString = IORegistryEntryCreateCFProperty(platformExpert,CFSTR(kIOPlatformSerialNumberKey), kCFAllocatorDefault, 0);
//     if (serialNumberAsCFString) {
//         serialNumber = serialNumberAsCFString;
//         return CFStringGetCStringPtr(serialNumber, kCFStringEncodingUTF8);
//     }
//     IOObjectRelease(platformExpert);
//   }
//   return NULL;
// }
import "C"
import (
	"strings"

	"github.com/mitchellh/go-ps"
)

// darwinUtils implements Darwin-specific platform utilities.
type darwinUtils struct {
	*unixUtils
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

// NewPlatformUtils returns a platform-independent interface to an
// instance of darwinUtils.
func NewPlatformUtils() PlatformUtils {
	return &darwinUtils{&unixUtils{}}
}

// GetSerialNumber returns the serial number of the Darwin device.
func (u *darwinUtils) GetSerialNumber() (string, error) {
	serialNumber := C.GoString(C.SerialNumber())

	if serialNumber == "" {
		return "", ErrInvalidSerialNumber
	}

	return serialNumber, nil
}
