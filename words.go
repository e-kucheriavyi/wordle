package main

import (
	_ "embed"
	"math/rand"
	"strings"
	"time"
)

//go:embed russian.txt
var f []byte

func GetWord(t time.Time) string {
	year, month, day := t.Date()
	s := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	r := rand.New(rand.NewSource(s.Unix()))

	var words []string = strings.Split(strings.ToLower(string(f)), "\n")
	i := r.Intn(len(words))

	return words[i]
}

func ValidateWord(guess string) bool {
	var words []string = strings.Split(strings.ToLower(string(f)), "\n")

	for i := range len(words) {
		word := words[i]

		// file word's len is 6 for some reason
		for j := range 6 {
			if j == 5 {
				return true
			}
			if guess[j] != word[j] {
				break
			}
		}
	}

	return false
}
