package daemon

import (
	"os"
	"encoding/json"
	"syscall"
)

var initialized = false

func (d *Context) reborn() (child *os.Process, err error) {



	return nil, nil
}

func (d *Context) search() (daemon *os.Process, err error) {
	if len(d.PidFileName) > 0 {
		var pid int
		if pid, err = ReadPidFile(d.PidFileName); err != nil {

			// todo pid 파일 읽기 실패 로그 넣기

			return
		}

		daemon, err = os.FindProcess(pid)
	}

	return
}

func (d *Context) release() (err error) {
	if !initialized {
		return
	}

	if d.pidFile != nil {
		err = d.pidFile.Remove()
	}

	return
}

func (d *Context) child() (err error) {
	if initialized {
		return os.ErrInvalid
	}

	initialized = true

	decoder := json.NewDecoder(os.Stdin)
	if err = decoder.Decode(d); err != nil {
		return
	}

	if err = syscall.Close(0); err != nil {
		return
	}

	if err = syscall.Dup2(3, 0); err != nil {
		return
	}

	if len(d.PidFileName) > 0 {
		d.pidFile = NewLockFile(os.NewFile(4, d.PidFileName))
		if err = d.pidFile.WritePid(); err != nil {
			return
		}
	}

	if d.Umask != 0 {
		syscall.Umask(int(d.Umask))
	}

	if len(d.Chroot) > 0 {
		err = syscall.Chroot(d.Chroot)
	}

	return
}