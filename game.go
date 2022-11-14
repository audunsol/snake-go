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

const initialNumberOfFruits = 10

type Game struct {
	Screen                tcell.Screen
	Width                 int
	Height                int
	Snake                 Snake
	Fruits                []Fruit
	IsGameOver            bool
	EatableFruitsPerLevel int
	StartTime             time.Time
	Lives                 int
	PreviousPoints        int
	HighScoreList         HighScoreList
	FinishPoint           FinishPoint
}

func NewGame(screen tcell.Screen) Game {
	game := Game{
		Screen:         screen,
		Snake:          NewSnake(),
		IsGameOver:     false,
		StartTime:      time.Now(),
		Lives:          3,
		PreviousPoints: 0,
	}
	game.ResizeScreen()
	game.Fruits = game.GenerateFruit(initialNumberOfFruits)
	game.HighScoreList = ReadHighScoresFromFile()
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

func (g *Game) renderSnakeWithRune(r rune) {
	s := g.Screen
	s.SetContent(g.Snake.X, g.Snake.Y, r, nil, defStyle)
	snake := g.Snake
	for i := 0; i < snake.Length; i++ {
		if i < len(snake.Body) {
			part := snake.Body[i]
			s.SetContent(part.X, part.Y, r, nil, defStyle)
			if i < (len(snake.Body)-1) && snake.Body[i+1].Y == part.Y {
				s.SetContent(part.X+1, part.Y, r, nil, defStyle)
			}
		}
	}
}

func (g *Game) RenderSnake() {
	g.renderSnakeWithRune(g.Snake.Display())
}

func (g *Game) ClearSnake() {
	g.renderSnakeWithRune(' ')
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

func (g *Game) RenderFinishPoint() {
	s := g.Screen
	if g.FinishPoint.Show {
		s.SetContent(g.FinishPoint.X, g.FinishPoint.Y, g.FinishPoint.Display(), nil, defStyle)
	}
}
func (g *Game) CalculatePoints() int {
	return g.PreviousPoints + (g.Snake.Length-StartLength)*100
}

func (g *Game) RemoveLife() {
	g.Lives--
	if g.Lives <= 0 {
		g.IsGameOver = true
	} else {
		g.PreviousPoints = g.CalculatePoints()
		g.ClearSnake()
		g.Snake = NewSnake()
		g.Fruits = g.GenerateFruit(initialNumberOfFruits)
	}
}

func fmtDuration(d time.Duration) string {
	sa := d.Round(time.Second)
	m := sa / time.Minute
	sa -= m * time.Minute
	s := sa / time.Second
	return fmt.Sprintf("%02d:%02d", m, s)
}

func (g *Game) RenderPanel() {
	x := g.Width + 2
	hearts := make([]rune, 5)
	if g.Lives < 6 {
		for i := 0; i < 5; i++ {
			if i < g.Lives {
				hearts[i] = '\U0001F9E1'
			} else {
				hearts[i] = ' '
			}
		}
		g.RenderText(x, 3, fmt.Sprintf("Lives: %v", string(hearts)))
	} else {
		g.RenderText(x, 3, fmt.Sprintf("Lives: %v x %s     ", g.Lives, string('\U0001F9E1')))
	}
	g.RenderText(x, 4, fmt.Sprintf("Points: %v", g.CalculatePoints()))
	duration := time.Since(g.StartTime)
	g.RenderText(x, 5, fmt.Sprintf("Duration: %v", fmtDuration(duration)))
}

func (g *Game) RenderText(startX int, startY int, text string) {
	s := g.Screen
	for pos, char := range text {
		s.SetContent(startX+pos, startY, rune(char), nil, defStyle)
	}
}

func (g *Game) CenterText(startY int, text string) int {
	startX := (g.Width / 2) - (len(text) / 2)
	g.RenderText(startX, startY, text)
	return startX + len(text)
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

func (g *Game) RenderHighScore(startY int) {
	g.RenderText(g.Width+2, startY, "========================")
	g.RenderText(g.Width+2, startY+1, "HIGH SCORE:")
	g.RenderText(g.Width+2, startY+2, "-----------")
	for i, hs := range g.HighScoreList {
		g.RenderText(g.Width+2, startY+4+i, fmt.Sprintf("%v: %v", i+1, hs.Name))
		g.RenderText(g.Width+rightPanelWidth-5, startY+4+i, fmt.Sprintf("%v", hs.Points))
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
			g.RemoveLife()
		}
		if f.Type == heart {
			g.Lives++
		}
		g.Fruits = append(g.Fruits[:i], g.Fruits[i+1:]...)
	}
}

func (g *Game) ReadInput(w int, h int, ch chan Action, inputCh chan rune) string {
	input := []rune{}
	s := g.Screen
	for {
		select {
		case answer := <-ch:
			if answer == Yes {
				return string(input)
			} else if answer == Quit {
				g.Exit()
			}
		case inputRune := <-inputCh:
			input = append(input, inputRune)
			s.SetContent(w+len(input), h, inputRune, nil, defStyle)
			g.Screen.Show()
		}
	}
}

func (g *Game) RenderGameOver(ch chan Action, input chan rune) {
	g.CenterText(7, "Game Over")
	g.CenterText(11, fmt.Sprintf("%v points", g.CalculatePoints()))
	points := g.CalculatePoints()
	highscorelist := ReadHighScoresFromFile()
	rank := highscorelist.IsNewHighScore(points)
	if rank > 0 {
		g.CenterText(15, fmt.Sprintf("New high score, rank %v!", rank))
		startX := g.CenterText(17, "Enter your name:")
		g.Screen.Show()
		name := g.ReadInput(startX+1, 17, ch, input)
		newList := highscorelist.Add(name, points)
		newList.WriteToFile()
		g.RenderHighScore(7)
		g.Screen.Show()
	}
	g.CenterText(20, "Hit ENTER to restart or ESC to quit")

	g.Screen.Show()

	for {
		select {
		case answer := <-ch:
			if answer == Yes {
				newGame := NewGame(g.Screen)
				g = &newGame
				g.Run(ch, input)
			} else if answer == Quit {
				g.Exit()
			}
		case <-input:
			// Ignore other runes when awaiting if start new game or not...
		}
	}
}

func (g *Game) Run(ch chan Action, input chan rune) {
	s := g.Screen
	s.Clear()
	s.SetStyle(defStyle)
	g.RenderBorders()
	g.RenderHighScore(7)

	tick := time.Tick(80 * time.Millisecond)

	// Main loop is here:
	for !g.IsGameOver {
		select {
		case <-tick:
			if !g.Snake.CheckEdges(g.Width, g.Height, borderSize) || !g.Snake.CheckSelfCollision() {
				g.RemoveLife()
				s.Clear()
				g.RenderBorders()
				g.RenderHighScore(7)
			}
			g.EatFruit()
			g.ClearSnake()
			g.Snake.Update()

			// Render:
			g.RenderPanel()
			// g.RenderCoordinates()
			g.RenderSnake()
			g.RenderFruits()
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
				s.Clear()
				g.RenderBorders()
				g.RenderHighScore(7)
			case Quit:
				g.Exit()
			}
		case <-input:
			// Ignore typing while playing in main loop for now...
		}
	}

	g.RenderGameOver(ch, input)
}

func (g *Game) Exit() {
	g.Screen.Fini()
	os.Exit(0)
}
