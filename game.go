package main

import (
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

var defStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorYellow)

type Game struct {
	Screen tcell.Screen
	Snake  Snake
	Fruits []Fruit
}

func (g *Game) RenderSnake() {
	s := g.Screen
	s.SetContent(g.Snake.X, g.Snake.Y, g.Snake.Display(), nil, defStyle)
	snake := g.Snake
	for i := 0; i < snake.Length; i++ {
		part := snake.Body[i]
		s.SetContent(part.X, part.Y, g.Snake.Display(), nil, defStyle)
	}
}

func (g *Game) RenderFruits() {
	s := g.Screen
	f := g.Fruits
	for i := 0; i < len(f); i++ {
		fruit := f[i]
		s.SetContent(fruit.X, fruit.Y, fruit.Display(), nil, defStyle)
	}
}

func (g *Game) Run() {
	s := g.Screen
	s.SetStyle(defStyle)

	// Main loop is here:
	for {
		s.Clear()

		width, height := s.Size()
		if !g.Snake.CheckEdges(width, height) {
			g.Exit()
		}
		// TODO: check if hit fruit
		g.Snake.Update()
		g.RenderSnake()
		g.RenderFruits()
		s.SetContent(g.Snake.X, g.Snake.Y, g.Snake.Display(), nil, defStyle)

		time.Sleep(40 * time.Millisecond)
		s.Show()
	}
}

func (g *Game) Exit() {
	g.Screen.Fini()
	log.Println("Game over!")
	os.Exit(0)
}