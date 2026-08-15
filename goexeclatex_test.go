package goexeclatex_test

import (
	"errors"
	"math"
	"testing"

	"github.com/stephen-mcelhose/goexeclatex"
)

func TestEvalBasic(t *testing.T) {
	tests := []struct {
		expr string
		vars map[string]float64
		want float64
	}{
		{`\frac{1}{2}`, nil, 0.5},
		{`\sqrt{9}`, nil, 3},
		{`1+2*3`, nil, 7},
		{`x^2 + 2x + 1`, map[string]float64{"x": 3}, 16},
		{`\sin{\pi/6}`, nil, 0.5},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			got, err := goexeclatex.Eval(tt.expr, tt.vars)
			if err != nil {
				t.Fatalf("Eval(%q): %v", tt.expr, err)
			}
			if math.Abs(got-tt.want) > 1e-10 {
				t.Errorf("Eval(%q) = %v, want %v", tt.expr, got, tt.want)
			}
		})
	}
}

func TestEvalNilVars(t *testing.T) {
	got, err := goexeclatex.Eval(`1+1`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != 2 {
		t.Errorf("got %v, want 2", got)
	}
}

func TestEvalSyntaxError(t *testing.T) {
	_, err := goexeclatex.Eval(`\frac{1}`, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	var syn *goexeclatex.SyntaxError
	if !errors.As(err, &syn) {
		t.Fatalf("want *SyntaxError, got %T: %v", err, err)
	}
	var ev *goexeclatex.EvalError
	if errors.As(err, &ev) {
		t.Fatal("syntax failure must not also be *EvalError")
	}
}

func TestEvalEvalError(t *testing.T) {
	_, err := goexeclatex.Eval(`\frac{1}{0}`, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	var ev *goexeclatex.EvalError
	if !errors.As(err, &ev) {
		t.Fatalf("want *EvalError, got %T: %v", err, err)
	}
	var syn *goexeclatex.SyntaxError
	if errors.As(err, &syn) {
		t.Fatal("eval failure must not also be *SyntaxError")
	}
}

func TestEvalUndefinedSymbol(t *testing.T) {
	_, err := goexeclatex.Eval(`xyz`, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	var ev *goexeclatex.EvalError
	if !errors.As(err, &ev) {
		t.Fatalf("want *EvalError, got %T: %v", err, err)
	}
}

// library.md §4.4 — empty input is left to lexer/parser; currently *SyntaxError (parser EOF).
func TestEvalEmptyExpression(t *testing.T) {
	for _, expr := range []string{"", "   "} {
		_, err := goexeclatex.Eval(expr, nil)
		if err == nil {
			t.Fatalf("Eval(%q) = nil error, want syntax error", expr)
		}
		var syn *goexeclatex.SyntaxError
		if !errors.As(err, &syn) {
			t.Fatalf("Eval(%q): want *SyntaxError, got %T: %v", expr, err, err)
		}
	}
}

// library.md §4.2 / §9 — user keys MAY override builtins.
func TestEvalOverrideBuiltin(t *testing.T) {
	got, err := goexeclatex.Eval(`\pi`, map[string]float64{"pi": 3})
	if err != nil {
		t.Fatal(err)
	}
	if got != 3 {
		t.Errorf("override pi: got %v, want 3", got)
	}
}

// library.md §4.3 — ±Inf is success (ADR-010).
func TestEvalInftySuccess(t *testing.T) {
	got, err := goexeclatex.Eval(`\infty`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !math.IsInf(got, 1) {
		t.Errorf("got %v, want +Inf", got)
	}
}
