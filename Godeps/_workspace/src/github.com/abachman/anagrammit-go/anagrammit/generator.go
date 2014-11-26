package anagrammit

import (
	"fmt"
)

type GeneratorArgs struct {
	WordsFile   string
	WordLength  int
	Shuffle     bool
	ResultLimit int
}

type Generator struct {
	WordsFile   string
	WordLength  int
	Shuffle     bool
	ResultLimit int

	baseLexicon *Lexicon
}

type generatorContext struct {
	counter int
	output  chan string
	stopped bool
}

func NewGenerator(args *GeneratorArgs) *Generator {
	generator := &Generator{
		WordsFile:   args.WordsFile,
		WordLength:  args.WordLength,
		ResultLimit: args.ResultLimit,
		Shuffle:     args.Shuffle,
	}
	generator.initialize()
	return generator
}

func (g *Generator) initialize() {
	// Base
	g.baseLexicon = NewLexicon(g.WordsFile, g.WordLength)
}

func (g *Generator) Generate(inputPhrase string, output chan string) {

	go func() {
		inputWord := NewWord(inputPhrase)

		// Calculated
		initialLexicon := g.baseLexicon.Generate(0, inputWord)

		if g.Shuffle {
			ShuffleLexicon(initialLexicon)
		}

		// hold per-process variables
		context := &generatorContext{
			counter: 0,
			output:  output,
			stopped: false,
		}

		if initialLexicon.Length == 0 {
			fmt.Println("[anagrammit] initial lexicon contains no words, no anagrams possible")
			g.Stop(context)
			return
		}

		// make sure we Stop if result limit isn't hit
		defer g.Stop(context)
		g.mainloop(initialLexicon, inputWord, NewPhrase(), context)
	}()
}

func (g *Generator) mainloop(lex *Lexicon, inpt *Word, phrase *Phrase, context *generatorContext) {
	for i := 0; i < lex.Length; i++ {

		// try the next word in the lexicon
		nextWord := lex.Words[i]

		phrase.Add(nextWord)

		// Decrement inpt's LetterCount by phrase's
		for i := 0; i < LETTER_COUNT; i++ {
			inpt.LetterCount[i] -= nextWord.LetterCount[i]
		}

		if inpt.LetterCount[LETTER_TOTAL] == 0 {
			// result found!
			context.output <- phrase.DisplayString()

			context.counter++
			if context.counter >= g.ResultLimit && g.ResultLimit != 0 {
				g.Stop(context)
				return
			}
		} else {
			nextLexicon := lex.Generate(i, inpt)

			if len(nextLexicon.Words) > 0 {
				// recurse, there's still hope
				g.mainloop(nextLexicon, inpt, phrase, context)

				if context.counter >= g.ResultLimit && g.ResultLimit != 0 {
					g.Stop(context)
					return
				}
			}
		}

		last := phrase.Pop()
		for i := 0; i < LETTER_COUNT; i++ {
			inpt.LetterCount[i] += last.LetterCount[i]
		}
	}
}

func (g *Generator) Stop(context *generatorContext) {
	if !context.stopped {
		close(context.output)
		context.stopped = true
	}
}
