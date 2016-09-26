package token

import (
	"net/http"
	"view"
	"env"
)

func Issue(w http.ResponseWriter, r *http.Request)  {
	slog := env.GetLogger()
	slog.Info("request token issue")


	view.RenderSimple(w, r, "index")
}
