package spellcheck

import (
	"math"
	"sort"
	"spellio/levenshtein"
	"strings"
)

type Candidate struct {
	Word      string
	Distance  int
	Frequency int
}

type Correction struct {
	Word       string
	Distance   int
	Frequency  int
	Confidence float64
}

func (wt *WordTrie) FindCandidates(word string, maxDist, N int) []Candidate {
	var candidates []Candidate
	wordLen := len(word)
	wt.collectWords(func(candidate string, frequency int) {
		candidateLen := len(candidate)
		if int(math.Abs(float64(wordLen-candidateLen))) > maxDist {
			return
		}
		dist := levenshtein.DistanceWithThreshold(word, candidate, maxDist)
		if dist <= maxDist {
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

	if contraction, exists := contractions[word]; exists {
		return []Correction{{
			Word:       wt.preserveCase(originalWord, contraction),
			Distance:   0,
			Frequency:  1000000,
			Confidence: 1.0,
		}}
	}

	// Check for pattern-based correction but don't return immediately -
	// let it be prioritized in the full candidate search
	var patternCorrection string
	if correction, exists := commonMisspellings[word]; exists && wt.IsWord(correction) {
		patternCorrection = correction
	}

	if wt.isPossessive(word) {
		baseWord := strings.TrimSuffix(word, "'s")
		if wt.IsWord(baseWord) {
			correctedPossessive := wt.preserveCase(originalWord, baseWord+"'s")
			return []Correction{{
				Word:       correctedPossessive,
				Distance:   0,
				Frequency:  wt.GetWordFrequency(baseWord),
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

		// Boost confidence for pattern-based corrections
		if patternCorrection != "" && c.Word == patternCorrection {
			confidence = 0.98
		}

		corrections = append(corrections, Correction{
			Word:       c.Word,
			Distance:   c.Distance,
			Frequency:  c.Frequency,
			Confidence: confidence,
		})
	}

	sort.Slice(corrections, func(i, j int) bool {
		// Primary sort: High-confidence pattern corrections first
		if corrections[i].Confidence >= 0.98 && corrections[j].Confidence < 0.98 {
			return true
		}
		if corrections[j].Confidence >= 0.98 && corrections[i].Confidence < 0.98 {
			return false
		}

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
		if math.Abs(scoreI-scoreJ) > 0.001 { // Use a small threshold for float comparison
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
