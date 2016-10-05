package main

import (
	"env"
	"github.com/gorilla/mux"
	"net/http"
	"token"
)


const (
	configFile	= "config"
	configDir	= "/home/dhkim/funny-web"
)

const (
	URL_ROOT	= "/"
)

func main() {
	/*
	env.Init(configDir, configFile)
	slog := env.GetLogger()

	slog.Info("===== start funny web =====")

	mx := mux.NewRouter()

	mx.HandleFunc(URL_ROOT, token.Issue)

	port := env.GetConfig("server.port")

	slog.Info("http server listen port [%s]", port)

	if err := http.ListenAndServe(":" + port, mx); err != nil {
		slog.Err("failed http server listen port [%s]", port)
	}
	*/
}
