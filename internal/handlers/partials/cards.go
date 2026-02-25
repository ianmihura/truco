package partials

import (
	"fmt"
	"net/http"
	"truco/pkg/fsm"
	"truco/pkg/truco"
)

func (h *Handler) GetCards(w http.ResponseWriter, r *http.Request) {
	match := GetMatch(r)

	data := struct {
		Cards  []CardUI
		State  string
		Action fsm.ValidAction
	}{
		Cards:  GetAvailableCards(match.GetStatsFilter()),
		State:  string(match.Encode()),
		Action: fsm.PLAY,
	}

	err := h.tmpl.ExecuteTemplate(w, "cards", data)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) GetLowerCards(w http.ResponseWriter, r *http.Request) {
	card1 := r.URL.Query().Get("card1")
	card2 := r.URL.Query().Get("card2")
	t1 := truco.RANKS[card1]
	t2 := truco.RANKS[card2]

	data := struct {
		Ranks []string
	}{
		Ranks: truco.CardsLowerEqual(min(t1, t2)),
	}

	fmt.Println(truco.CardsLowerEqual(min(t1, t2)))

	err := h.tmpl.ExecuteTemplate(w, "cards_lower", data)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
	}
}
