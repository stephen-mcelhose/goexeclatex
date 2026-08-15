package lexer

import (
	"fmt"
	"regexp"
	"strings"
)

// commandArities maps normalised command names (backslash-stripped, lowercased)
// to the number of char-mode arguments they consume.
// See spec §6.4.
var commandArities = map[string]int{
	"frac":  2,
	"dfrac": 2,
	"tfrac": 2,
	"cfrac": 2,
	"sqrt":  1,
	"sin":   1,
	"cos":   1,
	"tan":   1,
	"asin":  1,
	"acos":  1,
	"atan":  1,
	"sec":   1,
	"csc":   1,
	"cot":   1,
	"asec":  1,
	"acsc":  1,
	"acot":  1,
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
	// Priority 11: COMMAND before SYMBOL so \sin isn't SYMBOL(sin) — spec §5.2
	{COMMAND, regexp.MustCompile(`^\\[A-Za-z]+`)},
	{SYMBOL, regexp.MustCompile(`^[A-Za-z_][A-Za-z_0-9]*`)},
	// Priority 13: NUMBER — spec §5.2
	{NUMBER, regexp.MustCompile(`^\d+(?:\.\d+)?(?:[eE][+\-]?\d+)?`)},
}

type lexer struct {
	input  string
	pos    int
	tokens []Token
}

// Lex tokenises input and returns a flat []Token terminated by an EOF token.
//
// It implements spec §4 (pre-pass), §5 (pattern matching), §6 (char-mode),
// §7 (post-lex remapping), and §8 (EOF sentinel).
func Lex(input string) ([]Token, error) {
	// §4 Pre-pass: \_ → literal underscore.
	input = strings.ReplaceAll(input, `\_`, "_")

	l := &lexer{input: input}
	if err := l.lexExpression(false); err != nil {
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
// Algorithm mirrors evaluatex lexExpression(charMode) from
// evaluatex/src/lexer.js, elevated to normative in spec §6.
func (l *lexer) lexExpression(charMode bool) error {
	for {
		// Skip whitespace before checking whether there is anything left,
		// so trailing spaces do not cause a spurious end-of-input error.
		l.skipWhitespace()
		if l.pos >= len(l.input) {
			if charMode {
				return fmt.Errorf("lexer: unexpected end of input at position %d", l.pos)
			}
			break
		}

		var (
			tok Token
			err error
		)
		if charMode {
			tok, err = l.nextCharToken()
		} else {
			tok, err = l.next()
		}
		if err != nil {
			return err
		}
		l.tokens = append(l.tokens, tok)

		switch tok.Type {
		case POWER:
			// §6.2: POWER reads exactly 1 char-mode argument.
			if err := l.lexExpression(true); err != nil {
				return err
			}

		case COMMAND:
			// §6.2: COMMAND reads N char-mode arguments where N is the arity.
			// Normalise now (strip \, lowercase) to look up the arity table.
			// Final normalisation of the token value happens in postLexRemap (§7.3),
			// but we need the name here too.
			name := strings.ToLower(tok.Value[1:]) // strip leading backslash
			arity := commandArities[name]           // 0 if not in table — spec §6.4
			for i := 0; i < arity; i++ {
				if err := l.lexExpression(true); err != nil {
					return err
				}
			}

		case LPAREN:
			if tok.Value == "{" {
				// §6.3: { opens a group — recurse in normal mode until matching }.
				if err := l.lexExpression(false); err != nil {
					return err
				}
			}
			// Other LPAREN values ((, [, \left() are not group openers for
			// the purposes of char-mode expansion.
		}

		// Exit conditions (spec §6):
		//   - char mode: return after emitting exactly one token (plus its args).
		//   - group mode: return after the closing } is emitted.
		if charMode {
			return nil
		}
		if tok.Type == RPAREN && tok.Value == "}" {
			return nil
		}
	}
	return nil
}

// next reads the next non-whitespace token from the remaining input.
func (l *lexer) next() (Token, error) {
	return l.nextN(-1)
}

// nextCharToken implements the char-mode read rule from spec §6.1:
//   - If the next non-whitespace character is \, read a full COMMAND token.
//   - Otherwise read exactly one byte.
func (l *lexer) nextCharToken() (Token, error) {
	l.skipWhitespace()
	if l.pos < len(l.input) && l.input[l.pos] == '\\' {
		return l.next()
	}
	return l.nextN(1)
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
