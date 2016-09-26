package token

import (
	"net/http"
	"fmt"
	"view"
)

func Issue(w http.ResponseWriter, r *http.Request)  {

	fmt.Println("-----token Issue")

	view.RenderSimple(w, r, "index")
}
