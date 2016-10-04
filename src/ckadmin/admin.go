package main

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

var (
	cmd = flag.String("c", "", `usage command to the daemon
		start - startup
		stop - fast shutdown
		status - check status
		reload - reloading the configuration file`)
	stop = make(chan  struct{})
	done = make(chan struct{})
)

// NOTE 다른 데몬 프로세스를 만들때 변경 할 곳
// 환경 파일과 환경 파일에서 프로세스명 읽어 오기 전에 사용할 프로세스명
// 환경 파일을 읽어 온 이후에는 프로세스명은 환경 설정에 지정된 값을 사용 한다.
const (
	processName	= "ckadmin"
	confFileName	= "/Users/dhkim/Documents/Develop/project.my/funny-web/ckadmin.properties"
)

// 데몬 프로세스 만들기 전에 환경을 prepare 한다.
// 데몬 프로세스는 start, stop, status, reload 명령어를 지원 한다.
func prepare() bool {
	flag.Parse()

	// reborn 일때는 daemon start 이므로 강제로 start 로 설정 한다.
	if daemon.WasReborn() {
		*cmd = "start"
	}

	if len(*cmd) <= 0 || (*cmd != "start" && *cmd != "stop" && *cmd != "status" && *cmd != "reload"){
		fmt.Printf("Usage: %s -c <start|stop|reload>\n", processName)
		return false
	}

	if !env.Init(confFileName) {
		return false
	}

	pidFileName := env.GetValue("daemon.pidfile.name")
	daemonName := env.GetValue("daemon.name")

	switch {
	case *cmd == "stop" || *cmd == "status":
		if _, err := daemon.ReadPidFile(pidFileName); err != nil {
			fmt.Printf("[OK] %s already stopped\n", daemonName);
			return false
		}
	case *cmd == "reload":
		if _, err := daemon.ReadPidFile(pidFileName); err != nil {
			fmt.Printf("[FAIL] %s not running\n", daemonName);
			return false
		}
	case *cmd == "start":
		daemon.AddCommand(daemon.StringFlag(cmd, "stop"), syscall.SIGTERM, terminateHandler)
		daemon.AddCommand(daemon.StringFlag(cmd, "reload"), syscall.SIGHUP, reloadHandler)
		daemon.AddCommand(daemon.StringFlag(cmd, "status"), syscall.SIGINT, statusHandler)
	}

	return true
}

// 데몬 프로세스 로직
func worker() {
	for {
		time.Sleep(time.Second)

		fmt.Println("---loop")
	}

}

// todo terminate 걸어야 한다.

// ckadmin -c <start|stop|reload>
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
			slog.Err("Unable send signal to the %s\n", dname)
			fmt.Println(err)
		}
		daemon.SendCommands(d)
		return
	}

	d, err := cxt.Reborn()
	if err != nil {

		// todo stop, reload 명령어를 쳐도 여기 들어온다...음...걸러 줘야 하네...

		fmt.Printf("[FAIL] %s already running\n", dname)
		//fmt.Println(err)
		return
	}
	if d != nil {
		return
	}
	defer cxt.Release()

	slog.Info("===== start %s =====\n", dname)

	go worker()

	err = daemon.ServeSignals()
	if err != nil {
		slog.Err("server signals error")
		fmt.Println(err)
	}

	slog.Info("===== terminated %s =====\n", dname)
}


// 일단 만들어 두고 추후 수정 할 것임
func terminateHandler(sig os.Signal) error {

	fmt.Println("-----terminate")

	return daemon.ErrStop
}

func reloadHandler(sig os.Signal) error {

	fmt.Println("-----reloading")

	return nil
}

func statusHandler(sig os.Signal) error {

	fmt.Println("-----status")

	return nil
}