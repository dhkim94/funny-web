package cklog

import (
	"sync"
	"io"
	"fmt"
	"cktime"
	"time"
	"os"
	"log"
	"bufio"
	"bytes"
)

type Cklogger struct {
	mu	sync.Mutex
	out	io.Writer
	level	int
	path	string
}

const (
	DebugLevel	= iota
	InfoLevel
	WarnLevel
	ErrLevel
)

var bufferPool *sync.Pool

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

	bufferPool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
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

	if len(l.path) > 0 {
		ctm.SetFormat("YYYYMMDD")
		_path := l.path + "." + ctm.ToString()

		fd, err := os.OpenFile(_path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			log.Printf("failed to open log file [%s], error [%s]", _path, err)
			return err
		}
		defer fd.Close()

		wb := bufio.NewWriter(fd)

		buffer := bufferPool.Get().(*bytes.Buffer)
		buffer.Reset()
		defer bufferPool.Put(buffer)

		if len(v) > 0 {
			fmt.Fprintf(buffer, format+"\n", v...)
			buffer.WriteTo(wb)

			//wb.WriteString(fmt.Sprintf(format+"\n", v...))
		} else {
			fmt.Fprintf(buffer, format+"\n")
			buffer.WriteTo(wb)

			//wb.WriteString(fmt.Sprint(format+"\n"))
		}

		wb.Flush()

	} else {
		buffer := bufferPool.Get().(*bytes.Buffer)
		buffer.Reset()
		defer bufferPool.Put(buffer)

		if len(v) > 0 {
			fmt.Fprintf(buffer, format+"\n", v...)
			buffer.WriteTo(l.out)

			//l.out.Write([]byte(fmt.Sprintf(format+"\n", v...)))
		} else {
			fmt.Fprintf(buffer, format+"\n")
			buffer.WriteTo(l.out)

			//l.out.Write([]byte(fmt.Sprint(format+"\n")))
		}
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