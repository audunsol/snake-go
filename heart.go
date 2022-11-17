package main

import (
	"fmt"
	"time"
)

type Heart struct {
	Point
	ShowSeconds int
	startTime   time.Time
}

func NewHeart(minX, minY, maxX, maxY, showSeconds int) Heart {
	return Heart{
		Point:       NewPoint(minX, minY, maxX, maxY),
		ShowSeconds: showSeconds,
		startTime:   time.Now(),
	}
}

func (h *Heart) DidHit(s *Snake) bool {
	if h.Show() {
		return h.Point.DidHit(s)
	}
	return false
}

func (h *Heart) Show() bool {
	return h.secondsLeft() > 0
}

func (h *Heart) secondsLeft() int {
	diff := h.ShowSeconds - int(time.Since(h.startTime).Seconds())
	if diff > 0 {
		return diff
	}
	return 0
}

func (h *Heart) Display() []rune {
	secondsLeftStr := fmt.Sprint(h.secondsLeft())
	runes := []rune{'\U0001F9E1'}
	for _, r := range secondsLeftStr {
		runes = append(runes, r)
	}
	return runes
}
