package parser_test

import (
	"testing"

	"github.com/stephen-mcelhose/goexeclatex/internal/parser"
)

// Tests for docs/specs/parser-extensions.md §3.2 (floor / ceil).

func TestFloorNode(t *testing.T) {
	node := parse(t, `\lfloor 3.2 \rfloor`)
	fn, ok := node.(*parser.FunctionNode)
	if !ok {
		t.Fatalf(`parse(\lfloor 3.2 \rfloor) = %T, want *FunctionNode`, node)
	}
	if fn.Name != "floor" {
		t.Errorf("Name = %q, want floor", fn.Name)
	}
	if len(fn.Args) != 1 {
		t.Fatalf("Args len = %d, want 1", len(fn.Args))
	}
}

func TestCeilNode(t *testing.T) {
	node := parse(t, `\lceil 3.2 \rceil`)
	fn, ok := node.(*parser.FunctionNode)
	if !ok {
		t.Fatalf(`parse = %T, want *FunctionNode`, node)
	}
	if fn.Name != "ceil" {
		t.Errorf("Name = %q, want ceil", fn.Name)
	}
}

func TestFloorNotImplicitMultiply(t *testing.T) {
	node := parse(t, `\lfloor 3.2 \rfloor`)
	if _, ok := node.(*parser.BinaryNode); ok {
		t.Fatalf(`\lfloor 3.2 \rfloor must not parse as BinaryNode (implicit multiply)`)
	}
}

func TestFloorEmptyError(t *testing.T) {
	mustFailWith(t, `\lfloor\rfloor`, "empty")
}

func TestFloorUnmatchedError(t *testing.T) {
	mustFailWith(t, `\lfloor 3.2`, "lfloor")
}

func TestFloorMixedCeilError(t *testing.T) {
	mustFailWith(t, `\lfloor 3.2 \rceil`, "mismatch")
}

func TestCeilMixedFloorError(t *testing.T) {
	mustFailWith(t, `\lceil 3.2 \rfloor`, "mismatch")
}

func TestFloorExpression(t *testing.T) {
	node := parse(t, `\lfloor 1 + 2.5 \rfloor`)
	fn, ok := node.(*parser.FunctionNode)
	if !ok || fn.Name != "floor" {
		t.Fatalf("want floor FunctionNode, got %T %#v", node, node)
	}
}
