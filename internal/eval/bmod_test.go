package eval_test

import "testing"

func TestEvalBmod(t *testing.T) {
	tests := []struct {
		in   string
		want float64
	}{
		{`10 \bmod 3`, 1},
		{`15 \bmod 4`, 3},
		{`7 \bmod 7`, 0},
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

func TestEvalBmodZero(t *testing.T) {
	mustError(t, `5 \bmod 0`, "eval: division by zero")
}
