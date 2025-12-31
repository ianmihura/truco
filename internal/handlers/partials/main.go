package partials

import (
	"html/template"
	"net/http"
	"truco/pkg/fsm"
)

type TrackerData struct {
	ActionTitle string
	Actions     []fsm.ValidAction
	DoneActions []fsm.ValidAction
	State       string
	PlayedCard  string
	Stats       template.JS
}

type Handler struct {
	tmpl *template.Template
}

func NewHandler(tmpl *template.Template) *Handler {
	return &Handler{tmpl: tmpl}
}

// Returns a fsm.Match object from the state in queryparams,
// or a new empty fsm.Match
func GetMatch(r *http.Request) *fsm.Match {
	stateParam := r.URL.Query().Get("state")
	if stateParam == "" {
		return fsm.NewMatch()
	} else {
		return fsm.Decode([]byte(stateParam))
	}
}
