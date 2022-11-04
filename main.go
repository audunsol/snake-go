package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
)

func main() {
	screen, err := tcell.NewScreen()

	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	snake := NewSnake()
	game := NewGame(screen, snake)

	ch := make(chan Action)

	go game.Run(ch)

	for {
		switch event := screen.PollEvent().(type) {
		case *tcell.EventResize:
			screen.Sync()
			ch <- Resize
		case *tcell.EventKey:
			if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC {
				game.Exit()
			} else if event.Key() == tcell.KeyLeft {
				ch <- TurnLeft
			} else if event.Key() == tcell.KeyRight {
				ch <- TurnRight
			} else if event.Key() == tcell.KeyUp {
				ch <- TurnUp
			} else if event.Key() == tcell.KeyDown {
				ch <- TurnDown
			} else if event.Key() == tcell.KeyCtrlSpace {
				ch <- Pause
			}
		}
	}
}
