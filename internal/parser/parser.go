package parser

import (
	"fmt"
	"strconv"

	"github.com/stephen-mcelhose/goexeclatex/internal/lexer"
)

// commandArities maps normalised command names to their argument counts (spec §6.1).
// Commands not in this table have arity 0 and produce a SymbolNode.
var commandArities = map[string]int{
	// Arithmetic
	"frac": 2, "dfrac": 2, "tfrac": 2, "cfrac": 2,
	"sqrt": 1,
	// Trigonometric
	"sin": 1, "cos": 1, "tan": 1,
	"sec": 1, "csc": 1, "cot": 1,
	// Inverse trigonometric (abbreviations — ADR-003; canonical arcX also supported)
	"asin": 1, "acos": 1, "atan": 1,
	"asec": 1, "acsc": 1, "acot": 1,
	"arcsin": 1, "arccos": 1, "arctan": 1,
	// Logarithmic and exponential
	"ln": 1, "log": 1, "exp": 1,
	// Hyperbolic
	"sinh": 1, "cosh": 1, "tanh": 1,
	"coth": 1, "sech": 1, "csch": 1,
	// Combinatorics
	"binom": 2, "dbinom": 2, "tbinom": 2,
}

// Parse transforms a lexer token stream into an AST (spec §9).
func Parse(tokens []lexer.Token) (Node, error) {
	if len(tokens) == 0 || tokens[len(tokens)-1].Type != lexer.EOF {
		return nil, fmt.Errorf("parser: invalid token stream")
	}
	p := &parser{tokens: tokens}
	node, err := p.parseSum()
	if err != nil {
		return nil, err
	}
	if p.peek().Type != lexer.EOF {
		tok := p.peek()
		return nil, fmt.Errorf("parser: unexpected %s %q at position %d", tok.Type, tok.Value, tok.Pos)
	}
	return node, nil
}

// parser holds the token stream and current cursor position.
// absDepth tracks how many |…| absolute-value groups are currently open,
// so the product rule does not greedily consume a closing | as an implicit-
// multiply trigger (spec §4.2 ambiguity note).
type parser struct {
	tokens    []lexer.Token
	pos       int
	absDepth  int // counts nesting of |…| groups
	normDepth int // counts nesting of ‖…‖ groups
}

func (p *parser) peek() lexer.Token {
	return p.tokens[p.pos]
}

func (p *parser) consume() lexer.Token {
	tok := p.tokens[p.pos]
	p.pos++
	return tok
}

// expect consumes the next token if it matches the given type and value.
// Pass value="" to match any value for that type.
func (p *parser) expect(typ lexer.TokenType, value string) (lexer.Token, error) {
	tok := p.peek()
	if tok.Type != typ {
		if tok.Type == lexer.EOF {
			return tok, fmt.Errorf("parser: unexpected EOF")
		}
		return tok, fmt.Errorf("parser: unexpected %s %q at position %d", tok.Type, tok.Value, tok.Pos)
	}
	if value != "" && tok.Value != value {
		return tok, fmt.Errorf("parser: expected %q at position %d, got %q", value, tok.Pos, tok.Value)
	}
	return p.consume(), nil
}

// closerFor returns the expected closing bracket character for a given opener
// token value. The lexer may emit full forms like `\left(` or bare `(`;
// we key on the last character only (always the actual bracket).
func closerFor(opener string) string {
	if len(opener) == 0 {
		return ""
	}
	switch opener[len(opener)-1] {
	case '(':
		return ")"
	case '[':
		return "]"
	case '{':
		return "}"
	}
	return ""
}

// canBeginPower reports whether tok can start a power expression and therefore
// trigger implicit multiplication (spec §5). PIPE is only allowed when we are
// not already inside an |…| group — otherwise the closing | would be consumed
// as the start of a new implicit-multiply abs-value.
func (p *parser) canBeginPower(tok lexer.Token) bool {
	switch tok.Type {
	case lexer.NUMBER, lexer.SYMBOL, lexer.COMMAND, lexer.LPAREN:
		return true
	case lexer.NORM:
		return p.normDepth == 0
	case lexer.PIPE:
		return p.absDepth == 0
	}
	return false
}

// ── Grammar rules ─────────────────────────────────────────────────────────────

// parseSum → parseProduct { (PLUS | MINUS) parseProduct }
func (p *parser) parseSum() (Node, error) {
	left, err := p.parseProduct()
	if err != nil {
		return nil, err
	}
	for {
		tok := p.peek()
		if tok.Type != lexer.PLUS && tok.Type != lexer.MINUS {
			break
		}
		p.consume()
		op := tok.Value
		right, err := p.parseProduct()
		if err != nil {
			return nil, err
		}
		left = &BinaryNode{Op: op, Left: left, Right: right}
	}
	return left, nil
}

// parseProduct → parsePower { ('*' | '/') parsePower | <implicit multiply> }
func (p *parser) parseProduct() (Node, error) {
	left, err := p.parsePower()
	if err != nil {
		return nil, err
	}
	for {
		tok := p.peek()
		if tok.Type == lexer.TIMES || tok.Type == lexer.DIVIDE {
			p.consume()
			right, err := p.parsePower()
			if err != nil {
				return nil, err
			}
			left = &BinaryNode{Op: tok.Value, Left: left, Right: right}
			continue
		}
		// Implicit multiply (spec §5, ADR-006)
		if p.canBeginPower(tok) {
			right, err := p.parsePower()
			if err != nil {
				return nil, err
			}
			left = &BinaryNode{Op: "*", Left: left, Right: right}
			continue
		}
		break
	}
	return left, nil
}

// parsePower → parseUnary [POWER parseSuperArg]  (right-associative)
func (p *parser) parsePower() (Node, error) {
	base, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	if p.peek().Type == lexer.POWER {
		p.consume()
		// Use parseSuperArg (not recursive parsePower) so that a bare subscript
		// following the exponent is not silently absorbed into the exponent.
		// e.g. x^2_i → exp=2, then peek=UNDERSCORE → error (spec §4.1).
		exp, err := p.parseSuperArg()
		if err != nil {
			return nil, err
		}
		if p.peek().Type == lexer.UNDERSCORE {
			return nil, fmt.Errorf("parser: superscript before subscript is not allowed; write x_{i}^2, not x^2_i")
		}
		return &BinaryNode{Op: "^", Left: base, Right: exp}, nil
	}
	return base, nil
}

// parseSuperArg parses a superscript argument (spec §4.1), right-associatively.
//
// A brace group `{…}` is parsed as a full expression (subscripts allowed inside).
// A bare token chain is parsed without attaching subscripts at this level —
// this is what lets parsePower detect the x^2_i order violation.
//
// Right-associativity is preserved: 2^3^4 → parseSuperArg sees 3, then recursively
// handles ^4, yielding BinaryNode(^, 3, 4) as the exponent of 2.
func (p *parser) parseSuperArg() (Node, error) {
	if p.peek().Type == lexer.LPAREN && p.peek().Value == "{" {
		p.consume() // consume '{'
		arg, err := p.parseSum()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(lexer.RPAREN, "}"); err != nil {
			return nil, fmt.Errorf("parser: superscript {…} missing closing }")
		}
		return arg, nil
	}
	// Bare: single atom (no subscript), then optionally chain more ^ (right-recursive).
	base, err := p.parseAtom()
	if err != nil {
		return nil, err
	}
	if p.peek().Type == lexer.POWER {
		p.consume()
		exp, err := p.parseSuperArg()
		if err != nil {
			return nil, err
		}
		if p.peek().Type == lexer.UNDERSCORE {
			return nil, fmt.Errorf("parser: superscript before subscript is not allowed; write x_{i}^2, not x^2_i")
		}
		return &BinaryNode{Op: "^", Left: base, Right: exp}, nil
	}
	return base, nil
}

// parseUnary → MINUS parsePower | parsePostfix
func (p *parser) parseUnary() (Node, error) {
	if p.peek().Type == lexer.MINUS {
		p.consume()
		operand, err := p.parsePower()
		if err != nil {
			return nil, err
		}
		return &UnaryNode{Op: "-", Operand: operand}, nil
	}
	return p.parsePostfix()
}

// parsePostfix → parseAtom [UNDERSCORE sub_arg] [BANG]
func (p *parser) parsePostfix() (Node, error) {
	base, err := p.parseAtom()
	if err != nil {
		return nil, err
	}
	// spec §4.1: optional subscript (binds tighter than power).
	if p.peek().Type == lexer.UNDERSCORE {
		p.consume()
		sub, err := p.parseSubArg()
		if err != nil {
			return nil, err
		}
		if p.peek().Type == lexer.UNDERSCORE {
			return nil, fmt.Errorf("parser: chained subscripts are not allowed: x_{i}_{j}")
		}
		base = &SubscriptNode{Base: base, Sub: sub}
	}
	if p.peek().Type == lexer.BANG {
		p.consume()
		return &UnaryNode{Op: "!", Operand: base}, nil
	}
	return base, nil
}

// parseSubArg parses one subscript argument.
// If the next token is LPAREN('{'), it parses the full brace group.
// Otherwise it parses a single atom (char-mode).
func (p *parser) parseSubArg() (Node, error) {
	if p.peek().Type == lexer.LPAREN && p.peek().Value == "{" {
		p.consume() // consume '{'
		arg, err := p.parseSum()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(lexer.RPAREN, "}"); err != nil {
			return nil, fmt.Errorf("parser: subscript {…} missing closing }")
		}
		return arg, nil
	}
	return p.parseAtom()
}

// parseAtom → NUMBER | SYMBOL | COMMAND args | LPAREN sum RPAREN | PIPE sum PIPE
func (p *parser) parseAtom() (Node, error) {
	tok := p.peek()

	switch tok.Type {
	case lexer.NUMBER:
		p.consume()
		v, err := strconv.ParseFloat(tok.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("parser: invalid number %q at position %d", tok.Value, tok.Pos)
		}
		return &NumberNode{Value: v, Pos: tok.Pos}, nil

	case lexer.SYMBOL:
		p.consume()
		return &SymbolNode{Name: tok.Value, Pos: tok.Pos}, nil

	case lexer.COMMAND:
		return p.parseCommand()

	case lexer.LPAREN:
		return p.parseGroup()

	case lexer.PIPE:
		return p.parseAbsValue()

	case lexer.NORM:
		return p.parseNorm()

	case lexer.EOF:
		return nil, fmt.Errorf("parser: unexpected EOF")

	default:
		return nil, fmt.Errorf("parser: unexpected %s %q at position %d", tok.Type, tok.Value, tok.Pos)
	}
}

// parseCommand handles a COMMAND token and its arity-N arguments (spec §6).
func (p *parser) parseCommand() (Node, error) {
	tok := p.consume() // consume COMMAND
	name := tok.Value  // already normalised by lexer (no backslash, lowercased)

	// spec §4.2: large operators require their own grammar.
	if name == "sum" || name == "prod" {
		return p.parseLargeOp(name)
	}

	arity := commandArities[name]

	if arity == 0 {
		return &SymbolNode{Name: name, Pos: tok.Pos}, nil
	}

	args := make([]Node, 0, arity)
	for i := 0; i < arity; i++ {
		arg, err := p.parseCommandArg(name)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	return &FunctionNode{Name: name, Args: args, Pos: tok.Pos}, nil
}

// parseCommandArg parses one argument for a command (spec §6.2).
// If the next token is LPAREN('{'), it consumes the brace group.
// Otherwise (char-mode single-token arg) it calls parseAtom.
func (p *parser) parseCommandArg(name string) (Node, error) {
	if p.peek().Type == lexer.LPAREN && p.peek().Value == "{" {
		p.consume() // consume '{'
		arg, err := p.parseSum()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(lexer.RPAREN, "}"); err != nil {
			return nil, fmt.Errorf("parser: %s expects {…} argument, missing closing }", name)
		}
		return arg, nil
	}
	// Char-mode fallback: single atom without braces
	if !p.canBeginPower(p.peek()) {
		tok := p.peek()
		return nil, fmt.Errorf("parser: %s expects {…} argument at position %d", name, tok.Pos)
	}
	return p.parseAtom()
}

// parseGroup handles LPAREN sum RPAREN — transparent grouping (spec §4.3, ADR-005).
func (p *parser) parseGroup() (Node, error) {
	open := p.consume() // consume LPAREN
	closer := closerFor(open.Value)

	inner, err := p.parseSum()
	if err != nil {
		return nil, err
	}

	tok := p.peek()
	if tok.Type != lexer.RPAREN {
		if tok.Type == lexer.EOF {
			return nil, fmt.Errorf("parser: unexpected EOF")
		}
		return nil, fmt.Errorf("parser: expected %q at position %d, got %q", closer, tok.Pos, tok.Value)
	}
	// The RPAREN value may be a full form like `\right)` — compare on last char only.
	if len(tok.Value) == 0 || string(tok.Value[len(tok.Value)-1]) != closer {
		return nil, fmt.Errorf("parser: expected %q at position %d, got %q", closer, tok.Pos, tok.Value)
	}
	p.consume() // consume RPAREN
	return inner, nil
}

// parseAbsValue handles PIPE sum PIPE → FunctionNode("abs") (spec §4.2).
// absDepth is incremented while parsing the inner expression so that the
// closing | is not mistakenly consumed as an implicit-multiply trigger.
func (p *parser) parseAbsValue() (Node, error) {
	p.consume() // consume opening PIPE
	p.absDepth++
	inner, err := p.parseSum()
	p.absDepth--
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.PIPE, ""); err != nil {
		return nil, fmt.Errorf("parser: unexpected EOF: unclosed absolute value")
	}
	return &FunctionNode{Name: "abs", Args: []Node{inner}}, nil
}

// parseNorm handles NORM sum NORM → NormNode (spec §4.3).
// \lVert and \rVert both produce NORM tokens; any NORM closes the expression.
// normDepth is incremented while parsing the inner expression so that the
// closing \rVert is not mistakenly consumed as an implicit-multiply trigger.
func (p *parser) parseNorm() (Node, error) {
	p.consume() // consume opening NORM
	p.normDepth++
	if p.peek().Type == lexer.EOF {
		p.normDepth--
		return nil, fmt.Errorf("parser: empty \\lVert\\rVert body")
	}
	inner, err := p.parseSum()
	p.normDepth--
	if err != nil {
		return nil, err
	}
	if p.peek().Type != lexer.NORM {
		return nil, fmt.Errorf("parser: unmatched \\lVert: expected closing \\rVert")
	}
	p.consume() // consume closing NORM
	return &NormNode{Arg: inner}, nil
}

// parseLargeOp handles \sum and \prod with their bound and body (spec §4.2).
//
// Grammar:
//
//	largop → COMMAND("sum"|"prod")
//	         UNDERSCORE LPAREN('{') SYMBOL EQUALS parseSum RPAREN('}')
//	         POWER super_arg
//	         parsePower
func (p *parser) parseLargeOp(op string) (Node, error) {
	// Lower bound: _{var=from}
	if p.peek().Type != lexer.UNDERSCORE {
		return nil, fmt.Errorf("parser: \\%s requires _{var=from}^{to}: expected '_{' after \\%s", op, op)
	}
	p.consume() // consume UNDERSCORE

	if p.peek().Type != lexer.LPAREN || p.peek().Value != "{" {
		return nil, fmt.Errorf("parser: \\%s requires _{var=from}^{to}: expected '{' after '_'", op)
	}
	p.consume() // consume '{'

	if p.peek().Type != lexer.SYMBOL {
		return nil, fmt.Errorf("parser: \\%s iteration variable must be a symbol, got %q", op, p.peek().Value)
	}
	varTok := p.consume()

	if _, err := p.expect(lexer.EQUALS, "="); err != nil {
		return nil, fmt.Errorf("parser: \\%s requires _{var=from}: expected '=' after variable", op)
	}

	from, err := p.parseSum()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(lexer.RPAREN, "}"); err != nil {
		return nil, fmt.Errorf("parser: \\%s lower bound missing closing '}'", op)
	}

	// Upper bound: ^{to}
	if p.peek().Type != lexer.POWER {
		return nil, fmt.Errorf("parser: \\%s requires ^{to} after lower bound: expected '^'", op)
	}
	p.consume() // consume '^'

	to, err := p.parseCommandArg(op)
	if err != nil {
		return nil, err
	}

	// Body: parsed at parsePower level so that + and * are outside the sum.
	body, err := p.parsePower()
	if err != nil {
		return nil, err
	}

	return &LargeOpNode{Op: op, Var: varTok.Value, From: from, To: to, Body: body}, nil
}
