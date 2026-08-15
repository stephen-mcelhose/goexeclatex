package goexeclatex

// SyntaxError is a lex or parse failure (CLI exit 1 when used by the CLI).
type SyntaxError struct {
	Err error
}

func (e *SyntaxError) Error() string {
	if e == nil || e.Err == nil {
		return "syntax error"
	}
	return e.Err.Error()
}

func (e *SyntaxError) Unwrap() error { return e.Err }

// EvalError is an evaluation failure (CLI exit 2 when used by the CLI).
type EvalError struct {
	Err error
}

func (e *EvalError) Error() string {
	if e == nil || e.Err == nil {
		return "eval error"
	}
	return e.Err.Error()
}

func (e *EvalError) Unwrap() error { return e.Err }
