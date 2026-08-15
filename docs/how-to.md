---
type: how-to
title: How to use goexeclatex
description: User-facing guide covering stdin and -e flag usage, -p precision, -v variable binding, supported syntax, error handling, and shell composition patterns.
tags: [goexeclatex, cli, how-to]
timestamp: 2026-08-15T03:04:00Z
---

# How to use goexeclatex

`goexeclatex` evaluates a LaTeX math expression and prints the numeric result.
It is designed to be composed with shell pipelines.

## Install

```sh
go install github.com/stephen-mcelhose/goexeclatex/cmd/goexeclatex@latest
```

Or build from source:

```sh
git clone https://github.com/stephen-mcelhose/goexeclatex.git
cd goexeclatex
go build -o goexeclatex ./cmd/goexeclatex
```

## Evaluate an expression from stdin

```sh
echo '\frac{1}{2} + \sqrt{9}' | goexeclatex
# 3.5

echo '2^{10}' | goexeclatex
# 1024

echo '\sin{\pi/2}' | goexeclatex
# 1
```

## Evaluate an expression with the `-e` flag

```sh
goexeclatex -e '\sqrt{2}'
# 1.4142135623730951

goexeclatex -e '\frac{22}{7}'
# 3.142857142857143
```

## Control output precision with `-p`

`-p N` prints exactly N decimal places. The default (`-p -1`) uses the
shortest representation that round-trips.

```sh
goexeclatex -e '\sqrt{2}' -p 4
# 1.4142

goexeclatex -e '\pi' -p 2
# 3.14

goexeclatex -e '\pi' -p 10
# 3.1415926536
```

## Bind variables with `-v`

Use `-v name=value` to pass numeric values into the expression. The flag can
be repeated for multiple variables.

```sh
goexeclatex -v x=3 -e 'x^2'
# 9

goexeclatex -v x=3 -v y=4 -e '\sqrt{x^2+y^2}'
# 5

goexeclatex -v r=5 -e '\pi * r^2'
# 78.53981633974483
```

Variables can override built-in constants if needed:

```sh
goexeclatex -v pi=3 -e '\pi * r^2' -v r=1
# 3
```

## Supported syntax

| Construct                     | Example                                              |
| ----------------------------- | ---------------------------------------------------- |
| Arithmetic operators          | `1+2`, `10-3`, `3*4`, `10/4`                         |
| Exponentiation (right-assoc)  | `2^3`, `2^{10}`, `2^{3^2}`                           |
| Grouping                      | `(1+2)*3`, `[a+b]`, `{c+d}`                          |
| Fractions                     | `\frac{1}{2}`, `\dfrac{1}{2}`, `\tfrac{1}{2}`        |
| Square root                   | `\sqrt{4}`, `\sqrt{x^2+y^2}`                         |
| Absolute value                | `\|-3\|`, `\lvert -5 \rvert`                         |
| Factorial                     | `5!`, `0!`                                           |
| Implicit multiply             | `2\pi`, `3\sqrt{4}`                                  |
| Trig functions                | `\sin`, `\cos`, `\tan`, `\sec`, `\csc`, `\cot`       |
| Inverse trig                  | `\arcsin`, `\arccos`, `\arctan` (also `\asin` etc.)  |
| Hyperbolic trig               | `\sinh`, `\cosh`, `\tanh`, `\coth`, `\sech`, `\csch` |
| Log / exp                     | `\ln`, `\log` (base 10), `\exp`                      |
| Binomial coefficient          | `\binom{n}{k}`, `\dbinom{n}{k}`, `\tbinom{n}{k}`    |
| Built-in constants            | `\pi`, `e`, `\tau`, `\phi`, `\infty`                 |

Arguments to functions use `{}` braces:

```sh
goexeclatex -e '\sin{\pi/6}'        # 0.5 (approximately)
goexeclatex -e '\arccos{0}'         # 1.5707963267948966  (π/2)
goexeclatex -e '\ln{\e}'            # 1
goexeclatex -e '\binom{5}{2}'       # 10
```

See [[examples]] for more complete invocation patterns.

## Error handling

Parse errors (bad syntax) exit with code **1**:

```sh
echo '\frac{1}{2' | goexeclatex
# error: unexpected end of input: unclosed '{' group at position 10
# exit 1
```

Evaluation errors (domain violations, undefined symbols) exit with code **2**:

```sh
echo '1/0' | goexeclatex
# error: division by zero
# exit 2

echo '\sqrt{-1}' | goexeclatex
# error: domain error: sqrt of negative
# exit 2

echo 'undefined_var' | goexeclatex
# error: undefined symbol: undefined_var
# exit 2
```

Success exits with code **0**. You can use the exit code in shell pipelines:

```sh
result=$(goexeclatex -e '\sqrt{x^2+y^2}' -v x=3 -v y=4) && echo "Distance: $result"
# Distance: 5
```

## Compose with other tools

```sh
# Feed results into further computation
A=$(echo '3^2' | goexeclatex)
B=$(echo '4^2' | goexeclatex)
echo "$A + $B" | goexeclatex
# 25

# Use in a loop
for n in 1 2 3 4 5; do
  goexeclatex -v n="$n" -e 'n^2'
done
# 1
# 4
# 9
# 16
# 25
```
