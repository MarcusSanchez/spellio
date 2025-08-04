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

	sentence := []string{"cant", "dont", "wont", "its", "youre", "Johns", "book's", "thats", "well", "Im"}
	for _, word := range sentence {
		if wt.IsWord(word) {
			fmt.Printf("%s ", word)
		} else {
			suggestions := wt.AutocorrectMultiple(word, 3)
			if len(suggestions) > 0 {
				if len(suggestions) == 1 {
					fmt.Printf("%s ", suggestions[0].Word)
				} else {
					fmt.Printf("%s[", suggestions[0].Word)
					for i := 1; i < len(suggestions); i++ {
						if i > 1 {
							fmt.Printf(",")
						}
						fmt.Printf("%s(%.2f)", suggestions[i].Word, suggestions[i].Confidence)
					}
					fmt.Printf("] ")
				}
			} else {
				fmt.Printf("[%s] ", word)
			}
		}
	}
}
