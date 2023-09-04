# Snake Go

This is a very simple project, just made for me to:

- learn golang
- have fun
- show some programming stuff to my kids

Please don't take it too seriously.

## How to run

```bash
git clone <this repo>
cd snake-go
go run .
```

## Remaining work

- [x] Don't mess up terminal colors etc when terminating
- [x] Get fruits to render properly
- [x] Increase length of snake when eating fruit
- [x] Remove fruit when eaten
- [x] Fix bug where it is not possible to eat some bugs
- [x] Add Game Over screen
- [x] Add start over/exit question input
- [x] Increase speed when holding arrows (Fast Forward)
- [x] Set col/row size to make speed consistent in x/y direction
- [x] Add border/frame and point counter
- [x] Add menu to start new game, exit, set player name
      _Kinda added this, but just upon game over, and name only if on highscore_
- ~~Add sqlite DB to hold high-scores and user profiles~~
- [x] Store highscores in local JSON file
- [x] Accept name input when getting a new high score
- [x] Add multiple lives before game over
- [x] Add heart emojis that increase number of lives
- [x] Make hearts less probable than other fruits, and for a limited time
- [ ] Add settings
- [ ] Add snake head (other unicode/emoji symbol)
- [x] Change snake bodyparts (other unicode/emoji symbol that better aligns with fruit)
- [x] Add pause-functionality
- [x] Game over when hitting itself
- [x] Different points for different fruits
- [x] Negative points for some fruits (poo)?
- [x] Animate vomiting or nauseated face when eating poo
- [x] Animate explosion when eating bombs
- [ ] Bonus points for fruits (diamonds) showing up just for a short period (like hearts)
- [ ] Implement wormhole
- [x] Add bombs (fruits that one has to avoid)
- [x] Add new level when eaten a defined number of fruits
- [x] Walls with increasing challenges as leveling up
