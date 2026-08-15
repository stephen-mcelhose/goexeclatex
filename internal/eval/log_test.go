package eval_test

import "testing"

// Tests for docs/specs/parser-extensions.md §6.3 (log base eval).

func TestEvalLogBase(t *testing.T) {
	tests := []struct {
		in   string
		want float64
	}{
		{`\log{100}`, 2},
		{`\log_{2}(8)`, 3},
		{`\log_{10}(1000)`, 3},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := evaluate(t, tt.in)
			if !approx(got, tt.want) {
				t.Errorf("Eval(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestEvalLogBaseDomain(t *testing.T) {
	mustError(t, `\log_{2}(0)`, "eval: domain error:")
	mustError(t, `\log_{1}(8)`, "eval: domain error:")
	mustError(t, `\log_{-2}(8)`, "eval: domain error:")
}
