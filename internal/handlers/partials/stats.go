package partials

import (
	"encoding/json"
	"net/http"
	"truco/pkg/ar"
)

func (h *Handler) TrackStats(w http.ResponseWriter, r *http.Request) {
	fmatrixParam := r.URL.Query().Get("fmatrix")
	match := GetMatch(r)

	// Recalculate stats dynamically based on the current matrix mode
	stats, err := ar.ComputePairStats(fmatrixParam == "true", match.GetStatsFilter())
	if err != nil {
		http.Error(w, "Failed to compute stats: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, "Failed to marshal stats: "+err.Error(), http.StatusInternalServerError)
	}
}
