package truco

import (
	"testing"
)

func TestFlorUY(t *testing.T) {
	tests := []struct {
		handStr string
		muestra string
		want    uint8
	}{
		// 3 cards of same suit (no pieces)
		{"1e 2e 3e", "1b", 226},    // 1+2+3 + 220 = 226
		{"10e 11e 12e", "1b", 220}, // 0+0+0 + 220 = 220
		{"7e 6e 5e", "1b", 238},    // 7+6+5 + 220 = 238
		{"7e 6e 1e", "1b", 234},    // 7+6+1 + 220 = 234

		// 1 piece + 2 same suit
		{"2e 1b 7b", "1e", 238},   // 2e is pieza(10) + 1b(1) + 7b(7) + 220 = 238
		{"4e 10b 11b", "1e", 229}, // 4e is pieza(9) + 10b(0) + 11b(0) + 220 = 229
		{"12e 1b 2b", "2e", 233},  // 12e is pieza(10) + 1b(1) + 2b(2) + 220 = 233

		// 2 pieces + 1 card
		{"2e 4e 1c", "1e", 240},   // 2e(10) + 4e(9) + 1c(1) + 220 = 240
		{"11e 10e 1b", "1e", 235}, // 11e(7) + 10e(7) + 1b(1) + 220 = 235
		{"12e 5e 1b", "2e", 239},  // 12e(10) + 5e(8) + 1b(1) + 220 = 239

		// 3 pieces
		{"2e 4e 5e", "1e", 247},   // 10+9+8 + 220 = 247
		{"11e 10e 2e", "1e", 244}, // 7+7+10 + 220 = 244

		// Not flor
		{"1e 2b 3c", "1o", 3},
		{"2e 1b 1c", "1o", 2},  // 1 piece but others not same suit
		{"1e 2e 1b", "1b", 23}, // 2 same suit but no piece
	}

	for _, tt := range tests {
		t.Run(tt.handStr+"_m_"+tt.muestra, func(t *testing.T) {
			h := NewHand(tt.handStr)
			m := NewCard(tt.muestra)
			got := h.EnvidoUY(m)
			if got != tt.want {
				t.Errorf("Hand(%s).EnvidoUY(%s) = %d, want %d", tt.handStr, tt.muestra, got, tt.want)
			}
		})
	}
}
