package daemon

import (
	"syscall"
	"fmt"
)

// Non-block 파일 쓰기 Lock
func lockFile(fd uintptr) error {
	fmt.Println("----lockFile linux_adm64")

	err := syscall.Flock(int(fd), syscall.LOCK_EX|syscall.LOCK_NB)
	if err == syscall.EWOULDBLOCK {
		err = EWouldBlock
	}
	return err
}

func unlockFile(fd uintptr) error {
	fmt.Println("----unlockFile linux_adm64")


	err := syscall.Flock(int(fd), syscall.LOCK_UN)
	if err == syscall.EWOULDBLOCK {
		err = EWouldBlock
	}
	return err
}