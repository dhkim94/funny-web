package daemon

import (
	"testing"
	"fmt"
	"os"
	"io/ioutil"
	"runtime"
)

var (
	filename string		= os.TempDir() + "test.lock"
	fileperm os.FileMode	= 0644
	invalidname string	= "/a/b/c/d"
)

func setdata() {
	if runtime.GOOS == "darwin" {
		filename = os.TempDir() + "test.lock"
	} else {
		filename = os.TempDir() + "/test.lock"
	}
}

func TestCreatePidFile(t *testing.T) {
	setdata()

	if _, err := CreatePidFile(invalidname, fileperm); err == nil {
		t.Fatal("err1: ", err)
	}

	fmt.Println("> filename [", filename, "]")

	lock, err := CreatePidFile(filename, fileperm)
	if err != nil {
		t.Fatal("err2: ", err)
	}
	defer lock.Remove()

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal("err3: ", err)
	}

	if string(data) != fmt.Sprint(os.Getpid()) {
		t.Fatal("err4: pid not equal")
	}

	file, err := os.OpenFile(filename, os.O_RDONLY, fileperm)
	if err != nil {
		t.Fatal("err5: ", err)
	}

	if err = NewLockFile(file).WritePid(); err == nil {
		t.Fatal("err6: ", err)
	}
}


func TestNewLockFile(t *testing.T) {
	setdata()

	fmt.Println("-----test NewFileLock")

	file := NewLockFile(os.NewFile(1001, "/Users/dhkim/tmp/known_hosts"))

	err := file.Remove()



	if err != nil {
		t.Fatal("Remove(): error invalid fd")
	}






}


func TestGetFdName(t *testing.T) {
	setdata()

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
		t.Errorf("GetFdName(): detected invalid fd. name [%s]", name)
	}
}