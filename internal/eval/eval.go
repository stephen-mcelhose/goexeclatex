package eval

import (
	"fmt"
	"math"

	"github.com/stephen-mcelhose/goexeclatex/internal/parser"
)

// Eval reduces node to a float64 using scope for symbol lookup.
// It performs a depth-first walk of the AST (spec §4).
func Eval(node parser.Node, scope Scope) (float64, error) {
	switch n := node.(type) {

	case *parser.NumberNode:
		// §4.1: return value directly.
		return n.Value, nil

	case *parser.SymbolNode:
		// §4.2: look up in scope; error if absent.
		v, ok := scope[n.Name]
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
