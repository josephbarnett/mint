package system

import (
	"fmt"
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

func newSystemInfo() SystemInfo {
	var sysInfo SystemInfo

	sysInfo.Sysname = runtime.GOOS
	sysInfo.Nodename, _ = os.Hostname()

	// Get architecture
	sysInfo.Machine = runtime.GOARCH

	// Retrieve OS version information
	var versionInfo syscall.RtlOsVersionInfoEx
	versionInfo.DwOSVersionInfoSize = uint32(unsafe.Sizeof(versionInfo))

	if err := syscall.RtlGetVersion(&versionInfo); err == nil {
		sysInfo.Release = fmt.Sprintf("%d.%d.%d", versionInfo.DwMajorVersion, versionInfo.DwMinorVersion, versionInfo.DwBuildNumber)
	}

	return sysInfo
}

var defaultSysInfo = newSystemInfo()

func GetSystemInfo() SystemInfo {
	return defaultSysInfo
}
