package lexer

// Tests for spec §3.1 (UNDERSCORE), §3.2 (EQUALS), §3.3 (NORM).
// All tests in this file are RED until the corresponding implementation lands.

import "testing"

// ── §3.1 UNDERSCORE token ─────────────────────────────────────────────────────

// TestUnderscoreAlone — bare _ → UNDERSCORE with value "_"
func TestUnderscoreAlone(t *testing.T) {
	tokens := lex(t, "_")
	if tokens[0] != tok(UNDERSCORE, "_") {
		t.Errorf(`lex("_")[0] = %v, want UNDERSCORE("_")`, tokens[0])
	}
}

// TestUnderscoreCharModeNoBrace — x_i → SYMBOL(x) UNDERSCORE SYMBOL(i) EOF
// The _ triggers char-mode and reads the next single token as its argument.
func TestUnderscoreCharModeNoBrace(t *testing.T) {
	want := []Token{
		tok(SYMBOL, "x"),
		tok(UNDERSCORE, "_"),
		tok(SYMBOL, "i"),
		tok(EOF, ""),
	}
	tokens := lex(t, "x_i")
	if len(tokens) != len(want) {
		t.Fatalf("lex(%q) = %v (len %d), want len %d", "x_i", tokens, len(tokens), len(want))
	}
	for i, w := range want {
		if tokens[i] != w {
			t.Errorf("lex(%q)[%d] = %v, want %v", "x_i", i, tokens[i], w)
		}
	}
}

// TestUnderscoreBraceMode — x_{0} → SYMBOL(x) UNDERSCORE LPAREN({) NUMBER(0) RPAREN(}) EOF
func TestUnderscoreBraceMode(t *testing.T) {
	want := []Token{
		tok(SYMBOL, "x"),
		tok(UNDERSCORE, "_"),
		tok(LPAREN, "{"),
		tok(NUMBER, "0"),
		tok(RPAREN, "}"),
		tok(EOF, ""),
	}
	tokens := lex(t, "x_{0}")
	if len(tokens) != len(want) {
		t.Fatalf("lex(%q) = %v (len %d), want len %d", "x_{0}", tokens, len(tokens), len(want))
	}
	for i, w := range want {
		if tokens[i] != w {
			t.Errorf("lex(%q)[%d] = %v, want %v", "x_{0}", i, tokens[i], w)
		}
	}
}

// TestUnderscoreBraceExpression — x_{n+1} → full group expansion
func TestUnderscoreBraceExpression(t *testing.T) {
	want := []Token{
		tok(SYMBOL, "x"),
		tok(UNDERSCORE, "_"),
		tok(LPAREN, "{"),
		tok(SYMBOL, "n"),
		tok(PLUS, "+"),
		tok(NUMBER, "1"),
		tok(RPAREN, "}"),
		tok(EOF, ""),
	}
	tokens := lex(t, "x_{n+1}")
	if len(tokens) != len(want) {
		t.Fatalf("lex(%q) = %v (len %d), want len %d", "x_{n+1}", tokens, len(tokens), len(want))
	}
	for i, w := range want {
		if tokens[i] != w {
			t.Errorf("lex(%q)[%d] = %v, want %v", "x_{n+1}", i, tokens[i], w)
		}
	}
}

// ── §3.2 EQUALS token ─────────────────────────────────────────────────────────

// TestEqualsAlone — bare = → EQUALS with value "="
func TestEqualsAlone(t *testing.T) {
	tokens := lex(t, "=")
	if tokens[0] != tok(EQUALS, "=") {
		t.Errorf(`lex("=")[0] = %v, want EQUALS("=")`, tokens[0])
	}
}

// TestEqualsInContext — i=0 → SYMBOL(i) EQUALS NUMBER(0) EOF
func TestEqualsInContext(t *testing.T) {
	want := []Token{
		tok(SYMBOL, "i"),
		tok(EQUALS, "="),
		tok(NUMBER, "0"),
		tok(EOF, ""),
	}
	tokens := lex(t, "i=0")
	if len(tokens) != len(want) {
		t.Fatalf("lex(%q) = %v (len %d), want len %d", "i=0", tokens, len(tokens), len(want))
	}
	for j, w := range want {
		if tokens[j] != w {
			t.Errorf("lex(%q)[%d] = %v, want %v", "i=0", j, tokens[j], w)
		}
	}
}

// ── §3.3 NORM token ───────────────────────────────────────────────────────────

// TestLVert — \lVert → NORM (must NOT be remapped to PIPE like \lvert)
func TestLVert(t *testing.T) {
	tokens := lex(t, `\lVert`)
	if tokens[0].Type != NORM {
		t.Errorf(`lex("\\lVert")[0].Type = %v, want NORM`, tokens[0].Type)
	}
}

// TestRVert — \rVert → NORM
func TestRVert(t *testing.T) {
	tokens := lex(t, `\rVert`)
	if tokens[0].Type != NORM {
		t.Errorf(`lex("\\rVert")[0].Type = %v, want NORM`, tokens[0].Type)
	}
}

// TestNormDistinctFromPipe — \lVert produces NORM, \lvert produces PIPE (not the same)
func TestNormDistinctFromPipe(t *testing.T) {
	normTok := lex(t, `\lVert`)[0]
	pipeTok := lex(t, `\lvert`)[0]
	if normTok.Type == pipeTok.Type {
		t.Errorf(`\lVert and \lvert must produce different token types, both = %v`, normTok.Type)
	}
	if normTok.Type != NORM {
		t.Errorf(`\lVert token type = %v, want NORM`, normTok.Type)
	}
	if pipeTok.Type != PIPE {
		t.Errorf(`\lvert token type = %v, want PIPE`, pipeTok.Type)
	}
}

// TestNormSequence — \lVert x \rVert → NORM SYMBOL(x) NORM EOF
func TestNormSequence(t *testing.T) {
	tokens := lex(t, `\lVert x \rVert`)
	if len(tokens) != 4 {
		t.Fatalf(`lex("\\lVert x \\rVert") = %v (len %d), want len 4`, tokens, len(tokens))
	}
	if tokens[0].Type != NORM {
		t.Errorf("tokens[0].Type = %v, want NORM", tokens[0].Type)
	}
	if tokens[1] != tok(SYMBOL, "x") {
		t.Errorf("tokens[1] = %v, want SYMBOL(x)", tokens[1])
	}
	if tokens[2].Type != NORM {
		t.Errorf("tokens[2].Type = %v, want NORM", tokens[2].Type)
	}
	if tokens[3].Type != EOF {
		t.Errorf("tokens[3].Type = %v, want EOF", tokens[3].Type)
	}
}

// ── §3 Large-op full token stream ───────────────────────────────────────────────

// TestLargeOpTokenStream — \sum_{i=0}^{3} i → complete expected stream
func TestLargeOpTokenStream(t *testing.T) {
	want := []Token{
		tok(COMMAND, "sum"),
		tok(UNDERSCORE, "_"),
		tok(LPAREN, "{"),
		tok(SYMBOL, "i"),
		tok(EQUALS, "="),
		tok(NUMBER, "0"),
		tok(RPAREN, "}"),
		tok(POWER, "^"),
		tok(LPAREN, "{"),
		tok(NUMBER, "3"),
		tok(RPAREN, "}"),
		tok(SYMBOL, "i"),
		tok(EOF, ""),
	}
	tokens := lex(t, `\sum_{i=0}^{3} i`)
	if len(tokens) != len(want) {
		t.Fatalf("lex(%q) = %v (len %d), want len %d",
			`\sum_{i=0}^{3} i`, tokens, len(tokens), len(want))
	}
	for j, w := range want {
		if tokens[j] != w {
			t.Errorf("token[%d] = %v, want %v", j, tokens[j], w)
		}
	}
}
