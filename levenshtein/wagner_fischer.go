// Package levenshtein is a wagner-fischer implementation of the levenshtein distance algorithm
// This implementation uses dynamic programming to compute the Levenshtein distance
// between two strings, which is the minimum number of single-character edits (insertions,
// deletions, or substitutions) required to change one string into the other.
package levenshtein

func Distance(a, b string) int {
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
