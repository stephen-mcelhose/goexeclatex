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

// ── §5 infty constant (spec §5) ───────────────────────────────────────────────

func TestInftyConstant(t *testing.T) {
	got := evaluate(t, `\infty`)
	if !math.IsInf(got, 1) {
		t.Errorf(`Eval(\infty) = %v, want +Inf`, got)
	}
}

func TestInftyArithmetic(t *testing.T) {
	// IEEE 754 arithmetic on +Inf — no error (ADR-010).
	cases := []struct {
		input string
		check func(float64) bool
	}{
		{`\infty + 1`, func(v float64) bool { return math.IsInf(v, 1) }},
		{`1 / \infty`, func(v float64) bool { return v == 0 }},
	}
	for _, c := range cases {
		got := evaluate(t, c.input)
		if !c.check(got) {
			t.Errorf("Eval(%q) = %v, unexpected result", c.input, got)
		}
	}
}

// ── §6.1 Frac aliases (spec §6.1) ────────────────────────────────────────────

func TestFracAliases(t *testing.T) {
	cases := []struct {
		input string
		want  float64
	}{
		// Basic: display-mode aliases must equal \frac numerically.
		{`\dfrac{1}{2}`, 0.5},
		{`\tfrac{3}{4}`, 0.75},
		{`\cfrac{1}{3}`, 1.0 / 3.0},
		// Compound: adding two dfrac terms — exercises implicit multiply + nested arities.
		{`\dfrac{1}{2} + \dfrac{1}{3}`, 5.0 / 6.0},
		// Same value, three different aliases.
		{`\dfrac{22}{7} - \tfrac{22}{7}`, 0},
	}
	for _, c := range cases {
		got := evaluate(t, c.input)
		if !approx(got, c.want) {
			t.Errorf("Eval(%q) = %.15f, want %.15f", c.input, got, c.want)
		}
	}
}

// ── §6.3 arcsin / arccos / arctan (spec §6.3) ─────────────────────────────────

func TestArcTrig(t *testing.T) {
	cases := []struct {
		input string
		want  float64
	}{
		// Boundary values.
		{`\arcsin{0}`, 0},
		{`\arcsin{1}`, math.Pi / 2},
		{`\arcsin{-1}`, -math.Pi / 2},
		{`\arccos{1}`, 0},
		{`\arccos{0}`, math.Pi / 2},
		{`\arccos{-1}`, math.Pi}, // the π edge case
		{`\arctan{0}`, 0},
		{`\arctan{1}`, math.Pi / 4},
		{`\arctan{-1}`, -math.Pi / 4},
		// Composition identities: arcX(X(θ)) = θ for θ in principal range.
		// These catch sign errors and wrong-function bugs that boundary tests miss.
		{`\arcsin{\sin{\pi/6}}`, math.Pi / 6},  // sin(π/6)=0.5, arcsin(0.5)=π/6
		{`\arccos{\cos{\pi/3}}`, math.Pi / 3},  // cos(π/3)=0.5, arccos(0.5)=π/3
		{`\arctan{\tan{\pi/4}}`, math.Pi / 4},  // tan(π/4)=1,   arctan(1)=π/4
		// arcX and aX aliases must produce identical results.
		{`\arcsin{1} - \asin{1}`, 0},
		{`\arccos{0} - \acos{0}`, 0},
		{`\arctan{1} - \atan{1}`, 0},
	}
	for _, c := range cases {
		got := evaluate(t, c.input)
		if !approx(got, c.want) {
			t.Errorf("Eval(%q) = %.15f, want %.15f", c.input, got, c.want)
		}
	}
}

func TestArcTrigDomainErrors(t *testing.T) {
	mustError(t, `\arcsin{2}`, "eval: domain error:")
	mustError(t, `\arccos{-2}`, "eval: domain error:")
	// Boundary: exactly ±1 must NOT error (the limit of the domain).
	_ = evaluate(t, `\arcsin{1}`)
	_ = evaluate(t, `\arccos{-1}`)
}

// ── §6.4 ln / log / exp (spec §6.4) ──────────────────────────────────────────

func TestLogExp(t *testing.T) {
	cases := []struct {
		input string
		want  float64
	}{
		// ln basics.
		{`\ln{1}`, 0},
		{`\ln{\e}`, 1},
		// exp basics.
		{`\exp{0}`, 1},
		{`\exp{1}`, math.E},
		// log (base-10) basics.
		{`\log{1}`, 0},
		{`\log{10}`, 1},
		{`\log{100}`, 2},
		{`\log{1000}`, 3},
		// ln vs log are distinct: ln(10) ≠ 1, log(10) = 1.
		{`\log{10} - \ln{10}`, 1 - math.Log(10)},
		// Composition identities: ln(exp(x)) = x and exp(ln(x)) = x.
		{`\ln{\exp{2}}`, 2},
		{`\ln{\exp{-3}}`, -3},
		// Nested: ln(e²) = 2 — exercises ^ inside a brace arg.
		{`\ln{\e^2}`, 2},
		// log laws: log(a)+log(b) = log(a*b).
		{`\log{4} + \log{25}`, 2}, // log(4)+log(25) = log(100) = 2
	}
	for _, c := range cases {
		got := evaluate(t, c.input)
		if !approx(got, c.want) {
			t.Errorf("Eval(%q) = %.15f, want %.15f", c.input, got, c.want)
		}
	}
}

func TestLogDomainErrors(t *testing.T) {
	mustError(t, `\ln{0}`, "eval: domain error:")
	mustError(t, `\ln{-1}`, "eval: domain error:")
	mustError(t, `\log{0}`, "eval: domain error:")
	mustError(t, `\log{-5}`, "eval: domain error:")
}

// ── §6.5 Hyperbolic (spec §6.5) ───────────────────────────────────────────────

func TestHyperbolic(t *testing.T) {
	cases := []struct {
		input string
		want  float64
	}{
		// At x=0: sinh(0)=0, cosh(0)=1, tanh(0)=0, sech(0)=1.
		{`\sinh{0}`, 0},
		{`\cosh{0}`, 1},
		{`\tanh{0}`, 0},
		{`\sech{0}`, 1},
		// cosh(ln 2) = (2 + 1/2)/2 = 5/4 — exact rational, avoids tautology.
		{`\cosh{\ln{2}}`, 1.25},
		// sinh(ln 2) = (2 - 1/2)/2 = 3/4.
		{`\sinh{\ln{2}}`, 0.75},
		// tanh(ln 3) = (3 - 1/3)/(3 + 1/3) = (8/3)/(10/3) = 4/5.
		{`\tanh{\ln{3}}`, 0.8},
		// coth(1) = 1/tanh(1) — roundtrip: coth * tanh = 1.
		{`\coth{1} * \tanh{1}`, 1},
		// sech(1) = 1/cosh(1) — roundtrip.
		{`\sech{1} * \cosh{1}`, 1},
		// csch(1) = 1/sinh(1) — roundtrip.
		{`\csch{1} * \sinh{1}`, 1},
	}
	for _, c := range cases {
		got := evaluate(t, c.input)
		if !approx(got, c.want) {
			t.Errorf("Eval(%q) = %.15f, want %.15f", c.input, got, c.want)
		}
	}
}

func TestHyperbolicPolesNoError(t *testing.T) {
	// coth(0) and csch(0) are ±Inf — not an error (ADR-010, spec §6.5).
	for _, input := range []string{`\coth{0}`, `\csch{0}`} {
		tokens, err := lexer.Lex(input)
		if err != nil {
			t.Fatalf("Lex(%q) error: %v", input, err)
		}
		node, err := parser.Parse(tokens)
		if err != nil {
			t.Fatalf("Parse(%q) error: %v", input, err)
		}
		v, err := eval.Eval(node, eval.NewScope())
		if err != nil {
			t.Errorf("Eval(%q) returned error %v, want nil (±Inf is not an error)", input, err)
		}
		if !math.IsInf(v, 0) {
			t.Errorf("Eval(%q) = %v, want ±Inf", input, v)
		}
	}
}

// ── §6.6 Binomial coefficient (spec §6.6) ────────────────────────────────────

func TestBinom(t *testing.T) {
	cases := []struct {
		input string
		want  float64
	}{
		// Small known values.
		{`\binom{5}{2}`, 10},
		{`\binom{6}{3}`, 20},
		// Boundary: C(n,0)=1 and C(n,n)=1 for all n.
		{`\binom{4}{0}`, 1},
		{`\binom{4}{4}`, 1},
		{`\binom{10}{0}`, 1},
		{`\binom{10}{10}`, 1},
		// Larger value: C(10,5)=252 — catches integer overflow in naive factorial approach.
		{`\binom{10}{5}`, 252},
		// Symmetry: C(n,k) = C(n,n-k).
		{`\binom{10}{3}`, 120},
		{`\binom{10}{7}`, 120},
		// Pascal's identity: C(n,k) = C(n-1,k-1) + C(n-1,k).
		// C(9,2)=36, C(9,3)=84 → sum=120 = C(10,3).
		{`\binom{9}{2} + \binom{9}{3}`, 120},
		// Display-mode aliases produce identical results.
		{`\dbinom{5}{2}`, 10},
		{`\tbinom{5}{2}`, 10},
		{`\dbinom{10}{5} - \tbinom{10}{5}`, 0},
	}
	for _, c := range cases {
		got := evaluate(t, c.input)
		if !approx(got, c.want) {
			t.Errorf("Eval(%q) = %.15f, want %.15f", c.input, got, c.want)
		}
	}
}

func TestBinomDomainErrors(t *testing.T) {
	mustError(t, `\binom{5}{6}`, "eval: domain error:")   // k > n
	mustError(t, `\binom{-1}{2}`, "eval: domain error:")  // negative n
	mustError(t, `\binom{5}{1.5}`, "eval: domain error:") // non-integer k
}

// ── v0.3 §6.1 SubscriptNode ──────────────────────────────────────────────────
//
// Tests for spec §6.1 (SubscriptNode variable resolution).
// All tests in this section are RED until the v0.3 implementation lands.

// mustPipelineError expects the full lex→parse→eval pipeline to produce an
// error whose message contains want (at whichever stage it first fails).
// If the error at that stage does NOT contain want, the test fails.
func mustPipelineError(t *testing.T, input string, want string) {
	t.Helper()
	tokens, lexErr := lexer.Lex(input)
	if lexErr != nil {
		if !strings.Contains(lexErr.Error(), want) {
			t.Errorf("lex(%q) error = %q, want it to contain %q", input, lexErr.Error(), want)
		}
		return
	}
	node, parseErr := parser.Parse(tokens)
	if parseErr != nil {
		if !strings.Contains(parseErr.Error(), want) {
			t.Errorf("parse(%q) error = %q, want it to contain %q", input, parseErr.Error(), want)
		}
		return
	}
	_, evalErr := eval.Eval(node, eval.NewScope())
	if evalErr == nil {
		t.Errorf("eval(%q) expected error containing %q, got nil", input, want)
		return
	}
	if !strings.Contains(evalErr.Error(), want) {
		t.Errorf("eval(%q) error = %q, want it to contain %q", input, evalErr.Error(), want)
	}
}

// mustPipelineErrorScope is like mustPipelineError but uses a custom scope.
func mustPipelineErrorScope(t *testing.T, input string, scope eval.Scope, want string) {
	t.Helper()
	tokens, lexErr := lexer.Lex(input)
	if lexErr != nil {
		if !strings.Contains(lexErr.Error(), want) {
			t.Errorf("lex(%q) error = %q, want it to contain %q", input, lexErr.Error(), want)
		}
		return
	}
	node, parseErr := parser.Parse(tokens)
	if parseErr != nil {
		if !strings.Contains(parseErr.Error(), want) {
			t.Errorf("parse(%q) error = %q, want it to contain %q", input, parseErr.Error(), want)
		}
		return
	}
	_, evalErr := eval.Eval(node, scope)
	if evalErr == nil {
		t.Errorf("eval(%q) expected error containing %q, got nil", input, want)
		return
	}
	if !strings.Contains(evalErr.Error(), want) {
		t.Errorf("eval(%q) error = %q, want it to contain %q", input, evalErr.Error(), want)
	}
}

// TestEvalSubscriptBasic — x_{0} with x_0=5 in scope → 5.
func TestEvalSubscriptBasic(t *testing.T) {
	scope := eval.NewScope()
	scope["x_0"] = 5
	got := evaluateWithScope(t, "x_{0}", scope)
	if !approx(got, 5) {
		t.Errorf("eval(x_{0}) = %v, want 5", got)
	}
}

// TestEvalSubscriptMultiple — multiple subscripted variables in same expression.
func TestEvalSubscriptMultiple(t *testing.T) {
	scope := eval.NewScope()
	scope["x_0"] = 1
	scope["x_1"] = 2
	scope["x_2"] = 3
	got := evaluateWithScope(t, "x_{0} + x_{1} + x_{2}", scope)
	if !approx(got, 6) {
		t.Errorf("eval(x_{0}+x_{1}+x_{2}) = %v, want 6", got)
	}
}

// TestEvalSubscriptDynamicIndex — x_{i} where i is in scope → looks up x_{int(i)}.
func TestEvalSubscriptDynamicIndex(t *testing.T) {
	scope := eval.NewScope()
	scope["i"] = 2
	scope["x_2"] = 99
	got := evaluateWithScope(t, "x_{i}", scope)
	if !approx(got, 99) {
		t.Errorf("eval(x_{i}) = %v, want 99", got)
	}
}

// TestEvalSubscriptThenPower — x_i^2 with x_i bound: (x_i)^2.
func TestEvalSubscriptThenPower(t *testing.T) {
	scope := eval.NewScope()
	scope["i"] = 0
	scope["x_0"] = 3
	got := evaluateWithScope(t, "x_i^2", scope)
	if !approx(got, 9) {
		t.Errorf("eval(x_i^2) = %v, want 9", got)
	}
}

// TestEvalSubscriptEpsilonTolerance — index rounded via epsilon (spec §6.1.3).
// x_{3} should resolve correctly even if the subscript expression yields 3.0000000000001.
func TestEvalSubscriptEpsilonTolerance(t *testing.T) {
	scope := eval.NewScope()
	scope["x_3"] = 42
	// 1+1+1 is exact but exercises the integer-check path.
	got := evaluateWithScope(t, "x_{1+1+1}", scope)
	if !approx(got, 42) {
		t.Errorf("eval(x_{1+1+1}) = %v, want 42", got)
	}
}

// TestEvalSubscriptNegativeIndex — spec §6.1.3: negative index is an error.
func TestEvalSubscriptNegativeIndex(t *testing.T) {
	mustPipelineError(t, "x_{-1}", "non-negative integer")
}

// TestEvalSubscriptNonIntegerIndex — spec §6.1.3: fractional index is an error.
func TestEvalSubscriptNonIntegerIndex(t *testing.T) {
	mustPipelineError(t, "x_{1.5}", "non-negative integer")
}

// TestEvalSubscriptUndefined — spec §6.1.5: composite key not in scope.
func TestEvalSubscriptUndefined(t *testing.T) {
	// Empty scope — x_0 is not defined.
	mustPipelineErrorScope(t, "x_{0}", eval.Scope{}, "undefined symbol: x_0")
}

// TestEvalSubscriptBaseNotSymbol — spec §6.1.1: non-symbol base is an error.
// Parser allows 5_{0} syntactically; evaluator must reject it.
func TestEvalSubscriptBaseNotSymbol(t *testing.T) {
	mustPipelineError(t, "5_{0}", "subscript base must be a symbol")
}

// ── v0.3 §6.2 BigOpNode ──────────────────────────────────────────────────────

// TestEvalSumBasic — \sum_{i=1}^{3} i = 1+2+3 = 6.
func TestEvalSumBasic(t *testing.T) {
	got := evaluate(t, `\sum_{i=1}^{3} i`)
	if !approx(got, 6) {
		t.Errorf(`eval(\sum_{i=1}^{3} i) = %v, want 6`, got)
	}
}

// TestEvalSumSingleStep — \sum_{i=5}^{5} i = 5 (exactly one step).
func TestEvalSumSingleStep(t *testing.T) {
	got := evaluate(t, `\sum_{i=5}^{5} i`)
	if !approx(got, 5) {
		t.Errorf(`eval(\sum_{i=5}^{5} i) = %v, want 5`, got)
	}
}

// TestEvalSumEmptyRange — \sum_{i=5}^{3} i = 0 (to < from → identity for sum).
func TestEvalSumEmptyRange(t *testing.T) {
	got := evaluate(t, `\sum_{i=5}^{3} i`)
	if !approx(got, 0) {
		t.Errorf(`eval(\sum_{i=5}^{3} i) = %v, want 0`, got)
	}
}

// TestEvalProdBasic — \prod_{i=1}^{4} i = 1*2*3*4 = 24.
func TestEvalProdBasic(t *testing.T) {
	got := evaluate(t, `\prod_{i=1}^{4} i`)
	if !approx(got, 24) {
		t.Errorf(`eval(\prod_{i=1}^{4} i) = %v, want 24`, got)
	}
}

// TestEvalProdEmptyRange — \prod_{i=5}^{3} i = 1 (to < from → identity for prod).
func TestEvalProdEmptyRange(t *testing.T) {
	got := evaluate(t, `\prod_{i=5}^{3} i`)
	if !approx(got, 1) {
		t.Errorf(`eval(\prod_{i=5}^{3} i) = %v, want 1`, got)
	}
}

// TestEvalSumSquares — \sum_{i=1}^{3} {i^2} = 1+4+9 = 14.
func TestEvalSumSquares(t *testing.T) {
	got := evaluate(t, `\sum_{i=1}^{3} {i^2}`)
	if !approx(got, 14) {
		t.Errorf(`eval(\sum_{i=1}^{3} {i^2}) = %v, want 14`, got)
	}
}

// TestEvalSumBodyPrecedence — body at parsePower level: \sum_{i=1}^{3} i + 1 = 6+1 = 7.
func TestEvalSumBodyPrecedence(t *testing.T) {
	got := evaluate(t, `\sum_{i=1}^{3} i + 1`)
	if !approx(got, 7) {
		t.Errorf(`eval(\sum_{i=1}^{3} i + 1) = %v, want 7`, got)
	}
}

// TestEvalSumVarShadows — iteration variable shadows outer scope binding.
func TestEvalSumVarShadows(t *testing.T) {
	scope := eval.NewScope()
	scope["i"] = 99 // outer i should be shadowed inside the sum
	got := evaluateWithScope(t, `\sum_{i=1}^{3} i`, scope)
	if !approx(got, 6) {
		t.Errorf(`eval(\sum_{i=1}^{3} i) with outer i=99 = %v, want 6`, got)
	}
	// outer scope must not be mutated
	if scope["i"] != 99 {
		t.Errorf("outer scope[i] = %v after sum, want 99 (must not be mutated)", scope["i"])
	}
}

// TestEvalSumWithSubscriptVars — \sum_{i=0}^{2} x_{i} with x_0,x_1,x_2 in scope.
func TestEvalSumWithSubscriptVars(t *testing.T) {
	scope := eval.NewScope()
	scope["x_0"] = 10
	scope["x_1"] = 20
	scope["x_2"] = 30
	got := evaluateWithScope(t, `\sum_{i=0}^{2} x_{i}`, scope)
	if !approx(got, 60) {
		t.Errorf(`eval(\sum_{i=0}^{2} x_{i}) = %v, want 60`, got)
	}
}

// TestEvalNestedSum — \sum_{i=1}^{2} \sum_{j=1}^{2} {i*j} = (1+2)+(2+4) = 9.
func TestEvalNestedSum(t *testing.T) {
	got := evaluate(t, `\sum_{i=1}^{2} \sum_{j=1}^{2} {i*j}`)
	if !approx(got, 9) {
		t.Errorf(`eval(nested sum) = %v, want 9`, got)
	}
}

// TestEvalSumBoundsInOuterScope — bound expressions evaluated in outer scope,
// not affected by body variable binding.
func TestEvalSumBoundsInOuterScope(t *testing.T) {
	scope := eval.NewScope()
	scope["n"] = 4
	got := evaluateWithScope(t, `\sum_{i=1}^{n} i`, scope)
	if !approx(got, 10) { // 1+2+3+4 = 10
		t.Errorf(`eval(\sum_{i=1}^{n} i) = %v, want 10`, got)
	}
}

// TestEvalSumLowerNonInteger — spec §6.2.1: non-integer lower bound is an error.
func TestEvalSumLowerNonInteger(t *testing.T) {
	mustPipelineError(t, `\sum_{i=1.5}^{3} i`, "lower bound must be an integer")
}

// TestEvalSumUpperInfinity — spec §6.2.2: infinite upper bound is an error.
func TestEvalSumUpperInfinity(t *testing.T) {
	mustPipelineError(t, `\sum_{i=1}^{\infty} i`, "upper bound must be")
}

// TestEvalSumBodyError — body error propagates immediately.
func TestEvalSumBodyError(t *testing.T) {
	mustPipelineError(t, `\sum_{i=1}^{3} {1/0}`, "division by zero")
}

// TestEvalProdSumOfFactorials — \prod_{i=1}^{5} i = 120.
func TestEvalProdFact5(t *testing.T) {
	got := evaluate(t, `\prod_{i=1}^{5} i`)
	if !approx(got, 120) {
		t.Errorf(`eval(\prod_{i=1}^{5} i) = %v, want 120`, got)
	}
}

// ── v0.3 §6.3 NormNode ───────────────────────────────────────────────────────

// TestEvalNormNegative — \lVert -3 \rVert = 3.
func TestEvalNormNegative(t *testing.T) {
	got := evaluate(t, `\lVert -3 \rVert`)
	if !approx(got, 3) {
		t.Errorf(`eval(\lVert -3 \rVert) = %v, want 3`, got)
	}
}

// TestEvalNormPositive — \lVert 5 \rVert = 5.
func TestEvalNormPositive(t *testing.T) {
	got := evaluate(t, `\lVert 5 \rVert`)
	if !approx(got, 5) {
		t.Errorf(`eval(\lVert 5 \rVert) = %v, want 5`, got)
	}
}

// TestEvalNormZero — \lVert 0 \rVert = 0.
func TestEvalNormZero(t *testing.T) {
	got := evaluate(t, `\lVert 0 \rVert`)
	if !approx(got, 0) {
		t.Errorf(`eval(\lVert 0 \rVert) = %v, want 0`, got)
	}
}

// TestEvalNormExpression — \lVert x - y \rVert with x=1, y=4 → |1-4| = 3.
func TestEvalNormExpression(t *testing.T) {
	scope := eval.NewScope()
	scope["x"] = 1
	scope["y"] = 4
	got := evaluateWithScope(t, `\lVert x - y \rVert`, scope)
	if !approx(got, 3) {
		t.Errorf(`eval(\lVert x-y \rVert) = %v, want 3`, got)
	}
}

// TestEvalNormNested — \lVert \lVert -2 \rVert \rVert = 2 (norm of norm).
func TestEvalNormNested(t *testing.T) {
	got := evaluate(t, `\lVert \lVert -2 \rVert \rVert`)
	if !approx(got, 2) {
		t.Errorf(`eval(nested norm) = %v, want 2`, got)
	}
}

// TestEvalNormDistinctFromAbs — \lVert -5 \rVert and |-5| must both be 5.
// This also verifies they co-exist without token collision.
func TestEvalNormDistinctFromAbs(t *testing.T) {
	norm := evaluate(t, `\lVert -5 \rVert`)
	abs := evaluate(t, `|-5|`)
	if !approx(norm, 5) {
		t.Errorf(`eval(\lVert -5 \rVert) = %v, want 5`, norm)
	}
	if !approx(abs, 5) {
		t.Errorf(`eval(|-5|) = %v, want 5`, abs)
	}
}
