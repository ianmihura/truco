package truco

// list of all cards
var ALL_CARDS = []Card{
	{1, 'e'}, {1, 'b'},
	{7, 'e'}, {7, 'o'},
	{3, 'e'}, {3, 'b'}, {3, 'o'}, {3, 'c'},
	{2, 'e'}, {2, 'b'}, {2, 'o'}, {2, 'c'},
	{1, 'o'}, {1, 'c'},
	{12, 'e'}, {12, 'b'}, {12, 'o'}, {12, 'c'},
	{11, 'e'}, {11, 'b'}, {11, 'o'}, {11, 'c'},
	{10, 'e'}, {10, 'b'}, {10, 'o'}, {10, 'c'},
	{7, 'b'}, {7, 'c'},
	{6, 'e'}, {6, 'b'}, {6, 'o'}, {6, 'c'},
	{5, 'e'}, {5, 'b'}, {5, 'o'}, {5, 'c'},
	{4, 'e'}, {4, 'b'}, {4, 'o'}, {4, 'c'},
}

// list of all figures
var FIGURES = []Card{
	{10, 'e'}, {10, 'b'}, {10, 'o'}, {10, 'c'},
	{11, 'e'}, {11, 'b'}, {11, 'o'}, {11, 'c'},
	{12, 'e'}, {12, 'b'}, {12, 'o'}, {12, 'c'},
}

// RANKS_AR =
// "1e",
// "1b",
// "7e",
// "7o",
// "3",
// "2",
// "1f",
// "12",
// "11",
// "10",
// "7f",
// "6",
// "5",
// "4"

var PIEZAS = []Card{
	{2, 'p'}, {4, 'p'}, {5, 'p'}, {11, 'p'}, {10, 'p'},
}

// RANKS_UY =
// "2p",
// "4p",
// "5p",
// "11p",
// "10p",
// "1e",
// "1b",
// "7e",
// "7o",
// "3",
// "2",
// "1f",
// "12",
// "11",
// "10",
// "7f",
// "6",
// "5",
// "4"

// maps every card to a relative score for game truco:
// - if TRUCO[i] > TRUCO[j], then i beats j in truco
// - same score tie
var TRUCO = map[Card]uint8{
	{2, 'p'}:  19,
	{4, 'p'}:  18,
	{5, 'p'}:  17,
	{11, 'p'}: 16,
	{10, 'p'}: 15,
	{1, 'e'}:  14,
	{1, 'b'}:  13,
	{7, 'e'}:  12,
	{7, 'o'}:  11,
	{3, 'e'}:  10, {3, 'b'}: 10, {3, 'o'}: 10, {3, 'c'}: 10,
	{2, 'e'}: 9, {2, 'b'}: 9, {2, 'o'}: 9, {2, 'c'}: 9,
	{1, 'c'}: 8, {1, 'o'}: 8,
	{12, 'e'}: 7, {12, 'b'}: 7, {12, 'o'}: 7, {12, 'c'}: 7,
	{11, 'e'}: 6, {11, 'b'}: 6, {11, 'o'}: 6, {11, 'c'}: 6,
	{10, 'e'}: 5, {10, 'b'}: 5, {10, 'o'}: 5, {10, 'c'}: 5,
	{7, 'c'}: 4, {7, 'b'}: 4,
	{6, 'e'}: 3, {6, 'b'}: 3, {6, 'o'}: 3, {6, 'c'}: 3,
	{5, 'e'}: 2, {5, 'b'}: 2, {5, 'o'}: 2, {5, 'c'}: 2,
	{4, 'e'}: 1, {4, 'b'}: 1, {4, 'o'}: 1, {4, 'c'}: 1,
}

// Slightly more efficient way to get truco value of a card
func GetTruco(c Card) uint8 {
	switch c.N {
	case 1:
		switch c.S {
		case 'e':
			return 14
		case 'b':
			return 13
		default:
			return 8
		}
	case 2:
		return 9
	case 3:
		return 10
	case 4:
		return 1
	case 5:
		return 2
	case 6:
		return 3
	case 7:
		switch c.S {
		case 'e':
			return 12
		case 'o':
			return 11
		default:
			return 4
		}
	case 10:
		return 5
	case 11:
		return 6
	case 12:
		return 7
	default:
		return 0
	}
}

// Slightly more efficient way to get truco value of a card
func GetTrucoUY(c, m Card) uint8 {
	switch c.N {
	case 1:
		switch c.S {
		case 'e':
			return 14
		case 'b':
			return 13
		default:
			return 8
		}
	case 2:
		if c.S == m.S {
			return 19
		} else {
			return 9
		}
	case 3:
		return 10
	case 4:
		if c.S == m.S {
			return 18
		} else {
			return 1
		}
	case 5:
		if c.S == m.S {
			return 17
		} else {
			return 2
		}
	case 6:
		return 3
	case 7:
		switch c.S {
		case 'e':
			return 12
		case 'o':
			return 11
		default:
			return 4
		}
	case 10:
		if c.S == m.S {
			return 15
		} else {
			return 5
		}
	case 11:
		if c.S == m.S {
			return 16
		} else {
			return 6
		}
	case 12:
		if c.S == m.S {
			switch m.N {
			case 2:
				return 19
			case 4:
				return 18
			case 5:
				return 17
			case 11:
				return 16
			case 10:
				return 15
			default:
				return 7
			}
		}
		return 7
	default:
		return 0
	}
}

const MAX_ENVIDO_AR = 33
const MAX_ENVIDO_UY = 37

// Possible cards needed for an envido,
// indexed by the envido amount,
// mapped to a list of possible card number (Card.N)
//
//	1-7
//	f -> 10,11,12
var ENVIDOS_AR = map[uint8][][]uint8{
	0:  {{'f'}},
	1:  {{1}},
	2:  {{2}},
	3:  {{3}},
	4:  {{4}},
	5:  {{5}},
	6:  {{6}},
	7:  {{7}},
	20: {{'f', 'f'}},
	21: {{'f', 1}},
	22: {{'f', 2}},
	23: {{'f', 3}, {1, 2}},
	24: {{'f', 4}, {1, 3}},
	25: {{'f', 5}, {1, 4}, {2, 3}},
	26: {{'f', 6}, {1, 5}, {2, 4}},
	27: {{'f', 7}, {1, 6}, {2, 5}, {3, 4}},
	28: {{1, 7}, {2, 6}, {3, 5}},
	29: {{2, 7}, {3, 6}, {4, 5}},
	30: {{3, 7}, {4, 6}},
	31: {{4, 7}, {5, 6}},
	32: {{5, 7}},
	33: {{6, 7}},
}

// Possible cards needed for an envido,
// indexed by the envido amount,
// mapped to a list of possible card number (Card.N)
//
//		1-7
//		f -> 10,11,12
//	 0x2a -> 2p
//	 0x4a -> 4p
//	 0x5a -> 5p
//	 0x1a -> 11p,10p
var ENVIDOS_UY = map[uint8][][]uint8{
	0:  {{'f'}},
	1:  {{1}},
	2:  {{2}},
	3:  {{3}},
	4:  {{4}},
	5:  {{5}},
	6:  {{6}},
	7:  {{7}},
	20: {{'f', 'f'}},
	21: {{'f', 1}},
	22: {{'f', 2}},
	23: {{'f', 3}, {1, 2}},
	24: {{'f', 4}, {1, 3}},
	25: {{'f', 5}, {1, 4}, {2, 3}},
	26: {{'f', 6}, {1, 5}, {2, 4}},
	27: {{'f', 7}, {1, 6}, {2, 5}, {3, 4}, {0x1a}},
	28: {{1, 7}, {2, 6}, {3, 5}, {0x5a}, {0x1a, 1}},
	29: {{2, 7}, {3, 6}, {4, 5}, {0x4a}, {0x5a, 1}, {0x1a, 2}},
	30: {{3, 7}, {4, 6}, {0x2a}, {0x4a, 1}, {0x5a, 2}, {0x1a, 3}},
	31: {{4, 7}, {5, 6}, {0x2a, 1}, {0x4a, 2}, {0x5a, 3}, {0x1a, 4}},
	32: {{5, 7}, {0x2a, 2}, {0x4a, 3}, {0x5a, 4}, {0x1a, 5}},
	33: {{6, 7}, {0x2a, 3}, {0x4a, 4}, {0x5a, 5}, {0x1a, 6}},
	34: {{0x2a, 4}, {0x4a, 5}, {0x5a, 6}, {0x1a, 7}, {0x0a, 7}},
	35: {{0x2a, 5}, {0x4a, 6}, {0x5a, 7}},
	36: {{0x2a, 6}, {0x4a, 7}},
	37: {{0x2a, 7}},
}

// envido each pieza card brings
var ENVIDO_PIEZA = map[Card]uint8{
	{2, 'p'}:  10,
	{4, 'p'}:  9,
	{5, 'p'}:  8,
	{11, 'p'}: 7,
	{10, 'p'}: 7,
}
