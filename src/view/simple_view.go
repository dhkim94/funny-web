package view

import (
	"net/http"
	"github.com/yosssi/ace"
	"env"
)

func RenderSimple(w http.ResponseWriter, r *http.Request, tmpl string) {
	tmplPath := env.GetConfig("template.path")

	tpl, err := ace.Load(tmpl, "", &ace.Options{
		BaseDir: tmplPath,
		Extension: "jade",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tpl.Execute(w, map[string]string{"Msg": "Hello Ace"}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
