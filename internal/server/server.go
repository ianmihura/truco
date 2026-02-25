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
	matrixHandler := pages.NewMatrixHandler(s.Tmpl)
	homeHandler := pages.NewHomeHandler(s.Tmpl)
	handler := partials.NewHandler(s.Tmpl)

	// Serve static files
	fs := http.FileServer(http.Dir("web/static"))
	s.Handle("/static/", http.StripPrefix("/static/", fs))

	// Main app
	s.Handle("/", homeHandler)

	// Matrix
	s.Handle("/matrix", matrixHandler)
	s.HandleFunc("/get-cards", handler.GetCards)
	s.HandleFunc("/get-lower-cards", handler.GetLowerCards)
	s.HandleFunc("/track-act", handler.TrackAct)
	s.HandleFunc("/track-stats", handler.TrackStats)
}
