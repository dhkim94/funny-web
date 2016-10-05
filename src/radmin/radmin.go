package main

// radmin 은 다음 기능을 한다.
// 1. mqtt broker 상태 모니터링
// 2. 관리자 웹 서비스

import (
	"env"
	"flag"
	"fmt"
	"daemon"
	"syscall"
	"os"
	"strconv"
	"time"
)

// NOTE 다른 데몬 프로세스를 만들때 변경 할 곳
// 환경 파일과 환경 파일에서 프로세스명 읽어 오기 전에 사용할 프로세스명
// 환경 파일을 읽어 온 이후에는 프로세스명은 환경 설정에 지정된 값을 사용 한다.
const (
	processName	= "radmin"
	confFileName	= "/Users/dhkim/Documents/Develop/project.my/funny-web/radmin.properties"
)

var (
	cmd = flag.String("c", "", `usage command to the daemon
		start - startup
		stop - fast shutdown
		status - check status
		reload - reloading the configuration file`)
	stop = make(chan  struct{})
	done = make(chan struct{})
)

const (
	cmdStart	= "start"
	cmdStop		= "stop"
	cmdStatus	= "status"
	cmdReload	= "reload"
)

// 데몬 프로세스 만들기 전에 환경을 prepare 한다.
// 데몬 프로세스는 start, stop, status, reload 명령어를 지원 한다.
func prepare() bool {
	flag.Parse()

	// reborn 일때는 daemon start 이므로 강제로 start 로 설정 한다.
	if daemon.WasReborn() {
		*cmd = "start"
	}

	if len(*cmd) <= 0 || (*cmd != cmdStart && *cmd != cmdStop && *cmd != cmdStatus && *cmd != cmdReload){
		fmt.Printf("Usage: %s -c <%s|%s|%s|%s>\n", processName,
			cmdStart, cmdStop, cmdReload, cmdStatus)
		return false
	}

	if !env.Init(confFileName) {
		return false
	}

	pidFileName := env.GetValue("daemon.pidfile.name")
	daemonName := env.GetValue("daemon.name")

	switch {
	case *cmd == cmdStop:
		if pid, err := daemon.ReadPidFile(pidFileName); err != nil {
			fmt.Printf("[FAIL] %s already stopped\n", daemonName);
			return false
		} else {
			fmt.Printf("[OK] %s will stop. pid: %d\n", daemonName, pid)
		}
	case *cmd == cmdStatus:
		if pid, err := daemon.ReadPidFile(pidFileName); err == nil {
			fmt.Printf("[OK] %s already running. pid: %d\n", daemonName, pid)
		} else {
			fmt.Printf("[OK] %s not running\n", daemonName)
		}
		return false
	case *cmd == cmdReload:
		if _, err := daemon.ReadPidFile(pidFileName); err != nil {
			fmt.Printf("[FAIL] %s not running\n", daemonName);
			return false
		} else {
			fmt.Printf("[OK] %s will reload configuration.\n", daemonName)
		}
	case *cmd == cmdStart:
		if pid, err := daemon.ReadPidFile(pidFileName); err != nil {
			fmt.Printf("[OK] %s will startup\n", daemonName)
		} else {
			fmt.Printf("[FAIL] %s already running. pid: %d\n", daemonName, pid)
			return false
		}
	}

	daemon.AddCommand(daemon.StringFlag(cmd, cmdStop), syscall.SIGTERM, terminateHandler)
	daemon.AddCommand(daemon.StringFlag(cmd, cmdReload), syscall.SIGHUP, reloadHandler)

	return true
}

// NOTE 다른 데몬 프로세스를 만들때 변경 할 곳
// 데몬 프로세스 로직
func worker() {
	for {
		time.Sleep(time.Second)

		fmt.Println("---loop")
	}

}

// NOTE 다른 데몬 프로세스를 만들때 변경 할 곳
// 데몬 프로세스 종료 할때 호출 되는 함수.
// 종료 직전의 로직을 처리 하면 된다.
func terminateHandler(sig os.Signal) error {
	slog := env.GetLogger()
	slog.Debug("start terminate handler")


	slog.Debug("complete terminate handler")

	return daemon.ErrStop
}

// NOTE 다른 데몬 프로세스를 만들때 변경 할 곳
// 데몬 프로세스 reload 할때 호출 되는 함수.
// reload 명령어로 환경 값들을 다시 설정 할때 사용 하면 된다.
func reloadHandler(sig os.Signal) error {
	slog := env.GetLogger()

	slog.Info("reload configuration reload")

	return nil
}

// radmin -c <start|stop|reload|status>
func main() {
	//---------------------------------------
	// 1. prepare environment
	//---------------------------------------
	if !prepare() {
		return
	}


	//---------------------------------------
	// 2. prepare daemon
	//---------------------------------------
	pidFilePerm, _ := strconv.ParseUint(env.GetValue("daemon.pidfile.perm"), 8, 32)
	logFilePerm, _ := strconv.ParseUint(env.GetValue("log.file.perm"), 8, 32)
	umask, _ := strconv.ParseUint(env.GetValue("daemon.umask"), 8, 8)
	dname := env.GetValue("daemon.name")

	cxt := &daemon.Context{
		PidFileName: env.GetValue("daemon.pidfile.name"),
		PidFilePerm: os.FileMode(pidFilePerm),
		LogFileName: env.GetValue("log.file"),
		LogFilePerm: os.FileMode(logFilePerm),
		WorkDir:     env.GetValue("daemon.work.dir"),
		Umask:       int(umask),
		Args:        []string{dname},
	}

	slog := env.GetLogger()

	if len(daemon.ActiveFlags()) > 0 {
		d, err := cxt.Search()
		if err != nil {
			slog.Err("Unable send signal to the %s", dname)
			fmt.Println(err)
		}
		daemon.SendCommands(d)
		return
	}

	d, err := cxt.Reborn()
	if err != nil {
		fmt.Printf("[FAIL] %s already running\n", dname)
		return
	}
	if d != nil {
		return
	}
	defer cxt.Release()

	slog.Info("===== start %s =====", dname)

	go worker()

	err = daemon.ServeSignals()
	if err != nil {
		slog.Err("server signals error")
		fmt.Println(err)
	}

	slog.Info("===== terminated %s =====", dname)
}
