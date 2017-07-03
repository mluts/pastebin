package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/GeertJohan/go.rice"
	"github.com/julienschmidt/httprouter"
	"github.com/patrickmn/go-cache"
	"github.com/renstrom/shortuuid"
	"github.com/timewasted/go-accept-headers"
	"go.iondynamics.net/templice"
)

var AcceptedTypes = []string{
	"text/html",
	"text/plain",
}

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
		accepts, err := accept.Negotiate(
			r.Header.Get("Accept"), AcceptedTypes...,
		)
		if err != nil {
			log.Printf("error negotiating: %s", err)
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}

		switch accepts {
		case "text/html":
			s.render(w, "index", nil)
		case "text/plain":
		default:
		}
	}
}

// PasteHandler ...
func (s *Server) PasteHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var blob string

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}
		blob = string(body)

		err = r.ParseForm()
		if err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}

		if blob == "" {
			blob = r.Form.Get("blob")
		}

		if blob == "" {
			blob = r.URL.Query().Get("blob")
		}

		if blob == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		uuid := shortuuid.NewWithNamespace(s.config.fqdn)
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
		accepts, err := accept.Negotiate(
			r.Header.Get("Accept"), AcceptedTypes...,
		)
		if err != nil {
			log.Printf("error negotiating: %s", err)
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}

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

		switch accepts {
		case "text/html":
			s.render(w, "view", struct{ Blob string }{Blob: blob.(string)})
		case "text/plain":
			w.Write([]byte(blob.(string)))
		default:
			w.Write([]byte(blob.(string)))
		}
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
		store:     cache.New(cfg.expiry, cfg.expiry*2),
		templates: templice.New(rice.MustFindBox("templates")),
	}

	err := server.templates.Load()
	if err != nil {
		log.Panicf("error loading templates: %s", err)
	}

	server.initRoutes()

	return server
}
