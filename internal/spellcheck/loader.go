package spellcheck

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func New() (*WordTrie, error) {
	wt := NewWordTrie()
	if err := wt.loadWords("resources/english_words_freqs.txt"); err != nil {
		return nil, fmt.Errorf("failed to load words: %w", err)
	}
	return wt, nil
}

func (wt *WordTrie) loadWords(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pair := strings.Split(strings.TrimSpace(scanner.Text()), ",")
		word := pair[0]
		frequency, _ := strconv.Atoi(pair[1])

		if word != "" {
			wt.Insert(strings.ToLower(word), frequency)
		}
	}
	return scanner.Err()
}
