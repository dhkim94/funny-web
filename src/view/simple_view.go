package view

import (
	"net/http"
	"env"
	"time"
	"bufio"
	"os"
	"crypto/md5"
	"sync"
	"html/template"
)

type htmlFiles struct {
	modify time.Time
	content string
}

var rwMu *sync.RWMutex
var htmlTemplates map[string]*htmlFiles

func upsertHtmlText(filePath string) (out *string, oerr error) {
	// intialize
	if len(htmlTemplates) <= 0 {
		htmlTemplates = make(map[string]*htmlFiles)
		rwMu = new(sync.RWMutex)
	}

	slog := env.GetLogger()
	tmplPath := env.GetValue("template.path")

	fullPath := tmplPath + "/" + filePath

	fd, err := os.OpenFile(fullPath,
		os.O_RDONLY, os.FileMode(0644))
	if err != nil {
		slog.Err("not found template file [%s/%s]", tmplPath, filePath)
		out, oerr = nil, err
		return
	}
	defer fd.Close()

	info, _ := fd.Stat()

	_key := md5.Sum([]byte(fullPath))
	key := string(_key[:])

	rwMu.RLock()
	if htmlTemplates[key] != nil  && info.ModTime().Equal(htmlTemplates[key].modify) {
		out, oerr = &htmlTemplates[key].content, nil
		rwMu.RUnlock()
		return
	}
	rwMu.RUnlock()

	slog.Debug("read template file [%s]", fullPath)

	reader := bufio.NewReader(fd)
	buf := make([]byte, info.Size())
	reader.Read(buf)

	content := &htmlFiles{
		modify: info.ModTime(),
		content: string(buf),
	}

	rwMu.Lock()
	htmlTemplates[key] = content
	out, oerr = &htmlTemplates[key].content, nil
	rwMu.Unlock()

	return
}

// show simple html template
// data : nullable value
func RenderSimple(w http.ResponseWriter, r *http.Request, fname string, data interface{}, isComplete chan bool) {
	content, err := upsertHtmlText(fname)

	// todo header 설정 해 줘야 한다.

	if err != nil {

		// todo error 출력에 따른 헤더 설정 필요하면 여기서 추가 한다.

		w.Write([]byte("not found file: " + fname))
		isComplete <- false
		return
	}

	tmpl, err2 := template.New(fname).Parse(*content)
	if err2 != nil {

		// todo error 출력에 따른 헤더 설정 필요하면 여기서 추가 한다.

		w.Write([]byte("template error. file: " + fname))
		isComplete <- false
		return
	}

	//t1.Execute(w, nil)
	if err := tmpl.Execute(w, data); err != nil {
		// todo error 출력에 따른 헤더 설정 필요하면 여기서 추가 한다.

		w.Write([]byte("template execute error. file: " + fname))
		isComplete <- false
		return
	}

	isComplete <- true
}
