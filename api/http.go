package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Kane-Sendgrid/minefield/fs"
	"github.com/braintree/manners"
	"github.com/gorilla/mux"
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"os"
)

type HTTPServer struct {
	server  *manners.GracefulServer
	systems map[string]*fileSystem
}

type fileSystem struct {
	mountPoint string
	fs         *fs.Fs
	server     *fuse.Server
}

func newFileSystem(mountPoint string) (*fileSystem, error) {
	f := &fileSystem{
		mountPoint: mountPoint,
	}
	fs := &fs.Fs{FileSystem: pathfs.NewDefaultFileSystem(),
		Files: map[string]*fs.File{}}
	nfs := pathfs.NewPathNodeFs(fs, nil)
	os.MkdirAll(mountPoint, 0777)
	server, _, err := nodefs.MountRoot(mountPoint, nfs.Root(), nil)
	f.fs = fs
	f.server = server
	if err != nil {
		return nil, err
	}
	go f.server.Serve()
	return f, nil
}

func (f fileSystem) Unmount() error {
	return f.server.Unmount()
}

type jsonResponse map[string]interface{}

func NewHTTPServer() *HTTPServer {
	s := &HTTPServer{
		server:  manners.NewServer(),
		systems: map[string]*fileSystem{},
	}
	return s
}

func (s *HTTPServer) Serve() {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/mount", s.ApiMount)
	api.HandleFunc("/unmount", s.ApiUnmount)
	api.HandleFunc("/detonate", s.ApiDetonate)
	api.HandleFunc("/defuse", s.ApiDefuse)

	log.Println("starting api server on :7725...")
	s.server.ListenAndServe(":7725", r)
}

func (s *HTTPServer) Shutdown() {
	for _, s := range s.systems {
		fmt.Println("Unmounting", s.mountPoint)
		s.Unmount()
	}
	s.server.Shutdown <- true
}

func (s *HTTPServer) response(w http.ResponseWriter, j jsonResponse) {
	data, _ := json.Marshal(j)
	w.Write(data)
}

func (s *HTTPServer) ApiMount(w http.ResponseWriter, r *http.Request) {
	mountPoint := r.FormValue("mountpoint")
	if len(mountPoint) == 0 {
		s.response(w, jsonResponse{
			"success": false,
			"error":   "mountpoint is required",
		})
		return
	}
	fs, err := newFileSystem(mountPoint)
	if err != nil {
		s.response(w, jsonResponse{
			"success": false,
			"error":   err,
		})
		return
	}
	s.systems[mountPoint] = fs
	s.response(w, jsonResponse{
		"success": true,
	})
}

func (s *HTTPServer) ApiUnmount(w http.ResponseWriter, r *http.Request) {
	mountPoint := r.FormValue("mountpoint")
	if fs, ok := s.systems[mountPoint]; ok {
		err := fs.Unmount()
		if err == nil {
			delete(s.systems, mountPoint)
			s.response(w, jsonResponse{
				"success": true,
			})
			return
		} else {
			s.response(w, jsonResponse{
				"success": false,
				"error":   err,
			})
			return
		}
	}
	s.response(w, jsonResponse{
		"success": false,
		"error":   "not found",
	})
}

func (s *HTTPServer) ApiDetonate(w http.ResponseWriter, r *http.Request) {
	mountPoint := r.FormValue("mountpoint")
	if fs, ok := s.systems[mountPoint]; ok {
		fs.fs.Detonate("")
		s.response(w, jsonResponse{
			"success": true,
		})
		return
	}

	s.response(w, jsonResponse{
		"success": false,
		"error":   "not found",
	})
}

func (s *HTTPServer) ApiDefuse(w http.ResponseWriter, r *http.Request) {
	mountPoint := r.FormValue("mountpoint")
	if fs, ok := s.systems[mountPoint]; ok {
		fs.fs.Defuse("")
		s.response(w, jsonResponse{
			"success": true,
		})
		return
	}

	s.response(w, jsonResponse{
		"success": false,
		"error":   "not found",
	})
}