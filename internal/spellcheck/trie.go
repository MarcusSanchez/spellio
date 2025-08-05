package spellcheck

import (
	"strings"
	"unicode"
)

type LetterNode struct {
	Children  map[rune]*LetterNode
	IsWord    bool
	Frequency int
}

type WordTrie struct {
	Root *LetterNode
}

func NewWordTrie() *WordTrie {
	return &WordTrie{
		Root: &LetterNode{Children: make(map[rune]*LetterNode)},
	}
}

func (wt *WordTrie) Insert(word string, frequency int) {
	word = strings.ToLower(word)
	n := wt.Root
	for _, ch := range word {
		if _, ok := n.Children[ch]; !ok {
			n.Children[ch] = &LetterNode{Children: make(map[rune]*LetterNode)}
		}
		n = n.Children[ch]
	}
	n.IsWord = true
	n.Frequency = frequency
}

func (wt *WordTrie) IsWord(word string) bool {
	word = normalizeApostrophe(strings.ToLower(word))
	if _, ok := contractions[word]; ok {
		return false
	}
	if _, ok := commonMisspellings[word]; ok {
		return false
	}
	if _, ok := contractionCorrections[word]; ok {
		return true
	}
	if wt.isPossessive(word) {
		baseWord := strings.TrimSuffix(word, "'s")
		return wt.IsWord(baseWord)
	}

	n := wt.Root
	for _, ch := range word {
		if _, ok := n.Children[ch]; !ok {
			return false
		}
		n = n.Children[ch]
	}
	return n.IsWord
}

func (wt *WordTrie) GetWordFrequency(word string) int {
	word = strings.ToLower(word)
	n := wt.Root
	for _, ch := range word {
		if _, ok := n.Children[ch]; !ok {
			return 0 // Word not found
		}
		n = n.Children[ch]
	}
	if n.IsWord {
		return n.Frequency
	}
	return 0 // Not a valid word
}

func (wt *WordTrie) collectWords(fn func(string, int)) {
	var dfs func(node *LetterNode, prefix []rune)
	dfs = func(node *LetterNode, prefix []rune) {
		if node.IsWord {
			fn(string(prefix), node.Frequency)
		}
		for ch, child := range node.Children {
			prefix = append(prefix, ch)
			dfs(child, prefix)
			prefix = prefix[:len(prefix)-1]
		}
	}
	dfs(wt.Root, []rune{})
}

func (wt *WordTrie) isPossessive(word string) bool {
	return strings.HasSuffix(word, "'s") && len(word) > 2
}

func (wt *WordTrie) preserveCase(original, corrected string) string {
	if len(original) == 0 || len(corrected) == 0 {
		return corrected
	}

	if unicode.IsUpper(rune(original[0])) {
		if len(corrected) > 0 {
			runes := []rune(corrected)
			runes[0] = unicode.ToUpper(runes[0])
			return string(runes)
		}
	}

	return corrected
}

func (wt *WordTrie) calculateConfidence(distance, frequency, maxDistance int) float64 {
	distanceScore := 1.0 - (float64(distance) / float64(maxDistance+1))

	maxFreq := 1000000.0
	if frequency > int(maxFreq) {
		maxFreq = float64(frequency)
	}
	frequencyScore := float64(frequency) / maxFreq

	confidence := (distanceScore * 0.7) + (frequencyScore * 0.3)
	if confidence > 1.0 {
		confidence = 1.0
	}
	return confidence
}

func normalizeApostrophe(word string) string {
	result := make([]rune, 0, len(word))
	for _, r := range word {
		if r == '\u2018' || r == '\u2019' || r == '`' {
			result = append(result, '\'')
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}
