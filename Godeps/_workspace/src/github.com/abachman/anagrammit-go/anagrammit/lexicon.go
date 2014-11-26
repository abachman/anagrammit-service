package anagrammit

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 8196)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return count, err
		}

		count += bytes.Count(buf[:c], lineSep)

		if err == io.EOF {
			break
		}
	}

	return count, nil
}

type Lexicon struct {
	Words  []*Word
	Length int
}

func (l *Lexicon) Append(w *Word) {
	l.Words[l.Length] = w
	l.Length++
}

func NewLexicon(wordFile string, wordLen int) *Lexicon {
	// read
	file, err := os.Open(wordFile)
	if err != nil {
		fmt.Println("[anagrammit] unable to open file :(")
		log.Fatal(err)
	}
	defer file.Close()

	count, _ := lineCounter(file)
	file.Seek(0, 0)

	lexicon := &Lexicon{Words: make([]*Word, count)}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		if len(s) >= wordLen {
			lexicon.Append(NewWord(s))
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "[anagrammit] error reading standard input:", err)
	}

	return lexicon
}

// from existing, given a word, generate a new lexicon
func (existing *Lexicon) Generate(offset int, inputWord *Word) *Lexicon {
	lexicon := &Lexicon{Words: make([]*Word, existing.Length)}
	bad := false

	// check every word in existing lexicon
	for w := offset; w < existing.Length; w++ {
		word := existing.Words[w]

		bad = false
		// check every letter in lexicon word
		for l := 0; l < LETTER_COUNT; l++ {
			if word.LetterCount[l] > inputWord.LetterCount[l] {
				bad = true
				break
			}
		}

		if !bad {
			lexicon.Append(word)
		}
	}

	return lexicon
}
