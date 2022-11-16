package main

import (
	"github.com/common-nighthawk/go-figure"
)

const font = "block"

func GetText(t string) []string {
	fig := figure.NewFigure(t, font, true)
	return fig.Slicify()
}

func GetTextAndLength(t string) (text []string, length int) {
	splash := GetText(t)
	maxSplashLen := 0
	for _, line := range splash {
		if maxSplashLen < len(line) {
			maxSplashLen = len(line)
		}
	}
	return splash, maxSplashLen
}

func PrintLogo() {
	fig := figure.NewColorFigure("SNAKE-GO", font, "yellow", true)
	fig.Blink(2000, 200, 50)
	// time.Sleep(2 * time.Second)
}
