package server

import (
	"fmt"
	"os"
)

type Server struct {
	databases []*Database
	http      *HttpServer
}

func NewServer() *Server {
	server := &Server{}
	httpServer := &HttpServer{
		server: server,
	}
	server.http = httpServer
	return server
}

func (s *Server) Start() error {
	errChan := make(chan error)

	go func() {
		err := s.http.Serve()
		if err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
}

func (s *Server) GetDatabase(name string) *Database {
	for _, db := range s.databases {
		if db.name == name {
			return db
		}
	}

	return nil
}
