package spellcheck

import (
	"sort"
	"strings"
)

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
				suggestions = append(suggestions, &Suggestion{
					Word:      word,
					Frequency: wt.GetWordFrequency(word),
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
				suggestions = append(suggestions, Suggestion{
					Word:      word,
					Frequency: wt.GetWordFrequency(word),
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
