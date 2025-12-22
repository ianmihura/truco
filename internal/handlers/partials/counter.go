package partials

import (
	"fmt"
	"net/http"
	"sync"
)

type CounterHandler struct {
	mu      sync.Mutex
	counter int
}

func NewCounterHandler() *CounterHandler {
	return &CounterHandler{}
}

func (h *CounterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.mu.Lock()
	h.counter++
	currentVal := h.counter
	h.mu.Unlock()

	// Return just the new number as HTML (HTMX swaps this into the target)
	fmt.Fprintf(w, "%d", currentVal)
}
