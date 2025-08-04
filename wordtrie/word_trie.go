package wordtrie

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"spellio/levenshtein"
	"spellio/wordfreqs"
	"strings"
	"unicode"
)

type LetterNode struct {
	Children map[rune]*LetterNode
	IsWord   bool
}

type WordTrie struct {
	Root         *LetterNode
	WordFreqs    *wordfreqs.WordFrequencies
	contractions map[string]string
}

func New(f *wordfreqs.WordFrequencies) (*WordTrie, error) {
	wt := &WordTrie{
		Root:         &LetterNode{Children: make(map[rune]*LetterNode)},
		WordFreqs:    f,
		contractions: initContractions(),
	}
	if err := wt.loadWords("resources/words.txt"); err != nil {
		return nil, fmt.Errorf("failed to load words: %w", err)
	}
	return wt, nil
}

func initContractions() map[string]string {
	return map[string]string{
		"cant":     "can't",
		"wont":     "won't",
		"dont":     "don't",
		"isnt":     "isn't",
		"arent":    "aren't",
		"wasnt":    "wasn't",
		"werent":   "weren't",
		"hasnt":    "hasn't",
		"havent":   "haven't",
		"hadnt":    "hadn't",
		"wouldnt":  "wouldn't",
		"couldnt":  "couldn't",
		"shouldnt": "shouldn't",
		"mustnt":   "mustn't",
		"neednt":   "needn't",
		"oughtnt":  "oughtn't",
		"shant":    "shan't",
		"darent":   "daren't",
		"youre":    "you're",
		"theyre":   "they're",
		"were":     "we're",
		"youve":    "you've",
		"theyve":   "they've",
		"weve":     "we've",
		"ive":      "I've",
		"youll":    "you'll",
		"theyll":   "they'll",
		"well":     "we'll",
		"hell":     "he'll",
		"shell":    "she'll",
		"itll":     "it'll",
		"thatll":   "that'll",
		"wholl":    "who'll",
		"whatll":   "what'll",
		"wherll":   "where'll",
		"whenll":   "when'll",
		"whyll":    "why'll",
		"howll":    "how'll",
		"youd":     "you'd",
		"theyd":    "they'd",
		"wed":      "we'd",
		"hed":      "he'd",
		"shed":     "she'd",
		"itd":      "it'd",
		"thatd":    "that'd",
		"whod":     "who'd",
		"whatd":    "what'd",
		"whered":   "where'd",
		"whend":    "when'd",
		"whyd":     "why'd",
		"howd":     "how'd",
		"im":       "I'm",
		"lets":     "let's",
		"thats":    "that's",
		"whats":    "what's",
		"wheres":   "where's",
		"whens":    "when's",
		"whys":     "why's",
		"hows":     "how's",
		"whos":     "who's",
		"heres":    "here's",
		"theres":   "there's",
	}
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
	word = normalizeApostrophe(strings.ToLower(word))
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
	wordLen := len(word)
	wt.collectWords(func(candidate string) {
		candidateLen := len(candidate)
		if int(math.Abs(float64(wordLen-candidateLen))) > maxDist {
			return
		}
		dist := levenshtein.DistanceWithThreshold(word, candidate, maxDist)
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
	Word       string
	Distance   int
	Frequency  int
	Confidence float64
}

func (wt *WordTrie) Autocorrect(word string, md ...int) (Correction, bool) {
	suggestions := wt.AutocorrectMultiple(word, 1, md...)
	if len(suggestions) > 0 {
		return suggestions[0], true
	}
	return Correction{}, false
}

func (wt *WordTrie) AutocorrectMultiple(word string, maxSuggestions int, md ...int) []Correction {
	maxDist := 2
	if len(md) > 0 {
		maxDist = md[0]
	}

	originalWord := word
	word = normalizeApostrophe(strings.ToLower(word))

	if contraction, exists := wt.contractions[word]; exists {
		return []Correction{{
			Word:       wt.preserveCase(originalWord, contraction),
			Distance:   0,
			Frequency:  1000000,
			Confidence: 1.0,
		}}
	}

	if wt.isPossessive(word) {
		baseWord := strings.TrimSuffix(word, "'s")
		if wt.IsWord(baseWord) {
			correctedPossessive := wt.preserveCase(originalWord, baseWord+"'s")
			return []Correction{{
				Word:       correctedPossessive,
				Distance:   0,
				Frequency:  wt.WordFreqs.Frequencies[baseWord],
				Confidence: 0.95,
			}}
		}
	}

	candidates := wt.FindCandidates(word, maxDist, 1_000_000)

	var corrections []Correction
	for _, c := range candidates {
		if c.Word == word {
			continue
		}
		confidence := wt.calculateConfidence(c.Distance, c.Frequency, maxDist)
		corrections = append(corrections, Correction{
			Word:       c.Word,
			Distance:   c.Distance,
			Frequency:  c.Frequency,
			Confidence: confidence,
		})
	}

	sort.Slice(corrections, func(i, j int) bool {
		// Calculate composite scores that balance distance and frequency
		// Score = distance - log10(frequency) * scaling_factor
		// Lower scores rank higher
		scalingFactor := 0.6
		
		scoreI := float64(corrections[i].Distance)
		scoreJ := float64(corrections[j].Distance)
		
		// Apply frequency scaling if frequency > 0
		if corrections[i].Frequency > 0 {
			scoreI -= math.Log10(float64(corrections[i].Frequency)) * scalingFactor
		}
		if corrections[j].Frequency > 0 {
			scoreJ -= math.Log10(float64(corrections[j].Frequency)) * scalingFactor
		}
		
		// Primary sort by composite score
		if math.Abs(scoreI - scoreJ) > 0.001 { // Use small threshold for float comparison
			return scoreI < scoreJ
		}
		
		// Tie-breaker: keyboard distance
		keyboardDistI := levenshtein.KeyboardAwareDistance(word, corrections[i].Word)
		keyboardDistJ := levenshtein.KeyboardAwareDistance(word, corrections[j].Word)
		return keyboardDistI < keyboardDistJ
	})

	if len(corrections) > maxSuggestions {
		corrections = corrections[:maxSuggestions]
	}

	return corrections
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

func (wt *WordTrie) AutosuggestMultiple(prefix string, maxSuggestions int) []Suggestion {
	prefix = strings.ToLower(prefix)

	node := wt.Root
	for _, ch := range prefix {
		n, ok := node.Children[ch]
		if !ok {
			return nil
		}
		node = n
	}

	var suggestions []Suggestion
	var dfs func(n *LetterNode, current []rune)
	dfs = func(n *LetterNode, current []rune) {
		if n.IsWord {
			word := string(current)
			if word != prefix { // Skip the exact prefix match
				frequency := wt.WordFreqs.Frequencies[word]
				suggestions = append(suggestions, Suggestion{
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
		return nil
	}

	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Frequency > suggestions[j].Frequency
	})

	if len(suggestions) > maxSuggestions {
		suggestions = suggestions[:maxSuggestions]
	}

	return suggestions
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
