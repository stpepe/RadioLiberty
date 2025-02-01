package server

import (
	"net/http"
)

const StaticPath = "static/layouts/"

func (s *Server) Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, StaticPath+"index.html")
}
