package main_test

import (
	"bufio"
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

// bin is the path to the compiled binary, set once in TestMain.
var bin string

func TestMain(m *testing.M) {
	tmp, err := os.MkdirTemp("", "goexeclatex-test-*")
	if err != nil {
		panic("mktemp: " + err.Error())
	}

	bin = filepath.Join(tmp, "goexeclatex")
	out, err := exec.Command("go", "build", "-o", bin, ".").CombinedOutput()
	if err != nil {
		os.RemoveAll(tmp)
		panic("build failed:\n" + string(out))
	}

	code := m.Run()
	os.RemoveAll(tmp)
	os.Exit(code)
}

// run executes the binary with the given args, optionally piping stdin.
// Returns stdout, stderr, and the exit code. Aborts after 5 seconds.
func run(stdin string, args ...string) (stdout, stderr string, exit int) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin, args...)
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	exit = 0
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			exit = ee.ExitCode()
		}
	}
	return strings.TrimRight(outBuf.String(), "\n"), errBuf.String(), exit
}

// loadCases reads a tab-separated testdata file, skipping blank lines and
// lines starting with #. Returns rows of fields.
func loadCases(t *testing.T, name string) [][]string {
	t.Helper()
	f, err := os.Open(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("open %s: %v", name, err)
	}
	defer f.Close()
	var rows [][]string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		rows = append(rows, strings.Split(line, "\t"))
	}
	return rows
}

// TestStdinCases pipes each expression to the binary and checks stdout + exit.
func TestStdinCases(t *testing.T) {
	for _, row := range loadCases(t, "stdin.txt") {
		if len(row) != 3 {
			t.Fatalf("stdin.txt: bad row %v (want 3 fields)", row)
		}
		expr, wantOut, wantExitStr := row[0], row[1], row[2]
		wantExit, _ := strconv.Atoi(wantExitStr)
		t.Run(expr, func(t *testing.T) {
			got, _, exit := run(expr)
			if exit != wantExit {
				t.Errorf("exit = %d, want %d", exit, wantExit)
			}
			if got != wantOut {
				t.Errorf("stdout = %q, want %q", got, wantOut)
			}
		})
	}
}

// TestErrorCases pipes each expression and checks that stderr contains the
// expected substring and the exit code matches.
func TestErrorCases(t *testing.T) {
	for _, row := range loadCases(t, "errors.txt") {
		if len(row) != 3 {
			t.Fatalf("errors.txt: bad row %v (want 3 fields)", row)
		}
		expr, wantErrContains, wantExitStr := row[0], row[1], row[2]
		wantExit, _ := strconv.Atoi(wantExitStr)
		t.Run(expr, func(t *testing.T) {
			_, errOut, exit := run(expr)
			if exit != wantExit {
				t.Errorf("exit = %d, want %d (stderr: %q)", exit, wantExit, errOut)
			}
			if !strings.Contains(errOut, wantErrContains) {
				t.Errorf("stderr = %q, want it to contain %q", errOut, wantErrContains)
			}
		})
	}
}

// TestFlagCases covers -e, -v, and -p flags directly in Go (avoids
// shell-quoting complexity in testdata files).
func TestFlagCases(t *testing.T) {
	cases := []struct {
		name     string
		args     []string
		wantOut  string
		wantExit int
	}{
		// -e flag
		{"expr flag frac", []string{"-e", `\frac{1}{2}`}, "0.5", 0},
		{"expr flag sqrt", []string{"-e", `\sqrt{2}`}, "1.4142135623730951", 0},
		{"expr flag power", []string{"-e", `2^{10}`}, "1024", 0},

		// -p precision flag
		{"prec 4 sqrt2", []string{"-e", `\sqrt{2}`, "-p", "4"}, "1.4142", 0},
		{"prec 2 pi", []string{"-e", `\pi`, "-p", "2"}, "3.14", 0},
		{"prec 0 pi", []string{"-e", `\pi`, "-p", "0"}, "3", 0},
		{"prec 4 pi", []string{"-e", `\pi`, "-p", "4"}, "3.1416", 0},
		{"prec 0 half", []string{"-e", `\frac{1}{2}`, "-p", "0"}, "0", 0}, // round-half-to-even: 0.5 → 0

		// -v variable binding
		{"var x squared", []string{"-v", "x=3", "-e", "x^2"}, "9", 0},
		{"var x+y", []string{"-v", "x=3", "-v", "y=4", "-e", "x+y"}, "7", 0},
		{"var pythagorean", []string{"-v", "x=3", "-v", "y=4", "-e", `\sqrt{x^2+y^2}`}, "5", 0},
		{"var override pi", []string{"-v", "pi=3", "-e", `\pi`}, "3", 0},
		{"var float val", []string{"-v", "x=1.5", "-e", "x^2"}, "2.25", 0},

		// v0.2 functions with user-supplied variables.
		{"var greek sin", []string{"-v", "alpha=1.5", "-e", `\sin{\alpha}`}, "0.9974949866040544", 0},
		{"var ln of var", []string{"-v", "x=1", "-e", `\ln{x}`}, "0", 0},
		{"var binom n k", []string{"-v", "n=5", "-v", "k=2", "-e", `\binom{n}{k}`}, "10", 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, _, exit := run("", c.args...)
			if exit != c.wantExit {
				t.Errorf("exit = %d, want %d", exit, c.wantExit)
			}
			if got != c.wantOut {
				t.Errorf("stdout = %q, want %q", got, c.wantOut)
			}
		})
	}
}

// TestFlagErrors covers flag-level error paths.
func TestFlagErrors(t *testing.T) {
	cases := []struct {
		name            string
		args            []string
		stdin           string
		wantErrContains string
		wantExit        int
	}{
		{"empty -e", []string{"-e", ""}, "", "no expression provided", 1},
		{"whitespace -e", []string{"-e", "   "}, "", "no expression provided", 1},
		{"positional arg", []string{"1+2"}, "", "unknown command", 1},
		{"empty var name", []string{"-v", "=10", "-e", "1"}, "", "variable name cannot be empty", 1},
		{"bad -v format", []string{"-v", "x", "-e", "1"}, "", "expected name=value", 1},
		{"bad -v value", []string{"-v", "x=abc", "-e", "1"}, "", "invalid", 1},
		{"eval error exit 2", []string{"-e", "1/0"}, "", "division by zero", 2},
		{"parse error exit 1", []string{"-e", `\frac{1}{2`}, "", "unclosed '{'", 1},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, errOut, exit := run(c.stdin, c.args...)
			if exit != c.wantExit {
				t.Errorf("exit = %d, want %d (stderr: %q)", exit, c.wantExit, errOut)
			}
			if !strings.Contains(errOut, c.wantErrContains) {
				t.Errorf("stderr = %q, want it to contain %q", errOut, c.wantErrContains)
			}
		})
	}
}
