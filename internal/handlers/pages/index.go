package pages

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"truco/pkg/truco"
)

type HomeHandler struct {
	Tmpl *template.Template
}

func NewHomeHandler(tmpl *template.Template) *HomeHandler {
	return &HomeHandler{Tmpl: tmpl}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.handleCalculate(w, r)
		return
	}

	data := struct {
		AllCards []string
	}{
		AllCards: getAllCardStrings(),
	}

	if err := h.Tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getAllCardStrings() []string {
	res := make([]string, 0, len(truco.ALL_CARDS))
	for _, c := range truco.ALL_CARDS {
		res = append(res, c.ToString())
	}
	return res
}

func (h *HomeHandler) handleCalculate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	mHandStr := r.Form.Get("mHand")
	kCardsStr := r.Form.Get("kCards")
	envido := uint8(255) // Default to unknown
	if env := r.Form.Get("envido"); env != "" {
		fmt.Sscanf(env, "%d", &envido)
	}
	isMHandFirst := r.Form.Get("isMHandFirst") == "true"
	hasStrategy := true // Default to true for this simplified UI

	mHand := parseCards(mHandStr)
	kCards := parseCards(kCardsStr)

	if len(mHand) != 3 {
		http.Error(w, "Select exactly 3 cards for your hand", http.StatusBadRequest)
		return
	}

	stats := truco.Hand(mHand).TrucoStrengthStats(kCards, []truco.Card{}, envido, isMHandFirst, hasStrategy)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.Tmpl.ExecuteTemplate(w, "results_partial.html", stats); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func parseCards(s string) []truco.Card {
	if s == "" {
		return []truco.Card{}
	}
	parts := strings.Split(s, ",")
	cards := make([]truco.Card, 0, len(parts))
	for _, p := range parts {
		if p == "" {
			continue
		}
		cards = append(cards, truco.NewCard(p))
	}
	return cards
}
