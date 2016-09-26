package cklog

import (
	"sync"
	"io"
	"fmt"
	"cktime"
	"time"
	"os"
	"log"
)

type Cklogger struct {
	mu	sync.Mutex
	out	io.Writer
	pool	sync.Pool
	level	int
	path	string
}

const (
	DebugLevel	= iota
	InfoLevel
	WarnLevel
	ErrLevel
)

func NewLogger(level string, out string, path string) *Cklogger  {
	var iLevel int

	switch level {
	case "info":
		iLevel = InfoLevel
	case "warn":
		iLevel = WarnLevel
	case "err":
		iLevel = ErrLevel
	default:
		iLevel = DebugLevel
	}

	switch out {
	case "stderr":
		return &Cklogger{
			out: os.Stderr,
			level: iLevel,
		}
	case "file":
		return &Cklogger{
			path: path,
			level: iLevel,
		}
	default:
		return &Cklogger{
			out: os.Stdout,
			level: iLevel,
		}
	}
}

func New(out io.Writer, level int) *Cklogger {
	return &Cklogger{
		out: out,
		level: level,
	}
}

func (l *Cklogger) Output(level int, format string, v ...interface{}) error {
	if l.level > level {
		return nil
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	ctm := cktime.NewCktime(time.Now(), "hh:mm:ss.SSS")

	var slvl string
	switch level {
	case InfoLevel:
		slvl = "INFO"
	case WarnLevel:
		slvl = "WARN"
	case ErrLevel:
		slvl = "ERR"
	default:
		slvl = "DEBUG"
	}

	format = fmt.Sprintf("%s %s %s", ctm.ToString(), slvl, format)

	// TODO 이거 pool 로 바꾸는 것 찾아 보도록 하자.


	if len(l.path) > 0 {
		ctm.SetFormat("YYYYMMDD")
		_path := l.path + "." + ctm.ToString()

		fd, err := os.OpenFile(_path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			log.Printf("failed to open log file [%s], error [%s]", _path, err)
			return err
		}
		defer fd.Close()

		if len(v) > 0 {
			fd.WriteString(fmt.Sprintf(format+"\n", v...))
		} else {
			fd.WriteString(fmt.Sprint(format+"\n"))
		}

	} else {
		//var buffer *bytes.Buffer
		//var buffer = l.pool.Get().(*bytes.Buffer)
		//buffer.Reset()

		if len(v) > 0 {
			//l.pool.Put([]byte(fmt.Sprintf(format+"\n", v...)))

			l.out.Write([]byte(fmt.Sprintf(format+"\n", v...)))
		} else {
			//l.pool.Put([]byte(fmt.Sprint(format+"\n")))


			l.out.Write([]byte(fmt.Sprint(format+"\n")))
		}

		//fmt.Println("---[", buffer)


		//l.out.Write(buffer)

	}

	return nil
}

func (l *Cklogger) Info(format string, v ...interface{}) {
	l.Output(InfoLevel, format, v...)
}

func (l *Cklogger) Debug(format string, v ...interface{}) {
	l.Output(DebugLevel, format, v...)
}

func (l *Cklogger) Warn(format string, v ...interface{}) {
	l.Output(WarnLevel, format, v...)
}

func (l *Cklogger) Error(format string, v ...interface{}) {
	l.Output(ErrLevel, format, v...)
}