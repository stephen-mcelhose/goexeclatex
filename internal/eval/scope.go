package eval

import "math"

// Scope maps normalised symbol names (backslash-stripped, lowercased) to their
// float64 values. Callers may add user-defined variables before passing to Eval.
type Scope map[string]float64

// NewScope returns a Scope pre-seeded with the built-in constants (spec §5).
func NewScope() Scope {
	return Scope{
		"pi":  math.Pi,
		"e":   math.E,
		"tau": 2 * math.Pi,
		"phi": (1 + math.Sqrt(5)) / 2,
	}
}
