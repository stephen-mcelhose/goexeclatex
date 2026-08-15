package parser_test

// Tests for spec §4.1 (subscript rule), §4.2 (large operators), §4.3 (norm rule).
// All tests in this file are RED until the corresponding implementation lands.

import (
	"strings"
	"testing"

	"github.com/stephen-mcelhose/goexeclatex/internal/lexer"
	"github.com/stephen-mcelhose/goexeclatex/internal/parser"
)

// mustFailWith expects a lex or parse error whose message contains want.
// Unlike the existing mustFail helper, it verifies the error message.
func mustFailWith(t *testing.T, input, want string) {
	t.Helper()
	tokens, err := lexer.Lex(input)
	if err != nil {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("lex(%q) error = %q, want it to contain %q", input, err.Error(), want)
		}
		return
	}
	_, err = parser.Parse(tokens)
	if err == nil {
		t.Errorf("parse(%q) expected error, got nil", input)
		return
	}
	if !strings.Contains(err.Error(), want) {
		t.Errorf("parse(%q) error = %q, want it to contain %q", input, err.Error(), want)
	}
}

// ── §4.1 Subscript rule ───────────────────────────────────────────────────────

// TestSubscriptSymbolIndex — x_i → SubscriptNode{SymbolNode(x), SymbolNode(i)}
func TestSubscriptSymbolIndex(t *testing.T) {
	node := parse(t, "x_i")
	n, ok := node.(*parser.SubscriptNode)
	if !ok {
		t.Fatalf("parse(%q) = %T, want *SubscriptNode", "x_i", node)
	}
	base, ok := n.Base.(*parser.SymbolNode)
	if !ok || base.Name != "x" {
		t.Errorf("Base = %T %v, want SymbolNode(x)", n.Base, n.Base)
	}
	sub, ok := n.Sub.(*parser.SymbolNode)
	if !ok || sub.Name != "i" {
		t.Errorf("Sub = %T %v, want SymbolNode(i)", n.Sub, n.Sub)
	}
}

// TestSubscriptNumberIndex — x_{0} → SubscriptNode{SymbolNode(x), NumberNode(0)}
func TestSubscriptNumberIndex(t *testing.T) {
	node := parse(t, "x_{0}")
	n, ok := node.(*parser.SubscriptNode)
	if !ok {
		t.Fatalf("parse(%q) = %T, want *SubscriptNode", "x_{0}", node)
	}
	if _, ok := n.Base.(*parser.SymbolNode); !ok {
		t.Errorf("Base = %T, want *SymbolNode", n.Base)
	}
	num, ok := n.Sub.(*parser.NumberNode)
	if !ok || num.Value != 0 {
		t.Errorf("Sub = %T %v, want NumberNode(0)", n.Sub, n.Sub)
	}
}

// TestSubscriptExpressionIndex — x_{n+1} → SubscriptNode{x, BinaryNode(+, n, 1)}
func TestSubscriptExpressionIndex(t *testing.T) {
	node := parse(t, "x_{n+1}")
	n, ok := node.(*parser.SubscriptNode)
	if !ok {
		t.Fatalf("parse(%q) = %T, want *SubscriptNode", "x_{n+1}", node)
	}
	if _, ok := n.Sub.(*parser.BinaryNode); !ok {
		t.Errorf("Sub = %T, want *BinaryNode", n.Sub)
	}
}

// TestSubscriptThenPower — x_i^2 → BinaryNode(^, SubscriptNode(x,i), NumberNode(2))
// Subscript binds tighter than power: x_i^2 = (x_i)^2.
func TestSubscriptThenPower(t *testing.T) {
	node := parse(t, "x_i^2")
	pow, ok := node.(*parser.BinaryNode)
	if !ok || pow.Op != "^" {
		t.Fatalf("parse(%q) = %T, want BinaryNode(^)", "x_i^2", node)
	}
	if _, ok := pow.Left.(*parser.SubscriptNode); !ok {
		t.Errorf("Left of ^ = %T, want *SubscriptNode", pow.Left)
	}
	num, ok := pow.Right.(*parser.NumberNode)
	if !ok || num.Value != 2 {
		t.Errorf("Right of ^ = %T %v, want NumberNode(2)", pow.Right, pow.Right)
	}
}

// TestSubscriptChainError — x_{i}_{j} → parse error (chained subscripts forbidden)
func TestSubscriptChainError(t *testing.T) {
	mustFailWith(t, "x_{i}_{j}", "chained subscripts")
}

// TestSuperscriptBeforeSubscriptError — x^2_i → parse error (order constraint)
// Standard LaTeX allows this; goexeclatex requires _ before ^ (spec §4.1).
func TestSuperscriptBeforeSubscriptError(t *testing.T) {
	mustFailWith(t, "x^2_i", "superscript before subscript")
}

// TestStrayEquals — x = 1 → parse error (EQUALS only valid in large-op bound)
func TestStrayEquals(t *testing.T) {
	mustFailWith(t, "x = 1", "=")
}

// ── §4.2 Big operator rule ────────────────────────────────────────────────────

// TestSumNode — \sum_{i=0}^{3} i → LargeOpNode{sum, "i", 0, 3, SymbolNode(i)}
func TestSumNode(t *testing.T) {
	node := parse(t, `\sum_{i=0}^{3} i`)
	n, ok := node.(*parser.LargeOpNode)
	if !ok {
		t.Fatalf(`parse(\sum_{i=0}^{3} i) = %T, want *LargeOpNode`, node)
	}
	if n.Op != "sum" {
		t.Errorf("Op = %q, want %q", n.Op, "sum")
	}
	if n.Var != "i" {
		t.Errorf("Var = %q, want %q", n.Var, "i")
	}
	from, ok := n.From.(*parser.NumberNode)
	if !ok || from.Value != 0 {
		t.Errorf("From = %T %v, want NumberNode(0)", n.From, n.From)
	}
	to, ok := n.To.(*parser.NumberNode)
	if !ok || to.Value != 3 {
		t.Errorf("To = %T %v, want NumberNode(3)", n.To, n.To)
	}
	body, ok := n.Body.(*parser.SymbolNode)
	if !ok || body.Name != "i" {
		t.Errorf("Body = %T %v, want SymbolNode(i)", n.Body, n.Body)
	}
}

// TestProdNode — \prod_{i=1}^{4} i → LargeOpNode{prod, "i", 1, 4, SymbolNode(i)}
func TestProdNode(t *testing.T) {
	node := parse(t, `\prod_{i=1}^{4} i`)
	n, ok := node.(*parser.LargeOpNode)
	if !ok {
		t.Fatalf(`parse(\prod_{i=1}^{4} i) = %T, want *LargeOpNode`, node)
	}
	if n.Op != "prod" {
		t.Errorf("Op = %q, want %q", n.Op, "prod")
	}
	if n.Var != "i" {
		t.Errorf("Var = %q, want %q", n.Var, "i")
	}
	from, ok := n.From.(*parser.NumberNode)
	if !ok || from.Value != 1 {
		t.Errorf("From = %T %v, want NumberNode(1)", n.From, n.From)
	}
	to, ok := n.To.(*parser.NumberNode)
	if !ok || to.Value != 4 {
		t.Errorf("To = %T %v, want NumberNode(4)", n.To, n.To)
	}
}

// TestLargeOpBodyIsSingleAtom — body parsed at parsePower level, so
// \sum_{i=1}^{3} i + 1 → BinaryNode(+, LargeOpNode, 1) not LargeOpNode(body=i+1).
func TestLargeOpBodyIsSingleAtom(t *testing.T) {
	node := parse(t, `\sum_{i=1}^{3} i + 1`)
	add, ok := node.(*parser.BinaryNode)
	if !ok || add.Op != "+" {
		t.Fatalf(`parse(\sum_{i=1}^{3} i + 1) = %T, want outer BinaryNode(+)`, node)
	}
	if _, ok := add.Left.(*parser.LargeOpNode); !ok {
		t.Errorf("Left = %T, want *LargeOpNode", add.Left)
	}
	if _, ok := add.Right.(*parser.NumberNode); !ok {
		t.Errorf("Right = %T, want *NumberNode", add.Right)
	}
}

// TestLargeOpBracedBody — \sum_{i=1}^{3} {i+1} → LargeOpNode with body BinaryNode(+, i, 1)
func TestLargeOpBracedBody(t *testing.T) {
	node := parse(t, `\sum_{i=1}^{3} {i+1}`)
	n, ok := node.(*parser.LargeOpNode)
	if !ok {
		t.Fatalf(`parse(\sum_{i=1}^{3} {i+1}) = %T, want *LargeOpNode`, node)
	}
	if _, ok := n.Body.(*parser.BinaryNode); !ok {
		t.Errorf("Body = %T, want *BinaryNode (i+1)", n.Body)
	}
}

// TestLargeOpMissingSubscript — \sum^{3} → parse error (lower bound required first)
func TestLargeOpMissingSubscript(t *testing.T) {
	mustFailWith(t, `\sum^{3} i`, `_{`)
}

// TestLargeOpMissingBoth — \sum i → parse error (both bounds required)
func TestLargeOpMissingBoth(t *testing.T) {
	mustFailWith(t, `\sum i`, `_{`)
}

// TestLargeOpMissingSuperscript — \sum_{i=0} i → parse error (upper bound required)
func TestLargeOpMissingSuperscript(t *testing.T) {
	mustFailWith(t, `\sum_{i=0} i`, `^`)
}

// TestLargeOpVarNotSymbol — \sum_{1=0}^{3} 1 → parse error (variable must be symbol)
func TestLargeOpVarNotSymbol(t *testing.T) {
	mustFailWith(t, `\sum_{1=0}^{3} 1`, "iteration variable must be a symbol")
}

// TestLargeOpVarShadowedInBody — nested iteration variable correctly binds
// \sum_{i=1}^{2} \sum_{j=1}^{2} {i*j}: outer i=1..2, inner j=1..2
func TestLargeOpNestedSumNode(t *testing.T) {
	node := parse(t, `\sum_{i=1}^{2} \sum_{j=1}^{2} {i*j}`)
	outer, ok := node.(*parser.LargeOpNode)
	if !ok || outer.Op != "sum" {
		t.Fatalf("outer = %T, want *LargeOpNode(sum)", node)
	}
	if _, ok := outer.Body.(*parser.LargeOpNode); !ok {
		t.Errorf("outer.Body = %T, want *LargeOpNode", outer.Body)
	}
}

// ── §4.3 Norm rule ────────────────────────────────────────────────────────────

// TestNormNode — \lVert x \rVert → NormNode{SymbolNode(x)}
func TestNormNode(t *testing.T) {
	node := parse(t, `\lVert x \rVert`)
	n, ok := node.(*parser.NormNode)
	if !ok {
		t.Fatalf(`parse(\lVert x \rVert) = %T, want *NormNode`, node)
	}
	sym, ok := n.Arg.(*parser.SymbolNode)
	if !ok || sym.Name != "x" {
		t.Errorf("Arg = %T %v, want SymbolNode(x)", n.Arg, n.Arg)
	}
}

// TestNormNested — \lVert x + y \rVert → NormNode with BinaryNode body
func TestNormExpression(t *testing.T) {
	node := parse(t, `\lVert x + y \rVert`)
	n, ok := node.(*parser.NormNode)
	if !ok {
		t.Fatalf(`parse(\lVert x+y \rVert) = %T, want *NormNode`, node)
	}
	if _, ok := n.Arg.(*parser.BinaryNode); !ok {
		t.Errorf("Arg = %T, want *BinaryNode", n.Arg)
	}
}

// TestNormEmptyError — \lVert\rVert → parse error (empty body)
func TestNormEmptyError(t *testing.T) {
	mustFailWith(t, `\lVert\rVert`, "empty")
}

// TestNormUnmatchedError — \lVert x → parse error (unmatched \lVert)
func TestNormUnmatchedError(t *testing.T) {
	mustFailWith(t, `\lVert x`, `\lVert`)
}
