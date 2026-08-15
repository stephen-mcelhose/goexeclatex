package eval

import (
	"fmt"
	"math"

	"github.com/stephen-mcelhose/goexeclatex/internal/parser"
)

// Eval reduces node to a float64 using scope for symbol lookup.
// It performs a depth-first walk of the AST (spec §4; subscripts §6).
func Eval(node parser.Node, scope ScopeLookup) (float64, error) {
	switch n := node.(type) {
	case *parser.NumberNode:
		// §4.1: return value directly.
		return n.Value, nil
	case *parser.SymbolNode:
		return evalSymbol(n, scope)
	case *parser.BinaryNode:
		return evalBinary(n, scope)
	case *parser.UnaryNode:
		return evalUnary(n, scope)
	case *parser.FunctionNode:
		return evalFunction(n, scope)
	case *parser.SubscriptNode:
		return evalSubscript(n, scope)
	case *parser.LargeOpNode:
		return evalLargeOp(n, scope)
	case *parser.NormNode:
		return evalNorm(n, scope)
	default:
		return 0, fmt.Errorf("eval: unknown node type: %T", node)
	}
}

// evalSymbol looks up a symbol in scope (spec §4.2).
func evalSymbol(n *parser.SymbolNode, scope ScopeLookup) (float64, error) {
	v, ok := scope.Lookup(n.Name)
	if !ok {
		return 0, fmt.Errorf("eval: undefined symbol: %s", n.Name)
	}
	return v, nil
}

// evalBinary evaluates a binary infix node (spec §4.3).
func evalBinary(n *parser.BinaryNode, scope ScopeLookup) (float64, error) {
	left, err := Eval(n.Left, scope)
	if err != nil {
		return 0, err
	}
	right, err := Eval(n.Right, scope)
	if err != nil {
		return 0, err
	}
	return applyBinary(n.Op, left, right)
}

// applyBinary applies a binary operator to two already-evaluated operands (spec §4.3).
func applyBinary(op string, left, right float64) (float64, error) {
	switch op {
	case "+":
		return left + right, nil
	case "-":
		return left - right, nil
	case "*":
		return left * right, nil
	case "/":
		if right == 0 {
			return 0, fmt.Errorf("eval: division by zero")
		}
		return left / right, nil
	case "^":
		return math.Pow(left, right), nil
	case "bmod":
		if right == 0 {
			return 0, fmt.Errorf("eval: division by zero")
		}
		return math.Mod(left, right), nil
	default:
		return 0, fmt.Errorf("eval: unknown operator: %s", op)
	}
}

// evalUnary evaluates a prefix/postfix unary node (spec §4.4).
func evalUnary(n *parser.UnaryNode, scope ScopeLookup) (float64, error) {
	operand, err := Eval(n.Operand, scope)
	if err != nil {
		return 0, err
	}
	switch n.Op {
	case "-":
		return -operand, nil
	case "!":
		return factorial(operand)
	default:
		return 0, fmt.Errorf("eval: unknown unary operator: %s", n.Op)
	}
}

// evalFunction evaluates args then dispatches to the built-in table (spec §4.5 → §6).
func evalFunction(n *parser.FunctionNode, scope ScopeLookup) (float64, error) {
	args := make([]float64, len(n.Args))
	for i, arg := range n.Args {
		v, err := Eval(arg, scope)
		if err != nil {
			return 0, err
		}
		args[i] = v
	}
	return callBuiltin(n.Name, args)
}

// evalSubscript resolves x_{k} → scope key "x_<int(k)>" (spec subscripts-largeops §6.1).
func evalSubscript(n *parser.SubscriptNode, scope ScopeLookup) (float64, error) {
	base, ok := n.Base.(*parser.SymbolNode)
	if !ok {
		return 0, fmt.Errorf("eval: subscript base must be a symbol, got %T", n.Base)
	}
	idx, err := Eval(n.Sub, scope)
	if err != nil {
		return 0, err
	}
	key, err := subscriptKey(base.Name, idx)
	if err != nil {
		return 0, err
	}
	v, found := scope.Lookup(key)
	if !found {
		return 0, fmt.Errorf("eval: undefined symbol: %s", key)
	}
	return v, nil
}

// subscriptKey builds the composite scope key for a subscript index (spec §6.1.3–§6.1.5).
func subscriptKey(baseName string, idx float64) (string, error) {
	rounded := math.Round(idx)
	if math.IsNaN(idx) || math.IsInf(idx, 0) || math.Abs(idx-rounded) > 1e-9 || rounded < 0 {
		return "", fmt.Errorf("eval: subscript index must be a non-negative integer, got %v", idx)
	}
	return fmt.Sprintf("%s_%d", baseName, int(rounded)), nil
}

// evalLargeOp evaluates \sum or \prod over a discrete integer range (spec §6.2).
func evalLargeOp(n *parser.LargeOpNode, scope ScopeLookup) (float64, error) {
	from, to, err := largeOpBounds(n, scope)
	if err != nil {
		return 0, err
	}
	switch n.Op {
	case "sum":
		return foldRange(from, to, 0, scope, n.Var, n.Body, func(acc, v float64) float64 { return acc + v })
	case "prod":
		return foldRange(from, to, 1, scope, n.Var, n.Body, func(acc, v float64) float64 { return acc * v })
	default:
		return 0, fmt.Errorf("eval: unknown large operator: %s", n.Op)
	}
}

// largeOpBounds evaluates and validates From/To as finite integers (spec §6.2.1–§6.2.2).
func largeOpBounds(n *parser.LargeOpNode, scope ScopeLookup) (from, to int, err error) {
	fromVal, err := Eval(n.From, scope)
	if err != nil {
		return 0, 0, err
	}
	toVal, err := Eval(n.To, scope)
	if err != nil {
		return 0, 0, err
	}
	// Bounds must be finite integers (no epsilon tolerance — they are control values).
	if math.IsNaN(fromVal) || math.IsInf(fromVal, 0) || fromVal != math.Trunc(fromVal) {
		return 0, 0, fmt.Errorf("eval: \\%s lower bound must be an integer, got %v", n.Op, fromVal)
	}
	if math.IsNaN(toVal) || math.IsInf(toVal, 0) || toVal != math.Trunc(toVal) {
		return 0, 0, fmt.Errorf("eval: \\%s upper bound must be a finite integer, got %v", n.Op, toVal)
	}
	return int(fromVal), int(toVal), nil
}

// foldRange evaluates body once per integer i in [from, to], combining with op.
func foldRange(
	from, to int,
	acc float64,
	scope ScopeLookup,
	name string,
	body parser.Node,
	op func(acc, v float64) float64,
) (float64, error) {
	for i := from; i <= to; i++ {
		inner := &innerScope{parent: scope, name: name, value: float64(i)}
		v, err := Eval(body, inner)
		if err != nil {
			return 0, err
		}
		acc = op(acc, v)
	}
	return acc, nil
}

// evalNorm evaluates \lVert expr \rVert as scalar absolute value (spec §6.3).
func evalNorm(n *parser.NormNode, scope ScopeLookup) (float64, error) {
	v, err := Eval(n.Arg, scope)
	if err != nil {
		return 0, err
	}
	return math.Abs(v), nil
}

// factorial implements n! for non-negative integers (spec §4.4.1).
func factorial(n float64) (float64, error) {
	if n < 0 || n != math.Trunc(n) {
		return 0, fmt.Errorf("eval: factorial requires a non-negative integer")
	}
	result := 1.0
	for i := 2.0; i <= n; i++ {
		result *= i
	}
	return result, nil
}
