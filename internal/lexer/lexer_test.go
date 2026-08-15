package lexer

import (
	"testing"
)

// tok is a test helper that builds a Token without caring about Pos.
func tok(typ TokenType, value string) Token {
	return Token{Type: typ, Value: value}
}

// lex runs Lex and strips Pos from all tokens so table tests stay readable.
func lex(t *testing.T, input string) []Token {
	t.Helper()
	tokens, err := Lex(input)
	if err != nil {
		t.Fatalf("Lex(%q) unexpected error: %v", input, err)
	}
	out := make([]Token, len(tokens))
	for i, tk := range tokens {
		out[i] = Token{Type: tk.Type, Value: tk.Value}
	}
	return out
}

// ---- §3 Token Types --------------------------------------------------------

// TestNumbers covers spec §3 NUMBER pattern.
func TestNumbers(t *testing.T) {
	cases := []struct {
		input string
		want  Token
	}{
		{"42", tok(NUMBER, "42")},
		{"3.14", tok(NUMBER, "3.14")},
		{"0.5", tok(NUMBER, "0.5")},
		{"1e10", tok(NUMBER, "1e10")},
		{"2.5e-3", tok(NUMBER, "2.5e-3")},
		{"1E+2", tok(NUMBER, "1E+2")},
	}
	for _, c := range cases {
		tokens := lex(t, c.input)
		if tokens[0] != c.want {
			t.Errorf("Lex(%q)[0] = %v, want %v", c.input, tokens[0], c.want)
		}
	}
}

// TestSymbols covers spec §3 SYMBOL pattern.
func TestSymbols(t *testing.T) {
	cases := []struct {
		input string
		want  Token
	}{
		{"x", tok(SYMBOL, "x")},
		{"abc", tok(SYMBOL, "abc")},
		{"x1", tok(SYMBOL, "x1")},
		{"MyVar", tok(SYMBOL, "MyVar")},
	}
	for _, c := range cases {
		tokens := lex(t, c.input)
		if tokens[0] != c.want {
			t.Errorf("Lex(%q)[0] = %v, want %v", c.input, tokens[0], c.want)
		}
	}
}

// TestCommands covers spec §3 COMMAND type and §7.3 normalisation.
// Commands with arity > 0 must be given their arguments (spec §9: missing
// args in char mode are an error).  Arity-0 commands (\pi) can appear alone.
func TestCommands(t *testing.T) {
	cases := []struct {
		input string
		value string // normalised name: backslash-stripped, lowercased
	}{
		{`\sin x`, "sin"},   // arity 1 — x is the char-mode arg
		{`\frac{1}{2}`, "frac"}, // arity 2 — two group args
		{`\sqrt{4}`, "sqrt"}, // arity 1 — group arg
		{`\pi`, "pi"},        // arity 0 — no args needed
		{`\Sin x`, "sin"},   // §7.3: lowercased
		{`\FRAC{1}{2}`, "frac"}, // §7.3: lowercased
	}
	for _, c := range cases {
		tokens := lex(t, c.input)
		if tokens[0].Type != COMMAND {
			t.Errorf("Lex(%q)[0].Type = %v, want COMMAND", c.input, tokens[0].Type)
		}
		if tokens[0].Value != c.value {
			t.Errorf("Lex(%q)[0].Value = %q, want %q", c.input, tokens[0].Value, c.value)
		}
	}
}

// TestCommandTokenStreams verifies the full token sequence for command
// invocations, including correct {}-group wrapping and nested brace handling.
// TestCommands above only checks tokens[0]; this test checks the whole stream.
func TestCommandTokenStreams(t *testing.T) {
	cases := []struct {
		input string
		want  []Token
	}{
		{
			// arity-0: no brace args, just the command + EOF
			`\pi`,
			[]Token{tok(COMMAND, "pi"), tok(EOF, "")},
		},
		{
			// arity-1 with single-char arg (char-mode, not brace group)
			`\sin x`,
			[]Token{tok(COMMAND, "sin"), tok(SYMBOL, "x"), tok(EOF, "")},
		},
		{
			// arity-1 with brace group arg
			`\sqrt{4}`,
			[]Token{
				tok(COMMAND, "sqrt"),
				tok(LPAREN, "{"), tok(NUMBER, "4"), tok(RPAREN, "}"),
				tok(EOF, ""),
			},
		},
		{
			// arity-2: two sibling brace groups
			`\frac{1}{2}`,
			[]Token{
				tok(COMMAND, "frac"),
				tok(LPAREN, "{"), tok(NUMBER, "1"), tok(RPAREN, "}"),
				tok(LPAREN, "{"), tok(NUMBER, "2"), tok(RPAREN, "}"),
				tok(EOF, ""),
			},
		},
		{
			// nested: \sqrt inside first arg of \frac
			`\frac{\sqrt{2}}{3}`,
			[]Token{
				tok(COMMAND, "frac"),
				tok(LPAREN, "{"),
				tok(COMMAND, "sqrt"),
				tok(LPAREN, "{"), tok(NUMBER, "2"), tok(RPAREN, "}"),
				tok(RPAREN, "}"),
				tok(LPAREN, "{"), tok(NUMBER, "3"), tok(RPAREN, "}"),
				tok(EOF, ""),
			},
		},
		{
			// deeply nested: \frac inside \sqrt inside \frac
			`\frac{\sqrt{\frac{1}{2}}}{3}`,
			[]Token{
				tok(COMMAND, "frac"),
				tok(LPAREN, "{"),
				tok(COMMAND, "sqrt"),
				tok(LPAREN, "{"),
				tok(COMMAND, "frac"),
				tok(LPAREN, "{"), tok(NUMBER, "1"), tok(RPAREN, "}"),
				tok(LPAREN, "{"), tok(NUMBER, "2"), tok(RPAREN, "}"),
				tok(RPAREN, "}"),
				tok(RPAREN, "}"),
				tok(LPAREN, "{"), tok(NUMBER, "3"), tok(RPAREN, "}"),
				tok(EOF, ""),
			},
		},
		{
			// POWER with brace group
			`2^{10}`,
			[]Token{
				tok(NUMBER, "2"),
				tok(POWER, "^"),
				tok(LPAREN, "{"), tok(NUMBER, "10"), tok(RPAREN, "}"),
				tok(EOF, ""),
			},
		},
	}
	for _, c := range cases {
		tokens := lex(t, c.input)
		if !equalTokens(tokens, c.want) {
			t.Errorf("Lex(%q)\n got  %v\n want %v", c.input, tokens, c.want)
		}
	}
}

// TestOperators covers spec §3 PLUS/MINUS/TIMES/DIVIDE/POWER tokens.
// POWER requires a char-mode argument (spec §6.2), so it is tested with one.
func TestOperators(t *testing.T) {
	cases := []struct {
		input  string
		idx    int // index of the operator token in the stream
		want   TokenType
	}{
		{"2+3", 1, PLUS},
		{"2-3", 1, MINUS},
		{"2*3", 1, TIMES},
		{"2/3", 1, DIVIDE},
		{"2^3", 1, POWER}, // ^ consumes "3" as its char-mode arg
	}
	for _, c := range cases {
		tokens := lex(t, c.input)
		if tokens[c.idx].Type != c.want {
			t.Errorf("Lex(%q)[%d].Type = %v, want %v", c.input, c.idx, tokens[c.idx].Type, c.want)
		}
	}
}

// TestGrouping covers spec §3 LPAREN/RPAREN for (, [, {, ), ], }.
func TestGrouping(t *testing.T) {
	cases := []struct {
		input string
		want  TokenType
		value string
	}{
		{"(", LPAREN, "("},
		{"[", LPAREN, "["},
		{"{1}", LPAREN, "{"}, // bare { is an unclosed group — needs content + closing }
		{")", RPAREN, ")"},
		{"]", RPAREN, "]"},
		{"}", RPAREN, "}"},
	}
	for _, c := range cases {
		tokens := lex(t, c.input)
		if tokens[0].Type != c.want || tokens[0].Value != c.value {
			t.Errorf("Lex(%q)[0] = %v, want %v(%s)", c.input, tokens[0], c.want, c.value)
		}
	}
}

// TestMisc covers PIPE, BANG, COMMA, and EOF.
func TestMisc(t *testing.T) {
	cases := []struct {
		input string
		want  TokenType
	}{
		{"|", PIPE},
		{"!", BANG},
		{",", COMMA},
	}
	for _, c := range cases {
		tokens := lex(t, c.input)
		if tokens[0].Type != c.want {
			t.Errorf("Lex(%q)[0].Type = %v, want %v", c.input, tokens[0].Type, c.want)
		}
	}

	// EOF is always the last token (spec §8).
	tokens := lex(t, "1")
	last := tokens[len(tokens)-1]
	if last.Type != EOF {
		t.Errorf("last token = %v, want EOF", last)
	}
}

// ---- §5 Whitespace ----------------------------------------------------------

// TestWhitespaceSkipped covers spec §5.3: whitespace produces no token.
func TestWhitespaceSkipped(t *testing.T) {
	tokens := lex(t, "  2  +  3  ")
	want := []Token{
		tok(NUMBER, "2"),
		tok(PLUS, "+"),
		tok(NUMBER, "3"),
		tok(EOF, ""),
	}
	if !equalTokens(tokens, want) {
		t.Errorf("Lex(\"  2  +  3  \") = %v, want %v", tokens, want)
	}
}

// TestUnknownCharacter covers spec §5.4: error on unrecognised character.
func TestUnknownCharacter(t *testing.T) {
	_, err := Lex("2 @ 3")
	if err == nil {
		t.Errorf("Lex(\"2 @ 3\") expected error, got nil")
	}
}

// ---- §6 Char-Mode Rule ------------------------------------------------------

// TestCharModePower covers spec §6: 2^24 → NUMBER(2) POWER NUMBER(2) NUMBER(4).
// (LaTeX single-char superscript rule — spec §6.1/6.2)
func TestCharModePower(t *testing.T) {
	tokens := lex(t, "2^24")
	want := []Token{
		tok(NUMBER, "2"),
		tok(POWER, "^"),
		tok(NUMBER, "2"), // char mode: only "2" consumed, not "24"
		tok(NUMBER, "4"),
		tok(EOF, ""),
	}
	if !equalTokens(tokens, want) {
		t.Errorf("Lex(\"2^24\") = %v, want %v", tokens, want)
	}
}

// TestCharModePowerGroup covers spec §6.3: 2^{12} treats {12} as one group.
func TestCharModePowerGroup(t *testing.T) {
	tokens := lex(t, "2^{12}")
	want := []Token{
		tok(NUMBER, "2"),
		tok(POWER, "^"),
		tok(LPAREN, "{"),
		tok(NUMBER, "12"),
		tok(RPAREN, "}"),
		tok(EOF, ""),
	}
	if !equalTokens(tokens, want) {
		t.Errorf("Lex(\"2^{12}\") = %v, want %v", tokens, want)
	}
}

// TestCharModeCommand covers spec §6.2: \frac reads arity=2 char-mode args.
// \frac 4 2 → COMMAND(frac) NUMBER(4) NUMBER(2)
func TestCharModeCommand(t *testing.T) {
	tokens := lex(t, `\frac 4 2`)
	want := []Token{
		tok(COMMAND, "frac"),
		tok(NUMBER, "4"),
		tok(NUMBER, "2"),
		tok(EOF, ""),
	}
	if !equalTokens(tokens, want) {
		t.Errorf(`Lex(\frac 4 2) = %v, want %v`, tokens, want)
	}
}

// TestCharModeCommandGroup covers spec §6.3: \frac{1}{2} group expansion.
func TestCharModeCommandGroup(t *testing.T) {
	tokens := lex(t, `\frac{1}{2}`)
	want := []Token{
		tok(COMMAND, "frac"),
		tok(LPAREN, "{"),
		tok(NUMBER, "1"),
		tok(RPAREN, "}"),
		tok(LPAREN, "{"),
		tok(NUMBER, "2"),
		tok(RPAREN, "}"),
		tok(EOF, ""),
	}
	if !equalTokens(tokens, want) {
		t.Errorf(`Lex(\frac{1}{2}) = %v, want %v`, tokens, want)
	}
}

// TestCharModeCommandMixed covers spec §6: \frac 1{20}3 — first arg single char,
// second arg a group; trailing 3 is a separate token.
// Expected: COMMAND(frac) NUMBER(1) LPAREN({) NUMBER(20) RPAREN(}) NUMBER(3)
func TestCharModeCommandMixed(t *testing.T) {
	tokens := lex(t, `\frac 1{20}3`)
	want := []Token{
		tok(COMMAND, "frac"),
		tok(NUMBER, "1"),
		tok(LPAREN, "{"),
		tok(NUMBER, "20"),
		tok(RPAREN, "}"),
		tok(NUMBER, "3"),
		tok(EOF, ""),
	}
	if !equalTokens(tokens, want) {
		t.Errorf(`Lex(\frac 1{20}3) = %v, want %v`, tokens, want)
	}
}

// TestCharModeCommandInArg covers spec §6.1: if the char-mode character is \,
// the lexer MUST read a full COMMAND token.
// \sin \pi → COMMAND(sin) COMMAND(pi)  (sin reads \pi as its single char arg)
func TestCharModeCommandInArg(t *testing.T) {
	tokens := lex(t, `\sin \pi`)
	want := []Token{
		tok(COMMAND, "sin"),
		tok(COMMAND, "pi"),
		tok(EOF, ""),
	}
	if !equalTokens(tokens, want) {
		t.Errorf(`Lex(\sin \pi) = %v, want %v`, tokens, want)
	}
}

// TestCharModeUnknownCommand covers spec §6.4: commands not in the arity table
// have arity 0 — they consume no char-mode args.
func TestCharModeUnknownCommand(t *testing.T) {
	// \pi is not in the arity table (arity 0); 3 should be a separate token.
	tokens := lex(t, `\pi 3`)
	want := []Token{
		tok(COMMAND, "pi"),
		tok(NUMBER, "3"),
		tok(EOF, ""),
	}
	if !equalTokens(tokens, want) {
		t.Errorf(`Lex(\pi 3) = %v, want %v`, tokens, want)
	}
}

// ---- §7 Post-Lex Remapping --------------------------------------------------

// TestRemapLeftRight covers spec §7.1: \left( → LPAREN, \right) → RPAREN.
func TestRemapLeftRight(t *testing.T) {
	tokens := lex(t, `\left( x \right)`)
	if tokens[0].Type != LPAREN || tokens[0].Value != "\\left(" {
		t.Errorf("\\left( should produce LPAREN, got %v", tokens[0])
	}
	// find the last non-EOF token
	last := tokens[len(tokens)-2]
	if last.Type != RPAREN {
		t.Errorf("\\right) should produce RPAREN, got %v", last)
	}
}

// TestRemapLeftStandalone covers spec §7.1: standalone \left (before \lvert)
// MUST be dropped.
func TestRemapLeftStandalone(t *testing.T) {
	tokens := lex(t, `\left\lvert x \right\rvert`)
	// \left dropped, \lvert → PIPE, x, \right dropped, \rvert → PIPE
	want := []Token{
		tok(PIPE, "|"),
		tok(SYMBOL, "x"),
		tok(PIPE, "|"),
		tok(EOF, ""),
	}
	if !equalTokens(tokens, want) {
		t.Errorf(`Lex(\left\lvert x \right\rvert) = %v, want %v`, tokens, want)
	}
}

// TestRemapTimes covers spec §7.2: \times → TIMES.
func TestRemapTimes(t *testing.T) {
	tokens := lex(t, `2\times 3`)
	if tokens[1].Type != TIMES {
		t.Errorf(`\times should produce TIMES, got %v`, tokens[1])
	}
}

// TestRemapCdot covers spec §7.2: \cdot → TIMES.
func TestRemapCdot(t *testing.T) {
	tokens := lex(t, `2\cdot 3`)
	if tokens[1].Type != TIMES {
		t.Errorf(`\cdot should produce TIMES, got %v`, tokens[1])
	}
}

// TestRemapDiv covers spec §7.2: \div → DIVIDE.
func TestRemapDiv(t *testing.T) {
	tokens := lex(t, `6\div 2`)
	if tokens[1].Type != DIVIDE {
		t.Errorf(`\div should produce DIVIDE, got %v`, tokens[1])
	}
}

// TestRemapLvert covers spec §7.2: \lvert and \rvert → PIPE.
func TestRemapLvert(t *testing.T) {
	tokens := lex(t, `\lvert x \rvert`)
	if tokens[0].Type != PIPE {
		t.Errorf(`\lvert should produce PIPE, got %v`, tokens[0])
	}
	if tokens[2].Type != PIPE {
		t.Errorf(`\rvert should produce PIPE, got %v`, tokens[2])
	}
}

// ---- Integration: expression shapes from docs/examples.md ------------------

// TestExprFrac covers \frac{1}{2} + \sqrt{9} token stream.
func TestExprFrac(t *testing.T) {
	tokens := lex(t, `\frac{1}{2} + \sqrt{9}`)
	// COMMAND(frac) LPAREN NUMBER(1) RPAREN LPAREN NUMBER(2) RPAREN PLUS COMMAND(sqrt) LPAREN NUMBER(9) RPAREN EOF
	if tokens[0] != tok(COMMAND, "frac") {
		t.Errorf("token[0] = %v, want COMMAND(frac)", tokens[0])
	}
	plusIdx := -1
	for i, tk := range tokens {
		if tk.Type == PLUS {
			plusIdx = i
			break
		}
	}
	if plusIdx < 0 {
		t.Errorf("no PLUS token found in stream: %v", tokens)
	}
}

// TestExprPower covers 2^{10} → NUMBER(2) POWER LPAREN NUMBER(10) RPAREN.
func TestExprPower(t *testing.T) {
	tokens := lex(t, "2^{10}")
	want := []Token{
		tok(NUMBER, "2"),
		tok(POWER, "^"),
		tok(LPAREN, "{"),
		tok(NUMBER, "10"),
		tok(RPAREN, "}"),
		tok(EOF, ""),
	}
	if !equalTokens(tokens, want) {
		t.Errorf("Lex(\"2^{10}\") = %v, want %v", tokens, want)
	}
}

// TestExprSinPi covers \sin(\pi / 6).
func TestExprSinPi(t *testing.T) {
	tokens := lex(t, `\sin(\pi / 6)`)
	// COMMAND(sin) LPAREN(() COMMAND(pi) DIVIDE NUMBER(6) RPAREN()) EOF
	// Note: \sin has arity 1 — char-mode reads "(" as the single arg (an LPAREN).
	// But ( is not { so no group recursion; char-mode returns after (.
	// Then parser sees the closing ) and DIVIDE NUMBER(6) as separate tokens.
	if tokens[0] != tok(COMMAND, "sin") {
		t.Errorf("token[0] = %v, want COMMAND(sin)", tokens[0])
	}
	if tokens[1].Type != LPAREN {
		t.Errorf("token[1] = %v, want LPAREN", tokens[1])
	}
}

// ---- Helpers ----------------------------------------------------------------

// ---- Error paths (spec §9) --------------------------------------------------

// lexErr is a helper that asserts Lex returns a non-nil error and returns it.
func lexErr(t *testing.T, input string) error {
	t.Helper()
	_, err := Lex(input)
	if err == nil {
		t.Fatalf("Lex(%q): expected error, got nil", input)
	}
	return err
}

// TestErrorUnclosedBraceGroup verifies that unclosed { groups produce an error
// (spec §9). BLOCKING: previously these silently succeeded.
func TestErrorUnclosedBraceGroup(t *testing.T) {
	cases := []string{
		"2^{12",      // unclosed brace after POWER
		`\sqrt{4`,    // unclosed brace after 1-arg command
		`\frac{1}{2`, // unclosed second brace of 2-arg command
	}
	for _, input := range cases {
		_, err := Lex(input)
		if err == nil {
			t.Errorf("Lex(%q): expected error for unclosed '{' group, got nil", input)
		}
	}
}

// TestErrorEOFInCharMode verifies that EOF reached in char mode returns an error
// (spec §9). Covers POWER with no arg and a command with no arg.
func TestErrorEOFInCharMode(t *testing.T) {
	cases := []string{
		"2^",   // POWER with no argument
		`\sin`, // arity-1 command with no argument
	}
	for _, input := range cases {
		_, err := Lex(input)
		if err == nil {
			t.Errorf("Lex(%q): expected error for EOF in char mode, got nil", input)
		}
	}
}
// Known limitation: multi-byte UTF-8 runes in char-mode (e.g. 2^α) produce an
// error with an escaped byte in the message (\xce) rather than the full rune (α).
// LaTeX math expressions are ASCII-dominant so this is not fixed in Tier 1.
// See GAN review §D advisory "Broken UTF-8 Representation in Unexpected Character Error".

// ---- Helpers ----------------------------------------------------------------

// equalTokens compares two token slices ignoring Pos.
func equalTokens(a, b []Token) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Type != b[i].Type || a[i].Value != b[i].Value {
			return false
		}
	}
	return true
}
