package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

var defStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorYellow)

type Game struct {
	Screen     tcell.Screen
	Snake      Snake
	Fruits     []Fruit
	IsGameOver bool
}

func NewGame(screen tcell.Screen, snake Snake) Game {
	game := Game{
		Screen:     screen,
		Snake:      snake,
		IsGameOver: false,
	}
	fruits := game.GenerateFruit(10)
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

	if len(g.Fruits) == 0 {
		g.Fruits = g.GenerateFruit(3)
	}
}

func (g *Game) RenderText(startX int, startY int, text string) {
	s := g.Screen
	for pos, char := range text {
		s.SetContent(startX+pos, startY, rune(char), nil, defStyle)
	}
}

func (g *Game) CenterText(startY int, text string) {
	width, _ := g.Screen.Size()
	startX := (width / 2) - (len(text) / 2)
	g.RenderText(startX, startY, text)
}

func (g *Game) RenderCoordinates() {
	// f := g.Fruits
	sn := g.Snake
	g.RenderText(0, 0, fmt.Sprintf("Snake: (%v, %v)", sn.X, sn.Y))
	g.RenderText(0, 1, fmt.Sprintf("Snake len: %v", len(sn.Body)))
	// for i, fruit := range f {
	// 	g.RenderText(0, 1+i, fmt.Sprintf("Fruit %v: (%v, %v)", i, fruit.X, fruit.Y))
	// }
	for i, bp := range sn.Body {
		g.RenderText(0, 2+i, fmt.Sprintf("BodyPart %v: (%v, %v)", i, bp.X, bp.Y))
	}
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

func (g *Game) RenderGameOver() {
	g.CenterText(7, "Game Over")
	g.CenterText(11, fmt.Sprintf("%v points", g.Snake.Length))
	g.CenterText(15, "Hit ENTER to restart or ESC to quit")

	g.Screen.Show()
	// TODO: Make this Enter/ESC thing work...
	time.Sleep(2 * time.Second)
}

func (g *Game) Run(ch chan Action) {
	s := g.Screen
	s.SetStyle(defStyle)

	tick := time.Tick(80 * time.Millisecond)

	// Main loop is here:
	for !g.IsGameOver {
		select {
		case <-tick:

			width, height := s.Size()
			if !g.Snake.CheckEdges(width, height) || !g.Snake.CheckSelfCollision() {
				g.IsGameOver = true
			}
			g.EatFruit()
			g.Snake.Update()

			// Render:
			s.Clear()
			// g.RenderCoordinates()
			g.RenderSnake()
			g.RenderFruits()
			s.SetContent(g.Snake.X, g.Snake.Y, g.Snake.Display(), nil, defStyle)
			s.Show()
		case action := <-ch:
			switch action {
			case TurnLeft:
				g.Snake.TurnLeft()
			case TurnRight:
				g.Snake.TurnRight()
			case TurnUp:
				g.Snake.TurnUp()
			case TurnDown:
				g.Snake.TurnDown()
			case Pause:
				g.Snake.Pause()
			}
		}
	}

	g.RenderGameOver()
	g.Exit()
}

func (g *Game) Exit() {
	g.Screen.Fini()
	os.Exit(0)
}
