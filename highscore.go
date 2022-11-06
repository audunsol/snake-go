package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"
	"time"
)

const fileName = "./highscorelist.json"

type HighScore struct {
	Name   string
	Points int
	Time   time.Time
}

func NewHighScore(name string, points int) HighScore {
	return HighScore{
		Name:   name,
		Points: points,
		Time:   time.Now(),
	}
}

type HighScoreList []HighScore

func (list *HighScoreList) Sort() {
	l := *list
	sort.Slice(l, func(i, j int) bool {
		return l[i].Points > l[j].Points
	})
}

func ReadHighScoresFromFile() HighScoreList {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		return []HighScore{}
	}
	var highscorelist HighScoreList
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &highscorelist)
	highscorelist.Sort()
	defer jsonFile.Close()
	return highscorelist
}

func (list *HighScoreList) IsNewHighScore(points int) int {
	if len(*list) == 0 {
		return 1
	}
	for i, v := range *list {
		if v.Points < points {
			return i + 1
		}
	}
	return 0
}

func (list *HighScoreList) Add(name string, points int) HighScoreList {
	hs := NewHighScore(name, points)
	newList := append(*list, hs)
	newList.Sort()
	return newList
}

func (list *HighScoreList) WriteToFile() {
	file, _ := json.MarshalIndent(list, "", " ")

	_ = ioutil.WriteFile(fileName, file, 0644)
}
