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

var fruitRunes = []rune{
	apple,
	banana,
	pear,
	cherry,
	strawberry,
	eggplant,
}

type Fruit struct {
	X    int
	Y    int
	Type rune
}

func NewFruit(maxX int, maxY int) Fruit {
	randomFruit := fruitRunes[rand.Intn(len(fruitRunes))]
	return Fruit{
		X:    rand.Intn(maxX),
		Y:    rand.Intn(maxY),
		Type: randomFruit,
	}
}

func (f *Fruit) Display() rune {
	return f.Type
}

func (f *Fruit) DidHit(s *Snake) bool {
	return s.X == f.X && s.Y == f.Y
}
