package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

var defStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorYellow)

var rightPanelWidth = 50
var borderSize = 1

const initialNumberOfFruits = 10
const finishPointThreshold = 30

type Game struct {
	Screen                tcell.Screen
	Width                 int
	Height                int
	Snake                 Snake
	Fruits                []Fruit
	Wall                  Wall
	IsGameOver            bool
	EatableFruitsPerLevel int
	StartTime             time.Time
	Lives                 int
	Hearts                []Heart
	PreviousPoints        int
	HighScoreList         HighScoreList
	FinishPoint           FinishPoint
	Level                 int
}

func NewGame(screen tcell.Screen) Game {
	game := Game{
		Screen:         screen,
		Snake:          NewSnake(),
		IsGameOver:     false,
		StartTime:      time.Now(),
		Lives:          3,
		PreviousPoints: 0,
		Level:          1,
	}
	game.ResizeScreen()
	game.Fruits = game.GenerateFruit(initialNumberOfFruits)
	game.EatableFruitsPerLevel = eatableFruitLeft(game.Fruits)
	game.HighScoreList = ReadHighScoresFromFile()
	game.FinishPoint = NewFinishPoint(borderSize, borderSize, game.Width-borderSize, game.Height-borderSize, finishPointThreshold)
	game.Wall = game.GenerateWall()
	return game
}

func (g *Game) ClearAndRerenderFrame() {
	s := g.Screen
	s.Clear()
	g.RenderBorders()
	g.RenderHighScore(7)
	g.RenderWall()
	s.Show()
	s.Sync()
}

func (g *Game) GenerateWall() Wall {
	return WallPerLevel(borderSize, borderSize, g.Width-borderSize, g.Height-borderSize, g.Level)
}

func (g *Game) RenderSplashText(startLine int, t string) int {
	splash, maxSplashLen := GetTextAndLength(t)
	startX := (g.Width / 2) - (maxSplashLen / 2)
	for i, line := range splash {
		g.RenderText(startX, startLine+i, line)
	}
	g.Screen.Show()
	return len(splash)
}

func (g *Game) NextLevel() {
	g.ClearAndRerenderFrame()
	g.EatableFruitsPerLevel = 0
	g.Level++
	g.RenderSplashText(10, fmt.Sprintf("LEVEL: %v", g.Level))
	g.Screen.Show()
	time.Sleep(2 * time.Second)
	g.ClearAndRerenderFrame()
	g.PreviousPoints = g.CalculatePoints()
	// Add points just for completing the level:
	g.PreviousPoints += 1000
	g.Snake = NewSnake()
	g.Fruits = g.GenerateFruit(initialNumberOfFruits)
	g.FinishPoint = NewFinishPoint(borderSize, borderSize, g.Width-borderSize, g.Height-borderSize, finishPointThreshold)
	g.Wall = g.GenerateWall()
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

func (g *Game) GenerateHeart() {
	// Random, around every 500 tick, generate a new heart
	// that is visible for 10 seconds:
	if rand.Intn(500) == 0 {
		g.Hearts = append(g.Hearts, NewHeart(borderSize, borderSize, g.Width-borderSize, g.Height-borderSize, 10))
	}
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

func (g *Game) renderOrClearFruits(clear bool) {
	f := g.Fruits
	s := g.Screen
	for i := 0; i < len(f); i++ {
		fruit := f[i]
		r := fruit.Display()
		if clear {
			r = ' '
		}
		s.SetContent(fruit.X, fruit.Y, r, nil, defStyle)
	}
}

func (g *Game) ClearFruits() {
	g.renderOrClearFruits(true)
}

func eatableFruitLeft(f []Fruit) int {
	eatableFruitLeft := 0
	for _, fruit := range f {
		if fruit.IsEatable() {
			eatableFruitLeft++
		}
	}
	return eatableFruitLeft
}

func (g *Game) RenderFruits() {
	f := g.Fruits
	if eatableFruitLeft(f) == 0 {
		// Clear before regenerate:
		g.ClearAndRerenderFrame()
		g.Fruits = g.GenerateFruit(3)
		g.EatableFruitsPerLevel += eatableFruitLeft(g.Fruits)
	}
	g.renderOrClearFruits(false)
	g.FinishPoint.EvaluateShow(g.EatableFruitsPerLevel)
}

func (g *Game) renderOrClearHearts(clear bool) {
	h := g.Hearts
	s := g.Screen
	for i := 0; i < len(h); i++ {
		heart := h[i]
		runes := []rune{' ', ' ', ' '}
		if !clear && heart.Show() {
			runes = heart.Display()
		}
		for i, r := range runes {
			yOffset := 0
			if i > 0 {
				yOffset = -1
			}
			s.SetContent(heart.X+i, heart.Y+yOffset, r, nil, defStyle)
		}
	}
}

func (g *Game) ClearHearts() {
	g.renderOrClearHearts(true)
}

func (g *Game) RenderHearts() {
	g.renderOrClearHearts(false)
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

func (g *Game) RenderWall() {
	s := g.Screen
	w := g.Wall
	for _, p := range w.Points {
		s.SetContent(p.X, p.Y, w.Display(), nil, defStyle)
	}
}

func (g *Game) RenderFinishPoint() {
	s := g.Screen
	if g.FinishPoint.Show() {
		s.SetContent(g.FinishPoint.X, g.FinishPoint.Y, g.FinishPoint.Display(), nil, defStyle)
	}
}

// func (g *Game) RenderWormholes() {
// 	hole := '\U0001F573'
// }

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
	g.RenderText(x, 5, fmt.Sprintf("Level: %v", g.Level))
	duration := time.Since(g.StartTime)
	g.RenderText(x, 6, fmt.Sprintf("Duration: %v", fmtDuration(duration)))
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
		g.Fruits = append(g.Fruits[:i], g.Fruits[i+1:]...)
	}
}

func (g *Game) CheckHearts() {
	var i int
	for i = 0; i < len(g.Hearts); i++ {
		h := g.Hearts[i]
		if h.DidHit(&g.Snake) {
			g.Lives++
			break
		}
	}
	// Remove heart from list:
	if i < len(g.Hearts) {
		g.Hearts = append(g.Hearts[:i], g.Hearts[i+1:]...)
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
			} else if answer == Delete && len(input) > 0 {
				input = input[:len(input)-1]
				s.SetContent(w+len(input)+1, h, ' ', nil, defStyle)
				g.Screen.Show()
			}
		case inputRune := <-inputCh:
			input = append(input, inputRune)
			s.SetContent(w+len(input), h, inputRune, nil, defStyle)
			g.Screen.Show()
		}
	}
}

func (g *Game) RenderGameOver(ch chan Action, input chan rune) {
	lineCounter := 7
	lineCounter += g.RenderSplashText(lineCounter, "GAME OVER")
	lineCounter += 2
	g.CenterText(lineCounter, fmt.Sprintf("%v points", g.CalculatePoints()))
	points := g.CalculatePoints()
	highscorelist := ReadHighScoresFromFile()
	rank := highscorelist.IsNewHighScore(points)
	lineCounter++
	if rank > 0 {
		g.CenterText(lineCounter, fmt.Sprintf("New high score, rank %v!", rank))
		lineCounter += 2
		startX := g.CenterText(lineCounter, "Enter your name:")
		g.Screen.Show()
		name := g.ReadInput(startX+1, lineCounter, ch, input)
		newList := highscorelist.Add(name, points)
		newList.WriteToFile()
		g.RenderHighScore(7)
		g.Screen.Show()
	}
	lineCounter += 3
	g.CenterText(lineCounter, "Hit ENTER to restart or ESC to quit")

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
	s.SetStyle(defStyle)
	g.ClearAndRerenderFrame()

	splashHeight := g.RenderSplashText(10, "SNAKE-GO")
	g.CenterText(10+splashHeight+2, "Get ready...")
	s.Show()
	time.Sleep(2 * time.Second)
	tick := time.Tick(80 * time.Millisecond)

	// Main loop is here:
	for !g.IsGameOver {
		select {
		case <-tick:
			if !g.Snake.CheckEdges(g.Width, g.Height, borderSize) || !g.Snake.CheckSelfCollision() || g.Wall.DidHit(&g.Snake) {
				g.RemoveLife()
				g.ClearAndRerenderFrame()
			}
			g.GenerateHeart()
			g.ClearSnake()
			g.ClearFruits()
			g.ClearHearts()
			g.EatFruit()
			g.CheckHearts()
			if g.FinishPoint.DidHit(&g.Snake) {
				g.NextLevel()
			}
			g.Snake.Update()

			// Render:
			g.RenderPanel()
			g.RenderWall()
			// g.RenderCoordinates()
			g.RenderSnake()
			g.RenderFruits()
			g.RenderHearts()
			g.RenderFinishPoint()
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
				g.ClearAndRerenderFrame()
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
