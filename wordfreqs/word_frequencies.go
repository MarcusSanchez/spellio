package wordfreqs

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type WordFrequencies struct {
	Frequencies map[string]int
}

func New() (*WordFrequencies, error) {
	f := &WordFrequencies{Frequencies: make(map[string]int, 333_334)}
	if err := f.loadFrequencies("resources/english_words_freqs.txt"); err != nil {
		return nil, fmt.Errorf("failed to load frequencies: %w", err)
	}
	return f, nil
}

func (em *WordFrequencies) Insert(word string, frequency int) {
	word = strings.ToLower(word)
	em.Frequencies[word] = frequency
}

func (em *WordFrequencies) loadFrequencies(filename string) error {
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
			em.Insert(strings.ToLower(word), frequency)
		}
	}
	return scanner.Err()
}
