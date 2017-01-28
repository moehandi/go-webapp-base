package controller

import (
	"net/http"
	"strings"
)

// Static maps static files for julienschmidt/httpRouter
func Static(w http.ResponseWriter, r *http.Request) {
	// Disable listing directories
	if strings.HasSuffix(r.URL.Path, "/") {
		Error404(w, r)
		return
	}
	http.ServeFile(w, r, r.URL.Path[1:])
}

func GorillaStatic(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/") {
		Error404(w, r)
		return
	}
	http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
}