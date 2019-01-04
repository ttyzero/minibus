package lib

import (
	"os"
	"syscall"
)

func pidActive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err == nil {
		err := proc.Signal(syscall.Signal(0))
		if err == nil {
			return true
		}
	}
	return false
}
