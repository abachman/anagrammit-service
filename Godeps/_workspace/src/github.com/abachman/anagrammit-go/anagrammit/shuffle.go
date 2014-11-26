package anagrammit

import (
	"math/rand"
	"time"
)

// shuffle
// for i from n − 1 downto 1 do
//      j ← random integer with 0 ≤ j ≤ i
//      exchange a[j] and a[i]
func ShuffleLexicon(lex *Lexicon) {
	rand.Seed(time.Now().UnixNano())

	var j int
	for i := lex.Length - 1; i > 0; i-- {
		j = rand.Intn(i)
		w := lex.Words[j]
		lex.Words[j] = lex.Words[i]
		lex.Words[i] = w
	}
}
