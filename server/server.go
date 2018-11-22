package server

import (
	"github.com/eagledb/eagledb"
)

type Server struct {
	http ealgedb.Http
}

func (s *Server) Start() {
	s.http.server = s
	http.Serve()
}
