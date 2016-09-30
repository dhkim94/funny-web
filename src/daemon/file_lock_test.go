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

func TestGetFdName(t *testing.T) {
	name, err := GetFdName(0)
	if err != nil {
		t.Error(err)
	} else {
		if name != "/dev/null" {
			t.Errorf("Filename of fd 0: [%s]", name)
		}
	}

	name, err = GetFdName(1011)
	if err == nil {
		t.Errorf("detected invalid fd. name [%s]", name)
	}
}