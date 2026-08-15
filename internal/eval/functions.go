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
	case "frac", "dfrac", "tfrac", "cfrac":
		if args[1] == 0 {
			return 0, fmt.Errorf("eval: division by zero")
		}
		return args[0] / args[1], nil

	case "sqrt":
		if len(args) == 2 {
			x, n := args[0], args[1]
			if n == 0 {
				return 0, fmt.Errorf("eval: domain error: root index must not be zero")
			}
			if x < 0 {
				return 0, fmt.Errorf("eval: domain error: root of negative")
			}
			return math.Pow(x, 1/n), nil
		}
		if args[0] < 0 {
			return 0, fmt.Errorf("eval: domain error: sqrt of negative")
		}
		return math.Sqrt(args[0]), nil

	case "abs":
		return math.Abs(args[0]), nil

	case "floor":
		return math.Floor(args[0]), nil

	case "ceil":
		return math.Ceil(args[0]), nil

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

	// §6.3 Inverse trigonometric (aX abbreviations + canonical arcX)
	case "asin", "arcsin":
		if math.Abs(args[0]) > 1 {
			return 0, fmt.Errorf("eval: domain error: %s argument out of range", name)
		}
		return math.Asin(args[0]), nil

	case "acos", "arccos":
		if math.Abs(args[0]) > 1 {
			return 0, fmt.Errorf("eval: domain error: %s argument out of range", name)
		}
		return math.Acos(args[0]), nil

	case "atan", "arctan":
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

	// §6.4 Logarithmic and exponential
	case "ln":
		if args[0] <= 0 {
			return 0, fmt.Errorf("eval: domain error: ln argument out of range")
		}
		return math.Log(args[0]), nil

	case "log":
		if len(args) == 2 {
			x, b := args[0], args[1]
			if x <= 0 {
				return 0, fmt.Errorf("eval: domain error: log argument out of range")
			}
			if b <= 0 || b == 1 {
				return 0, fmt.Errorf("eval: domain error: log base out of range")
			}
			return math.Log(x) / math.Log(b), nil
		}
		if args[0] <= 0 {
			return 0, fmt.Errorf("eval: domain error: log argument out of range")
		}
		return math.Log10(args[0]), nil

	case "exp":
		return math.Exp(args[0]), nil

	// §6.5 Hyperbolic
	case "sinh":
		return math.Sinh(args[0]), nil
	case "cosh":
		return math.Cosh(args[0]), nil
	case "tanh":
		return math.Tanh(args[0]), nil
	case "coth":
		return 1 / math.Tanh(args[0]), nil
	case "sech":
		return 1 / math.Cosh(args[0]), nil
	case "csch":
		return 1 / math.Sinh(args[0]), nil

	// §6.6 Combinatorics
	case "binom", "dbinom", "tbinom":
		return binomCoeff(args[0], args[1])

	case "min":
		return minArgs(args), nil
	case "max":
		return maxArgs(args), nil
	case "gcd":
		return gcdArgs(args)

	default:
		return 0, fmt.Errorf("eval: unknown function: %s", name)
	}
}

func minArgs(args []float64) float64 {
	m := args[0]
	for _, v := range args[1:] {
		if v < m {
			m = v
		}
	}
	return m
}

func maxArgs(args []float64) float64 {
	m := args[0]
	for _, v := range args[1:] {
		if v > m {
			m = v
		}
	}
	return m
}

func gcdArgs(args []float64) (float64, error) {
	vals := make([]int64, len(args))
	for i, v := range args {
		if v < 0 || v != math.Trunc(v) {
			return 0, fmt.Errorf("eval: domain error: gcd requires non-negative integers")
		}
		vals[i] = int64(v)
	}
	g := vals[0]
	for _, v := range vals[1:] {
		g = gcd2(g, v)
	}
	return float64(g), nil
}

func gcd2(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

// binomCoeff computes C(n, k) = n! / (k! * (n-k)!) (spec §6.6.1).
func binomCoeff(n, k float64) (float64, error) {
	if n < 0 || n != math.Trunc(n) || k < 0 || k != math.Trunc(k) {
		return 0, fmt.Errorf("eval: domain error: binom requires non-negative integer arguments")
	}
	if k > n {
		return 0, fmt.Errorf("eval: domain error: binom requires k \u2264 n")
	}
	// Multiplicative formula: avoids computing large factorials.
	result := 1.0
	ki := k
	if n-k < ki {
		ki = n - k // use the smaller of k and n-k
	}
	for i := 0.0; i < ki; i++ {
		result = result * (n - i) / (i + 1)
	}
	return math.Round(result), nil
}
