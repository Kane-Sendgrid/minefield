package api

import (
	"net/http"
	"encoding/json"

	"github.com/braintree/manners"
	"github.com/gorilla/mux"
)

type HTTPServer struct {
	server *manners.GracefulServer
}

type jsonResponse map[string]interface{}

func NewHTTPServer() *HTTPServer {
	s := &HTTPServer{}
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/mount", s.ApiMount)

	return s
}

func (s *HTTPServer) response(w http.ResponseWriter, j jsonResponse) {
	data, _ := json.Marshal(j)
	w.Write(data)
}

func (s *HTTPServer) ApiMount(w http.ResponseWriter, r *http.Request) {
	s.response(w, jsonResponse{
		"success": true,
	})
}
