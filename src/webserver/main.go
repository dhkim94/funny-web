package main

import (
	"fmt"
	"env"
	"github.com/gorilla/mux"
	"net/http"
	"token"
)


const (
	configFile	= "config"
	configDir	= "/Users/dhkim/Documents/Develop/project.my/funny-web"
)

const (
	URL_ROOT	= "/"
)

func main() {
	env.Init(configDir, configFile)
	slog := env.GetLogger()

	slog.Info("===== start funny web =====")

	mx := mux.NewRouter()

	mx.HandleFunc(URL_ROOT, token.Issue)

	port := fmt.Sprintf("%s", env.GetConfig("server.port"))

	slog.Info("http server listen port [%s]", port)

	err := http.ListenAndServe(":" + port, mx)
	if err != nil {
		slog.Error("failed http server listen port [%s]", port)
	}
}
