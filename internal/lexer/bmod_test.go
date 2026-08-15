package lexer

import "testing"

func TestBmodRemap(t *testing.T) {
	tok := lex(t, `\bmod`)[0]
	if tok.Type != BMOD || tok.Value != "bmod" {
		t.Errorf(`lex("\\bmod")[0] = %v, want BMOD(bmod)`, tok)
	}
}

func TestPmodNotBmod(t *testing.T) {
	tok := lex(t, `\pmod`)[0]
	if tok.Type == BMOD {
		t.Errorf(`\pmod must not remap to BMOD`)
	}
	if tok.Type != COMMAND || tok.Value != "pmod" {
		t.Errorf(`lex("\\pmod")[0] = %v, want COMMAND(pmod)`, tok)
	}
}
