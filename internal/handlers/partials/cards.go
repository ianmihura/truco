package partials

import (
	"net/http"
	"slices"
	"truco/pkg/fsm"
	"truco/pkg/truco"
)

func (h *Handler) GetCards(w http.ResponseWriter, r *http.Request) {
	match := GetMatch(r)

	slices.SortFunc(truco.ALL_CARDS, truco.SortForTruco)
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
