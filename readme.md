# Solver for the game Truco, as played in Argentina and Uruguay

### Features
1. Whats the chance of each sub-hand (given two best cards, with or without envido) -- only relevant if we have UI
2. Given your hand
    - Calc your envido / flor points (easy)
    - Calc 'truco' strength
        - given sorted cards played against each other, against how many hands do you win
        - given unsorted, against how many hands do you win
    So we get:
    - Whats the chance that your hand is better than one (or more) oponents
    - Trailing: given we fix some piece of information, how do these probabilites change
3. Given your envido
    - Whats the chance your envido/flor is best (medium)
    - Whats the range of hands you could have (hard)
    - Bonus: what card reveals less, or contradictory, info (eg. 33 envido with 2m + 3, showing any 7 will bluff)
    - Trailing: given we fix some piece of information, how do these probabilites change
4. Guess range of other players given limited info
    - Flor (score)
    - Envido (score or partial score)
    - Individual cards

### TODO UI:
- choose your hand similar to poker GTO (see .odt)
- choose mCards and kCards out of a 4x10 matrix
    then also choose an envido
- Separate between arg & uru
- Make it easier for user to choose their hand
    - table of best 2 cards (ar)14x14 or (uy)14x19
    - sub-choose third (worst) card
    - sub-sub-choose envido
    - 2-3 click choose muestra
- How many players?

### TODO backend 
- chance that your hand is best in table, given known info
- chance that your envido is best in table, given known info
- TrucoStrength does not capture the practical stregth, because there are some permutations that will never be reasoanably played. Eg:
    - Tie the first round, you should play your strongest card right away.
    - If you lost the first round, you should not tie any other round, unless you're loosing anyway.
    Meaning, some hands seem stronger, or weaker, than they should. My intuition is that real strength polarizes scores even more: as strong hands would not loose is dumb ways, and viceversa (as score is calculated by averaging agains all possible hands). Mid-range hands would also tend to stay mid-range. I expect there to be no extreme cases where this drastically changes a hand's overall score.
- EV features (harder):
    1. EV of playing your envido. Considering
        - You may lose
        - You give away information
    2. EV of calling truco. Considering
        - information you have on others (their range)
        - information you gave away (your range)
    3. EV of each card played at each step
- Implement for Uruguayan truco

### References
https://quanam.com/todo-lo-que-siempre-quisiste-saber-del-truco-uruguayo/
