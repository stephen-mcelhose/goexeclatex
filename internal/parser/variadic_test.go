package parser_test

import (
	"testing"

	"github.com/stephen-mcelhose/goexeclatex/internal/parser"
)

// Tests for docs/specs/parser-extensions.md §5 (variadic paren args).

func TestMinTwoArgs(t *testing.T) {
	node := parse(t, `\min(3,1)`)
	fn, ok := node.(*parser.FunctionNode)
	if !ok || fn.Name != "min" || len(fn.Args) != 2 {
		t.Fatalf("want min/2, got %#v", node)
	}
}

func TestMaxThreeArgs(t *testing.T) {
	node := parse(t, `\max(3,1,2)`)
	fn, ok := node.(*parser.FunctionNode)
	if !ok || fn.Name != "max" || len(fn.Args) != 3 {
		t.Fatalf("want max/3, got %#v", node)
	}
}

func TestGcdTwoArgs(t *testing.T) {
	node := parse(t, `\gcd(12,8)`)
	fn, ok := node.(*parser.FunctionNode)
	if !ok || fn.Name != "gcd" || len(fn.Args) != 2 {
		t.Fatalf("want gcd/2, got %#v", node)
	}
}

func TestMinArityOneError(t *testing.T) {
	mustFailWith(t, `\min(1)`, "at least 2")
}

func TestSinParenStillWorks(t *testing.T) {
	node := parse(t, `\sin(x)`)
	fn, ok := node.(*parser.FunctionNode)
	if !ok || fn.Name != "sin" || len(fn.Args) != 1 {
		t.Fatalf("want sin/1, got %#v", node)
	}
}

func TestBareCommaStillErrors(t *testing.T) {
	mustFailWith(t, `1,2`, "COMMA")
}
