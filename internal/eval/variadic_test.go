package eval_test

import "testing"

// Tests for docs/specs/parser-extensions.md §5.4 (min/max/gcd eval).

func TestEvalMinMaxGcd(t *testing.T) {
	tests := []struct {
		in   string
		want float64
	}{
		{`\min(3,1)`, 1},
		{`\min(3,1,2)`, 1},
		{`\max(3,1)`, 3},
		{`\max(3,1,2)`, 3},
		{`\gcd(12,8)`, 4},
		{`\gcd(12,8,20)`, 4},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := evaluate(t, tt.in)
			if got != tt.want {
				t.Errorf("Eval(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestEvalGcdDomain(t *testing.T) {
	mustError(t, `\gcd(12,8.5)`, "eval: domain error:")
	mustError(t, `\gcd(-12,8)`, "eval: domain error:")
}
