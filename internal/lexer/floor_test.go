package lexer

import "testing"

// Tests for docs/specs/parser-extensions.md §3.1 (FLOOR / CEIL tokens).

func TestLfloorToken(t *testing.T) {
	tokens := lex(t, `\lfloor`)
	if tokens[0].Type != FLOOR {
		t.Errorf(`lex("\\lfloor")[0].Type = %v, want FLOOR`, tokens[0].Type)
	}
	if tokens[0].Value != "lfloor" {
		t.Errorf(`lex("\\lfloor")[0].Value = %q, want "lfloor"`, tokens[0].Value)
	}
}

func TestRfloorToken(t *testing.T) {
	tokens := lex(t, `\rfloor`)
	if tokens[0].Type != FLOOR || tokens[0].Value != "rfloor" {
		t.Errorf(`lex("\\rfloor")[0] = %v, want FLOOR("rfloor")`, tokens[0])
	}
}

func TestLceilToken(t *testing.T) {
	tokens := lex(t, `\lceil`)
	if tokens[0].Type != CEIL || tokens[0].Value != "lceil" {
		t.Errorf(`lex("\\lceil")[0] = %v, want CEIL("lceil")`, tokens[0])
	}
}

func TestRceilToken(t *testing.T) {
	tokens := lex(t, `\rceil`)
	if tokens[0].Type != CEIL || tokens[0].Value != "rceil" {
		t.Errorf(`lex("\\rceil")[0] = %v, want CEIL("rceil")`, tokens[0])
	}
}

func TestFloorSequence(t *testing.T) {
	tokens := lex(t, `\lfloor 3.2 \rfloor`)
	if len(tokens) != 4 {
		t.Fatalf(`lex len = %d, want 4: %v`, len(tokens), tokens)
	}
	if tokens[0].Type != FLOOR || tokens[0].Value != "lfloor" {
		t.Errorf("tokens[0] = %v, want FLOOR(lfloor)", tokens[0])
	}
	if tokens[1].Type != NUMBER {
		t.Errorf("tokens[1].Type = %v, want NUMBER", tokens[1].Type)
	}
	if tokens[2].Type != FLOOR || tokens[2].Value != "rfloor" {
		t.Errorf("tokens[2] = %v, want FLOOR(rfloor)", tokens[2])
	}
}

func TestFloorNotCommand(t *testing.T) {
	tok := lex(t, `\lfloor`)[0]
	if tok.Type == COMMAND {
		t.Errorf(`\lfloor must not be COMMAND, got %v`, tok)
	}
}
