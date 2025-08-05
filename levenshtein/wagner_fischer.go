// Package levenshtein is a wagner-fischer implementation of the levenshtein distance algorithm
// This implementation uses dynamic programming to compute the Levenshtein distance
// between two strings, which is the minimum number of single-character edits (insertions,
// deletions, or substitutions) required to change one string into the other.
package levenshtein

import "slices"

var keyboardLayout = map[rune][]rune{
	'q':  {'w', 'a', 's'},
	'w':  {'q', 'e', 'a', 's', 'd'},
	'e':  {'w', 'r', 's', 'd', 'f'},
	'r':  {'e', 't', 'd', 'f', 'g'},
	't':  {'r', 'y', 'f', 'g', 'h'},
	'y':  {'t', 'u', 'g', 'h', 'j'},
	'u':  {'y', 'i', 'h', 'j', 'k'},
	'i':  {'u', 'o', 'j', 'k', 'l'},
	'o':  {'i', 'p', 'k', 'l'},
	'p':  {'o', 'l'},
	'a':  {'q', 'w', 's', 'z'},
	's':  {'a', 'w', 'e', 'd', 'z', 'x'},
	'd':  {'s', 'e', 'r', 'f', 'x', 'c'},
	'f':  {'d', 'r', 't', 'g', 'c', 'v'},
	'g':  {'f', 't', 'y', 'h', 'v', 'b'},
	'h':  {'g', 'y', 'u', 'j', 'b', 'n'},
	'j':  {'h', 'u', 'i', 'k', 'n', 'm'},
	'k':  {'j', 'i', 'o', 'l', 'm'},
	'l':  {'k', 'o', 'p', 'm', '\''},
	'\'': {'l', ';'},
	';':  {'l', '\'', 'p'},
	'z':  {'a', 's', 'x'},
	'x':  {'z', 's', 'd', 'c'},
	'c':  {'x', 'd', 'f', 'v'},
	'v':  {'c', 'f', 'g', 'b'},
	'b':  {'v', 'g', 'h', 'n'},
	'n':  {'b', 'h', 'j', 'm'},
	'm':  {'n', 'j', 'k', 'l'},
}

func keyboardDistance(a, b rune) int {
	if a == b {
		return 0
	}

	if adjacent, exists := keyboardLayout[a]; exists {
		if slices.Contains(adjacent, b) {
			return 9
		}
	}

	return 10
}

func Distance(a, b string) int {
	return DistanceWithThreshold(a, b, -1)
}

func KeyboardAwareDistance(a, b string) int {
	return KeyboardAwareDistanceWithThreshold(a, b, -1)
}

func DistanceWithThreshold(a, b string, threshold int) int {
	la, lb := len(a), len(b)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}
	// Allocate 2 rows instead of the full matrix for optimization
	prev := make([]int, lb+1)
	curr := make([]int, lb+1)
	for j := 0; j <= lb; j++ {
		prev[j] = j
	}
	for i := 1; i <= la; i++ {
		curr[0] = i
		minInRow := curr[0]
		for j := 1; j <= lb; j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			curr[j] = minimum(
				prev[j]+1,      // deletion
				curr[j-1]+1,    // insertion
				prev[j-1]+cost, // substitution
			)
			if curr[j] < minInRow {
				minInRow = curr[j]
			}
		}
		if threshold >= 0 && minInRow > threshold {
			return threshold + 1
		}
		prev, curr = curr, prev
	}
	return prev[lb]
}

func KeyboardAwareDistanceWithThreshold(a, b string, threshold int) int {
	la, lb := len(a), len(b)
	if la == 0 {
		return lb * 10
	}
	if lb == 0 {
		return la * 10
	}

	prev := make([]int, lb+1)
	curr := make([]int, lb+1)
	for j := 0; j <= lb; j++ {
		prev[j] = j * 10
	}
	for i := 1; i <= la; i++ {
		curr[0] = i * 10
		minInRow := curr[0]
		for j := 1; j <= lb; j++ {
			cost := keyboardDistance(rune(a[i-1]), rune(b[j-1]))
			curr[j] = minimum(
				prev[j]+10,     // deletion
				curr[j-1]+10,   // insertion
				prev[j-1]+cost, // substitution
			)
			if curr[j] < minInRow {
				minInRow = curr[j]
			}
		}
		if threshold >= 0 && minInRow > threshold*10 {
			return (threshold + 1) * 10
		}
		prev, curr = curr, prev
	}
	return prev[lb]
}

func minimum(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	} else {
		if b < c {
			return b
		}
		return c
	}
}
