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

	go game.Run()

	for {
		switch event := screen.PollEvent().(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC {
				game.Exit()
			} else if event.Key() == tcell.KeyLeft {
				game.Snake.TurnLeft()
			} else if event.Key() == tcell.KeyRight {
				game.Snake.TurnRight()
			} else if event.Key() == tcell.KeyUp {
				game.Snake.TurnUp()
			} else if event.Key() == tcell.KeyDown {
				game.Snake.TurnDown()
			}
		}
	}
}
