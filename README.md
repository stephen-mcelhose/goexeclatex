# goexeclatex

**goexeclatex** — the Go execut(or) for LaTeX math expressions.

A Go library and CLI for parsing and evaluating LaTeX math syntax into numeric results.

## Library

```go
package main

import (
	"fmt"
	"log"

	"github.com/stephen-mcelhose/goexeclatex"
)

func main() {
	v, err := goexeclatex.Eval(`\frac{1}{2} + \sqrt{9}`, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(v) // 3.5

	v, err = goexeclatex.Eval(`x^2 + 2x + 1`, map[string]float64{"x": 3})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(v) // 16
}
```

`Eval` returns a full-precision `float64`. Errors are stage-discriminable
(`*goexeclatex.SyntaxError` vs `*goexeclatex.EvalError`). Normative API:
[`docs/specs/library.md`](docs/specs/library.md).

## CLI

```sh
go install github.com/stephen-mcelhose/goexeclatex/cmd/goexeclatex@latest

goexeclatex -e '\frac{1}{2} + \sqrt{9}'
# 3.5

goexeclatex -e '\sin{\pi/6}'
# 0.5

goexeclatex -v x=3 -e 'x^2 + 2x + 1'
# 16
```

See [`docs/how-to.md`](docs/how-to.md) for flags (`-e`, `-v`, `-p`), stdin/`printf` usage, and error exit codes.

## Status

Early development. Public library API: [issue #20](https://github.com/stephen-mcelhose/goexeclatex/issues/20).

## License

MIT
