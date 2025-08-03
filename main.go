package main

import (
	"fmt"
	"log"
	"os"
	"spellio/wordfreqs"
	"spellio/wordtrie"
)

func main() {
	wf, err := wordfreqs.New()
	if err != nil {
		log.Println("failed to load english frequency map: ", err)
		os.Exit(1)
	}

	wt, err := wordtrie.New(wf)
	if err != nil {
		log.Println("failed to load word-trie: ", err)
		os.Exit(1)
	}

	sentence := []string{"speling", "erors", "arre", "comon", "in", "impropppr", "sentenses"}
	for _, word := range sentence {
		if wt.IsWord(word) {
			fmt.Printf("%s ", word)
		} else {
			correction, ok := wt.Autocorrect(word)
			if ok {
				fmt.Printf("%s ", correction.Word)
			} else {
				fmt.Printf("[%s] ", word)
			}
		}
	}
}
