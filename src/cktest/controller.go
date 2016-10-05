package cktest

import (
	"net/http"
	"env"
	"view"
)

func Test1(w http.ResponseWriter, r *http.Request) {
	slog := env.GetLogger()
	slog.Info("request test1")

	isComplete := make(chan bool)

	go view.RenderSimple(w, r, "test1.html", nil, isComplete)

	if <-isComplete {
		slog.Info("response test1")
	}
}
