package daemon

import (
	"os"
	"fmt"
	"github.com/pkg/errors"
)

var EWouldBlock = errors.New("File already locked another process")

type FileLock struct {
	*os.File
}

// 파라미터로 주어진 파일로 부터 파일 Lock을 생성 한다.
func NewFileLock(file *os.File) *FileLock {
	return &FileLock{file}
}

// 파일 Lock을 삭제 한다.
func (fileLock *FileLock) Remove() error {
	defer fileLock.Close()

	fileLock.Unlock()



	return nil
}

func (fileLock *FileLock) Unlock() error {

	fmt.Println("----- file unlock")

	return unlockFile(fileLock.Fd())
}

// 파일 descriptor 를 이용하여 파일명을 구한다.
func GetFdName(fd uintptr) (name string, err error) {
	return getFdName(fd)
}



