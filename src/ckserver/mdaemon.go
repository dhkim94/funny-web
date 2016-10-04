package main

import (
	"flag"
	"daemon"
	"syscall"
	"os"
	"log"
	"time"
	"fmt"
)

var (
	signal = flag.String("s", "", `send signal to the daemon
		quit - graceful shutdown
		stop - fast shutdown
		reload - reloading the configuration file`)
	stop = make(chan  struct{})
	done = make(chan struct{})
)

func main() {
	flag.Parse()
	daemon.AddCommand(daemon.StringFlag(signal, "quit"), syscall.SIGQUIT, termHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "stop"), syscall.SIGTERM, termHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "reload"), syscall.SIGHUP, reloadHandler)

	cntxt := &daemon.Context{
		PidFileName: "pid",
		PidFilePerm: 0644,
		LogFileName: "log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"ckdaemon"},
	}

	if len(daemon.ActiveFlags()) > 0 {
		d, err := cntxt.Search()
		if err != nil {
			log.Fatalln("Unable send signal to the daemon:", err)
		}
		daemon.SendCommands(d)
		return
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatalln(err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Println("- - - - - - - - - - - - - - -")
	log.Println("daemon started")

	log.SetPrefix("INFO ")
	log.Println("-----info log")

	log.SetPrefix("DEBUG ")
	log.Println("-----debug log")

	fmt.Println("----print print")


	go worker()

	err = daemon.ServeSignals()
	if err != nil {
		log.Println("Error:", err)
	}
	log.Println("daemon terminated")
}

func worker() {
	for {
		time.Sleep(time.Second)
		log.Println("-----run worker 1")
		log.Println("-----run worker 1-1")

		// 이건 non-block code
		select {
		case ok := <-stop:
			log.Println("-----stop ok [", ok)
		default:
			log.Println("-----continue")
		}

		// 이건 block code
		//if _, ok := <-stop; ok {
		//	log.Println("-----run worker 2")
		//	break
		//}

		//log.Println(<-stop)

		log.Println("-----run worker 3")
	}
	log.Println("-----run worker 4")
	done <- struct{}{}
}

func termHandler(sig os.Signal) error {
	log.Println("terminating ...")
	stop <- struct {}{}

	if sig == syscall.SIGQUIT {
		<- done
	}

	return daemon.ErrStop
}

func reloadHandler(sig os.Signal) error {
	log.Println("configuration reloaded")
	return nil
}