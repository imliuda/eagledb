package server

import (
	"github.com/eagledb/eagledb/config"
	"github.com/eagledb/eagledb/point"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var decoder = schema.NewDecoder()

type HttpServer struct {
	ListenAddr string
	ListenPort int

	server *Server
}

func NewHttpServer() *HttpServer {
	return &HttpServer{
		ListenAddr: config.G.Http.ListenAddr,
		ListenPort: config.G.Http.ListenPort,
	}
}

func (s *HttpServer) Serve() error {
	r := mux.NewRouter()

	r.HandleFunc("/point/write", s.WritePoint).Methods("POST")
	r.HandleFunc("/point/query", s.QueryPoint).Methods("GET")

	go func() {
		err := http.ListenAndServe(s.ListenAddr+":"+strconv.Itoa(s.ListenPort), nil)
		if err != nil {
			os.Exit(1)
		}
	}()

	return nil
}

func (s *HttpServer) WritePoint(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	form := struct {
		Database string
	}{}

	err = decoder.Decode(form, r.PostForm)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	db := s.server.GetDatabase(form.Database)
	if db == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "write error", http.StatusInternalServerError)
		return
	}

	points, err := point.Parse(body)
	if err != nil {
		http.Error(w, "parse error", http.StatusBadRequest)
		return
	}

	err = db.WritePoints(points)
	if err != nil {
		http.Error(w, "write error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("ASF"))
}

func (s *HttpServer) QueryPoint(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ASF"))
}
