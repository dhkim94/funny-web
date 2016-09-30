package daemon

import (
	"testing"
	"fmt"
	"os"
)

func TestNewFileLock(t *testing.T) {
	fmt.Println("-----test NewFileLock")

	fileLock := NewFileLock(os.NewFile(1001, "/Users/dhkim/tmp/aa1"))
	fmt.Println(*fileLock.File)
	fmt.Println(*fileLock)

	fileLock.Remove()
}


func TestTest1(t *testing.T) {
	fmt.Println("-----test test1")
}