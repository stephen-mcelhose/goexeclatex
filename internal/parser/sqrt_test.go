package parser_test

import (
	"testing"

	"github.com/stephen-mcelhose/goexeclatex/internal/parser"
)

// Tests for docs/specs/parser-extensions.md §4 (\sqrt[n]{x}).

func TestSqrtPlain(t *testing.T) {
	node := parse(t, `\sqrt{9}`)
	fn, ok := node.(*parser.FunctionNode)
	if !ok || fn.Name != "sqrt" || len(fn.Args) != 1 {
		t.Fatalf("want sqrt/1, got %#v", node)
	}
}

func TestSqrtNthRoot(t *testing.T) {
	node := parse(t, `\sqrt[3]{27}`)
	fn, ok := node.(*parser.FunctionNode)
	if !ok {
		t.Fatalf("got %T, want FunctionNode", node)
	}
	if fn.Name != "sqrt" {
		t.Errorf("Name = %q, want sqrt", fn.Name)
	}
	if len(fn.Args) != 2 {
		t.Fatalf("Args len = %d, want 2 (radicand, index)", len(fn.Args))
	}
}

func TestSqrtNthRootNotImplicitMultiply(t *testing.T) {
	node := parse(t, `\sqrt[3]{27}`)
	if _, ok := node.(*parser.BinaryNode); ok {
		t.Fatal(`\sqrt[3]{27} must not be BinaryNode`)
	}
}

func TestSqrtCharModeArg(t *testing.T) {
	node := parse(t, `\sqrt9`)
	fn, ok := node.(*parser.FunctionNode)
	if !ok || fn.Name != "sqrt" || len(fn.Args) != 1 {
		t.Fatalf("want sqrt/1 for \\sqrt9, got %#v", node)
	}
}
