package handlers

import "net/http"

func FileServer(w http.ResponseWriter, r *http.Request) {
	dir := "./web"
	http.ServeFile(w, r, dir+r.URL.Path)
}
