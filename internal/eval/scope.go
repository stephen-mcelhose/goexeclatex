package eval

import "math"

// ScopeLookup is the read interface for symbol lookup used by Eval.
// Both Scope (flat map) and innerScope (chained for large-op iteration) satisfy it.
type ScopeLookup interface {
	Lookup(name string) (float64, bool)
}

// Scope maps normalised symbol names (backslash-stripped, lowercased) to their
// float64 values. Callers may add user-defined variables before passing to Eval.
type Scope map[string]float64

// Lookup implements ScopeLookup.
func (s Scope) Lookup(name string) (float64, bool) {
	v, ok := s[name]
	return v, ok
}

// NewScope returns a Scope pre-seeded with the built-in constants (spec §5).
func NewScope() Scope {
	return Scope{
		"pi":    math.Pi,
		"e":     math.E,
		"tau":   2 * math.Pi,
		"phi":   (1 + math.Sqrt(5)) / 2,
		"infty": math.Inf(1),
	}
}

// innerScope wraps a parent ScopeLookup and shadows exactly one variable.
// Used by parseLargeOp evaluation to bind the iteration variable without
// mutating the outer scope.
type innerScope struct {
	parent ScopeLookup
	name   string
	value  float64
}

func (s *innerScope) Lookup(name string) (float64, bool) {
	if name == s.name {
		return s.value, true
	}
	return s.parent.Lookup(name)
}
