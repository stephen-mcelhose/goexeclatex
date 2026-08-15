package eval

import (
	"strings"
	"testing"
)

// Defense-in-depth for library.md §6 / issue #20: under-arity must error, not panic.
func TestCallBuiltinMinMaxGcdArity(t *testing.T) {
	tests := []struct {
		name string
		args []float64
	}{
		{"min", nil},
		{"min", []float64{1}},
		{"max", nil},
		{"max", []float64{1}},
		{"gcd", nil},
		{"gcd", []float64{1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("callBuiltin(%q, %v) panicked: %v", tt.name, tt.args, r)
				}
			}()
			_, err := callBuiltin(tt.name, tt.args)
			if err == nil {
				t.Fatalf("callBuiltin(%q, %v) = nil error, want arity error", tt.name, tt.args)
			}
			if !strings.Contains(err.Error(), "requires at least 2") {
				t.Errorf("error %q does not mention arity", err)
			}
		})
	}
}
