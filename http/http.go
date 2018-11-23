package http

import (
	"github.com/eagledb/eagledb/config"
	"net/http"
	"strconv"
)

type Server struct {
	ListenAddr string
	ListenPort int
}

func NewServer() *Server {
	return &Server{
		ListenAddr: config.G.Http.ListenAddr,
		ListenPort: config.G.Http.ListenPort,
	}
}

func (s *Server) Serve() error {
	http.HandleFunc("/point/write", s.WritePoint)
	http.HandleFunc("/point/query", s.QueryPoint)

	return http.ListenAndServe(s.ListenAddr+":"+strconv.Itoa(s.ListenPort), nil)
}

func (s *Server) WritePoint(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ASF"))
}

func (s *Server) QueryPoint(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ASF"))
}
