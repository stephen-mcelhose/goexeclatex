package eval_test

import (
	"math"
	"testing"
)

// Tests for docs/specs/parser-extensions.md §3.3 (floor / ceil eval).

func TestEvalFloor(t *testing.T) {
	tests := []struct {
		in   string
		want float64
	}{
		{`\lfloor 3.2 \rfloor`, 3},
		{`\lfloor -3.2 \rfloor`, -4},
		{`\lfloor 5 \rfloor`, 5},
		{`\lfloor 1+2.7 \rfloor`, 3},
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

func TestEvalCeil(t *testing.T) {
	tests := []struct {
		in   string
		want float64
	}{
		{`\lceil 3.2 \rceil`, 4},
		{`\lceil -3.2 \rceil`, -3},
		{`\lceil 5 \rceil`, 5},
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

func TestEvalFloorNested(t *testing.T) {
	got := evaluate(t, `\lfloor \lceil 3.2 \rceil \rfloor`)
	if got != 4 {
		t.Errorf("got %v, want 4", got)
	}
}

func TestEvalFloorVsAbs(t *testing.T) {
	// sanity: floor of negative differs from abs
	got := evaluate(t, `\lfloor -2.1 \rfloor`)
	if got != -3 {
		t.Errorf("got %v, want -3", got)
	}
	_ = math.Floor // keep math import if needed for future float compares
}
