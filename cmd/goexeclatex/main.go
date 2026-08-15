package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stephen-mcelhose/goexeclatex/internal/eval"
	"github.com/stephen-mcelhose/goexeclatex/internal/lexer"
	"github.com/stephen-mcelhose/goexeclatex/internal/parser"
)

func main() {
	if err := rootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func rootCmd() *cobra.Command {
	var vars []string
	var prec int
	var expr string

	cmd := &cobra.Command{
		Use:   "goexeclatex",
		Short: "Evaluate a LaTeX math expression",
		Long: `goexeclatex reads a LaTeX math expression and prints its numeric result.

Reads from stdin by default:
  echo '\frac{1}{2} + \sqrt{9}' | goexeclatex

Or supply the expression directly with -e:
  goexeclatex -e '\sin{\pi/6}'`,
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			input, err := readInput(cmd, expr)
			if err != nil {
				die(1, err)
			}

			userVars, err := parseVars(vars)
			if err != nil {
				die(1, err)
			}

			tokens, err := lexer.Lex(input)
			if err != nil {
				die(1, err) // lex error → exit 1 (spec §8)
			}

			node, err := parser.Parse(tokens)
			if err != nil {
				die(1, err) // parse error → exit 1 (spec §8)
			}

			scope := eval.NewScope()
			for k, v := range userVars {
				scope[k] = v
			}

			result, err := eval.Eval(node, scope)
			if err != nil {
				die(2, err) // eval error → exit 2 (spec §8)
			}

			fmt.Println(formatResult(result, prec))
			return nil
		},
	}

	cmd.Flags().StringArrayVarP(&vars, "var", "v", nil, "bind a variable: -v x=3.14 (repeatable)")
	cmd.Flags().IntVarP(&prec, "prec", "p", -1, "decimal places in output (-1 = full precision)")
	cmd.Flags().StringVarP(&expr, "expr", "e", "", "expression to evaluate (skips stdin read)")

	return cmd
}

// die prints a user-facing error to stderr and exits with the given code.
// Package prefixes (eval:, lexer:, parser:) are stripped per ADR-011.
func die(code int, err error) {
	fmt.Fprintf(os.Stderr, "error: %s\n", userMessage(err))
	os.Exit(code)
}

// userMessage strips internal package prefixes from an error string (ADR-011).
func userMessage(err error) string {
	s := err.Error()
	for _, pfx := range []string{"eval: ", "lexer: ", "parser: "} {
		s = strings.TrimPrefix(s, pfx)
	}
	return s
}

// readInput returns the expression from -e or stdin.
func readInput(cmd *cobra.Command, expr string) (string, error) {
	if e := strings.TrimSpace(expr); e != "" {
		return e, nil
	}
	b, err := io.ReadAll(cmd.InOrStdin())
	if err != nil {
		return "", fmt.Errorf("reading stdin: %w", err)
	}
	input := strings.TrimSpace(string(b))
	if input == "" {
		return "", fmt.Errorf("no expression provided (use stdin or -e)")
	}
	return input, nil
}

// parseVars converts ["x=1.0", "y=2.0"] into a map[string]float64.
func parseVars(vars []string) (map[string]float64, error) {
	scope := make(map[string]float64, len(vars))
	for _, v := range vars {
		name, raw, ok := strings.Cut(v, "=")
		if !ok {
			return nil, fmt.Errorf("invalid -v value %q: expected name=value", v)
		}
		trimName := strings.TrimSpace(name)
		if trimName == "" {
			return nil, fmt.Errorf("invalid -v value %q: variable name cannot be empty", v)
		}
		val, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid -v value %q: %w", v, err)
		}
		scope[trimName] = val
	}
	return scope, nil
}

// formatResult formats a float64 per spec §7.
// Default (-1): 'g' format — shortest representation that round-trips.
// Fixed (≥0): 'f' format — exactly N decimal places.
func formatResult(v float64, prec int) string {
	if prec < 0 {
		return strconv.FormatFloat(v, 'g', -1, 64)
	}
	if prec > 100 {
		prec = 100
	}
	return strconv.FormatFloat(v, 'f', prec, 64)
}
