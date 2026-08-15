package eval_test

import (
	"math"
	"strings"
	"testing"

	"github.com/stephen-mcelhose/goexeclatex/internal/eval"
	"github.com/stephen-mcelhose/goexeclatex/internal/lexer"
	"github.com/stephen-mcelhose/goexeclatex/internal/parser"
)

// evaluate is a test helper: lex → parse → eval with a fresh scope.
func evaluate(t *testing.T, input string) float64 {
	t.Helper()
	tokens, err := lexer.Lex(input)
	if err != nil {
		t.Fatalf("Lex(%q) error: %v", input, err)
	}
	node, err := parser.Parse(tokens)
	if err != nil {
		t.Fatalf("Parse(%q) error: %v", input, err)
	}
	result, err := eval.Eval(node, eval.NewScope())
	if err != nil {
		t.Fatalf("Eval(%q) unexpected error: %v", input, err)
	}
	return result
}

// evaluateWithScope runs evaluate with a custom scope (for variable tests).
func evaluateWithScope(t *testing.T, input string, scope eval.Scope) float64 {
	t.Helper()
	tokens, err := lexer.Lex(input)
	if err != nil {
		t.Fatalf("Lex(%q) error: %v", input, err)
	}
	node, err := parser.Parse(tokens)
	if err != nil {
		t.Fatalf("Parse(%q) error: %v", input, err)
	}
	result, err := eval.Eval(node, scope)
	if err != nil {
		t.Fatalf("Eval(%q) unexpected error: %v", input, err)
	}
	return result
}

// mustError expects a non-nil eval error matching prefix.
func mustError(t *testing.T, input string, wantPrefix string) {
	t.Helper()
	tokens, err := lexer.Lex(input)
	if err != nil {
		return // lex error counts — pipeline fails before eval
	}
	node, err := parser.Parse(tokens)
	if err != nil {
		return // parse error counts
	}
	_, err = eval.Eval(node, eval.NewScope())
	if err == nil {
		t.Errorf("Eval(%q) expected error with prefix %q, got nil", input, wantPrefix)
		return
	}
	if !strings.HasPrefix(err.Error(), wantPrefix) {
		t.Errorf("Eval(%q) error = %q, want prefix %q", input, err.Error(), wantPrefix)
	}
}

// approx checks float64 equality within 1e-10.
func approx(a, b float64) bool {
	if math.IsNaN(a) && math.IsNaN(b) {
		return true
	}
	if math.IsInf(a, 0) || math.IsInf(b, 0) {
		return a == b
	}
	return math.Abs(a-b) < 1e-10
}

// ── §4.1 NumberNode ───────────────────────────────────────────────────────────

func TestNumberLiteral(t *testing.T) {
	cases := []struct {
		input string
		want  float64
	}{
		{"0", 0},
		{"1", 1},
		{"3.14", 3.14},
		{"42", 42},
		{"1.5", 1.5},
	}
	for _, c := range cases {
		got := evaluate(t, c.input)
		if !approx(got, c.want) {
			t.Errorf("Eval(%q) = %v, want %v", c.input, got, c.want)
		}
	}
}

// ── §4.2 SymbolNode ───────────────────────────────────────────────────────────

func TestUndefinedSymbolError(t *testing.T) {
	// Error message must contain the symbol name (spec §7).
	tokens, _ := lexer.Lex("x")
	node, _ := parser.Parse(tokens)
	_, err := eval.Eval(node, eval.NewScope())
	if err == nil {
		t.Fatal("Eval(\"x\") expected error, got nil")
	}
	if !strings.Contains(err.Error(), "x") {
		t.Errorf("Eval(\"x\") error = %q, want it to contain symbol name \"x\"", err.Error())
	}
	if !strings.HasPrefix(err.Error(), "eval: undefined symbol:") {
		t.Errorf("Eval(\"x\") error = %q, want prefix \"eval: undefined symbol:\"", err.Error())
	}
}

func TestUserDefinedVariable(t *testing.T) {
	scope := eval.NewScope()
	scope["x"] = 3.0
	got := evaluateWithScope(t, "x", scope)
	if !approx(got, 3.0) {
		t.Errorf("Eval(\"x\") with x=3 = %v, want 3", got)
	}
}

// ── §4.3 BinaryNode ───────────────────────────────────────────────────────────

func TestBinaryArithmetic(t *testing.T) {
	cases := []struct {
		input string
		want  float64
	}{
		{"1+2", 3},
		{"5-3", 2},
		{"3*4", 12},
		{"10/4", 2.5},
		{"2^{10}", 1024}, // braces required — char mode after ^ reads one digit only
		{"1+2*3", 7},    // precedence
		{"(1+2)*3", 9},  // grouping
		{"2^3^2", 512},  // right-assoc: 2^(3^2) = 2^9
	}
	for _, c := range cases {
		got := evaluate(t, c.input)
		if !approx(got, c.want) {
			t.Errorf("Eval(%q) = %v, want %v", c.input, got, c.want)
		}
	}
}

func TestDivisionByZeroError(t *testing.T) {
	mustError(t, "1/0", "eval: division by zero")
}

// ── §4.4 UnaryNode ────────────────────────────────────────────────────────────

func TestUnaryNegation(t *testing.T) {
	cases := []struct {
		input string
		want  float64
	}{
		{"-3", -3},
		{"-0", 0},
		{"--3", 3},
		{"-3^2", -9}, // -(3^2) not (-3)^2
	}
	for _, c := range cases {
		got := evaluate(t, c.input)
		if !approx(got, c.want) {
			t.Errorf("Eval(%q) = %v, want %v", c.input, got, c.want)
		}
	}
}

func TestFactorial(t *testing.T) {
	cases := []struct {
		input string
		want  float64
	}{
		{"0!", 1},
		{"1!", 1},
		{"5!", 120},
		{"6!", 720},
	}
	for _, c := range cases {
		got := evaluate(t, c.input)
		if !approx(got, c.want) {
			t.Errorf("Eval(%q) = %v, want %v", c.input, got, c.want)
		}
	}
}

func TestFactorialNonIntegerError(t *testing.T) {
	mustError(t, "1.5!", "eval: factorial requires a non-negative integer")
}

func TestFactorialNegativeError(t *testing.T) {
	// "-1!" parses as -(1!) = -1 (no error). Need grouping to pass -1 to factorial.
	mustError(t, "(-1)!", "eval: factorial requires a non-negative integer")
}

// ── §5 Built-in constants ─────────────────────────────────────────────────────

func TestBuiltinConstants(t *testing.T) {
	cases := []struct {
		input string
		want  float64
	}{
		{`\pi`, math.Pi},
		{`e`, math.E},
		{`\tau`, 2 * math.Pi},
		{`\phi`, (1 + math.Sqrt(5)) / 2},
	}
	for _, c := range cases {
		got := evaluate(t, c.input)
		if !approx(got, c.want) {
			t.Errorf("Eval(%q) = %.15f, want %.15f", c.input, got, c.want)
		}
	}
}

// ── §6.1 Arithmetic functions ─────────────────────────────────────────────────

func TestFrac(t *testing.T) {
	cases := []struct {
		input string
		want  float64
	}{
		{`\frac{1}{2}`, 0.5},
		{`\frac{22}{7}`, 22.0 / 7.0},
		{`\frac{3}{4}+\frac{1}{4}`, 1.0},
	}
	for _, c := range cases {
		got := evaluate(t, c.input)
		if !approx(got, c.want) {
			t.Errorf("Eval(%q) = %v, want %v", c.input, got, c.want)
		}
	}
}

func TestFracDivisionByZeroError(t *testing.T) {
	mustError(t, `\frac{1}{0}`, "eval: division by zero")
}

func TestSqrt(t *testing.T) {
	cases := []struct {
		input string
		want  float64
	}{
		{`\sqrt{4}`, 2},
		{`\sqrt{9}`, 3},
		{`\sqrt{2}`, math.Sqrt(2)},
		{`\sqrt{0}`, 0},
	}
	for _, c := range cases {
		got := evaluate(t, c.input)
		if !approx(got, c.want) {
			t.Errorf("Eval(%q) = %v, want %v", c.input, got, c.want)
		}
	}
}

func TestSqrtNegativeError(t *testing.T) {
	mustError(t, `\sqrt{-1}`, "eval: domain error: sqrt of negative")
}

func TestAbs(t *testing.T) {
	cases := []struct {
		input string
		want  float64
	}{
		{"|3|", 3},
		{"|-3|", 3},
		{"|0|", 0},
		{`\lvert -5 \rvert`, 5},
	}
	for _, c := range cases {
		got := evaluate(t, c.input)
		if !approx(got, c.want) {
			t.Errorf("Eval(%q) = %v, want %v", c.input, got, c.want)
		}
	}
}

// ── §6.2 Trig functions ───────────────────────────────────────────────────────

func TestTrigBasic(t *testing.T) {
	cases := []struct {
		input string
		want  float64
	}{
		{`\sin{0}`, 0},
		{`\cos{0}`, 1},
		{`\tan{0}`, 0},
		{`\sin{\pi}`, 0},           // sin(π) ≈ 0
		{`\cos{\pi}`, -1},          // cos(π) = -1
	}
	for _, c := range cases {
		got := evaluate(t, c.input)
		if !approx(got, c.want) {
			t.Errorf("Eval(%q) = %.15f, want %.15f", c.input, got, c.want)
		}
	}
}

func TestTrigReciprocal(t *testing.T) {
	cases := []struct {
		input string
		want  float64
	}{
		{`\sec{0}`, 1},                    // 1/cos(0)=1
		{`\csc{\pi/2}`, 1},               // 1/sin(π/2)=1
		{`\cot{\pi/4}`, 1},               // 1/tan(π/4)=1
	}
	for _, c := range cases {
		got := evaluate(t, c.input)
		if !approx(got, c.want) {
			t.Errorf("Eval(%q) = %.15f, want %.15f", c.input, got, c.want)
		}
	}
}

// ── §6.3 Inverse trig functions ───────────────────────────────────────────────

func TestInverseTrig(t *testing.T) {
	cases := []struct {
		input string
		want  float64
	}{
		{`\asin{0}`, 0},
		{`\asin{1}`, math.Pi / 2},
		{`\acos{1}`, 0},
		{`\acos{0}`, math.Pi / 2},
		{`\atan{0}`, 0},
		{`\atan{1}`, math.Pi / 4},
		{`\acot{1}`, math.Pi / 4},
	}
	for _, c := range cases {
		got := evaluate(t, c.input)
		if !approx(got, c.want) {
			t.Errorf("Eval(%q) = %.15f, want %.15f", c.input, got, c.want)
		}
	}
}

func TestInverseTrigHappyPath(t *testing.T) {
	// asec and acsc have only domain-error tests above — add happy-path cases (spec §6.3).
	cases := []struct {
		input string
		want  float64
	}{
		{`\asec{2}`, math.Acos(0.5)},   // asec(2) = acos(1/2) = π/3
		{`\acsc{2}`, math.Asin(0.5)},   // acsc(2) = asin(1/2) = π/6
		{`\asec{-1}`, math.Pi},         // asec(-1) = acos(-1) = π
	}
	for _, c := range cases {
		got := evaluate(t, c.input)
		if !approx(got, c.want) {
			t.Errorf("Eval(%q) = %.15f, want %.15f", c.input, got, c.want)
		}
	}
}

func TestInverseTrigDomainErrors(t *testing.T) {
	mustError(t, `\asin{2}`, "eval: domain error:")
	mustError(t, `\acos{-2}`, "eval: domain error:")
	mustError(t, `\asec{0}`, "eval: domain error:")
	mustError(t, `\acsc{0}`, "eval: domain error:")
}

// ── §6.2 Trig pole behaviour ──────────────────────────────────────────────────

func TestTrigNonTrivial(t *testing.T) {
	// tan at a non-trivial value (spec §6.2 — not just zero).
	got := evaluate(t, `\tan{\pi/4}`)
	if !approx(got, 1.0) {
		t.Errorf(`Eval(\tan{\pi/4}) = %.15f, want 1`, got)
	}
}

func TestTrigPolesNoError(t *testing.T) {
	// At poles, sec/csc/cot/tan produce ±Inf but MUST NOT return an error (spec §6.2).
	// Use a value very close to π/2 where cos≈0, making sec very large.
	// We verify no error is returned, not the exact value.
	poleInputs := []string{
		`\sec{\pi/2}`,
		`\csc{\pi}`,
	}
	for _, input := range poleInputs {
		tokens, err := lexer.Lex(input)
		if err != nil {
			t.Fatalf("Lex(%q) error: %v", input, err)
		}
		node, err := parser.Parse(tokens)
		if err != nil {
			t.Fatalf("Parse(%q) error: %v", input, err)
		}
		_, err = eval.Eval(node, eval.NewScope())
		if err != nil {
			t.Errorf("Eval(%q) at trig pole returned error %v, want nil (spec §6.2)", input, err)
		}
	}
}

// ── §4.5 Unknown function name ────────────────────────────────────────────────

func TestUnknownFunctionError(t *testing.T) {
	// Directly construct a FunctionNode with an unknown name to test the
	// defensive default in callBuiltin (spec §4.5). Unreachable via the
	// normal lex→parse pipeline (unknown commands become SymbolNodes), but
	// the spec says Eval MUST return this error if Name is not in the table.
	node := &parser.FunctionNode{Name: "unknownfn", Args: []parser.Node{&parser.NumberNode{Value: 1}}}
	_, err := eval.Eval(node, eval.NewScope())
	if err == nil {
		t.Fatal("Eval(FunctionNode{unknownfn}) expected error, got nil")
	}
	if !strings.HasPrefix(err.Error(), "eval: unknown function:") {
		t.Errorf("error = %q, want prefix \"eval: unknown function:\"", err.Error())
	}
}

// ── §7 Error message prefix ───────────────────────────────────────────────────

func TestErrorPrefix(t *testing.T) {
	// All eval errors must begin with "eval: "
	mustError(t, "x", "eval: ")
	mustError(t, "1/0", "eval: ")
	mustError(t, `\sqrt{-1}`, "eval: ")
	mustError(t, "1.5!", "eval: ")
}

// ── Integration: end-to-end expressions ──────────────────────────────────────

func TestIntegration(t *testing.T) {
	cases := []struct {
		input string
		want  float64
	}{
		{`\frac{1}{2} + \sqrt{9}`, 3.5},
		{`2^{10}`, 1024},
		{`\sin{\pi/6}`, 0.5},
		{`2\pi`, 2 * math.Pi},
		{`\frac{\sqrt{2}}{2}`, math.Sqrt(2) / 2},
		{`5!`, 120},
		{`|-3| + \sqrt{4}`, 5},
	}
	for _, c := range cases {
		got := evaluate(t, c.input)
		if !approx(got, c.want) {
			t.Errorf("Eval(%q) = %.15f, want %.15f", c.input, got, c.want)
		}
	}
}

func TestImplicitMultiplyEval(t *testing.T) {
	// 2x with x=3 → 6
	scope := eval.NewScope()
	scope["x"] = 3.0
	got := evaluateWithScope(t, "2x", scope)
	if !approx(got, 6.0) {
		t.Errorf("Eval(\"2x\", x=3) = %v, want 6", got)
	}
}
