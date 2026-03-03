#include <stdint.h>
#include <stdio.h>

// Binary structure: SS M NNNNN
// M = can it be pieza (is 2,4,5,11,10,12)

uint8_t M_N  = 0b00011111; // Number Mask
uint8_t M_P  = 0b00100000; // Can be pieza Mask
uint8_t _2p  = 0b00000000; // Number
uint8_t _4p  = 0b00000001; // Number
uint8_t _5p  = 0b00000010; // Number
uint8_t _11p = 0b00000011; // Number
uint8_t _10p = 0b00000100; // Number
uint8_t _1e  = 0b00000101; // Number
uint8_t _1b  = 0b00000110; // Number
uint8_t _7e  = 0b00000111; // Number
uint8_t _7o  = 0b00001000; // Number
uint8_t _3   = 0b00001001; // Number
uint8_t _2   = 0b00101010; // Number
uint8_t _1f  = 0b00001011; // Number
uint8_t _12  = 0b00101100; // Number
uint8_t _11  = 0b00101101; // Number
uint8_t _10  = 0b00101110; // Number
uint8_t _7f  = 0b00001111; // Number
uint8_t _6   = 0b00010000; // Number
uint8_t _5   = 0b00110001; // Number
uint8_t _4   = 0b00110010; // Number
uint8_t M_S  = 0b11000000; // Suit Mask
uint8_t e    = 0b00000000; // Suit
uint8_t b    = 0b01000000; // Suit
uint8_t o    = 0b10000000; // Suit
uint8_t c    = 0b11000000; // Suit

// Card binary structure: SS M NNNNN
// S = Suit
// N = Number
// M = can it be pieza (2,4,5,11,10,12)
uint8_t CARDS[40] = {
    uint8_t(_1e),
    uint8_t(_1b),
    uint8_t(_7e),
    uint8_t(_7o),
    uint8_t(_3|e), uint8_t(_3|b), uint8_t(_3|o), uint8_t(_3|c),
    uint8_t(_2|e), uint8_t(_2|b), uint8_t(_2|o), uint8_t(_2|c),
    uint8_t(_1f|o), uint8_t(_1f|c),
    uint8_t(_12|e), uint8_t(_12|b), uint8_t(_12|o), uint8_t(_12|c),
    uint8_t(_11|e), uint8_t(_11|b), uint8_t(_11|o), uint8_t(_11|c),
    uint8_t(_10|e), uint8_t(_10|b), uint8_t(_10|o), uint8_t(_10|c),
    uint8_t(_7f|b), uint8_t(_7f|c),
    uint8_t(_6|e), uint8_t(_6|b), uint8_t(_6|o), uint8_t(_6|c),
    uint8_t(_5|e), uint8_t(_5|b), uint8_t(_5|o), uint8_t(_5|c),
    uint8_t(_4|e), uint8_t(_4|b), uint8_t(_4|o), uint8_t(_4|c),
};

int main() {

    // TODO shuffle CARDS

    uint8_t PACK[32];
    int64_t i = 0;

    // separate full deck of 40 cards into 5 iterations of 32 cards each
    for (int pack = 0; pack < 5; ++pack) {
        // TODO construct pack with cards at index % pack = 0

        // pick my three cards
        for (int m0 = 0; m0 < 32; ++m0) {
            for (int m1 = m0+1; m1 < 32; ++m1) {
                for (int m2 = m1+1; m2 < 32; ++m2) {
                    // uint8_t mHand[3] = { PACK[m0], PACK[m1], PACK[m2] }
                    
                    for (int m = m2+1; m < 32; ++m) {
                        // uint8_t mHand[3] = { to_pieza(PACK[m0], PACK[m]), to_pieza(PACK[m1], PACK[m]), to_pieza(PACK[m2], PACK[m]) }
                        
                        for (int m3 = m+1; m3 < 32; ++m3) {
                            for (int m4 = m3+1; m4 < 32; ++m4) {
                                for (int m5 = m4+1; m5 < 32; ++m5) {
                                    // uint8_t partnerHand[3] = { PACK[m3], PACK[m4], PACK[m5] }

                                    for (int o0 = m5+1; o0 < 32; ++o0) {
                                        for (int o1 = o0+1; o1 < 32; ++o1) {
                                            for (int o2 = o1+1; o2 < 32; ++o2) {
                                                for (int o3 = o2+1; o3 < 32; ++o3) {
                                                    for (int o4 = o3+1; o4 < 32; ++o4) {
                                                        for (int o5 = o4+1; o5 < 32; ++o5) {
                                                            // TODO play permutations
                                                            ++i;

                                                        }
                                                    }
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
    }

    printf("%lld\n",i);

    return 0;
}

// win or tie
bool wins_tie(uint8_t m, uint8_t o) {
    return (m&M_N) >= (o&M_N);
}

bool wins(uint8_t m, uint8_t o) {
    return (m&M_N) > (o&M_N);
}

uint8_t to_pieza(uint8_t card, uint8_t pieza) {
    if ((M_S & card == M_S & pieza) && (M_P & card)) {
        // same suit and can be pieza
        
        // TODO convert to a single level switch
        uint8_t cardN = card & M_N;
        if (cardN == _4) {
            return _4p;
        } else if (cardN == _5) {
            return _5p;
        } else if (cardN == _11) {
            return _11p;
        } else if (cardN == _10) {
            return _10p;
        } else if (cardN == _12) {

            uint8_t piezaN = pieza & M_N;
            if (piezaN == _4) {
                return _4p;
            } else if (piezaN == _5) {
                return _5p;
            } else if (piezaN == _11) {
                return _11p;
            } else if (piezaN == _10) {
                return _10p;
            } else if (piezaN == _2) {
                return _2p;
            } else {
                throw;
            }

        } else if (cardN == _2) {
            return _2p;
        } else {
            throw;
        }
    } else {
        return card; // not pieza
    }
}
