package main

import "sort"

const WallRune = '#'

type Wall struct {
	Points []Point
}

func WallPerLevel(minX, minY, maxX, maxY, level int) Wall {
	switch level {
	case 1:
		return emptyWall()
	case 2:
		return horizMidWall(minX, minY, maxX, maxY, 4)
	case 3:
		return vertMidWall(minX, minY, maxX, maxY, 4)
	case 4:
		hori := horizMidWall(minX, minY, maxX, maxY, 4)
		verti := vertMidWall(minX, minY, maxX, maxY, 4)
		return merge(hori, verti)
	case 5:
		return horizMidWall(minX, minY, maxX, maxY, 2)
	case 6:
		return vertMidWall(minX, minY, maxX, maxY, 2)
	case 7:
		hori := horizMidWall(minX, minY, maxX, maxY, 2)
		verti := vertMidWall(minX, minY, maxX, maxY, 2)
		return merge(hori, verti)
	default:
		return emptyWall()
	}
}

func (w *Wall) Display() rune {
	return WallRune
}

func emptyWall() Wall {
	return Wall{
		Points: []Point{},
	}
}

func merge(a, b Wall) Wall {
	points := append(a.Points, b.Points...)
	sort.Slice(points, func(i, j int) bool {
		return points[i].X > points[j].X
	})
	var uniquePoints []Point
	for i := 0; i < len(points)-1; i++ {
		p1 := points[i]
		p2 := points[i+1]
		if i == 0 {
			uniquePoints = []Point{p1}
		}
		if p1.X != p2.X || p1.Y != p2.Y {
			uniquePoints = append(uniquePoints, p2)
		}
	}
	return Wall{Points: points}
}

func horizMidWall(minX, minY, maxX, maxY, fracOfScreen int) Wall {
	length := (maxX - minX) / fracOfScreen
	xPos := (maxX-minX)/2 - (length / 2)
	yPos := (maxY - minY) / 2

	points := []Point{}
	for i := 0; i < length; i++ {
		points = append(points, Point{X: xPos + i, Y: yPos})
	}
	return Wall{
		Points: points,
	}
}

func vertMidWall(minX, minY, maxX, maxY, fracOfScreen int) Wall {
	length := (maxY - minY) / fracOfScreen
	yPos := (maxY-minY)/2 - (length / 2)
	xPos := (maxX - minX) / 2

	points := []Point{}
	for i := 0; i < length; i++ {
		points = append(points, Point{X: xPos, Y: yPos + i})
	}
	return Wall{
		Points: points,
	}
}

func (w *Wall) DidHit(s *Snake) bool {
	for _, p := range w.Points {
		if p.DidHit(s) {
			return true
		}
	}
	return false
}
