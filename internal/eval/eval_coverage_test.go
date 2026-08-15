package eval_test

import (
	"math"
	"strings"
	"testing"

	"github.com/stephen-mcelhose/goexeclatex/internal/eval"
	"github.com/stephen-mcelhose/goexeclatex/internal/parser"
)

// Direct AST construction covers Eval error-propagation and defensive default
// arms that the lex→parse pipeline cannot reach (unknown ops / nil node).
// Part of #12 — raise Eval statement coverage before the complexity refactor.

func num(v float64) *parser.NumberNode {
	return &parser.NumberNode{Value: v}
}

func sym(name string) *parser.SymbolNode {
	return &parser.SymbolNode{Name: name}
}

func mustEvalErr(t *testing.T, node parser.Node, wantSubstr string) {
	t.Helper()
	_, err := eval.Eval(node, eval.NewScope())
	if err == nil {
		t.Fatalf("Eval(%T) expected error containing %q, got nil", node, wantSubstr)
	}
	if !strings.Contains(err.Error(), wantSubstr) {
		t.Fatalf("Eval(%T) error = %q, want substring %q", node, err.Error(), wantSubstr)
	}
}

// TestEvalBinaryLeftError — §4.3: Left evaluation failure propagates.
func TestEvalBinaryLeftError(t *testing.T) {
	mustEvalErr(t, &parser.BinaryNode{
		Op:    "+",
		Left:  sym("undefined_left"),
		Right: num(1),
	}, "undefined symbol: undefined_left")
}

// TestEvalBinaryRightError — §4.3: Right evaluation failure propagates.
func TestEvalBinaryRightError(t *testing.T) {
	mustEvalErr(t, &parser.BinaryNode{
		Op:    "+",
		Left:  num(1),
		Right: sym("undefined_right"),
	}, "undefined symbol: undefined_right")
}

// TestEvalBinaryUnknownOp — §4.3 default: unknown operator is an error.
func TestEvalBinaryUnknownOp(t *testing.T) {
	mustEvalErr(t, &parser.BinaryNode{
		Op:    "%",
		Left:  num(1),
		Right: num(2),
	}, "unknown operator: %")
}

// TestEvalUnaryOperandError — §4.4: operand evaluation failure propagates.
func TestEvalUnaryOperandError(t *testing.T) {
	mustEvalErr(t, &parser.UnaryNode{
		Op:      "-",
		Operand: sym("undefined_unary"),
	}, "undefined symbol: undefined_unary")
}

// TestEvalUnaryUnknownOp — §4.4 default: unknown unary operator is an error.
func TestEvalUnaryUnknownOp(t *testing.T) {
	mustEvalErr(t, &parser.UnaryNode{
		Op:      "~",
		Operand: num(1),
	}, "unknown unary operator: ~")
}

// TestEvalFunctionArgError — §4.5: argument evaluation failure propagates.
func TestEvalFunctionArgError(t *testing.T) {
	mustEvalErr(t, &parser.FunctionNode{
		Name: "sqrt",
		Args: []parser.Node{sym("undefined_arg")},
	}, "undefined symbol: undefined_arg")
}

// TestEvalSubscriptSubError — §6.1: Sub evaluation failure propagates.
func TestEvalSubscriptSubError(t *testing.T) {
	mustEvalErr(t, &parser.SubscriptNode{
		Base: sym("x"),
		Sub:  sym("undefined_idx"),
	}, "undefined symbol: undefined_idx")
}

// TestEvalSubscriptNaNIndex — §6.1.3: NaN index is rejected.
func TestEvalSubscriptNaNIndex(t *testing.T) {
	mustEvalErr(t, &parser.SubscriptNode{
		Base: sym("x"),
		Sub:  num(math.NaN()),
	}, "non-negative integer")
}

// TestEvalSubscriptInfIndex — §6.1.3: Inf index is rejected.
func TestEvalSubscriptInfIndex(t *testing.T) {
	mustEvalErr(t, &parser.SubscriptNode{
		Base: sym("x"),
		Sub:  num(math.Inf(1)),
	}, "non-negative integer")
}

// TestEvalLargeOpFromError — §6.2: From evaluation failure propagates.
func TestEvalLargeOpFromError(t *testing.T) {
	mustEvalErr(t, &parser.LargeOpNode{
		Op:   "sum",
		Var:  "i",
		From: sym("undefined_from"),
		To:   num(3),
		Body: sym("i"),
	}, "undefined symbol: undefined_from")
}

// TestEvalLargeOpToError — §6.2: To evaluation failure propagates.
func TestEvalLargeOpToError(t *testing.T) {
	mustEvalErr(t, &parser.LargeOpNode{
		Op:   "sum",
		Var:  "i",
		From: num(1),
		To:   sym("undefined_to"),
		Body: sym("i"),
	}, "undefined symbol: undefined_to")
}

// TestEvalLargeOpLowerNaN — §6.2.1: NaN lower bound is an error.
func TestEvalLargeOpLowerNaN(t *testing.T) {
	mustEvalErr(t, &parser.LargeOpNode{
		Op:   "sum",
		Var:  "i",
		From: num(math.NaN()),
		To:   num(3),
		Body: num(1),
	}, "lower bound must be an integer")
}

// TestEvalLargeOpUpperNaN — §6.2.2: NaN upper bound is an error.
func TestEvalLargeOpUpperNaN(t *testing.T) {
	mustEvalErr(t, &parser.LargeOpNode{
		Op:   "sum",
		Var:  "i",
		From: num(1),
		To:   num(math.NaN()),
		Body: num(1),
	}, "upper bound must be a finite integer")
}

// TestEvalLargeOpLowerInf — §6.2.1: Inf lower bound is an error.
func TestEvalLargeOpLowerInf(t *testing.T) {
	mustEvalErr(t, &parser.LargeOpNode{
		Op:   "sum",
		Var:  "i",
		From: num(math.Inf(1)),
		To:   num(3),
		Body: num(1),
	}, "lower bound must be an integer")
}

// TestEvalProdBodyError — §6.2: prod body error propagates (sum body already covered).
func TestEvalProdBodyError(t *testing.T) {
	mustPipelineError(t, `\prod_{i=1}^{3} {1/0}`, "division by zero")
}

// TestEvalLargeOpUnknown — §6.2 default: unknown large operator is an error.
func TestEvalLargeOpUnknown(t *testing.T) {
	mustEvalErr(t, &parser.LargeOpNode{
		Op:   "lim",
		Var:  "i",
		From: num(1),
		To:   num(3),
		Body: num(1),
	}, "unknown large operator: lim")
}

// TestEvalNormArgError — §6.3: Arg evaluation failure propagates.
func TestEvalNormArgError(t *testing.T) {
	mustEvalErr(t, &parser.NormNode{
		Arg: sym("undefined_norm"),
	}, "undefined symbol: undefined_norm")
}

// TestEvalUnknownNodeType — Eval default: nil Node is rejected.
func TestEvalUnknownNodeType(t *testing.T) {
	var node parser.Node // nil interface — no concrete type
	mustEvalErr(t, node, "unknown node type")
}
