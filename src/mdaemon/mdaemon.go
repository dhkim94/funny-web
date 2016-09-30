package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"
	"github.com/pkg/errors"
)

type Context struct {
	PidFileName string
	PidFilePerm os.FileMode
	LogFileName string
	LogFilePerm os.FileMode
	WorkDir string
	Umask int
	Args []string
}

type stringFlag struct {
	s *string
	v string
}

type Flag interface {
	IsSet() bool
}

type SignalHandlerFunc func(sig os.Signal) (err error)

func (f *stringFlag) IsSet() bool {
	if f == nil {
		return false
	}
	return *f.s == f.v
}

func StringFlag(f *string, v string) Flag {
	return &stringFlag(f, v)
}

func AddCommand(f Flag, sig os.Signal, handler SignalHandlerFunc)  {
	if f != nil {
		AddFlag(f, sig)
	}
	if handler != nil {
		SetSigHandler(handler, sig)
	}
}

func AddFlag(f Flag, sig os.Signal) {
	flags[f] = sig
}

func SetSigHandler(handler SignalHandlerFunc, signals ...os.Signal)  {
	for _, sig := range signals {
		handlers[sig] = handler
	}
}

func (d *Context) reborn() (child *os.Process, err error) {
	return nil, errNotSupported
}

func (d *Context) search() (daemon *os.Process, err error) {
	return nil, errNotSupported
}

func (d *Context) release() (err error) {
	return errNotSupported
}

func ActiveFlags() (ret []Flag) {
	ret = make([]Flag, 0, 1)
	for f := range flags {
		if f.IsSet() {
			ret = append(ret, f)
		}
	}
	return
}

func termHandler(sig os.Signal) (err error) {

	return nil
}

func reloadHandler(sig os.Signal) (err error) {
	return nil
}

var flags = make(map[Flag]os.Signal)
var handlers = make(map[os.Signal]SignalHandlerFunc)
var errNotSupported = errors.New("daemon: Non-POSIX OS is not supported")

func main() {
	signal := flag.String("s", "", `send signal to the daemon
		quit - gracefull shutdown
		stop - fast shutdown
		reload - reloading the configuration file`)

	flag.Parse()

	AddCommand(StringFlag(signal, "quit"), syscall.SIGQUIT, termHandler)
	AddCommand(StringFlag(signal, "stop"), syscall.SIGTERM, termHandler)
	AddCommand(StringFlag(signal, "reload"), syscall.SIGHUP, reloadHandler)

	fmt.Println("signal [", *signal, "]")

	//cxt := &Context{
	//	PidFileName: "mpid",
	//	PidFilePerm: 0644,
	//	LogFileName: "mlog",
	//	LogFilePerm: 0644,
	//	WorkDir: "./",
	//	Umask: 027,
	//	Args: []string{"mdaemon study]"},
	//}

	fmt.Println(len(ActiveFlags()))
}
