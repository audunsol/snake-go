package main

const BodyRune = '\U0001F7E7'
const NauseatedBodyRune = '\U0001F922'
const ExplodingBodyRune = '\U0001F4A5'

type BodyPart struct {
	X int
	Y int
}

type Snake struct {
	X         int
	Y         int
	Xspeed    int
	Yspeed    int
	Paused    bool
	Length    int
	Body      []BodyPart
	Nauseated bool
	Exploding bool
}

func NewSnake() Snake {
	s := Snake{}
	s.X = 5
	s.Y = 10
	s.Xspeed = 1
	s.Yspeed = 0
	s.Length = 4
	s.Paused = false
	s.Nauseated = false
	s.Exploding = false
	for i := 0; i < s.Length; i++ {
		s.Body = append(s.Body, BodyPart{
			X: s.X - i,
			Y: s.Y,
		})
	}
	return s
}

func (s *Snake) Display() rune {
	if s.Nauseated {
		return NauseatedBodyRune
	} else if s.Exploding {
		return ExplodingBodyRune
	} else {
		return BodyRune
	}
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
	if len(s.Body) >= s.Length {
		// Remove last item from body
		// if snake has its full length (nothing eaten recently):
		s.Body = s.Body[:len(s.Body)-1]
	}
	// Move head at the direction of speed:
	s.X += s.Xspeed
	s.Y += s.Yspeed
}

func (s *Snake) Eat(f Fruit) {
	if f.Lethal {
		s.Exploding = true
	} else if f.Points <= 0 {
		s.Nauseated = true
	} else if f.Points > 0 {
		s.Nauseated = false
	}
	s.Length += f.Points
	if s.Length < 1 {
		s.Length = 1
	}
}

func (s *Snake) CheckEdges(w int, h int) bool {
	if s.X > w || s.X < 0 || s.Y > h || s.Y < 0 {
		return false
	}
	return true
}

func (s *Snake) CheckSelfCollision() bool {
	for i := len(s.Body) - 1; i > 4; i-- {
		b := s.Body[i]
		if s.X == b.X && s.Y == b.Y {
			return false
		}
	}
	return true
}

func (s *Snake) TurnLeft() {
	if s.Xspeed == 0 {
		s.Xspeed = -1
		s.Yspeed = 0
	}
}

func (s *Snake) TurnRight() {
	if s.Xspeed == 0 {
		s.Xspeed = 1
		s.Yspeed = 0
	}
}

func (s *Snake) TurnUp() {
	if s.Yspeed == 0 {
		s.Xspeed = 0
		s.Yspeed = -1
	}
}

func (s *Snake) TurnDown() {
	if s.Yspeed == 0 {
		s.Xspeed = 0
		s.Yspeed = 1
	}
}

func (s *Snake) Pause() {
	s.Paused = !s.Paused
}
