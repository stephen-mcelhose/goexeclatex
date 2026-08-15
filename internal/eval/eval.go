package eval

import (
	"fmt"
	"math"

	"github.com/stephen-mcelhose/goexeclatex/internal/parser"
)

// Eval reduces node to a float64 using scope for symbol lookup.
// It performs a depth-first walk of the AST (spec §4).
func Eval(node parser.Node, scope ScopeLookup) (float64, error) {
	switch n := node.(type) {

	case *parser.NumberNode:
		// §4.1: return value directly.
		return n.Value, nil

	case *parser.SymbolNode:
		// §4.2: look up in scope; error if absent.
		v, ok := scope.Lookup(n.Name)
		if !ok {
			return 0, fmt.Errorf("eval: undefined symbol: %s", n.Name)
		}
		return v, nil

	case *parser.BinaryNode:
		// §4.3: evaluate both sides, then apply operator.
		left, err := Eval(n.Left, scope)
		if err != nil {
			return 0, err
		}
		right, err := Eval(n.Right, scope)
		if err != nil {
			return 0, err
		}
		switch n.Op {
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
		default:
			return 0, fmt.Errorf("eval: unknown operator: %s", n.Op)
		}

	case *parser.UnaryNode:
		// §4.4
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

	case *parser.FunctionNode:
		// §4.5: evaluate args then dispatch to built-in table.
		args := make([]float64, len(n.Args))
		for i, arg := range n.Args {
			v, err := Eval(arg, scope)
			if err != nil {
				return 0, err
			}
			args[i] = v
		}
		return callBuiltin(n.Name, args)

	case *parser.SubscriptNode:
		// §6.1: resolve x_{k} → look up "x_<int(k)>" in scope.
		base, ok := n.Base.(*parser.SymbolNode)
		if !ok {
			return 0, fmt.Errorf("eval: subscript base must be a symbol, got %T", n.Base)
		}
		idx, err := Eval(n.Sub, scope)
		if err != nil {
			return 0, err
		}
		rounded := math.Round(idx)
		if math.IsNaN(idx) || math.IsInf(idx, 0) || math.Abs(idx-rounded) > 1e-9 || rounded < 0 {
			return 0, fmt.Errorf("eval: subscript index must be a non-negative integer, got %v", idx)
		}
		key := fmt.Sprintf("%s_%d", base.Name, int(rounded))
		v, found := scope.Lookup(key)
		if !found {
			return 0, fmt.Errorf("eval: undefined symbol: %s", key)
		}
		return v, nil

	case *parser.LargeOpNode:
		// §6.2: evaluate \sum or \prod over a discrete integer range.
		fromVal, err := Eval(n.From, scope)
		if err != nil {
			return 0, err
		}
		toVal, err := Eval(n.To, scope)
		if err != nil {
			return 0, err
		}
		// Bounds must be finite integers (no epsilon tolerance — they are control values).
		if math.IsNaN(fromVal) || math.IsInf(fromVal, 0) || fromVal != math.Trunc(fromVal) {
			return 0, fmt.Errorf("eval: \\%s lower bound must be an integer, got %v", n.Op, fromVal)
		}
		if math.IsNaN(toVal) || math.IsInf(toVal, 0) || toVal != math.Trunc(toVal) {
			return 0, fmt.Errorf("eval: \\%s upper bound must be a finite integer, got %v", n.Op, toVal)
		}
		from, to := int(fromVal), int(toVal)

		switch n.Op {
		case "sum":
			acc := 0.0
			for i := from; i <= to; i++ {
				inner := &innerScope{parent: scope, name: n.Var, value: float64(i)}
				v, err := Eval(n.Body, inner)
				if err != nil {
					return 0, err
				}
				acc += v
			}
			return acc, nil
		case "prod":
			acc := 1.0
			for i := from; i <= to; i++ {
				inner := &innerScope{parent: scope, name: n.Var, value: float64(i)}
				v, err := Eval(n.Body, inner)
				if err != nil {
					return 0, err
				}
				acc *= v
			}
			return acc, nil
		default:
			return 0, fmt.Errorf("eval: unknown large operator: %s", n.Op)
		}

	case *parser.NormNode:
		// §6.3: scalar absolute value (Euclidean norm for scalars).
		v, err := Eval(n.Arg, scope)
		if err != nil {
			return 0, err
		}
		return math.Abs(v), nil

	default:
		return 0, fmt.Errorf("eval: unknown node type: %T", node)
	}
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
