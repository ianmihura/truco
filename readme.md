# Solver for the game Truco, as played in Argentina and Uruguay

![image](./screenshot.jpeg)

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
- settings on top
    - num players
    - Separate between arg & uru
- nicer styles
- choose mCards and kCards out of a 4x10 matrix
    then also choose an envido
- on click
    ### finish identifying your hand:
    - choose the third card
    - choose envido
    ### stats:
    - % chance you can get this hand
    - strength (truco and envido)
    - % hands you can win (similar to strength)
- main matrix: relevant info
    - que cartas puedo tener - my range
    - que chance hay que mis cartas sean las mejores de la mesa - strength
    - que cartas puede tener el otro
- translate error messages from fsm

### add stats
- chance that your hand is best in table, given known info
- chance that your envido is best in table, given known info

### TODO backend 
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

### Uruguay
- Make sure to reuse generic functions

### References
https://quanam.com/todo-lo-que-siempre-quisiste-saber-del-truco-uruguayo/

----

### Next steps
- return updated stats when FSM action is taken
- Track progress of a hand
    - scrollable action tracker
    - The matrix will constantly be updated every time some player makes an action
    - make a nicer card chooser (and reduce the options based on info available)
- remake matrix to be a triangle - we dont really care for separation between envido and non envido
    - pasa que el envido es solo relevante en la primera mano
        si ya anuncion que envido tiene, ya estamos mucho mas claros
        pero si no anuncio, lo mas relevante va a ser que cartas tiene, y separarlas en 2 (en especial por la cuenta) seria medio dificil de ver en lo obvio
    - Finish selecting a hand in with a third click (maybe 2 clicks: number, suit)
- Reduce the amount of js in the `matrix.html` file.
    Pass some of the logic to go, and allow generation of html in the backend (as htmx is meant to do). It will help that later we need to be able to change the `hand_stats.csv` file dynamically (remove some impossible hands, given known information). 
    This means that the function `CreatePairStatsCSV` will be called many times (every time we need to render the matrix)
- Define what metrics are encoded with color in the matrix
    We may need to show many stats, options:
    - click and show a side panel of relevant info (current, all info, as in GTOWizard)
    - color-code relevant stats for a first-glance intuition (use % bars in-hand to 2-3 dimensions)
        color: can vary, but if its too complex, it breaks vertical and horizontal
        vertical: must normalize
        horizontal: must normalize

        relevant metrics:
        % strength (vary between combined and truco alone)
        count (% chance that this is the hand)

### FSM
[Example in golang](https://refactoring.guru/design-patterns/state/go/example)

Mechanics:
- The frontend must:
    - Save the state
    - Send it to the backend on each request
    - UI interactions are easy: no state transition or interpretation
- The backend must:
    - Encode/decode state to and from the frontend
- Generalizer for uy

Envido scores:
```
tn| nq  quiero
1 | 1 - 2
1 | 1 - 3
1 | 1 - 255

2 | 2 - 4
2 | 2 - 5
2 | 2 - 255

3 | 4 - 7
3 | 5 - 255

4 | 7 - 255
```
