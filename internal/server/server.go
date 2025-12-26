package server

import (
	"html/template"
	"net/http"
	"truco/internal/handlers/pages"
	"truco/internal/handlers/partials"
)

type Server struct {
	*http.ServeMux
	Tmpl *template.Template
}

func NewServer(tmpl *template.Template) *Server {
	s := &Server{
		ServeMux: http.NewServeMux(),
		Tmpl:     tmpl,
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	// Initialize Handlers
	homeHandler := pages.NewHomeHandler(s.Tmpl)
	trackerHandler := partials.NewTrackerHandler(s.Tmpl)

	// Serve static files
	fs := http.FileServer(http.Dir("web/static"))
	s.Handle("/static/", http.StripPrefix("/static/", fs))

	// Register Routes
	s.Handle("/", homeHandler)
	s.Handle("/track-act", trackerHandler)
}
