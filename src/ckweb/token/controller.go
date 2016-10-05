package token

import (
	"net/http"
	"ckweb/view"
	"cklib/env"
)

type Sam struct {
	Name string
	Gender string
}

func Issue(w http.ResponseWriter, r *http.Request)  {
	slog := env.GetLogger()
	slog.Info("request token issue")

	data := &Sam{
		Name: "dhkim1",
		Gender: "sssss",
	}

	isComplete := make(chan bool)

	go view.RenderSimple(w, r, "hello.html", data, isComplete)

	if <-isComplete {
		slog.Info("response token issue")
	}
}
