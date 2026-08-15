// Package lexer implements the LaTeX math lexer as specified in wiki/specs/lexer.md.
package lexer

// TokenType identifies the category of a lexical token.
// See spec §3 for the full token type table.
type TokenType int

const (
	NUMBER     TokenType = iota // spec §3: \d+(\.\d+)?([eE][+-]?\d+)?
	SYMBOL                      // spec §3: [A-Za-z_][A-Za-z_0-9]*
	COMMAND                     // spec §3: \\[A-Za-z]+ normalised per §7.3
	PLUS                        // spec §3: +
	MINUS                       // spec §3: -
	TIMES                       // spec §3: * or remapped from \times, \cdot
	DIVIDE                      // spec §3: / or remapped from \div
	POWER                       // spec §3: ^
	LPAREN                      // spec §3: (, [, { and \left variants
	RPAREN                      // spec §3: ), ], } and \right variants
	PIPE                        // spec §3: | and remapped from \lvert, \rvert
	BANG                        // spec §3: !
	COMMA                       // spec §3: ,
	UNDERSCORE                  // spec §3: _ (Tier 0.3 — subscript)
	EQUALS                      // spec §3.2: = (big-op bound only)
	NORM                        // spec §3.3: \lVert / \rVert
	EOF                         // spec §8: synthetic end-of-input sentinel
)

var tokenTypeNames = [...]string{
	"NUMBER", "SYMBOL", "COMMAND", "PLUS", "MINUS", "TIMES", "DIVIDE",
	"POWER", "LPAREN", "RPAREN", "PIPE", "BANG", "COMMA", "UNDERSCORE",
	"EQUALS", "NORM", "EOF",
}

// String returns the name of the token type.
func (t TokenType) String() string {
	if int(t) >= 0 && int(t) < len(tokenTypeNames) {
		return tokenTypeNames[t]
	}
	return "UNKNOWN"
}

// Token is a single lexical unit produced by the lexer.
//
// Value contains the raw matched text for NUMBER, SYMBOL, LPAREN, RPAREN;
// for COMMAND it contains the normalised (backslash-stripped, lowercased) name
// per spec §7.3.
//
// Pos is the zero-based byte offset of the token in the original input.
type Token struct {
	Type  TokenType
	Value string
	Pos   int
}

// String returns a compact human-readable representation of the token.
func (t Token) String() string {
	switch t.Type {
	case PLUS, MINUS, TIMES, DIVIDE, POWER, PIPE, BANG, COMMA, UNDERSCORE, EQUALS, NORM, EOF:
		return t.Type.String()
	default:
		return t.Type.String() + "(" + t.Value + ")"
	}
}
