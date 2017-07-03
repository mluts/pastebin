package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/GeertJohan/go.rice"
	"github.com/julienschmidt/httprouter"
	"github.com/patrickmn/go-cache"
	"github.com/renstrom/shortuuid"
	"go.iondynamics.net/templice"
)

// Server ...
type Server struct {
	bind      string
	config    Config
	store     *cache.Cache
	templates *templice.Template
	router    *httprouter.Router
}

func (s *Server) render(w http.ResponseWriter, tmpl string, data interface{}) {
	err := s.templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// IndexHandler ...
func (s *Server) IndexHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		s.render(w, "index", nil)
	}
}

// PasteHandler ...
func (s *Server) PasteHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}

		uuid := shortuuid.NewWithNamespace(s.config.fqdn)
		blob := r.Form.Get("blob")
		s.store.Set(uuid, blob, cache.DefaultExpiration)

		u, err := url.Parse(fmt.Sprintf("./%s", uuid))
		if err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
		}
		http.Redirect(w, r, r.URL.ResolveReference(u).String(), http.StatusFound)
	}
}

// ViewHandler ...
func (s *Server) ViewHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		uuid := p.ByName("uuid")
		if uuid == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		blob, ok := s.store.Get(uuid)
		if !ok {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		s.render(w, "view", struct{ Blob string }{Blob: blob.(string)})
	}
}

// ListenAndServe ...
func (s *Server) ListenAndServe() {
	log.Fatal(http.ListenAndServe(s.bind, s.router))
}

func (s *Server) initRoutes() {
	s.router.GET("/", s.IndexHandler())
	s.router.POST("/", s.PasteHandler())
	s.router.GET("/:uuid", s.ViewHandler())
}

// NewServer ...
func NewServer(bind string, config Config) *Server {
	server := &Server{
		bind:      bind,
		config:    config,
		router:    httprouter.New(),
		store:     cache.New(5*time.Minute, 10*time.Minute),
		templates: templice.New(rice.MustFindBox("templates")),
	}

	err := server.templates.Load()
	if err != nil {
		log.Panicf("error loading templates: %s", err)
	}

	server.initRoutes()

	return server
}
