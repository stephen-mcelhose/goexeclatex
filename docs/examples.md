---
type: how-to
title: goexeclatex — Example Calls
description: Concrete CLI invocations showing common evaluation patterns.
tags: [goexeclatex, examples, cli]
timestamp: 2026-08-14T23:00:00Z
---

# goexeclatex — Example Calls

## Basic arithmetic

```sh
echo '2^{10}'                          | goexeclatex   # 1024
echo '\frac{22}{7}'                    | goexeclatex   # 3.142857...
goexeclatex -e '\sqrt{2}'                              # 1.414213...
goexeclatex -e '\sqrt[3]{27}'                          # 3  (v0.2)
goexeclatex -e '\frac{1}{2} + \sqrt{9}'               # 3.5
```

## Trigonometry

```sh
goexeclatex -e '\sin(\pi / 6)'                        # 0.5
goexeclatex -e '\cos(0)'                              # 1
goexeclatex -e '\tan(\pi / 4)'                        # 1
goexeclatex -e '\arctan(1) * 4'                       # 3.141592...  (v0.2)
```

## Variables

```sh
goexeclatex -v x=3 -e 'x^2 + 2x + 1'                 # 16
goexeclatex -v a=5 -v b=12 -e '\sqrt{a^2 + b^2}'     # 13
goexeclatex -v r=7 -e '\pi * r^2'                     # 153.938...
```

## Absolute value / factorial

```sh
goexeclatex -e '\lvert -42 \rvert'                    # 42
goexeclatex -e '|{-3.14}|'                            # 3.14
goexeclatex -e '10!'                                  # 3628800
```

## Area Under the Progress Curve (trapezoidal step)

One step of the accumulator: `AUPC_i = AUPC_prev + ½(val_cur + val_prev)(t_cur - t_prev)`

```sh
# v0.1 — plain variable names
goexeclatex \
  -v AUPC=100 -v val1=85 -v val0=90 -v t1=6 -v t0=3 \
  -e 'AUPC + \frac{1}{2} * (val1 + val0) * (t1 - t0)'
# => 362.5

# v0.3 — subscript variable names
goexeclatex \
  -v 'AUPC_{0}=100' -v 'val_{1}=85' -v 'val_{0}=90' \
  -v 't_{1}=6' -v 't_{0}=3' \
  -e 'AUPC_{0} + \frac{1}{2} * (val_{1} + val_{0}) * (t_{1} - t_{0})'
# => 362.5
```

## Weighted Index

`WI = Σ(x_i · w_i) / Σ(x_i)`

```sh
# v0.1 — plain variable names
goexeclatex \
  -v x0=80 -v x1=90 -v x2=70 \
  -v w0=1  -v w1=2  -v w2=3  \
  -e '\frac{(x0*w0) + (x1*w1) + (x2*w2)}{x0 + x1 + x2}'
# => 1.958333...

# v0.3 — subscript variable names
goexeclatex \
  -v 'x_{0}=80' -v 'x_{1}=90' -v 'x_{2}=70' \
  -v 'w_{0}=1'  -v 'w_{1}=2'  -v 'w_{2}=3'  \
  -e '\frac{(x_{0}*w_{0})+(x_{1}*w_{1})+(x_{2}*w_{2})}{x_{0}+x_{1}+x_{2}}'
# => 1.958333...
```

## Townsend-Heuberger

`TH = Σ(x_i · w_i) / (w_final · Σ(x_i)) × 100`

```sh
# v0.1 — plain variable names
goexeclatex \
  -v x0=80 -v x1=90 -v x2=70 \
  -v w0=1  -v w1=2  -v w2=3  \
  -e '\frac{(x0*w0) + (x1*w1) + (x2*w2)}{w2 * (x0 + x1 + x2)} * 100'
# => 65.277...

# v0.3 — subscript variable names
goexeclatex \
  -v 'x_{0}=80' -v 'x_{1}=90' -v 'x_{2}=70' \
  -v 'w_{0}=1'  -v 'w_{1}=2'  -v 'w_{2}=3'  \
  -e '\frac{(x_{0}*w_{0})+(x_{1}*w_{1})+(x_{2}*w_{2})}{w_{2}*(x_{0}+x_{1}+x_{2})}*100'
# => 65.277...
```

## Precision control

```sh
goexeclatex -p 2 -e '\pi'              # 3.14
goexeclatex -p 0 -e '\sqrt{2}'        # 1
goexeclatex -p 10 -e '\pi'            # 3.1415926536
```

## Milestone reference

| Example group          | Milestone | Key features needed            |
| ---------------------- | --------- | ------------------------------ |
| Basic arithmetic       | v0.1      | `\frac`, `^`, arithmetic       |
| Trig                   | v0.1      | `\sin/\cos/\tan`               |
| `\arctan`, `\sqrt[n]`  | v0.2      | arc trig, nth root             |
| Plain variable calls   | v0.1      | symbol table                   |
| Subscript variable calls | v0.3    | `_` token, subscript grammar   |
