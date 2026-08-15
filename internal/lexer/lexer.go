package lexer

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// commandArities maps normalised command names (backslash-stripped, lowercased)
// to the number of char-mode arguments they consume.
// See spec §6.4.
var commandArities = map[string]int{
	// Arithmetic
	"frac": 2, "dfrac": 2, "tfrac": 2, "cfrac": 2,
	// sqrt arity 0 in lexer — parser owns [n] and radicand (parser-extensions §4.1)
	// Trigonometric
	"sin": 1, "cos": 1, "tan": 1,
	"sec": 1, "csc": 1, "cot": 1,
	// Inverse trigonometric (evaluatex abbreviations — ADR-003; canonical arcX also supported)
	"asin": 1, "acos": 1, "atan": 1,
	"asec": 1, "acsc": 1, "acot": 1,
	"arcsin": 1, "arccos": 1, "arctan": 1,
	// Logarithmic and exponential
	"ln": 1, "exp": 1,
	// log arity 0 — parser handles optional _{b} and value (parser-extensions §6.1)
	// Hyperbolic
	"sinh": 1, "cosh": 1, "tanh": 1,
	"coth": 1, "sech": 1, "csch": 1,
	// Combinatorics
	"binom": 2, "dbinom": 2, "tbinom": 2,
}

type scanPattern struct {
	typ TokenType
	re  *regexp.Regexp
}

// patterns is the ordered pattern table from spec §5.2.
// Order is significant: first match wins.
var patterns = []scanPattern{
	// Priority 1: \left( \left[ \left{ or bare (, [, {  — spec §5.2
	{LPAREN, regexp.MustCompile(`^(?:\\left[(\[{]|[(\[{])`)},
	// Priority 2: \right) \right] \right} or bare ), ], }
	{RPAREN, regexp.MustCompile(`^(?:\\right[)\]}]|[)\]}])`)},
	{PLUS, regexp.MustCompile(`^\+`)},
	{MINUS, regexp.MustCompile(`^-`)},
	{TIMES, regexp.MustCompile(`^\*`)},
	{DIVIDE, regexp.MustCompile(`^/`)},
	{POWER, regexp.MustCompile(`^\^`)},
	{PIPE, regexp.MustCompile(`^\|`)},
	{BANG, regexp.MustCompile(`^!`)},
	{COMMA, regexp.MustCompile(`^,`)},
	{UNDERSCORE, regexp.MustCompile(`^_`)},
	{EQUALS, regexp.MustCompile(`^=`)},
	// Priority 11a: NORM — \lVert / \rVert (capital V) before generic COMMAND.
	// \lvert / \rvert (lowercase) fall through to COMMAND and are remapped to PIPE.
	{NORM, regexp.MustCompile(`^\\(?:lVert|rVert)`)},
	// Priority 11b: FLOOR / CEIL before COMMAND (parser-extensions §3.1).
	{FLOOR, regexp.MustCompile(`^\\(?:lfloor|rfloor)`)},
	{CEIL, regexp.MustCompile(`^\\(?:lceil|rceil)`)},
	// Priority 12: COMMAND before SYMBOL so \sin isn't SYMBOL(sin) — spec §5.2
	{COMMAND, regexp.MustCompile(`^\\[A-Za-z]+`)},
	// Priority 13: SYMBOL — spec §5.2 and ADR-013 (no underscore in names)
	{SYMBOL, regexp.MustCompile(`^[A-Za-z][A-Za-z0-9]*`)},
	// Priority 14: NUMBER — spec §5.2
	{NUMBER, regexp.MustCompile(`^\d+(?:\.\d+)?(?:[eE][+\-]?\d+)?`)},
}

type lexer struct {
	input  string
	pos    int
	tokens []Token
}

// Lex tokenises input and returns a flat []Token terminated by an EOF token.
//
// It implements spec §5 (pattern matching), §6 (char-mode),
// §7 (post-lex remapping), and §8 (EOF sentinel).
func Lex(input string) ([]Token, error) {
	l := &lexer{input: input}
	if err := l.lexExpression(false, false); err != nil {
		return nil, err
	}

	// §7 Post-lex remapping (three passes, applied in one scan).
	l.tokens = postLexRemap(l.tokens)

	// §8 EOF sentinel.
	l.tokens = append(l.tokens, Token{Type: EOF, Pos: len(input)})
	return l.tokens, nil
}

// lexExpression lexes the input recursively.
//
// When charMode is true the function reads exactly one char-mode token
// (spec §6.1) then returns.  When charMode is false it reads until the buffer
// is empty or a closing "}" ends a group (spec §6.3).
//
// inGroup must be true when the call was triggered by a "{" opener; in that
// case reaching EOF without a matching "}" is an error (spec §9).
//
// Algorithm mirrors evaluatex lexExpression(charMode) from
// evaluatex/src/lexer.js, elevated to normative in spec §6.
func (l *lexer) lexExpression(charMode, inGroup bool) error {
	for {
		// Skip whitespace before each token (spec §5.3). We do this here—not in
		// next()—so trailing spaces do not cause a spurious end-of-input error.
		l.skipWhitespace()
		if l.pos >= len(l.input) {
			return l.errAtEOF(charMode, inGroup)
		}

		tok, err := l.readToken(charMode)
		if err != nil {
			return err
		}
		l.tokens = append(l.tokens, tok)

		if err := l.expandArgs(tok); err != nil {
			return err
		}

		// Exit conditions (spec §6):
		//   - char mode: return after emitting exactly one token (plus its args).
		//   - group mode: return after the closing } is emitted.
		if charMode || isGroupCloser(tok) {
			return nil
		}
	}
}

// errAtEOF handles end-of-input for lexExpression (spec §9).
// Top-level EOF is success (nil); char-mode and unclosed groups are errors.
func (l *lexer) errAtEOF(charMode, inGroup bool) error {
	if charMode {
		return fmt.Errorf("lexer: unexpected end of input at position %d", l.pos)
	}
	if inGroup {
		return fmt.Errorf("lexer: unexpected end of input: unclosed '{' group at position %d", l.pos)
	}
	return nil
}

// readToken selects the next token using char-mode or normal scanning (spec §6.1).
func (l *lexer) readToken(charMode bool) (Token, error) {
	if charMode {
		return l.nextCharToken()
	}
	return l.next()
}

// expandArgs consumes char-mode / group arguments required by tok (spec §6.2–§6.3).
func (l *lexer) expandArgs(tok Token) error {
	switch tok.Type {
	case POWER:
		// §6.2: POWER reads exactly 1 char-mode argument.
		return l.lexExpression(true, false)

	case COMMAND:
		// §6.2: COMMAND reads N char-mode arguments where N is the arity.
		// Normalise now (strip \, lowercase) to look up the arity table.
		// Final normalisation of the token value happens in postLexRemap (§7.3),
		// but we need the name here too.
		name := strings.ToLower(tok.Value[1:])     // strip leading backslash
		return l.lexCharArgs(commandArities[name]) // 0 if not in table — spec §6.4

	case LPAREN:
		if tok.Value == "{" {
			// §6.3: { opens a group — recurse in normal mode until matching }.
			// inGroup=true so EOF inside returns an error (spec §9).
			return l.lexExpression(false, true)
		}
		// Other LPAREN values ((, [, \left() are not group openers for
		// the purposes of char-mode expansion.
	}
	return nil
}

// lexCharArgs reads n successive char-mode arguments (spec §6.2).
func (l *lexer) lexCharArgs(n int) error {
	for i := 0; i < n; i++ {
		if err := l.lexExpression(true, false); err != nil {
			return err
		}
	}
	return nil
}

// isGroupCloser reports whether tok closes a {…} group (spec §6.3).
func isGroupCloser(tok Token) bool {
	return tok.Type == RPAREN && tok.Value == "}"
}

// next reads the next non-whitespace token from the remaining input.
func (l *lexer) next() (Token, error) {
	return l.nextN(-1)
}

// nextCharToken implements the char-mode read rule from spec §6.1:
//   - If the next non-whitespace character is \, read a full COMMAND token.
//   - Otherwise read exactly one UTF-8 rune (spec §6.1.3).
func (l *lexer) nextCharToken() (Token, error) {
	l.skipWhitespace()
	if l.pos < len(l.input) && l.input[l.pos] == '\\' {
		return l.next()
	}
	_, size := utf8.DecodeRuneInString(l.input[l.pos:])
	return l.nextN(size)
}

// nextN reads the next token from a window of at most maxLen bytes.
// maxLen < 0 means no length limit (greedy).
func (l *lexer) nextN(maxLen int) (Token, error) {
	l.skipWhitespace()
	if l.pos >= len(l.input) {
		return Token{}, fmt.Errorf("lexer: unexpected end of input at position %d", l.pos)
	}

	end := len(l.input)
	if maxLen > 0 && l.pos+maxLen < end {
		end = l.pos + maxLen
	}
	sub := l.input[l.pos:end]
	startPos := l.pos

	for _, p := range patterns {
		loc := p.re.FindStringIndex(sub)
		if loc != nil { // anchor guarantees loc[0] == 0
			tok := Token{
				Type:  p.typ,
				Value: l.input[startPos : startPos+loc[1]],
				Pos:   startPos,
			}
			l.pos += loc[1]
			return tok, nil
		}
	}

	return Token{}, fmt.Errorf(
		"lexer: unexpected character %q at position %d",
		string(l.input[l.pos]), l.pos,
	)
}

// skipWhitespace advances pos past any ASCII whitespace (spec §5.3).
func (l *lexer) skipWhitespace() {
	for l.pos < len(l.input) {
		switch l.input[l.pos] {
		case ' ', '\t', '\n', '\r':
			l.pos++
		default:
			return
		}
	}
}

// postLexRemap applies the three post-lex passes from spec §7 in a single scan.
//
//  1. Drop standalone \left / \right COMMAND tokens (§7.1).
//  2. Remap operator COMMAND tokens to their target types (§7.2).
//  3. Normalise remaining COMMAND values: strip backslash, lowercase (§7.3).
func postLexRemap(in []Token) []Token {
	out := make([]Token, 0, len(in))
	for _, tok := range in {
		if tok.Type == FLOOR || tok.Type == CEIL {
			// parser-extensions §3.1: strip leading backslash from delimiter values.
			tok.Value = strings.TrimPrefix(tok.Value, `\`)
			out = append(out, tok)
			continue
		}
		if tok.Type != COMMAND {
			out = append(out, tok)
			continue
		}

		// Normalise name for matching (§7.3 applied inline).
		name := strings.ToLower(tok.Value[1:]) // strip backslash

		switch name {
		case "left", "right":
			// §7.1: drop standalone \left / \right decorators.

		case "times", "cdot":
			// §7.2: multiplication remaps.
			out = append(out, Token{Type: TIMES, Value: "*", Pos: tok.Pos})

		case "div":
			// §7.2: division remap.
			out = append(out, Token{Type: DIVIDE, Value: "/", Pos: tok.Pos})

		case "lvert", "rvert":
			// §7.2: absolute-value pipe remaps.
			out = append(out, Token{Type: PIPE, Value: "|", Pos: tok.Pos})

		default:
			// §7.3: normalise value; keep as COMMAND.
			tok.Value = name
			out = append(out, tok)
		}
	}
	return out
}
