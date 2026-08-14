Short Math Guide for LaTeX, version [2.0 (2017/12/22)]{.upright}

::: center
Version [2.0 (2017/12/22)]{.upright}, currently available from a link at `https://www.ams.org/tex/amslatex`
:::

# Acknowledgments and plans for future work {#acknowledgments-and-plans-for-future-work .unnumbered}

Thanks to all who contributed suggestions, assistance and encouragement. Special thanks to David Carlisle for repairing unruly macros and to Jennifer Wright Sharp for applying consistent editing in AMS style.

Plans for a future edition include addition of an index.

Reports concerning errors and suggestions for improvement should be sent to\
[`tech-support@ams.org`](mailto:tech-support@ams.org) .

# Introduction

This is a concise summary of recommended features in LaTeX and a couple of extension packages for **writing math formulas**. Readers needing greater depth of detail are referred to the sources listed in the bibliography, especially [@lamport], [@amsldoc], and [@fntguide]. A certain amount of familiarity with standard LaTeX terminology is assumed; if your memory needs refreshing on the LaTeX meaning of *command*, *optional argument*, *environment*, *package*, and so forth, see [@lamport].

Most of the features described here are available to you if you use LaTeX with two extension packages published by the American Mathematical Society: and . Thus, the source file for this document begins with

    \documentclass{article}
    \usepackage{amssymb,amsmath}

The package might be omissible for documents whose math symbol usage is relatively modest; in Section [3](#mathsymbols){reference-type="ref" reference="mathsymbols"}, the symbols that require are marked with ^a^ or ^b^ (font or ). In Section [3.3](#alpha-digit){reference-type="ref" reference="alpha-digit"}, a few additional fonts are included; the necessary packages are identified there.

Many noteworthy features found in other packages are not covered here; see Section [10](#other-packages){reference-type="ref" reference="other-packages"}. Regarding math symbols, please note especially that the list given here is not intended to be comprehensive, but to illustrate such symbols as users will normally find already present in their / system and usable without installing any additional fonts or doing other setup work.

If you have a need for a symbol not shown here, you will probably want to consult *The Comprehensive LaTeX Symbol List* [@comprehensive]. If your / installation is based on TeX Live, and includes documentation, the list can also be accessed by typing `texdoc comprehensive` at a system prompt.

:::::::::::: makeimage
::::::::::: minipage
::: eqxample
    \begin{equation}\label{xx}
    \begin{split}
    a& =b+c-d\\
     & \quad +e-f\\
     & =g+h\\
     & =i
    \end{split}
    \end{equation}

$$\begin{equation}
\label{xx}
\begin{split}
a& =b+c-d\\
 & \quad +e-f\\
 & =g+h\\
 & =i
\end{split}
\end{equation}$$
:::

::: eqxample
    \begin{multline}
    a+b+c+d+e+f\\
    +i+j+k+l+m+n\\
    +o+p+q+r+s
    \end{multline}

$$\begin{multline}
a+b+c+d+e+f\\
+i+j+k+l+m+n\\
+o+p+q+r+s
\end{multline}$$
:::

::: eqxample
    \begin{gather}
    a_1=b_1+c_1\\
    a_2=b_2+c_2-d_2+e_2
    \end{gather}

$$\begin{gather}
a_1=b_1+c_1\\
a_2=b_2+c_2-d_2+e_2
\end{gather}$$
:::

::: eqxample
    \begin{align}
    a_1& =b_1+c_1\\
    a_2& =b_2+c_2-d_2+e_2
    \end{align}

$$\begin{align}
a_1& =b_1+c_1\\
a_2& =b_2+c_2-d_2+e_2
\end{align}$$
:::

::: eqxample
    \begin{align}
    a_{11}& =b_{11}&
      a_{12}& =b_{12}\\
    a_{21}& =b_{21}&
      a_{22}& =b_{22}+c_{22}
    \end{align}

$$\begin{align}
a_{11}& =b_{11}&
  a_{12}& =b_{12}\\
a_{21}& =b_{21}&
  a_{22}& =b_{22}+c_{22}
\end{align}$$
:::

::: eqxample
    \begin{alignat}{2}
    a_1& =b_1+c_1&      &+e_1-f_1\\
    a_2& =b_2+c_2&{}-d_2&+e_2
    \end{alignat}

$$\begin{alignat}
{2}
a_1& =b_1+c_1&      &+e_1-f_1\\
a_2& =b_2+c_2&{}-d_2&+e_2
\end{alignat}$$
:::

::: eqxample
    \begin{flalign}
    a_{11}& =b_{11}&
      a_{12}& =b_{12}\\
    a_{21}& =b_{21}&
      a_{22}& =b_{22}+c_{22}
    \end{flalign}

$$\begin{flalign}
a_{11}& =b_{11}&
  a_{12}& =b_{12}\\
a_{21}& =b_{21}&
  a_{22}& =b_{22}+c_{22}
\end{flalign}$$
:::

::: notes
Applying to any primary environment will suppress the assignment of equation numbers. However, may be used to apply a visible label, and can be used to reference such manually tagged lines. Use of either or a on a subordinate environment is an error.

The environment is something of a special case. It is a subordinate environment that can be used as the contents of an environment or the contents of one in a multiple-equation structure such as or .

The primary environments , and have subordinate "" counterparts (, and ) that can be used as components of more complicated displays, or within in-line math. These "" environments can be positioned vertically using an optional argument `[t]`, `[c]` or `[b]`.

The name is meant as "full length", not "flush left" as often mistakenly reported. However, since a display occupying the full width will often begin at the left margin, this confusion is understandable. The indent applied to from both margins is set with .
:::
:::::::::::
::::::::::::

# Inline math formulas and displayed equations {#first-step}

## The fundamentals

Entering and leaving math mode in LaTeX is normally done with the following commands and environments.

:::: center
+-------------------+---+-------------------------------------------------+
|                   |   |                                                 |
+:=================:+:=:+:===============================================:+
| 1-1               |   | +:------------------------+:------------------+ |
|                   |   | |   -----------           | unnumbered        | |
|   --------------- |   | |   `\[...\]`             |                   | |
|     `$` ... `$`   |   | |   -----------           |                   | |
|    `\(` ... `\)`  |   | +-------------------------+-------------------+ |
|   --------------- |   | |   --------------------- | unnumbered        | |
|                   |   | |   `\begin{equation*}`   |                   | |
|                   |   | |   ...                   |                   | |
|                   |   | |   `\end{equation*}`     |                   | |
|                   |   | |   --------------------- |                   | |
|                   |   | +-------------------------+-------------------+ |
|                   |   | |   --------------------  |   --------------- | |
|                   |   | |   `\begin{equation}`    |   automatically   | |
|                   |   | |   ...                   |   numbered        | |
|                   |   | |   `\end{equation}`      |   --------------- | |
|                   |   | |   --------------------  |                   | |
|                   |   | +-------------------------+-------------------+ |
+-------------------+---+-------------------------------------------------+

::: notes
Do not leave a blank line between text and a displayed equation. This allows a page break at that location, which is bad style. It also causes the spacing between text and display to be incorrect, usually larger than it should be. If a visual break is desired in the input, insert a line containing only a `%` at the beginning. Leave a blank line between a display and following text only if a new paragraph is intended.

Do not group multiple display structures in the input (`\[...\]`, , etc.). Instead, use a multiline structure with substructures (, , etc.) as appropriate.

The alternative environments `math`  ... `math` and\
`displaymath`  ... `displaymath` are seldom needed in practice. Using the plain TeX notation `$$` ... `$$` for displayed equations is strongly discouraged. Although it is not expressly forbidden in LaTeX, it is not documented anywhere in the LaTeX book as being part of the LaTeX command set, and it interferes with the proper operation of various features such as the option.

The and environments described in [@lamport] are strongly discouraged because they produce inconsistent spacing of the equal signs and make no attempt to prevent overprinting of the equation body by the equation number.
:::
::::

Environments for handling equation groups and multiline equations are shown in Table [\[displays\]](#displays){reference-type="ref" reference="displays"}.

## Automatic numbering and cross-referencing

To get an auto-numbered equation, use the environment; to assign a label for cross-referencing, use the command:

    \begin{equation}\label{reio}
    ...
    \end{equation}

To get a cross-reference to an auto-numbered equation, use the command:

    ... using equations~\eqref{ax1} and~\eqref{bz2}, we
    can derive ...

The above example would produce something like

> using equations (3.2) and (3.5), we can derive

In other words, `\eqref{ax1}` is equivalent to `(\ref{ax1})`, but the parentheses produced by are always upright.

To give your equation numbers the form *m.n* (*section-number.equation-number*), use the command in the preamble of your document:

    \numberwithin{equation}{section}

For more details on custom numbering schemes see [@lamport §6.3, §C.8.4].

The environment provides a convenient way to number equations in a group with a subordinate numbering scheme. For example, supposing that the current equation number is , write

    \begin{equation}\label{first}
    a=b+c
    \end{equation}
    some intervening text
    \begin{subequations}\label{grp}
    \begin{align}
    a&=b+c\label{second}\\
    d&=e+f+g\label{third}\\
    h&=i+j\label{fourth}
    \end{align}
    \end{subequations}

to get $$\begin{equation}
\label{first}
a=b+c
\end{equation}$$ some intervening text $$\label{grp}
\begin{align}
a&=b+c\label{second}\\
d&=e+f+g\label{third}\\
h&=i+j\label{fourth}
\end{align}$$ By putting a command immediately after `\begin{subequations}` you can get a reference to the parent number; `\eqref{grp}` from the above example would produce [\[grp\]](#grp){reference-type="eqref" reference="grp"} while `\eqref{second}` would produce [\[second\]](#second){reference-type="eqref" reference="second"}.

An example at `https://tex.stackexchange.com/questions/220001/` shows a variant of the above example, with numbering like (2.1), (2.1a), ..., rather than (2.1), (2.2a), .... This is accomplished by using with a cross-reference to the principal component of the subequation number.

# Math symbols and math fonts {#mathsymbols}

## Classes of math symbols

The symbols in a math formula fall into different classes that correspond more or less to the part of speech each symbol would have if the formula were expressed in words. Certain spacing and positioning cues are traditionally used for the different symbol classes to increase the readability of formulas.

::: center
  --- ------- ------------------------------- --------------------------------
   0  Ord     simple/ordinary ()              $A\;0\;\Phi\;\infty$
   1  Op      prefix operator                 $\sum\;\prod\;\int$
   2  Bin     binary operator (conjunction)   ${+}\;{\cup}\;{\wedge}$
   3  Rel     relation/comparison (verb)      ${=}\;{<}\;{\subset}$
   4  Open    left/opening delimiter          $(\;{[}\;{\lbrace}\;{\langle}$
   5  Close   right/closing delimiter         $)\;{]}\;{\rbrace}\;{\rangle}$
   6  Punct   postfix/punctuation             ${.}\;{,}\;{;}\;{!}$
  --- ------- ------------------------------- --------------------------------
:::

::: notes
The distinction in TeX between class 0 and an additional class 7 has to do only with font selection issues, and it is immaterial here.

Symbols of class 2 (Bin), notably the minus sign $-$, are automatically printed by LaTeX as class 0 (no space) if they do not have a suitable left operande.g., at the beginning of a math formula or after an opening delimiter.
:::

The spacing for a few symbols follows tradition instead of the general rule: although $/$ is (semantically speaking) of class 2, we write $k/2$ with no space around the slash rather than $k\mathbin{/}2$. And compare `p|q` $p\vert q$ (no space) with `p\mid q` $p\mid q$ (class-3 spacing).

The proper way to define a new math symbol is discussed in *font selection* [@fntguide]. It is not really possible to give a useful synopsis here because one needs first to understand the ramifications of font specifications. But supposing one knows that a Cyrillic font named is available, here is a minimal example showing how to define a LaTeX command to print one letter from that font as a math symbol:

    % Declare that the combination of font attributes OT2/wncyr/m/n
    % should select the wncyr font.
    \DeclareFontShape{OT2}{wncyr}{m}{n}{<->wncyr10}{}
    % Declare that the symbolic math font name "cyr" should resolve to
    % OT2/wncyr/m/n.
    \DeclareSymbolFont{cyr}{OT2}{wncyr}{m}{n}
    % Declare that the command \Sh should print symbol 88 from the math font
    % "cyr", and that the symbol class is 0 (= alphabetic = Ord).
    \DeclareMathSymbol{\Sh}{\mathalpha}{cyr}{88}

## Some symbols intentionally omitted here

The following math symbols that are mentioned in the LaTeX book [@lamport] are intentionally omitted from this discussion because they are superseded by equivalent symbols when the package is loaded. If you are using the package anyway, the only thing that you are likely to gain by using the alternate name is an unnecessary increase in the number of fonts used by your document.

::: center
  -- ---------------------------------------
      $\csname square\endcsname$
      $\csname lozenge\endcsname$
      $\csname rightsquigarrow\endcsname$
      $\csname bowtie\endcsname$
      $\csname vartriangleleft\endcsname$
      $\csname trianglelefteq\endcsname$
      $\csname vartriangleright\endcsname$
      $\csname trianglerighteq\endcsname$
  -- ---------------------------------------
:::

Furthermore, there are available for / use above and beyond the ones included here. This list is not intended to be comprehensive. For a much more comprehensive list of symbols, including nonmathematically oriented ones, such as phonetic alphabetic or dingbats, see *The Comprehensive LaTeX Symbol List* [@comprehensive]. (Full font tables, ordered by font name, for all the fonts covered by the comprehensive list are included in the documentation provided by TeX Live: `texdoc rawtables`. These tables do not include symbol names.) Another source of symbol information is the package; see [@uc-math].

## Alphabets and digits {#alpha-digit}

### Latin letters and Arabic numerals

The Latin letters are simple symbols, class 0. The default font for them in math formulas is italic.

::: center
  ---------------------------------------------
     $A\,B\,C\,D\,E\,F\,G\,H\,I\,J\,K\,L\,M%
      \,N\,O\,P\,Q\,R\,S\,T\,U\,V\,W\,X\,Y\,Z$
     $a\,b\,c\,d\,e\,f\,g\,h\,i\,j\,k\,l\,m%
      \,n\,o\,p\,q\,r\,s\,t\,u\,v\,w\,x\,y\,z$
  ---------------------------------------------
:::

When adding an accent to an $i$ or $j$ in math, dotless variants can be obtained with and :

::: symlist
:::

Arabic numerals 0 are also of class 0. Their default font is upright/roman. $$0\,1\,2\,3\,4\,5\,6\,7\,8\,9$$

### Greek letters

Like the Latin letters, the Greek letters are simple symbols, class 0. For obscure historical reasons, the default font for lowercase Greek letters in math formulas is italic while the default font for capital Greek letters is upright/roman. (In other fields such as physics and chemistry, however, the typographical traditions are somewhat different.) The capital Greek letters not present in this list are the letters that have the same appearance as some Latin letter: A for Alpha, B for Beta, and so on. In the list of lowercase letters there is no omicron because it would be identical in appearance to Latin $o$. In practice, the Greek letters that have Latin look-alikes are seldom used in math formulas, to avoid confusion.

::: symlist
:::

### Other "basic" alphabetic symbols

These are also class 0.

::: symlist
:::

::: notes
:::

### Math font switches {#mathfonts}

Not all of the fonts necessary to support comprehensive math font switching are commonly available in a typical LaTeX setup. Here are the results of applying various font switches to a wide range of math symbols when the standard set of Computer Modern fonts is in use. It can be seen that the only symbols that respond correctly to all of the font switches are the uppercase Latin letters. In fact, *nearly all* math symbols apart from Latin letters remain unaffected by font switches; and although the lowercase Latin letters, capital Greek letters, and numerals do respond properly to some font switches, they produce bizarre results for other font switches. (Use of alternative math font sets such as Lucida New Math may ameliorate the situation somewhat.) $$\renewcommand{\arraystretch}{1.3}
\begin{array}{cccccccc}
%\text{default}& \cn{mathbf}& \cn{mathsf}& \cn{mathit}& \cn{mathcal}&
%  \cn{mathscr}&  \cn{mathbb}& \cn{mathfrak}\\
\text{default}& \cn{mathbf}& \cn{mathrm}& \cn{mathsf}& \cn{mathit}&
  \cn{mathcal}& \cn{mathbb}& \cn{mathfrak}\\
\hline
\symrow{X} \\
\symrow{x} \\
\symrow{0} \\
\symrow{[\,]} \\
\symrow{+} \\
\symrow{-} \\
\symrow{=} \\
\symrow{\Xi} \\
\symrow{\xi} \\
\symrow{\infty} \\
\symrow{\aleph} \\
\symrow{\sum}\\
\symrow{\amalg} \\
\symrow{\Re} \\
\end{array}$$ A common desire is to get a bold version of a particular math symbol. For those symbols where is not applicable, the or commands can be used. $$\begin{equation}
A_\infty + \pi A_0
\sim \mathbf{A}_{\boldsymbol{\infty}} \boldsymbol{+}
  \boldsymbol{\pi} \mathbf{A}_{\boldsymbol{0}}
\sim\pmb{A}_{\pmb{\infty}} \pmb{+}\pmb{\pi} \pmb{A}_{\pmb{0}}
\end{equation}$$

    A_\infty + \pi A_0
    \sim \mathbf{A}_{\boldsymbol{\infty}} \boldsymbol{+}
      \boldsymbol{\pi} \mathbf{A}_{\boldsymbol{0}}
    \sim\pmb{A}_{\pmb{\infty}} \pmb{+}\pmb{\pi} \pmb{A}_{\pmb{0}}

The command is obtained preferably by using the package, which provides a newer, more powerful version than the one provided by the package. It is usually ill-advised to apply to more than one symbol at a time; if such a need seems to arise, it more likely means that there is another, better way of going about it.

### Blackboard Bold letters (; no lowercase)

Usage: `\mathbb{R}`. Requires . $$\mathbb{A}\,\mathbb{B}\,\mathbb{C}\,\mathbb{D}\,\mathbb{E}\,\mathbb{F}
\,\mathbb{G}\,\mathbb{H}\,\mathbb{I}\,\mathbb{J}\,\mathbb{K}\,\mathbb{L}
\,\mathbb{M}\,\mathbb{N}\,\mathbb{O}\,\mathbb{P}\,\mathbb{Q}\,\mathbb{R}
\,\mathbb{S}\,\mathbb{T}\,\mathbb{U}\,\mathbb{V}\,\mathbb{W}\,\mathbb{X}
\,\mathbb{Y}\,\mathbb{Z}$$ One lowercase letter is available with a distinct name: $\Bbbk$

### Calligraphic letters (; no lowercase)

Usage: `\mathcal{M}`. $$\mathcal{A}\,\mathcal{B}\,\mathcal{C}\,\mathcal{D}\,\mathcal{E}
\,\mathcal{F}\,\mathcal{G}\,\mathcal{H}\,\mathcal{I}\,\mathcal{J}
\,\mathcal{K}\,\mathcal{L}\,\mathcal{M}\,\mathcal{N}\,\mathcal{O}
\,\mathcal{P}\,\mathcal{Q}\,\mathcal{R}\,\mathcal{S}\,\mathcal{T}
\,\mathcal{U}\,\mathcal{V}\,\mathcal{W}\,\mathcal{X}\,\mathcal{Y}
\,\mathcal{Z}$$

### Non-CM calligraphic and script letters

(; no lowercase) Usage: `\usepackage{mathrsfs}` `\mathscr{B}`. $$\mathscr{A}\,\mathscr{B}\,\mathscr{C}\,\mathscr{D}\,\mathscr{E}
\,\mathscr{F}\,\mathscr{G}\,\mathscr{H}\,\mathscr{I}\,\mathscr{J}
\,\mathscr{K}\,\mathscr{L}\,\mathscr{M}\,\mathscr{N}\,\mathscr{O}
\,\mathscr{P}\,\mathscr{Q}\,\mathscr{R}\,\mathscr{S}\,\mathscr{T}
\,\mathscr{U}\,\mathscr{V}\,\mathscr{W}\,\mathscr{X}\,\mathscr{Y}
\,\mathscr{Z}$$

(; no lowercase) Usage: `\usepackage{euscript}` `\mathscr{E}`. $$\EuScript{A}\,\EuScript{B}\,\EuScript{C}\,\EuScript{D}\,\EuScript{E}
\,\EuScript{F}\,\EuScript{G}\,\EuScript{H}\,\EuScript{I}\,\EuScript{J}
\,\EuScript{K}\,\EuScript{L}\,\EuScript{M}\,\EuScript{N}\,\EuScript{O}
\,\EuScript{P}\,\EuScript{Q}\,\EuScript{R}\,\EuScript{S}\,\EuScript{T}
\,\EuScript{U}\,\EuScript{V}\,\EuScript{W}\,\EuScript{X}\,\EuScript{Y}
\,\EuScript{Z}$$

### Fraktur letters ()

Usage: `\mathfrak{S}`. Requires . $$\mathfrak{A}\,\mathfrak{B}\,\mathfrak{C}\,\mathfrak{D}\,\mathfrak{E}
\,\mathfrak{F}\,\mathfrak{G}\,\mathfrak{H}\,\mathfrak{I}\,\mathfrak{J}
\,\mathfrak{K}\,\mathfrak{L}\,\mathfrak{M}\,\mathfrak{N}\,\mathfrak{O}
\,\mathfrak{P}\,\mathfrak{Q}\,\mathfrak{R}\,\mathfrak{S}\,\mathfrak{T}
\,\mathfrak{U}\,\mathfrak{V}\,\mathfrak{W}\,\mathfrak{X}\,\mathfrak{Y}
\,\mathfrak{Z}$$ $$\mathfrak{a}\,\mathfrak{b}\,\mathfrak{c}\,\mathfrak{d}\,\mathfrak{e}
\,\mathfrak{f}\,\mathfrak{g}\,\mathfrak{h}\,\mathfrak{i}\,\mathfrak{j}
\,\mathfrak{k}\,\mathfrak{l}\,\mathfrak{m}\,\mathfrak{n}\,\mathfrak{o}
\,\mathfrak{p}\,\mathfrak{q}\,\mathfrak{r}\,\mathfrak{s}\,\mathfrak{t}
\,\mathfrak{u}\,\mathfrak{v}\,\mathfrak{w}\,\mathfrak{x}\,\mathfrak{y}
\,\mathfrak{z}$$

## Miscellaneous simple symbols

These symbols are also of class 0 (ordinary) which means they do not have any built-in spacing.

::: symlist
:::

::: notes
A common mistake in the use of the symbols $\square$ and $\#$ is to try to make them serve as binary operators or relation symbols without using a properly defined math symbol command. If you merely use the existing commands or the intersymbol spacing will be incorrect because those commands produce a class-0 symbol.

Synonyms:
:::

## Binary operator symbols

::: symlist
:::

::: notes
, , ,
:::

## Relation symbols: $<$ $=$ $>$ $\succ$ $\sim$ and variants

::: symlist
:::

::: notes
, , , , ,
:::

## Relation symbols: arrows

See also Section [4](#notations){reference-type="ref" reference="notations"}.

::: symlist
:::

::: notes
, ,
:::

## Relation symbols: miscellaneous

::: symlist
:::

::: notes
:::

## Cumulative (variable-size) operators

::: symlist
:::

## Punctuation

::: symlist
:::

::: notes
The `:` by itself produces a colon with class-3 (relation) spacing. The command produces special spacing for use in constructions such as `f\colon A\to B` $f\colon A\to B$.

Although the commands and are frequently used, we recommend the more semantically oriented commands for most purposes (see Section [4.6](#dots){reference-type="ref" reference="dots"}).
:::

## Pairing delimiters (extensible) {#pair-delims}

See Section [6](#delim){reference-type="ref" reference="delim"} for more information.

::: symlist
:::

## Nonpairing extensible symbols

::: symlist
:::

::: notes
Using , `|`, , or for paired delimiters is not recommended (see Section [6.2](#verts){reference-type="ref" reference="verts"}). Instead, use delimiters from the list in Section [3.11](#pair-delims){reference-type="ref" reference="pair-delims"}.
:::

## Extensible vertical arrows

::: symlist
:::

## Math accents {#accents}

::: symlist
:::

## Named operators

These operators are represented by a multiletter abbreviation.

::: symlist
:::

To define additional named operators outside the above list, use the command; for example, after

    \DeclareMathOperator{\rank}{rank}
    \DeclareMathOperator{\esssup}{ess\,sup}

one could write

::: center
  ---------------- -----------------------------------
        `\rank(x)` $\mathop{\mathrm{rank}}(x)$
    `\esssup(y,z)` $\mathop{\mathrm{ess\,sup}}(y,z)$
  ---------------- -----------------------------------
:::

The star form creates an operator that takes limits in a displayed formula, such as $\sup$ or $\max$.

When predefining such a named operator is problematic (e.g., when using one in the title or abstract of an article), there is an alternative form that can be used directly: $$\verb'\operatorname{rank}(x)'\quad
  \rightarrow\quad\operatorname{rank}(x)$$

# Notations

## Top and bottom embellishments

These are visually similar to accents but generally span multiple symbols rather than being applied to a single base symbol. For ease of reference, and are redundantly included here and in the table of math accents.

::: symlist
:::

## Extensible arrows

and produce arrows that extend automatically to accommodate unusually wide subscripts or superscripts. These commands take one optional argument (the subscript) and one mandatory argument (the superscript, possibly empty): $$\begin{equation}
A\xleftarrow{n+\mu-1}B \xrightarrow[T]{n\pm i-1}C
\end{equation}$$

      \xleftarrow{n+\mu-1}\quad \xrightarrow[T]{n\pm i-1}

## Affixing symbols to other symbols

In addition to the standard accents (Section [3.14](#accents){reference-type="ref" reference="accents"}), other symbols can be placed above or below a base symbol with the and commands. For example, writing `\overset{*}{X}` will place a superscript-size $*$ above the $X$, thus: $\overset{*}{X}$. See also the description of in Section [8.4](#sideset){reference-type="ref" reference="sideset"}.

## Matrices {#ss:matrix}

The environments , , , , and have (respectively) $(\,)$, $[\,]$, $\lbrace\,\rbrace$, $\lvert\,\rvert$, and $\lVert\,\rVert$ delimiters built in. There is also a environment without delimiters and an environment that can be used to obtain left alignment or other variations in the column specs.

::::: center
::: minipage
    \begin{pmatrix}
    \alpha& \beta^{*}\\
    \gamma^{*}& \delta
    \end{pmatrix}
:::

::: minipage
$$\begin{pmatrix}
\alpha& \beta^{*}\\
\gamma^{*}& \delta
\end{pmatrix}$$
:::
:::::

To produce a small matrix suitable for use in text, there is a environment (e.g., $\bigl( \begin{smallmatrix}
  a&b\\ c&d
\end{smallmatrix} \bigr)$) that comes closer to fitting within a single text line than a normal matrix. This example was produced by

    \bigl( \begin{smallmatrix}
      a&b\\ c&d
    \end{smallmatrix} \bigr)

By default, all elements in a matrix are centered horizontally. The package provides starred versions of all the matrix environments that facilitate other alignments. That package also provides fenced versions of with parallel names in both starred and nonstarred versions.

To produce a row of dots in a matrix spanning a given number of columns, use . For example, `\hdotsfor{3}` in the second column of a four-column matrix will print a row of dots across the final three columns.

For piecewise function definitions there is a environment:

    P_{r-j}=\begin{cases}
        0&  \text{if $r-j$ is odd},\\
        r!\,(-1)^{(r-j)/2}&  \text{if $r-j$ is even}.
      \end{cases}

Notice the use of and the embedded math.

::: notes
The plain TeX form `\matrix{...\cr...\cr}` and the related commands , should be avoided in LaTeX (and when the package is loaded they are disabled).
:::

## Math spacing commands

When the package is used, all of these math spacing commands can be used both in and out of math mode.

::: center
  Abbrev.   Spelled out   Example      Abbrev.   Spelled out   Example
  --------- ------------- ------------ --------- ------------- --------------------
            no space      $34$                   no space      $34$
                          $3\,4$                               $3\!4$
                          $3\:4$                               $3\negmedspace4$
                          $3\;4$                               $3\negthickspace4$
                          $3\quad4$                            
                          $3\qquad4$                           
:::

For finer control over math spacing, use and 'math units'. One math unit, or `mu`, is equal to 1/18 em. Thus to get a negative half write `\mspace{-9.0mu}`.

There are also three commands that leave a space equal to the height and/or width of a given fragment of / material:

::: center
  ------------------ ----------------------------------------
  `\phantom{XXX}`    space as wide and high as three X's
  `\hphantom{XXX}`   space as wide as three X's; height 0
  `\vphantom{X}`     space of width 0, height = height of X
  ------------------ ----------------------------------------
:::

## Dots

For preferred placement of ellipsis dots (raised or on-line) in various contexts there is no general consensus. It may therefore be considered a matter of taste. In most situations, the generic can be used, and will interpret it in the manner preferred by the AMS, namely low dots () between commas or raised dots () between binary operators and relations, etc. If what follows the dots is ambiguous as to the choice, the specific form of the command can be used. However, by using the semantically oriented commands

- for

- for

- for

- for

- for (none of the above)

instead of and , you make it possible for your document to be adapted to different conventions on the fly, in case (for example) you have to submit it to a publisher who insists on following house tradition in this respect. The default treatment for the various kinds follows American Mathematical Society conventions:

::: center
+:-------------------------------------------+:-------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| ::: minipage                               | ::: minipage                                                                                                                                                             |
|     We have the series $A_1,A_2,\dotsc$,   | We have the series $A_1,A_2,\dotsc$, the regional sum $A_1+A_2+\dotsb$, the orthogonal product $A_1A_2\dotsm$, and the infinite integral $$\int_{A_1}\int_{A_2}\dotsi.$$ |
|     the regional sum $A_1+A_2+\dotsb$,     | :::                                                                                                                                                                      |
|     the orthogonal product $A_1A_2\dotsm$, |                                                                                                                                                                          |
|     and the infinite integral              |                                                                                                                                                                          |
|     \[\int_{A_1}\int_{A_2}\dotsi\].        |                                                                                                                                                                          |
| :::                                        |                                                                                                                                                                          |
+--------------------------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
:::

## Nonbreaking dashes

The command suppresses the possibility of a linebreak after the following hyphen or dash. For example, if you write 'pages 1' as `pages 1\nobreakdash--9` then a linebreak will never occur between the dash and the 9. You can also use to prevent undesirable hyphenations in combinations like `$p$-adic`. For frequent use, it's advisable to make abbreviations, e.g.,

    \newcommand{\p}{$p$\nobreakdash}% for "\p adic" ("p-adic")
    \newcommand{\Ndash}{\nobreakdash\textendash}% for "pages 1\Ndash 9"
    %    For "\n dimensional" ("n-dimensional"):
    \newcommand{\n}{$n$\nobreakdash-\hspace{0pt}}

The last example shows how to prohibit a linebreak after the hyphen but allow normal hyphenation in the following word. (Add a zero-width space after the hyphen.)

## Roots

The command produces a square root. To specify an explicit radix, give it as an optional argument. $$\verb'\sqrt{\frac{n}{n-1} S}'\quad\sqrt{\frac{n}{n-1} S}, \qquad
\verb'\sqrt[3]{2}'\quad
\sqrt[3]{2}$$

## Boxed formulas

The command puts a box around its argument, like except that the contents are in math mode: $$\begin{equation}
\boxed{\eta \leq C(\delta(\eta) +\Lambda_M(0,\delta))}
\end{equation}$$

      \boxed{\eta \leq C(\delta(\eta) +\Lambda_M(0,\delta))}

If you need to box an equation including the equation number, it may be difficult, depending on the context; there are some suggestions in the AMS author FAQ; see the entry outlined in red on the page `https://www.ams.org/faq?faq_id=290`.

# Fractions and related constructions

## The , , and commands

The command takes two arguments numerator and denominatorand typesets them in normal fraction form. Use or to overrule LaTeX's guess about the proper size to use for the fraction's contents (t = text style, d = display style). $$\begin{equation}
\frac{1}{k}\log_2 c(f),\quad\dfrac{1}{k}\log_2 c(f),\quad\tfrac{1}{k}\log_2 c(f)
\end{equation}$$

    \begin{equation}
    \frac{1}{k}\log_2 c(f),\quad\dfrac{1}{k}\log_2 c(f),
        \quad\tfrac{1}{k}\log_2 c(f)
    \end{equation}

$$\begin{equation}
\Re{z} =\frac{n\pi \dfrac{\theta +\psi}{2}}{
        \left(\dfrac{\theta +\psi}{2}\right)^2 + \left( \dfrac{1}{2}
        \log \left\lvert\dfrac{B}{A}\right\rvert\right)^2}.
\end{equation}$$

    \begin{equation}
    \Re{z} =\frac{n\pi \dfrac{\theta +\psi}{2}}{
            \left(\dfrac{\theta +\psi}{2}\right)^2 + \left( \dfrac{1}{2}
            \log \left\lvert\dfrac{B}{A}\right\rvert\right)^2}.
    \end{equation}

## The , , and commands

For binomial expressions such as $\binom{n}{k}$ there are , and commands: $$\begin{equation}
2^k-\binom{k}{1}2^{k-1}+\binom{k}{2}2^{k-2}
\end{equation}$$

    2^k-\binom{k}{1}2^{k-1}+\binom{k}{2}2^{k-2}

## The command

The capabilities of , , and their variants are subsumed by a generalized fraction command with six arguments. The last two correspond to 's numerator and denominator; the first two are optional delimiters (as seen in ); the third is a line thickness override ( uses this to set the fraction line thickness to 0 pti.e., invisible); and the fourth argument is a mathstyle override: integer values 0 select, respectively, , , , and . If the third argument is left empty, the line thickness defaults to "normal".

::: cmdspec
:::

To illustrate, here is how , , and might be defined.

    \newcommand{\frac}[2]{\genfrac{}{}{}{}{#1}{#2}}
    \newcommand{\tfrac}[2]{\genfrac{}{}{}{1}{#1}{#2}}
    \newcommand{\binom}[2]{\genfrac{(}{)}{0pt}{}{#1}{#2}}

::: notes
For technical reasons, using the primitive fraction commands , , in a LaTeX document is not recommended (see, e.g., `https://www.ams.org/faq?faq_id=288`, the entry outlined in red).
:::

## Continued fractions

The continued fraction $$\begin{equation}
\cfrac{1}{\sqrt{2}+
 \cfrac{1}{\sqrt{2}+
  \cfrac{1}{\sqrt{2}+\cdots
}}}
\end{equation}$$ can be obtained by typing

    \cfrac{1}{\sqrt{2}+
     \cfrac{1}{\sqrt{2}+
      \cfrac{1}{\sqrt{2}+\dotsb
    }}}

This produces better-looking results than straightforward use of . Left or right placement of any of the numerators is accomplished by using `[l]` or `[r]` instead of .

# Delimiters {#delim}

## Delimiter sizes {#bigdel}

Unless you indicate otherwise, delimiters in math formulas will remain at the standard size regardless of the height of the enclosed material. To get larger sizes, you can either select a particular size using a prefix (see below), or you can use and prefixes for autosizing.

The automatic delimiter sizing done by and has two limitations: first, it is applied mechanically to produce delimiters large enough to encompass the largest contained item, and second, the range of sizes has fairly large quantum jumps. This means that an expression that is infinitesimally too large for a given delimiter size will get the next larger size, a jump of 6pt or so (3pt top and bottom) in normal-sized text. There are two or three situations where the delimiter size is commonly adjusted. These adjustments are done using the following commands:

::: center
  ------------------------------------------ --------------------------------- ------------------------------------------------------- ----------------------------------------------------- ----------------------------------------------------- --------------------------------------------------------- ---------------------------------------------------------
  Delimiter                                  no size                                                                                                                                                                                                                                                         
  size                                       specified                                                                                                                                                                                                                                                       
                                                                                                                                                                                                                                                                                                             
  Result $\vphantom{\Bigg|^{\frac{1}{2}}}$   $\displaystyle(b)(\frac{c}{d})$   $\displaystyle\left(b\right)\left(\frac{c}{d}\right)$   $\displaystyle\bigl(b\bigr)\bigl(\frac{c}{d}\bigr)$   $\displaystyle\Bigl(b\Bigr)\Bigl(\frac{c}{d}\Bigr)$   $\displaystyle\biggl(b\biggr)\biggl(\frac{c}{d}\biggr)$   $\displaystyle\Biggl(b\Biggr)\Biggl(\frac{c}{d}\Biggr)$
  ------------------------------------------ --------------------------------- ------------------------------------------------------- ----------------------------------------------------- ----------------------------------------------------- --------------------------------------------------------- ---------------------------------------------------------
:::

The first kind of adjustment is done for cumulative operators with limits, such as summation signs. With and the delimiters usually turn out larger than necessary, and using the `Big` or `bigg` sizes instead gives better results: $$\begin{equation*}
\left[\sum_i a_i\left\lvert\sum_j x_{ij}\right\rvert^p\right]^{1/p}
\quad\text{versus}\quad
\biggl[\sum_i a_i\Bigl\lvert\sum_j x_{ij}\Bigr\rvert^p\biggr]^{1/p}
\end{equation*}$$

    \biggl[\sum_i a_i\Bigl\lvert\sum_j x_{ij}\Bigr\rvert^p\biggr]^{1/p}

The second kind of situation is clustered pairs of delimiters, where and make them all the same size (because that is adequate to cover the encompassed material), but what you really want is to make some of the delimiters slightly larger to make the nesting easier to see. $$\begin{equation*}
\left((a_1 b_1) - (a_2 b_2)\right)
\left((a_2 b_1) + (a_1 b_2)\right)
\quad\text{versus}\quad
\bigl((a_1 b_1) - (a_2 b_2)\bigr)
\bigl((a_2 b_1) + (a_1 b_2)\bigr)
\end{equation*}$$

    \left((a_1 b_1) - (a_2 b_2)\right)
    \left((a_2 b_1) + (a_1 b_2)\right)
    \quad\text{versus}\quad
    \bigl((a_1 b_1) - (a_2 b_2)\bigr)
    \bigl((a_2 b_1) + (a_1 b_2)\bigr)

The third kind of situation is a slightly oversize object in running text, such as $\left\lvert\frac{b'}{d'}\right\rvert$ where the delimiters produced by and cause too much line spreading. In that case and can be used to produce delimiters that are larger than the base size but still able to fit within the normal line spacing: $\bigl\lvert\frac{b'}{d'}\bigr\rvert$.

The package provides a feature that can simplify sizing; see the package documentation for details.

## Vertical bar notations {#verts}

The use of the `|` character to produce paired delimiters is not recommended. There is an ambiguity about the directionality of the symbol that will in rare cases produce incorrect spacinge.g., `|k|=|-k|` produces $|k|=|-k|$, and `|\sin x|` produces $|\sin x|$ instead of the correct $\lvert\sin x\rvert$. Using for a and for a whenever they are used in pairs will prevent this problem; compare $\lvert -k\rvert$, produced by `\lvert -k\rvert`. For double bars there are analogous , commands. Recommended practice is to define suitable commands in the document preamble for any paired-delimiter use of vert bar symbols:

    \providecommand{\abs}[1]{\lvert#1\rvert}
    \providecommand{\norm}[1]{\lVert#1\rVert}

whereupon `\abs{z}` would produce $\lvert z\rvert$ and `\norm{v}` would produce $\lVert v\rVert$.

# The command

The main use of the command is for words or phrases in a display. It is similar to in its effects but, unlike , automatically produces subscript-size text if used in a subscript. $$\begin{equation}
f_{[x_{i-1},x_i]} \text{ is monotonic,}
\quad i = 1,\dots,c+1
\end{equation}$$

    f_{[x_{i-1},x_i]} \text{ is monotonic,}
    \quad i = 1,\dots,c+1

##  and its relatives

Commands , , , deal with the special spacing conventions of notation. and are variants of preferred by some authors; omits the parentheses, whereas omits the and retains the parentheses. $$\begin{equation}
\gcd(n,m\bmod n) ;\quad x\equiv y\pmod b
;\quad x\equiv y\mod c ;\quad x\equiv y\pod d
\end{equation}$$

    \gcd(n,m\bmod n) ;\quad x\equiv y\pmod b
    ;\quad x\equiv y\mod c ;\quad x\equiv y\pod d

# Integrals and sums

## Altering the placement of limits

The limits on integrals, sums, and similar symbols are placed either to the side of or above and below the base symbol, depending on convention and context. / has rules for automatically choosing one or the other, and most of the time the results are satisfactory. In the event they are not, there are three / commands that can be used to influence the placement of the limits: , , . Compare

::::: center
::: minipage
$$\int_{\lvert x-x_z(t)\rvert<X_0} z^6(t)\phi(x)$$

    \int_{\abs{x-x_z(t)}<X_0} ...
:::

and

::: minipage
$$\int\limits_{\lvert x-x_z(t)\rvert<X_0} z^6(t)\phi(x)$$

    \int\limits_{\abs{x-x_z(t)}<X_0} ...
:::
:::::

The command should follow immediately after the base symbol to which it applies, and its meaning is: shift the following subscript and/or superscript to the limits position, regardless of the usual convention for this symbol. means to shift them to the side instead, and , which might be used in defining a new kind of base symbol, means to use standard positioning as for the command.

See also the description of the options and in [@amsldoc].

## Multiple integral signs

, , and give multiple integral signs with the spacing between them nicely adjusted, in both text and display style. is an extension of the same idea that gives two integral signs with dots between them. Notice the use of thin space () before $dx$ and friends to clarify the meaning. $$\begin{gather}
\iint\limits_A f(x,y)\,dx\,dy\qquad\iiint\limits_A
f(x,y,z)\,dx\,dy\,dz\\
\iiiint\limits_A
f(w,x,y,z)\,dw\,dx\,dy\,dz\qquad\idotsint\limits_A f(x_1,\dots,x_k)
\end{gather}$$

    \iint\limits_A f(x,y)\,dx\,dy\qquad\iiint\limits_A
    f(x,y,z)\,dx\,dy\,dz\\
    \iiiint\limits_A
    f(w,x,y,z)\,dw\,dx\,dy\,dz\qquad\idotsint\limits_A f(x_1,\dots,x_k)

## Multiline subscripts and superscripts

The command can be used to produce a multiline subscript or superscript: for example

::: center
+:---------------------------+:----------------------------------------------+
| ::: minipage               | $\displaystyle                                |
|     \sum_{\substack{       | \sum_{\substack{0\le i\le m\\ 0<j<n}} P(i,j)$ |
|              0\le i\le m\\ |                                               |
|              0<j<n}}       |                                               |
|       P(i,j)               |                                               |
| :::                        |                                               |
+----------------------------+-----------------------------------------------+
:::

## The command {#sideset}

There's also a command called , for a rather special purpose: putting symbols at the subscript and superscript corners of a symbol like $\sum$ or $\prod$. *Note: The command is only designed for use with large operator symbols; with ordinary symbols the results are unreliable.* With , you can write

::: center
+:---------------------------------------+:-----------------------------------------------+
| ::: minipage                           | $\displaystyle                                 |
|     \sideset{}{'}                      | \sideset{}{'}\sum_{n<k,\;\text{$n$ odd}} nE_n$ |
|       \sum_{n<k,\;\text{$n$ odd}} nE_n |                                                |
| :::                                    |                                                |
+----------------------------------------+------------------------------------------------+
:::

The extra pair of empty braces is explained by the fact that has the capability of putting an extra symbol or symbols at each corner of a large operator; to put an asterisk at each corner of a product symbol, you would type

::: center
+:------------------------------+:---------------------------+
| ::: minipage                  | $\displaystyle             |
|     \sideset{_*^*}{_*^*}\prod | \sideset{_*^*}{_*^*}\prod$ |
| :::                           |                            |
+-------------------------------+----------------------------+
:::

# Changing the size of elements in a formula

The / mechanisms for changing font size inside a math formula are completely different from the ones used outside math formulas. If you try to make something larger in a formula with one of the text commands such as or : $$\text{\large \#}\qquad\verb'{\large \#}'$$ you will get a warning message

    Command \large invalid in math mode

Such an attempt, however, often indicates a misunderstanding of how / math symbols work. If you want a \# symbol analogous to a summation sign in its typographical properties, then in principle the best way to achieve that is to define it as a symbol of type "mathop" with the standard / command (see [@fntguide]). (This entails, however, getting hold of a math font with a suitable text-size/display-size pair, which may not be so easy.)

Consider the expression: $$\frac{\sum_{n > 0} z^n}{\prod_{1\leq k\leq n} (1-q^k)}
\qquad\begin{minipage}{.5\columnwidth}
\begin{verbatim}
\frac{\sum_{n > 0} z^n}
     {\prod_{1\leq k\leq n} (1-q^k)}
\end{verbatim}
\end{minipage}$$ Using instead of wouldn't change anything in this case; if you want the sum and product symbols to appear full size, you need the command: $$\frac{{\displaystyle\sum_{n > 0} z^n}}
     {{\displaystyle\prod_{1\leq k\leq n} (1-q^k)}}
\qquad\begin{minipage}{.7\columnwidth}
\begin{verbatim}
\frac{{\displaystyle\sum_{n > 0} z^n}}
     {{\displaystyle\prod_{1\leq k\leq n} (1-q^k)}}
\end{verbatim}
\end{minipage}$$ And if you want full-size symbols but with limits on the side, use the command also: $$\frac{{\displaystyle\sum\nolimits_{n> 0} z^n}}
  {{\displaystyle\prod\nolimits_{1\leq k\leq n} (1-q^k)}}
\qquad\begin{minipage}{.76\columnwidth}
\begin{verbatim}
\frac{{\displaystyle\sum\nolimits_{n> 0} z^n}}
  {{\displaystyle\prod\nolimits_{1\leq k\leq n} (1-q^k)}}
\end{verbatim}
\end{minipage}$$ There are similar commands , , and , to force / to use the symbol size and spacing that would be applied in (respectively) inline math, first-order subscript, or second-order subscript, even when the current context would normally yield some other size.

**Note:** These commands belong to a special class of commands referred to in the / book as "declarations". In particular, notice where the braces fall that delimit the effect of the command:

::: center
**Right:** `{\displaystyle ...}` **Wrong:** `\displaystyle{...}`
:::

# Other packages of interest {#other-packages}

Many other LaTeX packages that address some aspect of mathematical formulas are available from CTAN (the Comprehensive TeX Archive Network). To recommend a few examples:

::: description
Additional features extending ; loads .

General theorem and proof setup.

Defines and , and provides access to many additional symbols (without names; provides the names).

Under accents and accents using arbitrary symbols.

Bold math package, provides a more general and more robust implementation of .

Ralph Smith's Formal Script, font setup.

Apply a large brace to two or more equations without losing the individual equation numbers.

Delimiters spanning multiple rows of an array.

Commutative diagrams and other diagrams.

Comprehensive graphical facilities, including features for drawing diagrams.
:::

The TeX Catalogue,\
`http://mirror.ctan.org/help/Catalogue/alpha.html`,\
is a good place to look if you know a package's name.

Questions and answers on specific TeX-related topics are the *raison d'être* of this forum:\
`https://tex.stackexchange.com/questions`\
Check the archives for existing answers; pointers to selected topics may expedite your search:\
`https://tex.meta.stackexchange.com/a/2425`\
If nothing useful turns up, ask your own question.

# Other documentation of interest

::: thebibliography
AMUG

American Mathematical Society and the LaTeX3 Project: *User's Guide for the [amsmath]{.nodecor} package*, Version 2.$+$, `http://mirror.ctan.org/macros/latex/required/amsmath/amsldoc.tex` and `http://mirror.ctan.org/macros/latex/required/amsmath/amsldoc.pdf`, 2017.

American Mathematical Society: *User's Guide, AMSFonts*, `http://mirror.ctan.org/fonts/amsfonts/amsfndoc.pdf`, 2002.

Scott Pakin: *The Comprehensive LaTeX Symbol List*, `http://mirror.ctan.org/tex-archive/info/symbols/comprehensive/`, January 2017. Raw font tables, without symbol names, are shown alphabetically by font name in the files in the same area of CTAN and from TeX Live with `texdoc rawtables`.

Leslie Lamport: *LaTeX: A document preparation system*, 2nd edition, Addison-Wesley, 1994.

Frank Mittelbach and Michel Goossens, with Johannes Braams, David Carlisle, and Chris Rowley: *The LaTeX Companion*, 2nd edition, Addison-Wesley, 2004.

LaTeX3 Project Team: *font selection*, `http://mirror.ctan.org/macros/latex/doc/fntguide.pdf`, 2005.

Michel Goossens, Frank Mittelbach, Sebastian Rahtz, Denis Roegel, and Herbert Voß: *The LaTeX Graphics Companion*, 2nd edition, Addison-Wesley, 2008.

D. P. Carlisle, LaTeX3 Project: *Packages in the 'graphics' bundle*, `http://mirror.ctan.org/macros/latex/required/graphics/grfguide.pdf`, 2017.

LaTeX3 Project Team: * for authors*, `http://mirror.ctan.org/macros/latex/doc/usrguide.pdf`, 2015.

George Grätzer: *More Math into LaTeX*, 5th edition, Springer, New York, 2016.

Will Robertson: *Every symbol [(]{.upright}most symbols[)]{.upright} defined by* , `http://mirror.ctan.org/macros/latex/contrib/unicode-math/unimath-symbols.pdf`, 2017; and\
Will Robertson, Philipp Stephani, Joseph Wright, and Khaled Hosny: *Experimental Unicode mathematical typesetting: The package*, `http://mirror.ctan.org/macros/latex/contrib/unicode-math/unicode-math.pdf`, 2017.
:::
