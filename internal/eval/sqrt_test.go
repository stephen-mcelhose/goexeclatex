package eval_test

import "testing"

// Tests for docs/specs/parser-extensions.md §4.3 (nth-root eval).

func TestEvalSqrtNthRoot(t *testing.T) {
	tests := []struct {
		in   string
		want float64
	}{
		{`\sqrt{9}`, 3},
		{`\sqrt[3]{27}`, 3},
		{`\sqrt[2]{16}`, 4},
		{`\sqrt[4]{16}`, 2},
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

func TestEvalSqrtNthRootDomain(t *testing.T) {
	mustError(t, `\sqrt[0]{8}`, "eval: domain error:")
	mustError(t, `\sqrt[3]{-8}`, "eval: domain error:")
}
