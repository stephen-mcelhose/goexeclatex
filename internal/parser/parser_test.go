package parser_test

import (
	"math"
	"testing"

	"github.com/stephen-mcelhose/goexeclatex/internal/lexer"
	"github.com/stephen-mcelhose/goexeclatex/internal/parser"
)

// parse is a test helper: lex then parse, failing immediately on lex error.
func parse(t *testing.T, input string) parser.Node {
	t.Helper()
	tokens, err := lexer.Lex(input)
	if err != nil {
		t.Fatalf("Lex(%q) unexpected error: %v", input, err)
	}
	node, err := parser.Parse(tokens)
	if err != nil {
		t.Fatalf("Parse(%q) unexpected error: %v", input, err)
	}
	return node
}

// mustFail is a test helper: expects Parse to return a non-nil error.
func mustFail(t *testing.T, input string) error {
	t.Helper()
	tokens, err := lexer.Lex(input)
	if err != nil {
		// lex error is fine — the point is it doesn't succeed
		return err
	}
	_, err = parser.Parse(tokens)
	if err == nil {
		t.Errorf("Parse(%q) expected error, got nil", input)
	}
	return err
}

// ── §3.1 Invalid token stream ─────────────────────────────────────────────────

func TestInvalidTokenStream(t *testing.T) {
	_, err := parser.Parse(nil)
	if err == nil {
		t.Error("Parse(nil) expected error, got nil")
	}
	_, err = parser.Parse([]lexer.Token{})
	if err == nil {
		t.Error("Parse([]) expected error, got nil")
	}
}

func TestOnlyEOF(t *testing.T) {
	// A stream containing only EOF has no expression — must error.
	_, err := parser.Parse([]lexer.Token{{Type: lexer.EOF}})
	if err == nil {
		t.Error("Parse([EOF]) expected error, got nil")
	}
}

// ── §3.2 Post-lex invariants ──────────────────────────────────────────────────

func TestLeftRightDecoratorStripped(t *testing.T) {
	// \left and \right are dropped by the lexer; the brackets remain as LPAREN/RPAREN.
	// The parser must treat \left( x \right) identically to ( x ).
	node := parse(t, `\left( 1+2 \right)`)
	if _, ok := node.(*parser.BinaryNode); !ok {
		t.Errorf(`Parse(\left( 1+2 \right)) = %T, want *BinaryNode (transparent)`, node)
	}
}

func TestCdotIsTimes(t *testing.T) {
	// \cdot is remapped to TIMES by the lexer (spec §3.2).
	node := parse(t, `3\cdot4`)
	n, ok := node.(*parser.BinaryNode)
	if !ok || n.Op != "*" {
		t.Errorf(`Parse(3\cdot4) = %T op=%q, want BinaryNode(*)`, node, "")
	}
}

// ── §7.1 NumberNode ───────────────────────────────────────────────────────────

func TestNumber(t *testing.T) {
	cases := []struct {
		input string
		want  float64
	}{
		{"1", 1},
		{"3.14", 3.14},
		{"0", 0},
		{"100", 100},
	}
	for _, c := range cases {
		node := parse(t, c.input)
		n, ok := node.(*parser.NumberNode)
		if !ok {
			t.Errorf("Parse(%q) = %T, want *NumberNode", c.input, node)
			continue
		}
		if n.Value != c.want {
			t.Errorf("Parse(%q).Value = %v, want %v", c.input, n.Value, c.want)
		}
	}
}

// ── §7.2 SymbolNode ───────────────────────────────────────────────────────────

func TestSymbol(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"x", "x"},
		{"y", "y"},
		{`\pi`, "pi"},
		{`\tau`, "tau"},
		{`\phi`, "phi"},
	}
	for _, c := range cases {
		node := parse(t, c.input)
		n, ok := node.(*parser.SymbolNode)
		if !ok {
			t.Errorf("Parse(%q) = %T, want *SymbolNode", c.input, node)
			continue
		}
		if n.Name != c.want {
			t.Errorf("Parse(%q).Name = %q, want %q", c.input, n.Name, c.want)
		}
	}
}

// ── §7.3 BinaryNode — arithmetic ─────────────────────────────────────────────

func TestBinaryAddSub(t *testing.T) {
	cases := []struct {
		input string
		op    string
	}{
		{"1+2", "+"},
		{"1-2", "-"},
	}
	for _, c := range cases {
		node := parse(t, c.input)
		n, ok := node.(*parser.BinaryNode)
		if !ok {
			t.Errorf("Parse(%q) = %T, want *BinaryNode", c.input, node)
			continue
		}
		if n.Op != c.op {
			t.Errorf("Parse(%q).Op = %q, want %q", c.input, n.Op, c.op)
		}
	}
}

func TestBinaryMulDiv(t *testing.T) {
	cases := []struct {
		input string
		op    string
	}{
		{"2*3", "*"},
		{"6/2", "/"},
		{`2\times3`, "*"},
		{`6\div2`, "/"},
	}
	for _, c := range cases {
		node := parse(t, c.input)
		n, ok := node.(*parser.BinaryNode)
		if !ok {
			t.Errorf("Parse(%q) = %T, want *BinaryNode", c.input, node)
			continue
		}
		if n.Op != c.op {
			t.Errorf("Parse(%q).Op = %q, want %q", c.input, n.Op, c.op)
		}
	}
}

// ── §4.1 Associativity ────────────────────────────────────────────────────────

// Left-associativity of division: 6/2/3 → (6/2)/3
func TestLeftAssociativityDivision(t *testing.T) {
	node := parse(t, "6/2/3")
	outer, ok := node.(*parser.BinaryNode)
	if !ok || outer.Op != "/" {
		t.Fatalf("Parse(\"6/2/3\") = %T, want BinaryNode(/)", node)
	}
	inner, ok := outer.Left.(*parser.BinaryNode)
	if !ok || inner.Op != "/" {
		t.Errorf("Parse(\"6/2/3\").Left = %T, want BinaryNode(/) — left-associativity", outer.Left)
	}
}

// Left-associativity: 1-2-3 → (1-2)-3
func TestLeftAssociativity(t *testing.T) {
	node := parse(t, "1-2-3")
	outer, ok := node.(*parser.BinaryNode)
	if !ok || outer.Op != "-" {
		t.Fatalf("Parse(\"1-2-3\") = %T(%v), want BinaryNode(-)", node, node)
	}
	inner, ok := outer.Left.(*parser.BinaryNode)
	if !ok || inner.Op != "-" {
		t.Errorf("Parse(\"1-2-3\").Left = %T, want BinaryNode(-) — left-associativity", outer.Left)
	}
}

// Right-associativity: 2^3^4 → 2^(3^4)
func TestRightAssociativityPower(t *testing.T) {
	node := parse(t, "2^3^4")
	outer, ok := node.(*parser.BinaryNode)
	if !ok || outer.Op != "^" {
		t.Fatalf("Parse(\"2^3^4\") = %T, want BinaryNode(^)", node)
	}
	_, ok = outer.Right.(*parser.BinaryNode)
	if !ok {
		t.Errorf("Parse(\"2^3^4\").Right = %T, want BinaryNode(^) — right-associativity", outer.Right)
	}
}

// ── §4.1 Precedence ───────────────────────────────────────────────────────────

// 1+2*3 → 1+(2*3): multiply binds tighter than add
func TestPrecedenceAddMul(t *testing.T) {
	node := parse(t, "1+2*3")
	outer, ok := node.(*parser.BinaryNode)
	if !ok || outer.Op != "+" {
		t.Fatalf("Parse(\"1+2*3\") top = %T, want BinaryNode(+)", node)
	}
	if _, ok := outer.Right.(*parser.BinaryNode); !ok {
		t.Errorf("Parse(\"1+2*3\").Right = %T, want BinaryNode(*)", outer.Right)
	}
}

// 2*3^4 → 2*(3^4): power binds tighter than multiply
func TestPrecedenceMulPow(t *testing.T) {
	node := parse(t, "2*3^4")
	outer, ok := node.(*parser.BinaryNode)
	if !ok || outer.Op != "*" {
		t.Fatalf("Parse(\"2*3^4\") top = %T, want BinaryNode(*)", node)
	}
	if _, ok := outer.Right.(*parser.BinaryNode); !ok {
		t.Errorf("Parse(\"2*3^4\").Right = %T, want BinaryNode(^)", outer.Right)
	}
}

// ── §7.4 UnaryNode ────────────────────────────────────────────────────────────

func TestUnaryNegation(t *testing.T) {
	node := parse(t, "-3")
	n, ok := node.(*parser.UnaryNode)
	if !ok || n.Op != "-" {
		t.Errorf("Parse(\"-3\") = %T, want UnaryNode(-)", node)
	}
	if _, ok := n.Operand.(*parser.NumberNode); !ok {
		t.Errorf("Parse(\"-3\").Operand = %T, want NumberNode", n.Operand)
	}
}

func TestUnaryNegationOfGroup(t *testing.T) {
	// -(1+2) — negation wraps a grouped sub-expression
	node := parse(t, "-(1+2)")
	n, ok := node.(*parser.UnaryNode)
	if !ok || n.Op != "-" {
		t.Errorf("Parse(\"-(1+2)\") = %T, want UnaryNode(-)", node)
	}
	if _, ok := n.Operand.(*parser.BinaryNode); !ok {
		t.Errorf("Parse(\"-(1+2)\").Operand = %T, want BinaryNode", n.Operand)
	}
}

// -3^2 → -(3^2): unary minus has lower precedence than power (spec grammar order:
// unary → postfix; power → unary [POWER power])
func TestUnaryNegationPowerPrecedence(t *testing.T) {
	node := parse(t, "-3^2")
	outer, ok := node.(*parser.UnaryNode)
	if !ok || outer.Op != "-" {
		t.Fatalf("Parse(\"-3^2\") = %T, want UnaryNode(-)", node)
	}
	if _, ok := outer.Operand.(*parser.BinaryNode); !ok {
		t.Errorf("Parse(\"-3^2\").Operand = %T, want BinaryNode(^)", outer.Operand)
	}
}

func TestUnaryFactorial(t *testing.T) {
	node := parse(t, "5!")
	n, ok := node.(*parser.UnaryNode)
	if !ok || n.Op != "!" {
		t.Errorf("Parse(\"5!\") = %T, want UnaryNode(!)", node)
	}
	if _, ok := n.Operand.(*parser.NumberNode); !ok {
		t.Errorf("Parse(\"5!\").Operand = %T, want NumberNode", n.Operand)
	}
}

func TestFactorialOfGroup(t *testing.T) {
	// (5+1)! — factorial of a grouped sub-expression
	node := parse(t, "(5+1)!")
	n, ok := node.(*parser.UnaryNode)
	if !ok || n.Op != "!" {
		t.Errorf("Parse(\"(5+1)!\") = %T, want UnaryNode(!)", node)
	}
	if _, ok := n.Operand.(*parser.BinaryNode); !ok {
		t.Errorf("Parse(\"(5+1)!\").Operand = %T, want BinaryNode", n.Operand)
	}
}

// ── §4.3 Bracket matching — transparent grouping ─────────────────────────────

func TestGroupingParen(t *testing.T) {
	// (1+2) should produce the inner BinaryNode directly, no wrapper
	node := parse(t, "(1+2)")
	if _, ok := node.(*parser.BinaryNode); !ok {
		t.Errorf("Parse(\"(1+2)\") = %T, want *BinaryNode (transparent)", node)
	}
}

func TestGroupingSquare(t *testing.T) {
	node := parse(t, "[1+2]")
	if _, ok := node.(*parser.BinaryNode); !ok {
		t.Errorf("Parse(\"[1+2]\") = %T, want *BinaryNode (transparent)", node)
	}
}

func TestGroupingBrace(t *testing.T) {
	node := parse(t, "{1+2}")
	if _, ok := node.(*parser.BinaryNode); !ok {
		t.Errorf("Parse(\"{1+2}\") = %T, want *BinaryNode (transparent)", node)
	}
}

func TestMismatchedBracketError(t *testing.T) {
	mustFail(t, "(1+2]")
	mustFail(t, "[1+2)")
	mustFail(t, "{1+2)")
	mustFail(t, "{1+2]")
	mustFail(t, "(1+2}")
}

func TestUnclosedBracketError(t *testing.T) {
	// An opened bracket with no closer must produce an error.
	mustFail(t, "(1+2")
	mustFail(t, "[1+2")
	// Note: "{1+2" is caught at lex time (unclosed brace — spec §3.1), not parser time.
}

// ── §4.2 Absolute value ───────────────────────────────────────────────────────

func TestAbsoluteValue(t *testing.T) {
	node := parse(t, "|3+4|")
	n, ok := node.(*parser.FunctionNode)
	if !ok || n.Name != "abs" {
		t.Errorf("Parse(\"|3+4|\") = %T, want FunctionNode(abs)", node)
	}
	if len(n.Args) != 1 {
		t.Errorf("FunctionNode(abs).Args len = %d, want 1", len(n.Args))
	}
}

func TestAbsoluteValueLvert(t *testing.T) {
	node := parse(t, `\lvert 3+4 \rvert`)
	n, ok := node.(*parser.FunctionNode)
	if !ok || n.Name != "abs" {
		t.Errorf(`Parse(\lvert 3+4 \rvert) = %T, want FunctionNode(abs)`, node)
	}
}

func TestUnclosedAbsValueError(t *testing.T) {
	// A PIPE with no closing PIPE must produce an error.
	mustFail(t, "|3+4")
}

func TestImplicitMultiplyWithAbs(t *testing.T) {
	// 2|x| → 2 * abs(x): PIPE triggers implicit multiply at absDepth == 0 (ADR-007)
	node := parse(t, "2|x|")
	outer, ok := node.(*parser.BinaryNode)
	if !ok || outer.Op != "*" {
		t.Fatalf("Parse(\"2|x|\") = %T, want BinaryNode(*)", node)
	}
	fn, ok := outer.Right.(*parser.FunctionNode)
	if !ok || fn.Name != "abs" {
		t.Errorf("Parse(\"2|x|\").Right = %T, want FunctionNode(abs)", outer.Right)
	}
}

// ── §5 Implicit multiplication ────────────────────────────────────────────────

func TestImplicitMultiplyNumberSymbol(t *testing.T) {
	node := parse(t, "2x")
	n, ok := node.(*parser.BinaryNode)
	if !ok || n.Op != "*" {
		t.Errorf("Parse(\"2x\") = %T, want BinaryNode(*)", node)
	}
}

func TestImplicitMultiplyNumberParen(t *testing.T) {
	node := parse(t, "2(3+4)")
	n, ok := node.(*parser.BinaryNode)
	if !ok || n.Op != "*" {
		t.Errorf("Parse(\"2(3+4)\") = %T, want BinaryNode(*)", node)
	}
}

func TestImplicitMultiplyWithCommand(t *testing.T) {
	// 2\sin x → 2 * sin(x) — spec §5 example table
	node := parse(t, `2\sin x`)
	outer, ok := node.(*parser.BinaryNode)
	if !ok || outer.Op != "*" {
		t.Fatalf(`Parse(2\sin x) = %T, want BinaryNode(*)`, node)
	}
	fn, ok := outer.Right.(*parser.FunctionNode)
	if !ok || fn.Name != "sin" {
		t.Errorf(`Parse(2\sin x).Right = %T, want FunctionNode(sin)`, outer.Right)
	}
}

func TestImplicitMultiplySymbols(t *testing.T) {
	// "xy" (no space) is a single multi-char identifier per evaluatex convention.
	// Implicit multiply between two symbols requires whitespace: "x y".
	node := parse(t, "x y")
	n, ok := node.(*parser.BinaryNode)
	if !ok || n.Op != "*" {
		t.Errorf("Parse(\"x y\") = %T, want BinaryNode(*)", node)
	}
}

// Implicit multiply is at product level: 2x^2 → 2*(x^2), not (2x)^2
func TestImplicitMultiplyNotBindingToPower(t *testing.T) {
	node := parse(t, "2x^2")
	outer, ok := node.(*parser.BinaryNode)
	if !ok || outer.Op != "*" {
		t.Fatalf("Parse(\"2x^2\") top = %T(%v), want BinaryNode(*)", node, node)
	}
	// Right side must be x^2 (BinaryNode), not just x (SymbolNode)
	if _, ok := outer.Right.(*parser.BinaryNode); !ok {
		t.Errorf("Parse(\"2x^2\").Right = %T, want BinaryNode(^)", outer.Right)
	}
}

// ── §6 Command arity dispatch ─────────────────────────────────────────────────

func TestCommandArity0Symbol(t *testing.T) {
	cases := []struct{ input, want string }{
		{`\pi`, "pi"},
		{`\tau`, "tau"},
		{`\phi`, "phi"},
	}
	for _, c := range cases {
		node := parse(t, c.input)
		n, ok := node.(*parser.SymbolNode)
		if !ok || n.Name != c.want {
			t.Errorf("Parse(%q) = %T(%v), want SymbolNode(%s)", c.input, node, node, c.want)
		}
	}
}

func TestUnknownCommandIsSymbol(t *testing.T) {
	// An unknown command (not in the arity table) → arity 0 → SymbolNode.
	node := parse(t, `\unknowncmd`)
	n, ok := node.(*parser.SymbolNode)
	if !ok || n.Name != "unknowncmd" {
		t.Errorf(`Parse(\unknowncmd) = %T, want SymbolNode("unknowncmd")`, node)
	}
}

func TestArcsinIsFunction(t *testing.T) {
	// \arcsin is a known arity-1 command → FunctionNode (spec §6.3).
	node := parse(t, `\arcsin{1}`)
	n, ok := node.(*parser.FunctionNode)
	if !ok || n.Name != "arcsin" || len(n.Args) != 1 {
		t.Errorf(`Parse(\arcsin{1}) = %T, want FunctionNode("arcsin", 1 arg)`, node)
	}
}

func TestCommandArity1BraceArg(t *testing.T) {
	// \sqrt{4} — brace-wrapped argument
	node := parse(t, `\sqrt{4}`)
	n, ok := node.(*parser.FunctionNode)
	if !ok || n.Name != "sqrt" {
		t.Errorf(`Parse(\sqrt{4}) = %T, want FunctionNode(sqrt)`, node)
		return
	}
	if len(n.Args) != 1 {
		t.Errorf("sqrt args = %d, want 1", len(n.Args))
	}
}

func TestCommandArity1CharModeArg(t *testing.T) {
	// \sin x — char-mode single-token arg (no braces)
	node := parse(t, `\sin x`)
	n, ok := node.(*parser.FunctionNode)
	if !ok || n.Name != "sin" {
		t.Errorf(`Parse(\sin x) = %T, want FunctionNode(sin)`, node)
		return
	}
	if len(n.Args) != 1 {
		t.Errorf("sin args = %d, want 1", len(n.Args))
	}
}

func TestCommandArity2(t *testing.T) {
	// \frac{1}{2}
	node := parse(t, `\frac{1}{2}`)
	n, ok := node.(*parser.FunctionNode)
	if !ok || n.Name != "frac" {
		t.Errorf(`Parse(\frac{1}{2}) = %T, want FunctionNode(frac)`, node)
		return
	}
	if len(n.Args) != 2 {
		t.Errorf("frac args = %d, want 2", len(n.Args))
	}
}

func TestCommandAllTrig(t *testing.T) {
	cmds := []string{
		`\sin{x}`, `\cos{x}`, `\tan{x}`,
		`\sec{x}`, `\csc{x}`, `\cot{x}`,
		`\asin{x}`, `\acos{x}`, `\atan{x}`,
		`\asec{x}`, `\acsc{x}`, `\acot{x}`,
	}
	for _, c := range cmds {
		node := parse(t, c)
		n, ok := node.(*parser.FunctionNode)
		if !ok || len(n.Args) != 1 {
			t.Errorf("Parse(%q) = %T, want FunctionNode with 1 arg", c, node)
		}
	}
}

// ── §6.2 Missing command argument errors ─────────────────────────────────────

func TestMissingArgSqrt(t *testing.T) {
	// \sqrt with no argument — should error, not silently return
	err := mustFail(t, `\sqrt`)
	if err == nil {
		t.Error(`Parse(\sqrt) expected error, got nil`)
	}
}

func TestMissingSecondArgFrac(t *testing.T) {
	// \frac{1} with no second argument — should error
	err := mustFail(t, `\frac{1}`)
	if err == nil {
		t.Error(`Parse(\frac{1}) expected error, got nil`)
	}
}

// ── §8 Error conditions + message prefix ─────────────────────────────────────

func TestUnexpectedEOF(t *testing.T) {
	// "1+" — trailing operator with no right operand
	mustFail(t, "1+")
}

func TestUnexpectedToken(t *testing.T) {
	// "+" at start — nothing to be a left operand for a binary op, and
	// unary minus is allowed but unary plus is not
	mustFail(t, "+1")
}

func TestErrorMessagePrefix(t *testing.T) {
	// All parser errors must begin with "parser: " (spec §8).
	cases := []string{"1+", "+1", "(1+2]", `\sqrt`}
	for _, c := range cases {
		tokens, lexErr := lexer.Lex(c)
		if lexErr != nil {
			continue // lex error — skip, not testing parser here
		}
		_, err := parser.Parse(tokens)
		if err == nil {
			t.Errorf("Parse(%q) expected error, got nil", c)
			continue
		}
		msg := err.Error()
		if len(msg) < 8 || msg[:8] != "parser: " {
			t.Errorf("Parse(%q) error = %q, want prefix \"parser: \"", c, msg)
		}
	}
}

// ── Integration: compound expressions ────────────────────────────────────────

func TestNestedFrac(t *testing.T) {
	// \frac{\sqrt{2}}{3}
	node := parse(t, `\frac{\sqrt{2}}{3}`)
	n, ok := node.(*parser.FunctionNode)
	if !ok || n.Name != "frac" || len(n.Args) != 2 {
		t.Errorf(`Parse(\frac{\sqrt{2}}{3}) = %T, want FunctionNode(frac,2)`, node)
		return
	}
	if _, ok := n.Args[0].(*parser.FunctionNode); !ok {
		t.Errorf("frac arg[0] = %T, want FunctionNode(sqrt)", n.Args[0])
	}
}

func TestPowerWithBraceGroup(t *testing.T) {
	// 2^{10} — brace group as power exponent
	node := parse(t, "2^{10}")
	n, ok := node.(*parser.BinaryNode)
	if !ok || n.Op != "^" {
		t.Errorf("Parse(\"2^{10}\") = %T, want BinaryNode(^)", node)
		return
	}
	exp, ok := n.Right.(*parser.NumberNode)
	if !ok || exp.Value != 10 {
		t.Errorf("Parse(\"2^{10}\").Right = %T(%v), want NumberNode(10)", n.Right, n.Right)
	}
}

func TestComplexExpression(t *testing.T) {
	// \frac{1}{2} + \sqrt{9}
	node := parse(t, `\frac{1}{2} + \sqrt{9}`)
	n, ok := node.(*parser.BinaryNode)
	if !ok || n.Op != "+" {
		t.Fatalf(`Parse(\frac{1}{2} + \sqrt{9}) = %T, want BinaryNode(+)`, node)
	}
	if _, ok := n.Left.(*parser.FunctionNode); !ok {
		t.Errorf("Left = %T, want FunctionNode(frac)", n.Left)
	}
	if _, ok := n.Right.(*parser.FunctionNode); !ok {
		t.Errorf("Right = %T, want FunctionNode(sqrt)", n.Right)
	}
}

// NaN-safe float comparison
func approxEqual(a, b float64) bool {
	if math.IsNaN(a) && math.IsNaN(b) {
		return true
	}
	return math.Abs(a-b) < 1e-12
}
