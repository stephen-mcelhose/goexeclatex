---
type: how-to
title: goexeclatex — Example Calls
description: Concrete CLI invocations and pointer to the public library API.
tags: [goexeclatex, examples, cli, library]
timestamp: 2026-08-15T19:30:00Z
---

# goexeclatex — Example Calls

For the importable Go API (`goexeclatex.Eval`), see the [README](../README.md)
and [[specs/library]]. The rest of this page is CLI-oriented ([[how-to]]).

## Basic arithmetic

```sh
printf '%s\n' '2^{10}'                 | goexeclatex   # 1024
goexeclatex -e '\frac{22}{7}'                          # 3.142857...
goexeclatex -e '\sqrt{2}'                              # 1.414213...
goexeclatex -e '\frac{1}{2} + \sqrt{9}'               # 3.5
goexeclatex -e '\sqrt[3]{27}'                          # ~3 (v0.4)
goexeclatex -e '\sqrt[3]{-8}'                          # ~-2 (v0.4, ADR-015)
goexeclatex -e '\lfloor 3.2 \rfloor'                   # 3 (v0.4)
goexeclatex -e '10 \bmod 3'                            # 1 (v0.4)
```

## Trigonometry

Brace args (LaTeX command style) and unary paren args both work:

```sh
goexeclatex -e '\sin{\pi/6}'                           # 0.5
goexeclatex -e '\sin(\pi/6)'                           # 0.5
goexeclatex -e '\cos{0}'                               # 1
goexeclatex -e '\tan{\pi/4}'                           # 1
goexeclatex -e '\arctan{1} * 4'                        # 3.141592...
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

## Area Under the Progress Curve (AUPC)

Trapezoidal accumulation:

\[
\mathrm{AUPC}_i = \mathrm{AUPC}_{i-1}
  + \tfrac12(v_i + v_{i-1})(t_i - t_{i-1})
\]

Equivalently as a discrete sum (brace the body so `+`/`*` stay inside):

\[
\mathrm{AUPC} = \sum_{i=2}^{n}
  \tfrac12\,(v_i + v_{i-1})\,(t_i - t_{i-1})
\]

### IncidenceRamp — one trapezoid step → 362.5

```sh
# v0.1 — plain variable names
goexeclatex \
  -v AUPC=100 -v val1=85 -v val0=90 -v t1=6 -v t0=3 \
  -e 'AUPC + \frac{1}{2} * (val1 + val0) * (t1 - t0)'
# => 362.5

# v0.3 — subscript variable names (both x_0 and x_{0} accepted in -v)
goexeclatex \
  -v 'AUPC_{0}=100' -v 'val_{1}=85' -v 'val_{0}=90' \
  -v 't_{1}=6' -v 't_{0}=3' \
  -e 'AUPC_{0} + \frac{1}{2} * (val_{1} + val_{0}) * (t_{1} - t_{0})'
# => 362.5
```

### PestProgress — four points → 101

Series: \((t,v) = (0,0),\ (4,3),\ (9,9),\ (14,17)\).

Step form (final step; prior AUPC already 36):

```sh
goexeclatex \
  -v 'AUPC_{0}=36' -v 'v_{2}=9' -v 'v_{3}=17' -v 't_{2}=9' -v 't_{3}=14' \
  -e 'AUPC_{0} + \frac{1}{2} * (v_{3} + v_{2}) * (t_{3} - t_{2})'
# => 101
```

Sum form (full series in one expression; uses dynamic subscripts `v_{i}`, `v_{i-1}`):

```sh
goexeclatex \
  -v 'v_{1}=0' -v 'v_{2}=3' -v 'v_{3}=9' -v 'v_{4}=17' \
  -v 't_{1}=0' -v 't_{2}=4' -v 't_{3}=9' -v 't_{4}=14' \
  -e '\sum_{i=2}^{4}{\frac{1}{2}*(v_{i}+v_{i-1})*(t_{i}-t_{i-1})}'
# => 101
```

### EpidemicCurve — seven points → 387.5

\((t,v) = (0,20),\ (1,32),\ (2,50),\ (3,78),\ (4,92),\ (5,88),\ (6,75)\)

```sh
goexeclatex \
  -v 'v_{1}=20' -v 'v_{2}=32' -v 'v_{3}=50' -v 'v_{4}=78' \
  -v 'v_{5}=92' -v 'v_{6}=88' -v 'v_{7}=75' \
  -v 't_{1}=0' -v 't_{2}=1' -v 't_{3}=2' -v 't_{4}=3' \
  -v 't_{5}=4' -v 't_{6}=5' -v 't_{7}=6' \
  -e '\sum_{i=2}^{7}{\frac{1}{2}*(v_{i}+v_{i-1})*(t_{i}-t_{i-1})}'
# => 387.5
```

### SpikeAndFade — ten points → 514

\((t,v) = (0,0),\ (1,2),\ (2,4),\ (3,8),\ (6,16),\ (9,32),\ (12,64),\ (15,32),\ (18,16),\ (21,8)\)

```sh
goexeclatex \
  -v 'v_{1}=0' -v 'v_{2}=2' -v 'v_{3}=4' -v 'v_{4}=8' -v 'v_{5}=16' \
  -v 'v_{6}=32' -v 'v_{7}=64' -v 'v_{8}=32' -v 'v_{9}=16' -v 'v_{10}=8' \
  -v 't_{1}=0' -v 't_{2}=1' -v 't_{3}=2' -v 't_{4}=3' -v 't_{5}=6' \
  -v 't_{6}=9' -v 't_{7}=12' -v 't_{8}=15' -v 't_{9}=18' -v 't_{10}=21' \
  -e '\sum_{i=2}^{10}{\frac{1}{2}*(v_{i}+v_{i-1})*(t_{i}-t_{i-1})}'
# => 514
```

## Weighted Index

`WI = Σ(x_i · w_i) / Σ(x_i)`

### ThreeClassTextbook → ≈1.9583

```sh
# v0.1 — plain variable names
goexeclatex \
  -v x0=80 -v x1=90 -v x2=70 \
  -v w0=1  -v w1=2  -v w2=3  \
  -e '\frac{(x0*w0) + (x1*w1) + (x2*w2)}{x0 + x1 + x2}'
# => 1.958333...

# v0.3 — subscript variable names (both x_0 and x_{0} accepted in -v)
goexeclatex \
  -v 'x_{0}=80' -v 'x_{1}=90' -v 'x_{2}=70' \
  -v 'w_{0}=1'  -v 'w_{1}=2'  -v 'w_{2}=3'  \
  -e '\frac{(x_{0}*w_{0})+(x_{1}*w_{1})+(x_{2}*w_{2})}{x_{0}+x_{1}+x_{2}}'
# => 1.958333...
```

### FourClassSeverity → 1.25

\(v=(8,4,3,5)\), \(w=(0,1,2,3)\)

```sh
goexeclatex \
  -v 'x_{0}=8' -v 'x_{1}=4' -v 'x_{2}=3' -v 'x_{3}=5' \
  -v 'w_{0}=0' -v 'w_{1}=1' -v 'w_{2}=2' -v 'w_{3}=3' \
  -e '\frac{(x_{0}*w_{0})+(x_{1}*w_{1})+(x_{2}*w_{2})+(x_{3}*w_{3})}{x_{0}+x_{1}+x_{2}+x_{3}}'
# => 1.25
```

### TwoClassSimple → ≈0.6667

```sh
goexeclatex \
  -v 'x_{0}=10' -v 'x_{1}=20' -v 'w_{0}=0' -v 'w_{1}=1' \
  -e '\frac{(x_{0}*w_{0})+(x_{1}*w_{1})}{x_{0}+x_{1}}'
# => 0.666...
```

## Townsend-Heuberger (TH Weighted Index)

`TH = Σ(x_i · w_i) / (w_final · Σ(x_i)) × 100`

Same inputs as the Weighted Index cases; \(w_{final}\) is the last class weight.

### ThreeClassTextbook → ≈65.277

```sh
# v0.1 — plain variable names
goexeclatex \
  -v x0=80 -v x1=90 -v x2=70 \
  -v w0=1  -v w1=2  -v w2=3  \
  -e '\frac{(x0*w0) + (x1*w1) + (x2*w2)}{w2 * (x0 + x1 + x2)} * 100'
# => 65.277...

# v0.3 — subscript variable names (both x_0 and x_{0} accepted in -v)
goexeclatex \
  -v 'x_{0}=80' -v 'x_{1}=90' -v 'x_{2}=70' \
  -v 'w_{0}=1'  -v 'w_{1}=2'  -v 'w_{2}=3'  \
  -e '\frac{(x_{0}*w_{0})+(x_{1}*w_{1})+(x_{2}*w_{2})}{w_{2}*(x_{0}+x_{1}+x_{2})}*100'
# => 65.277...
```

### FourClassSeverity → ≈41.667

```sh
goexeclatex \
  -v 'x_{0}=8' -v 'x_{1}=4' -v 'x_{2}=3' -v 'x_{3}=5' \
  -v 'w_{0}=0' -v 'w_{1}=1' -v 'w_{2}=2' -v 'w_{3}=3' \
  -e '\frac{(x_{0}*w_{0})+(x_{1}*w_{1})+(x_{2}*w_{2})+(x_{3}*w_{3})}{w_{3}*(x_{0}+x_{1}+x_{2}+x_{3})}*100'
# => 41.666...
```

### TwoClassSimple → ≈66.667

```sh
goexeclatex \
  -v 'x_{0}=10' -v 'x_{1}=20' -v 'w_{0}=0' -v 'w_{1}=1' \
  -e '\frac{(x_{0}*w_{0})+(x_{1}*w_{1})}{w_{1}*(x_{0}+x_{1})}*100'
# => 66.666...
```

## Subsample Aggregation

Named dataset **PlotStripFive**: \(1, 19, 2, 8, 5\).

```sh
goexeclatex -e '1+19+2+8+5'                 # Sum → 35
goexeclatex -e '\frac{1+19+2+8+5}{5}'        # Average → 7
goexeclatex -e '\max(1,19,2,8,5)'            # Maximum → 19
goexeclatex -e '\min(1,19,2,8,5)'            # Minimum → 1
```

Median (5) and sample standard deviation (≈7.2457) need host-side stats;
coefficient of variation is StdDev / Average (≈1.0351).

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
| `\arctan`              | v0.2      | arc trig aliases               |
| `\sqrt[n]`, floor, `\bmod`, `\log_{b}` | v0.4 | [[specs/parser-extensions]] |
| Plain variable calls   | v0.1      | symbol table                   |
| Subscript variable calls | v0.3    | `_` token, subscript grammar   |
| AUPC via `\sum_{i=a}^{b}` | v0.3  | large-op engine + dynamic `v_{i}` ([[specs/subscripts-largeops]]) |
| `\max` / `\min` (n≥2)  | v0.4      | variadic paren args            |

## Sources

- [[roadmap]], [[how-to]], [[specs/parser-extensions]], [[specs/subscripts-largeops]]
