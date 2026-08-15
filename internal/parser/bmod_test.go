package parser_test

import (
	"testing"

	"github.com/stephen-mcelhose/goexeclatex/internal/parser"
)

func TestBmodBinary(t *testing.T) {
	node := parse(t, `10 \bmod 3`)
	bn, ok := node.(*parser.BinaryNode)
	if !ok || bn.Op != "bmod" {
		t.Fatalf("want BinaryNode bmod, got %#v", node)
	}
}

func TestBmodNotImplicitMultiply(t *testing.T) {
	node := parse(t, `a \bmod b`)
	bn, ok := node.(*parser.BinaryNode)
	if !ok || bn.Op != "bmod" {
		t.Fatalf("want bmod BinaryNode, got %#v", node)
	}
}
