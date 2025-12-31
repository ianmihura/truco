package partials

import (
	"html/template"
	"net/http"
	"truco/pkg/ar"
	"truco/pkg/fsm"
	"truco/pkg/truco"
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

// Derived from truco.Card.
// We use different cards for UI
// to track selectable and unselectable cards
type CardUI struct {
	truco.Card
	OK bool
}

// Returns a copy of ALL_CARDS, excluding all cards in
// filter.MCards and filter.KCards
func GetAvailableCards(filter ar.FilterHands) []CardUI {
	excluded := make(map[truco.Card]bool, len(filter.KCards)*2)
	for _, c := range filter.KCards {
		excluded[c] = true
	}
	for _, c := range filter.MCards {
		excluded[c] = true
	}

	res := make([]CardUI, 0, len(truco.ALL_CARDS))
	for _, c := range truco.ALL_CARDS {
		res = append(res, CardUI{Card: c, OK: !excluded[c]})
	}
	return res
}
