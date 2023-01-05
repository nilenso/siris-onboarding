package server

import (
	"net/http"
)

func (s *Server) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
