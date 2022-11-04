package main

type Action int

const (
	TurnLeft = iota
	TurnRight
	TurnUp
	TurnDown
	Pause
	Quit
	Resize
)
