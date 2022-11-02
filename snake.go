package main

const Xspeed = 1
const Yspeed = 1
const BodyRune = '*'

type BodyPart struct {
	X int
	Y int
}

type Snake struct {
	X      int
	Y      int
	Xspeed int
	Yspeed int
	Paused bool
	Length int
	Body   []BodyPart
}

func NewSnake() Snake {
	s := Snake{}
	s.X = 5
	s.Y = 10
	s.Xspeed = Xspeed
	s.Yspeed = 0
	s.Length = 4
	s.Paused = false
	for i := 0; i < s.Length; i++ {
		s.Body = append(s.Body, BodyPart{
			X: s.X - i,
			Y: s.Y,
		})
	}
	return s
}

func (s *Snake) Display() rune {
	return BodyRune
}

func (s *Snake) Update() {
	if s.Paused {
		return
	}
	// Add new body part where head was now
	b := BodyPart{
		X: s.X,
		Y: s.Y,
	}
	s.Body = append([]BodyPart{b}, s.Body...)
	if len(s.Body) == s.Length {
		// Remove last item from body
		// if snake has its full length (nothing eaten recently):
		s.Body = s.Body[:len(s.Body)-1]
	}
	// Move head at the direction of speed:
	s.X += s.Xspeed
	s.Y += s.Yspeed
}

func (s *Snake) Eat() {
	s.Length += 3
}

func (s *Snake) CheckEdges(w int, h int) bool {
	if s.X > w || s.X < 0 || s.Y > h || s.Y < 0 {
		return false
	}
	return true
}

func (s *Snake) TurnLeft() {
	if s.Xspeed == 0 {
		s.Xspeed = -1 * Xspeed
		s.Yspeed = 0
	}
}

func (s *Snake) TurnRight() {
	if s.Xspeed == 0 {
		s.Xspeed = Xspeed
		s.Yspeed = 0
	}
}

func (s *Snake) TurnUp() {
	if s.Yspeed == 0 {
		s.Xspeed = 0
		s.Yspeed = -1 * Yspeed
	}
}

func (s *Snake) TurnDown() {
	if s.Yspeed == 0 {
		s.Xspeed = 0
		s.Yspeed = Yspeed
	}
}

func (s *Snake) Pause() {
	s.Paused = !s.Paused
}
