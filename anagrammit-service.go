// A basic anagram generating web service

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/abachman/anagrammit-go/anagrammit"
)

var (
	generator  *anagrammit.Generator
	validChars = regexp.MustCompile("[^a-z]")
)

func filtered(word string) string {
	word = strings.ToLower(word)
	word = validChars.ReplaceAllString(word, "")

	// only use the first 32 characters
	if len(word) > 32 {
		return word[:32]
	} else {
		return word
	}
}

func anagramHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	phrase := r.FormValue("phrase")

	phrase = filtered(phrase)
	log.Println("responding to input phrase", phrase)

	if len(phrase) > 0 {
		output := make(chan string)
		generator.Generate(phrase, output)
		for msg := range output {
			fmt.Fprintln(w, msg)
			// flush phrases as they are generated
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		}
	} else {
		fmt.Fprintln(w, "")
	}
}

func main() {
	wordLen := flag.Int("wordlength", 3, "minimum word length")
	limit := flag.Int("limit", 100, "result limit")
	shuffle := flag.Bool("shuffle", true, "shuffle lexicon")
	port := flag.Int("port", 8080, "port to listen on")
	flag.Parse()

	args := &anagrammit.GeneratorArgs{
		WordLength:  *wordLen,
		ResultLimit: *limit,
		Shuffle:     *shuffle,
		WordsFile:   "words/common-word-list.txt",
	}

	generator = anagrammit.NewGenerator(args)

	http.HandleFunc("/generate", anagramHandler)
	fmt.Println("Listening on", fmt.Sprintf(":%v", *port))
	err := http.ListenAndServe(fmt.Sprintf(":%v", *port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
