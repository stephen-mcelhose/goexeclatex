package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	if err := rootCmd().Execute(); err != nil {
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
  goexeclatex -e '\sin(\pi / 6)'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			input, err := readInput(cmd, expr)
			if err != nil {
				return err
			}

			scope, err := parseVars(vars)
			if err != nil {
				return err
			}

			result, err := evaluate(input, scope)
			if err != nil {
				return err
			}

			fmt.Println(formatResult(result, prec))
			return nil
		},
		SilenceUsage: true,
	}

	cmd.Flags().StringArrayVarP(&vars, "var", "v", nil, "bind a variable: -v x=3.14 (repeatable)")
	cmd.Flags().IntVarP(&prec, "prec", "p", -1, "decimal places in output (-1 = full precision)")
	cmd.Flags().StringVarP(&expr, "expr", "e", "", "expression to evaluate (skips stdin read)")

	return cmd
}

// readInput returns the expression string from -e or stdin.
func readInput(cmd *cobra.Command, expr string) (string, error) {
	if expr != "" {
		return strings.TrimSpace(expr), nil
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

// parseVars converts ["x=1.0", "y=2.0"] into a map.
func parseVars(vars []string) (map[string]float64, error) {
	scope := make(map[string]float64, len(vars))
	for _, v := range vars {
		name, raw, ok := strings.Cut(v, "=")
		if !ok {
			return nil, fmt.Errorf("invalid -v value %q: expected name=value", v)
		}
		val, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid -v value %q: %w", v, err)
		}
		scope[strings.TrimSpace(name)] = val
	}
	return scope, nil
}

// evaluate is the stub that will call lexer → parser → eval.
// Returns an error until those packages are implemented.
func evaluate(_ string, _ map[string]float64) (float64, error) {
	return 0, fmt.Errorf("not yet implemented")
}

// formatResult prints result with the requested precision.
func formatResult(v float64, prec int) string {
	if prec < 0 {
		return strconv.FormatFloat(v, 'f', -1, 64)
	}
	return strconv.FormatFloat(v, 'f', prec, 64)
}
