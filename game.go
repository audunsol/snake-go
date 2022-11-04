package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

var defStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorYellow)

var rightPanelWidth = 50
var borderSize = 1

type Game struct {
	Screen     tcell.Screen
	Width      int
	Height     int
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
	game.ResizeScreen()
	fruits := game.GenerateFruit(10)
	game.Fruits = fruits
	return game
}

func (g *Game) ResizeScreen() {
	width, height := g.Screen.Size()
	g.Width = width - rightPanelWidth - borderSize*2
	g.Height = height - borderSize*2
}

func (g *Game) GenerateFruit(n int) []Fruit {
	fruits := []Fruit{}
	for i := 0; i < n; i++ {
		fruits = append(fruits, NewFruit(borderSize, borderSize, g.Width-borderSize, g.Height-borderSize))
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
			if i < (len(snake.Body)-1) && snake.Body[i+1].Y == part.Y {
				s.SetContent(part.X+1, part.Y, g.Snake.Display(), nil, defStyle)
			}
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

	eatableFruitLeft := 0
	for _, fruit := range f {
		if fruit.IsEatable() {
			eatableFruitLeft++
		}
	}
	if eatableFruitLeft == 0 {
		g.Fruits = g.GenerateFruit(3)
	}
}

func (g *Game) RenderBorders() {
	s := g.Screen
	for i := 0; i < g.Width; i++ {
		s.SetContent(i, 0, '=', nil, defStyle)
		s.SetContent(i, g.Height, '=', nil, defStyle)
	}
	for j := 0; j < g.Height; j++ {
		s.SetContent(0, j, '#', nil, defStyle)
		s.SetContent(g.Width, j, '#', nil, defStyle)
	}
}

func (g *Game) CalculatePoints() int {
	return (g.Snake.Length - StartLength) * 100
}

func (g *Game) RenderPanel() {
	g.RenderText(g.Width+2, 2, fmt.Sprintf("Points: %v", g.CalculatePoints()))
}

func (g *Game) RenderText(startX int, startY int, text string) {
	s := g.Screen
	for pos, char := range text {
		s.SetContent(startX+pos, startY, rune(char), nil, defStyle)
	}
}

func (g *Game) CenterText(startY int, text string) {
	startX := (g.Width / 2) - (len(text) / 2)
	g.RenderText(startX, startY, text)
}

func (g *Game) RenderCoordinates() {
	// f := g.Fruits
	sn := g.Snake
	g.RenderText(g.Width+2, 5, "========================")
	g.RenderText(g.Width+2, 6, "DEBUG INFO:")
	g.RenderText(g.Width+2, 7, "-----------")
	g.RenderText(g.Width+2, 8, fmt.Sprintf("Snake: (%v, %v)", sn.X, sn.Y))
	g.RenderText(g.Width+2, 9, fmt.Sprintf("Snake len: %v", len(sn.Body)))
	// for i, fruit := range f {
	// 	g.RenderText(0, 1+i, fmt.Sprintf("Fruit %v: (%v, %v)", i, fruit.X, fruit.Y))
	// }
	for i, bp := range sn.Body {
		g.RenderText(g.Width+2, 10+i, fmt.Sprintf("BodyPart %v: (%v, %v)", i, bp.X, bp.Y))
	}
}

func (g *Game) EatFruit() {
	var i = 0
	var f Fruit
	for i = 0; i < len(g.Fruits); i++ {
		f = g.Fruits[i]
		if f.DidHit(&g.Snake) {
			break
		}
	}
	if i < len(g.Fruits) {
		g.Snake.Eat(f)
		if f.Lethal {
			g.IsGameOver = true
		}
		g.Fruits = append(g.Fruits[:i], g.Fruits[i+1:]...)
	}
}

func (g *Game) RenderGameOver() {
	g.CenterText(7, "Game Over")
	g.CenterText(11, fmt.Sprintf("%v points", g.CalculatePoints()))
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
			if !g.Snake.CheckEdges(g.Width, g.Height, borderSize) || !g.Snake.CheckSelfCollision() {
				g.IsGameOver = true
			}
			g.EatFruit()
			g.Snake.Update()

			// Render:
			s.Clear()
			g.RenderBorders()
			g.RenderPanel()
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
			case Resize:
				g.ResizeScreen()
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
