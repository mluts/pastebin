package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	// Logging
	"github.com/unrolled/logger"

	// Stats/Metrics
	"github.com/rcrowley/go-metrics"
	"github.com/rcrowley/go-metrics/exp"
	"github.com/thoas/stats"

	"github.com/GeertJohan/go.rice"
	"github.com/julienschmidt/httprouter"
	"github.com/patrickmn/go-cache"
	"github.com/renstrom/shortuuid"
	"github.com/timewasted/go-accept-headers"
	"go.iondynamics.net/templice"
)

// AcceptedTypes ...
var AcceptedTypes = []string{
	"text/html",
	"text/plain",
}

// Counters ...
type Counters struct {
	r metrics.Registry
}

func NewCounters() *Counters {
	counters := &Counters{
		r: metrics.NewRegistry(),
	}
	return counters
}

func (c *Counters) Inc(name string) {
	metrics.GetOrRegisterCounter(name, c.r).Inc(1)
}

func (c *Counters) Dec(name string) {
	metrics.GetOrRegisterCounter(name, c.r).Dec(1)
}

func (c *Counters) IncBy(name string, n int64) {
	metrics.GetOrRegisterCounter(name, c.r).Inc(n)
}

func (c *Counters) DecBy(name string, n int64) {
	metrics.GetOrRegisterCounter(name, c.r).Dec(n)
}

// Server ...
type Server struct {
	bind      string
	config    Config
	store     *cache.Cache
	templates *templice.Template
	router    *httprouter.Router

	// Logger
	logger *logger.Logger

	// Stats/Metrics
	counters *Counters
	stats    *stats.Stats
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
		s.counters.Inc("n_index")

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
		s.counters.Inc("n_paste")

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

		u, err := url.Parse(fmt.Sprintf("./view/%s", uuid))
		if err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
		}
		http.Redirect(w, r, r.URL.ResolveReference(u).String(), http.StatusFound)
	}
}

// ViewHandler ...
func (s *Server) ViewHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		s.counters.Inc("n_view")

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

// StatsHandler ...
func (s *Server) StatsHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		bs, err := json.Marshal(s.stats.Data())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write(bs)
	}
}

// ListenAndServe ...
func (s *Server) ListenAndServe() {
	log.Fatal(
		http.ListenAndServe(
			s.bind,
			s.logger.Handler(
				s.stats.Handler(s.router),
			),
		),
	)
}

func (s *Server) initRoutes() {
	s.router.Handler("GET", "/debug/metrics", exp.ExpHandler(s.counters.r))
	s.router.GET("/debug/stats", s.StatsHandler())

	s.router.GET("/", s.IndexHandler())
	s.router.POST("/", s.PasteHandler())
	s.router.GET("/view/:uuid", s.ViewHandler())
}

// NewServer ...
func NewServer(bind string, config Config) *Server {
	server := &Server{
		bind:      bind,
		config:    config,
		router:    httprouter.New(),
		store:     cache.New(cfg.expiry, cfg.expiry*2),
		templates: templice.New(rice.MustFindBox("templates")),

		// Logger
		logger: logger.New(logger.Options{
			Prefix:               "pastebin",
			RemoteAddressHeaders: []string{"X-Forwarded-For"},
			OutputFlags:          log.LstdFlags,
		}),

		// Stats/Metrics
		counters: NewCounters(),
		stats:    stats.New(),
	}

	err := server.templates.Load()
	if err != nil {
		log.Panicf("error loading templates: %s", err)
	}

	server.initRoutes()

	return server
}
