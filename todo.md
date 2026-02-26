# Solver for the game Truco, as played in Argentina and Uruguay

### Features
1. Given your hand
    - Calc your envido / flor points (easy)
    - Calc 'truco' strength
2. Given your envido
    - Whats the chance your envido/flor is best (medium)
    - Whats the range of hands you could have (hard)
    - Bonus: what card reveals less, or contradictory, info (eg. 33 envido with 2m + 3, showing any 7 will bluff)
3. Guess range of other players given limited info
    - Flor (score)
    - Envido (score or partial score)
    - Individual cards

settings:
- uy
- num players (2, 4)
- hints, explain

test
11e 7b 11b
1o 5o (less than 11)

example:
1e 7c 7b
kEnvido = 33
hasStrategy is key to understand the results

### TODO UI:
- select past turns:
    - without breaking history
    - enable undo logic (change card or action)
    - once you click on any player action, the matrix should reflect the changes of stats relevant to this user
    """
    save every state-change action (and params) to a map
    send this past actions to the frontend separately, with the turn id
    on-request past action: recreate the match with the past actions
    frontend: on click any action in a tracker box, delete all next (if any) and concat incoming
    """
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
    - toast
- make 'choose card' easier

### update stats
- truco: only median strength
    - click (not hover), optionally provide third card: more info
    - decile of remaining hands: 'chance its the best'
- envido: only in first turns
    - allow me to dig deeper into envido: provide info on mCards
    - chance that your envido is best in table

### Truco analytics
- 1v1 problem: TrucoStrength is 1v1; 2v2 should have a different beatness model
    - the problem is the search space is too big. 
    - mCards, kCards, pCards (partner cards): I may have a value, range or unknown a pCard
    - we can model the match as 6_cards vs 6_cards, each turn of 4 cards, keeping first card of each turn invariant (mCard)
        mkpk|mkpk|mkpk
- Ideal features: How to play: given your hand (and known cards)
    - what's the best card to play
    - what's the chance you win (ask for truco)

### Uruguay
- Make sure to reuse generic functions

### Next steps
- host in a free tier server
- hover on different elements also shows hints in the bottom, explaining UI
    - tutorial?
- Track progress of a hand
    - scrollable action tracker
    - The matrix will constantly be updated every time some player makes an action
    - make a nicer card chooser (and reduce the options based on info available)

Envido scores:
turno | no q | quiero | canto
------|------|--------|------
1     | 1    | 2      | env
1     | 1    | 3      | real
1     | 1    | 255    | falta
2     | 2    | 4      | env + env
2     | 2    | 5      | env + real
2     | 2    | 255    | env + falta
3     | 4    | 7      | env + env + real
3     | 5    | 255    | env + env + falta
4     | 7    | 255    | env + env + real + falta

### References
https://quanam.com/todo-lo-que-siempre-quisiste-saber-del-truco-uruguayo/

https://railway.com/pricing
https://render.com/pricing
https://www.alwaysdata.com/en/offers/
