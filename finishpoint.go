package main

type FinishPoint struct {
	Point
	Show bool
}

func NewFinishPoint(minX, minY, maxX, maxY int) FinishPoint {
	return FinishPoint{
		Point: NewPoint(minX, minY, maxX, maxY),
		Show:  false,
	}
}

func (f *FinishPoint) Display() rune {
	return '\U0001F3C1'
}
