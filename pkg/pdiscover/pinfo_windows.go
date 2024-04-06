package pdiscover

import (
	"os"
)

const (
	CTLKern      = 1
	KernProcArgs = 38
)

func GetOwnProcPath() (string, error) {
	return GetProcPath(os.Getpid())
}

func GetProcPath(pid int) (string, error) {
	// stubbed
	return "", nil
}

func getSysCtlInfo(name []int32, value *byte, valueLen *uintptr) error {
	return nil
}

func GetProcInfo(pid int) map[string]string {
	return nil
}
