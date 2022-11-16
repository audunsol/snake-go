package main

type FinishPoint struct {
	Point
	show                   bool
	EatableFruitsThreshold int
}

func NewFinishPoint(minX, minY, maxX, maxY, treshold int) FinishPoint {
	return FinishPoint{
		Point:                  NewPoint(minX, minY, maxX, maxY),
		show:                   false,
		EatableFruitsThreshold: treshold,
	}
}

func (f *FinishPoint) Show() bool {
	return f.show
}

func (f *FinishPoint) DidHit(s *Snake) bool {
	if f.show {
		return f.Point.DidHit(s)
	}
	return false
}

func (f *FinishPoint) EvaluateShow(eatableFruitsGenerated int) bool {
	if eatableFruitsGenerated > f.EatableFruitsThreshold {
		f.show = true
	} else {
		f.show = false
	}
	return f.show
}

func (f *FinishPoint) Display() rune {
	return '\U0001F3C1'
}
