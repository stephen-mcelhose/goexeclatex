package eval

import (
	"fmt"
	"math"
)

// callBuiltin dispatches a FunctionNode name to its implementation.
// args are the already-evaluated argument values.
// Returns (result, error).
func callBuiltin(name string, args []float64) (float64, error) {
	switch name {
	// §6.1 Arithmetic
	case "frac":
		if args[1] == 0 {
			return 0, fmt.Errorf("eval: division by zero")
		}
		return args[0] / args[1], nil

	case "sqrt":
		if args[0] < 0 {
			return 0, fmt.Errorf("eval: domain error: sqrt of negative")
		}
		return math.Sqrt(args[0]), nil

	case "abs":
		return math.Abs(args[0]), nil

	// §6.2 Trigonometric
	case "sin":
		return math.Sin(args[0]), nil
	case "cos":
		return math.Cos(args[0]), nil
	case "tan":
		return math.Tan(args[0]), nil
	case "sec":
		return 1 / math.Cos(args[0]), nil
	case "csc":
		return 1 / math.Sin(args[0]), nil
	case "cot":
		return 1 / math.Tan(args[0]), nil

	// §6.3 Inverse trigonometric
	case "asin":
		if math.Abs(args[0]) > 1 {
			return 0, fmt.Errorf("eval: domain error: asin argument out of range")
		}
		return math.Asin(args[0]), nil

	case "acos":
		if math.Abs(args[0]) > 1 {
			return 0, fmt.Errorf("eval: domain error: acos argument out of range")
		}
		return math.Acos(args[0]), nil

	case "atan":
		return math.Atan(args[0]), nil

	case "asec":
		if math.Abs(args[0]) < 1 {
			return 0, fmt.Errorf("eval: domain error: asec argument out of range")
		}
		return math.Acos(1 / args[0]), nil

	case "acsc":
		if math.Abs(args[0]) < 1 {
			return 0, fmt.Errorf("eval: domain error: acsc argument out of range")
		}
		return math.Asin(1 / args[0]), nil

	case "acot":
		return math.Atan(1 / args[0]), nil

	default:
		return 0, fmt.Errorf("eval: unknown function: %s", name)
	}
}
