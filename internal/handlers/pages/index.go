package pages

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
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
	mode := r.Form.Get("mode")
	muestraStr := r.Form.Get("muestra")
	isMHandFirst := r.Form.Get("isMHandFirst") == "true"
	hasStrategy := r.Form.Get("hasStrategy") == "true"
	sonBuenas := r.Form.Get("sonBuenas") == "true"
	flor := r.Form.Get("flor") == "true"
	kEnvido, err := strconv.Atoi(r.Form.Get("envido"))
	if err != nil {
		kEnvido = 255
	}

	mHand := truco.NewHand(mHandStr)
	if len(mHand) != 3 {
		http.Error(w, "Select exactly 3 cards for your hand", http.StatusBadRequest)
		return
	}
	kCards := []truco.Card(truco.NewHand(kCardsStr))

	var muestra truco.Card
	if mode == "UY" {
		if muestraStr == "" {
			http.Error(w, "Select a muestra", http.StatusBadRequest)
			return
		}
		muestra = truco.NewHand(muestraStr)[0]
	}

	if sonBuenas {
		if mode == "UY" {
			kEnvido = 100 + int(mHand.EnvidoUY(muestra))
		} else {
			kEnvido = 100 + int(mHand.Envido())
		}
	} else if flor {
		kEnvido = 200
	}

	var stats truco.TrucoStats
	if mode == "UY" {
		stats = mHand.TrucoStrengthStatsUY(kCards, []truco.Card{muestra}, uint8(kEnvido), isMHandFirst, hasStrategy)
	} else {
		stats = mHand.TrucoStrengthStats(kCards, []truco.Card{}, uint8(kEnvido), isMHandFirst, hasStrategy)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.Tmpl.ExecuteTemplate(w, "results_partial.html", stats); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
