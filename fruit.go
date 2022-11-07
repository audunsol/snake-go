package main

import (
	"math/rand"
)

var apple = '\U0001F34E'
var banana = '\U0001F34C'
var pear = '\U0001F350'
var cherry = '\U0001F352'
var strawberry = '\U0001F353'
var eggplant = '\U0001F346'
var poo = '\U0001F4A9'
var firecracker = '\U0001F9E8'
var bomb = '\U0001F4A3'
var heart = '\U0001F9E1'

var availableFruits = []struct {
	Rune   rune
	Points int
	Lethal bool
}{
	{apple, 3, false},
	{banana, 4, false},
	{pear, 3, false},
	{cherry, 2, false},
	{strawberry, 5, false},
	{eggplant, 3, false},
	{poo, -10, false},
	{firecracker, 0, true},
	{bomb, 0, true},
	{heart, 0, false},
}

type Fruit struct {
	Point
	Type   rune
	Points int
	Lethal bool
}

func NewFruit(minX, minY, maxX, maxY int) Fruit {
	f := availableFruits[rand.Intn(len(availableFruits))]
	return Fruit{
		Point: Point{
			X: rand.Intn(maxX-minX) + minX,
			Y: rand.Intn(maxY-minY) + minY,
		},
		Type:   f.Rune,
		Points: f.Points,
		Lethal: f.Lethal,
	}
}

func (f *Fruit) Display() rune {
	return f.Type
}

func (f *Fruit) IsEatable() bool {
	return !f.Lethal && f.Points > 0
}
