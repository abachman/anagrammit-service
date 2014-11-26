package anagrammit

import (
	"strings"
)

type Word struct {
	Display     string
	LetterCount []int
}

func letterFrequency(instr string) []int {
	// last cell in letter frequency list is sum of whole list
	out := make([]int, LETTER_COUNT)
	letters := "abcdefghijklmnopqrstuvwxyz"
	for i := 0; i < len(letters); i++ {
		c := strings.Count(instr, string(letters[i]))
		if c > 0 {
			out[i] += c
			out[LETTER_TOTAL] += c
		}
	}
	return out
}

func NewWord(inWord string) *Word {
	return &Word{inWord, letterFrequency(inWord)}
}
