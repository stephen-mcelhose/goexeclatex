// Package goexeclatex evaluates LaTeX math expressions to float64 results.
//
// The public surface is a thin facade over internal lexer, parser, and
// evaluator packages. See docs/specs/library.md and ADR-017.
package goexeclatex

import (
	"github.com/stephen-mcelhose/goexeclatex/internal/eval"
	"github.com/stephen-mcelhose/goexeclatex/internal/lexer"
	"github.com/stephen-mcelhose/goexeclatex/internal/parser"
)

// Eval parses and evaluates a LaTeX math expression.
// Variables supplied via vars are merged onto the built-in scope before
// evaluation. A nil vars map is treated as empty (no user bindings).
// Returns the full-precision numeric result or a stage-discriminable error
// (*SyntaxError or *EvalError).
func Eval(expr string, vars map[string]float64) (float64, error) {
	tokens, err := lexer.Lex(expr)
	if err != nil {
		return 0, &SyntaxError{Err: err}
	}

	node, err := parser.Parse(tokens)
	if err != nil {
		return 0, &SyntaxError{Err: err}
	}

	scope := eval.NewScope()
	for k, v := range vars {
		scope[k] = v
	}

	result, err := eval.Eval(node, scope)
	if err != nil {
		return 0, &EvalError{Err: err}
	}
	return result, nil
}
