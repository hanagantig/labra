package speller

import (
	"math"
)

type Speller struct {
	dict      []string
	threshold int
}

func New(dict []string, threshold int) *Speller {
	return &Speller{
		dict:      dict,
		threshold: threshold,
	}
}

func (s *Speller) CorrectSpelling(w string) string {
	for _, d := range s.dict {
		if w == d {
			return w
		}

		diff := len(w) - len(d)
		if math.Abs(float64(diff)) > float64(len(w)+s.threshold) {
			continue
		}

		dist := s.LevenshteinDistance(w, d)
		if dist <= 2 {
			return d
		}
	}

	return ""
}

// LevenshteinDistance measures the difference between two strings.
// The Levenshtein distance between two words is the minimum number of
// single-character edits (i.e. insertions, deletions or substitutions)
// required to change one word into the other.
//
// This implemention is optimized to use O(min(m,n)) space and is based on the
// optimized C version found here:
// http://en.wikibooks.org/wiki/Algorithm_implementation/Strings/Levenshtein_distance#C
func (s *Speller) LevenshteinDistance(src, t string) int {
	r1, r2 := []rune(src), []rune(t)
	column := make([]int, 1, 64)

	for y := 1; y <= len(r1); y++ {
		column = append(column, y)
	}

	for x := 1; x <= len(r2); x++ {
		column[0] = x

		for y, lastDiag := 1, x-1; y <= len(r1); y++ {
			oldDiag := column[y]
			cost := 0
			if r1[y-1] != r2[x-1] {
				cost = 1
			}
			column[y] = min(column[y]+1, column[y-1]+1, lastDiag+cost)
			lastDiag = oldDiag
		}
	}

	return column[len(r1)]
}

func min2(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func min(a, b, c int) int {
	return min2(min2(a, b), c)
}
