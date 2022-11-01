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

func NewGame(screen tcell.Screen, snake Snake) Game {
	game := Game{
		Screen: screen,
		Snake:  snake,
	}
	fruits := game.GenerateFruit(3)
	game.Fruits = fruits
	return game
}

func (g *Game) GenerateFruit(n int) []Fruit {
	width, height := g.Screen.Size()
	fruits := []Fruit{}
	for i := 0; i < n; i++ {
		fruits = append(fruits, NewFruit(width, height))
	}
	return fruits
}

func (g *Game) RenderSnake() {
	s := g.Screen
	s.SetContent(g.Snake.X, g.Snake.Y, g.Snake.Display(), nil, defStyle)
	snake := g.Snake
	for i := 0; i < snake.Length; i++ {
		if i < len(snake.Body) {
			part := snake.Body[i]
			s.SetContent(part.X, part.Y, g.Snake.Display(), nil, defStyle)
		}
	}
}

func (g *Game) RenderFruits() {
	s := g.Screen
	f := g.Fruits
	for i := 0; i < len(f); i++ {
		fruit := f[i]
		s.SetContent(fruit.X, fruit.Y, fruit.Display(), nil, defStyle)
	}
func (g *Game) EatFruit() {
	var i = 0
	for i = 0; i < len(g.Fruits); i++ {
		f := g.Fruits[i]
		if f.DidHit(&g.Snake) {
			break
		}
	}
	if i < len(g.Fruits) {
		g.Fruits = append(g.Fruits[:i], g.Fruits[i+1:]...)
		g.Snake.Length += 3
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
		g.EatFruit()
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
