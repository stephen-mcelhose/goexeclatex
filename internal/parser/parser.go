package parser

import (
	"fmt"
	"strconv"

	"github.com/stephen-mcelhose/goexeclatex/internal/lexer"
)

// commandArities maps normalised command names to their argument counts (spec §6.1).
// Commands not in this table have arity 0 and produce a SymbolNode.
var commandArities = map[string]int{
	"frac": 2,
	"sqrt": 1,
	// trig
	"sin": 1, "cos": 1, "tan": 1,
	"sec": 1, "csc": 1, "cot": 1,
	// inverse trig (evaluatex abbreviations — see ADR-003)
	"asin": 1, "acos": 1, "atan": 1,
	"asec": 1, "acsc": 1, "acot": 1,
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
	tokens   []lexer.Token
	pos      int
	absDepth int
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

// closerFor returns the expected closing bracket value for a given opener.
func closerFor(opener string) string {
	switch opener {
	case "(":
		return ")"
	case "[":
		return "]"
	case "{":
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

// parsePower → parseUnary [POWER parsePower]  (right-associative)
func (p *parser) parsePower() (Node, error) {
	base, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	if p.peek().Type == lexer.POWER {
		p.consume()
		exp, err := p.parsePower() // right-recursive
		if err != nil {
			return nil, err
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

// parsePostfix → parseAtom [BANG]
func (p *parser) parsePostfix() (Node, error) {
	base, err := p.parseAtom()
	if err != nil {
		return nil, err
	}
	if p.peek().Type == lexer.BANG {
		p.consume()
		return &UnaryNode{Op: "!", Operand: base}, nil
	}
	return base, nil
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
	if tok.Value != closer {
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
