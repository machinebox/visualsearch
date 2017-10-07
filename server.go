package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/machinebox/sdk-go/tagbox"
	"github.com/matryer/way"
)

// Server is the app server.
type Server struct {
	assets string
	tagbox *tagbox.Client
	items  map[string]Item
	router *way.Router
}

// NewServer makes a new Server.
func NewServer(assets string, tagbox *tagbox.Client, items map[string]Item) *Server {
	srv := &Server{
		assets: assets,
		tagbox: tagbox,
		items:  items,
		router: way.NewRouter(),
	}
	srv.router.Handle(http.MethodGet, "/assets/", Static("/assets/", assets))
	srv.router.HandleFunc(http.MethodGet, "/api/random-images", srv.handleRandomImages)
	srv.router.HandleFunc(http.MethodGet, "/api/similar-images", srv.handleSimilarImages)
	srv.router.HandleFunc(http.MethodGet, "/", srv.handleIndex)
	return srv
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(s.assets, "index.html"))
}

func (s *Server) handleSimilarImages(w http.ResponseWriter, r *http.Request) {
	urlStr := r.URL.Query().Get("url")
	u, err := url.Parse(urlStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tags, err := s.tagbox.SimilarURL(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var res struct {
		Items []Item `json:"items"`
	}
	for _, tag := range tags {
		item := s.items[tag.ID]
		item.Confidence = tag.Confidence
		res.Items = append(res.Items, item)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleRandomImages(w http.ResponseWriter, r *http.Request) {
	var res struct {
		Items []Item `json:"items"`
	}
	var count int
	for _, v := range s.items {
		res.Items = append(res.Items, v)
		count++
		if count == 20 {
			break
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Static gets a static file server for the specified path.
func Static(stripPrefix, dir string) http.Handler {
	h := http.StripPrefix(stripPrefix, http.FileServer(http.Dir(dir)))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}
