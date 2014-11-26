package anagrammit

import (
	"strings"
)

type Phrase struct {
	Words       []*Word
	LetterCount []int
	Next        int
}

func NewPhrase() *Phrase {
	return &Phrase{make([]*Word, MAX_PHRASE), letterFrequency(""), 0}
}

func (p *Phrase) Last() *Word {
	return p.Words[p.Next-1]
}

// add Word to Phrase
func (p *Phrase) Add(other *Word) {
	p.Words[p.Next] = other
	for i := 0; i < LETTER_COUNT; i++ {
		p.LetterCount[i] += other.LetterCount[i]
	}
	p.Next++
}

// remove Word from Phrase
func (p *Phrase) Pop() *Word {
	p.Next -= 1
	removed := p.Words[p.Next]
	p.Words[p.Next] = nil
	for i := 0; i < LETTER_COUNT; i++ {
		p.LetterCount[i] -= removed.LetterCount[i]
	}
	return removed
}

func (p *Phrase) Display() []string {
	out := make([]string, p.Next)
	for i := 0; i < p.Next; i++ {
		out[i] = p.Words[i].Display
	}
	return out
}

func (p Phrase) DisplayString() string {
	return strings.Join(p.Display(), " ")
}
