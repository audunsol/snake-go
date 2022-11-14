package main

import (
	"testing"
)

type checkEdgesTestCase struct {
	X          int
	Y          int
	Width      int
	Height     int
	BorderSize int
	Want       bool
}

func TestCheckEdgesSimple(t *testing.T) {
	cases := []checkEdgesTestCase{
		{1, 1, 5, 5, 0, true},
		{6, 1, 5, 5, 0, false},
		{100, 1, 5, 5, 1, false},
		{5, 1, 5, 5, 1, false},
		{1, 5, 5, 5, 1, false},
		{1, 1, 5, 5, 2, false},
	}

	for _, c := range cases {
		snake := Snake{
			Point: Point{
				X: c.X,
				Y: c.Y,
			},
		}

		got := snake.CheckEdges(c.Width, c.Height, c.BorderSize)

		if c.Want != got {
			t.Fail()
		}
	}
}

func TestUpdateDoNothingIfPaused(t *testing.T) {
	snake := Snake{
		Point: Point{
			X: 5,
			Y: 1,
		},
		Xspeed: 2,
		Yspeed: 0,
		Paused: true,
		Body:   []Point{{4, 1}, {3, 1}},
		Length: 3,
	}

	snake.Update()

	if snake.X != 5 || snake.Y != 1 || len(snake.Body) != 2 {
		t.Errorf("Update while paused still changed something!")
	}
}
