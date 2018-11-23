package server

import (
	"github.com/eagledb/eagledb/http"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start() error {
	http_server := http.NewServer()

	err := http_server.Serve()
	if err != nil {
		return err
	}

	return nil
}
