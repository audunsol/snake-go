package main

import "math/rand"

type Point struct {
	X int
	Y int
}

func NewPoint(minX, minY, maxX, maxY int) Point {
	return Point{
		X: rand.Intn(maxX-minX) + minX,
		Y: rand.Intn(maxY-minY) + minY,
	}
}

func (p *Point) didFastForwardOver(s *Snake) bool {
	if s.Y == p.Y {
		return (s.Xspeed == 4 && ((s.X-2) == p.X || (s.X-3) == p.X)) || (s.Xspeed == -4 && ((s.X+2) == p.X || s.X+3 == p.X))
	}
	if s.X == p.X {
		return (s.Yspeed == 2 && (s.Y-1) == p.Y) || (s.Yspeed == -2 && (s.Y+1) == p.Y)
	}
	return false
}

func (p *Point) DidHit(s *Snake) bool {
	if (s.X == p.X || s.X+1 == p.X) && s.Y == p.Y {
		return true
	}

	return p.didFastForwardOver(s)
}
