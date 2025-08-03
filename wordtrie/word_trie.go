package wordtrie

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"spellio/levenshtein"
	"spellio/wordfreqs"
	"strings"
)

type LetterNode struct {
	Children map[rune]*LetterNode
	IsWord   bool
}

type WordTrie struct {
	Root      *LetterNode
	WordFreqs *wordfreqs.WordFrequencies
}

func New(f *wordfreqs.WordFrequencies) (*WordTrie, error) {
	wt := &WordTrie{
		Root:      &LetterNode{Children: make(map[rune]*LetterNode)},
		WordFreqs: f,
	}
	if err := wt.loadWords("resources/words.txt"); err != nil {
		return nil, fmt.Errorf("failed to load words: %w", err)
	}
	return wt, nil
}

func (wt *WordTrie) Insert(word string) {
	word = strings.ToLower(word)
	n := wt.Root
	for _, ch := range word {
		if _, ok := n.Children[ch]; !ok {
			n.Children[ch] = &LetterNode{Children: make(map[rune]*LetterNode)}
		}
		n = n.Children[ch]
	}
	n.IsWord = true
}

func (wt *WordTrie) IsWord(word string) bool {
	word = strings.ToLower(word)
	n := wt.Root
	for _, ch := range word {
		if _, ok := n.Children[ch]; !ok {
			return false
		}
		n = n.Children[ch]
	}
	return n.IsWord
}

type Candidate struct {
	Word      string
	Distance  int
	Frequency int
}

func (wt *WordTrie) FindCandidates(word string, maxDist, N int) []Candidate {
	var candidates []Candidate
	wt.collectWords(func(candidate string) {
		dist := levenshtein.Distance(word, candidate)
		if dist <= maxDist {
			frequency := wt.WordFreqs.Frequencies[candidate]
			candidates = append(candidates, Candidate{
				Word:      candidate,
				Distance:  dist,
				Frequency: frequency,
			})
		}
	})
	// Sort by Distance, then Frequency
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Distance == candidates[j].Distance {
			return candidates[i].Frequency > candidates[j].Frequency
		}
		return candidates[i].Distance < candidates[j].Distance
	})
	if len(candidates) > N {
		candidates = candidates[:N]
	}
	return candidates
}

type Correction struct {
	Word      string
	Distance  int
	Frequency int
}

func (wt *WordTrie) Autocorrect(word string, md ...int) (Correction, bool) {
	maxDist := 2
	if len(md) > 0 {
		maxDist = md[0]
	}

	word = strings.ToLower(word)
	candidates := wt.FindCandidates(word, maxDist, 1_000_000)

	best := Correction{
		Word:      "",
		Distance:  maxDist + 1,
		Frequency: -1,
	}
	for _, c := range candidates {
		if c.Word == word {
			continue
		}
		if c.Distance < best.Distance || (c.Distance == best.Distance && c.Frequency > best.Frequency) {
			best.Word = c.Word
			best.Distance = c.Distance
			best.Frequency = c.Frequency
		}
	}
	return best, best.Word != ""
}

type Suggestion struct {
	Word      string
	Frequency int
}

func (wt *WordTrie) Autosuggest(prefix string) (*Suggestion, bool) {
	prefix = strings.ToLower(prefix)

	node := wt.Root
	for _, ch := range prefix {
		n, ok := node.Children[ch]
		if !ok {
			return nil, false
		}
		node = n
	}

	var suggestions []*Suggestion
	var dfs func(n *LetterNode, current []rune)
	dfs = func(n *LetterNode, current []rune) {
		if n.IsWord {
			word := string(current)
			if word != prefix { // Skip the exact prefix match
				frequency := wt.WordFreqs.Frequencies[word]
				suggestions = append(suggestions, &Suggestion{
					Word:      word,
					Frequency: frequency,
				})
			}
		}
		for ch, child := range n.Children {
			dfs(child, append(current, ch))
		}
	}
	dfs(node, []rune(prefix))

	if len(suggestions) == 0 {
		return nil, false
	}

	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Frequency > suggestions[j].Frequency
	})

	return suggestions[0], suggestions[0].Frequency > 1
}

func (wt *WordTrie) loadWords(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			wt.Insert(strings.ToLower(word))
		}
	}
	return scanner.Err()
}

func (wt *WordTrie) collectWords(fn func(string)) {
	var dfs func(node *LetterNode, prefix []rune)
	dfs = func(node *LetterNode, prefix []rune) {
		if node.IsWord {
			fn(string(prefix))
		}
		for ch, child := range node.Children {
			prefix = append(prefix, ch)
			dfs(child, prefix)
			prefix = prefix[:len(prefix)-1]
		}
	}
	dfs(wt.Root, []rune{})
}
