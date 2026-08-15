package parser_test

import (
	"testing"

	"github.com/stephen-mcelhose/goexeclatex/internal/parser"
)

// Tests for docs/specs/parser-extensions.md §6 (\log_{b}(x)).

func TestLogBase10Brace(t *testing.T) {
	node := parse(t, `\log{100}`)
	fn, ok := node.(*parser.FunctionNode)
	if !ok || fn.Name != "log" || len(fn.Args) != 1 {
		t.Fatalf("want log/1, got %#v", node)
	}
}

func TestLogBaseSubscriptParen(t *testing.T) {
	node := parse(t, `\log_{2}(8)`)
	fn, ok := node.(*parser.FunctionNode)
	if !ok || fn.Name != "log" || len(fn.Args) != 2 {
		t.Fatalf("want log/2, got %#v", node)
	}
}

func TestLogBaseRequiresParen(t *testing.T) {
	mustFailWith(t, `\log_{2}{8}`, "expects (x)")
}
