package handlers

import (
	"net/http"

	"forum/config"
)

func IndexHandler(rw http.ResponseWriter, r *http.Request) {
	config.TMPL.Render(rw, "index.html", nil)
}
