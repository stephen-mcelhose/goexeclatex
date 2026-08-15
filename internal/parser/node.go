package parser

// Node is the AST node interface (spec §7).
type Node interface {
	node() // unexported marker
}

// NumberNode represents a literal numeric value (spec §7.1).
type NumberNode struct {
	Value float64
	Pos   int
}

func (n *NumberNode) node() {}

// SymbolNode represents a variable or arity-0 constant command (spec §7.2).
type SymbolNode struct {
	Name string // normalised: backslash-stripped, lowercased
	Pos  int
}

func (n *SymbolNode) node() {}

// BinaryNode represents a binary infix operation (spec §7.3).
type BinaryNode struct {
	Op    string // "+", "-", "*", "/", "^"
	Left  Node
	Right Node
}

func (n *BinaryNode) node() {}

// UnaryNode represents a prefix or postfix unary operation (spec §7.4).
// Negation ("-") is prefix; factorial ("!") is postfix.
type UnaryNode struct {
	Op      string // "-" or "!"
	Operand Node
}

func (n *UnaryNode) node() {}

// FunctionNode represents an arity ≥ 1 command or absolute value (spec §7.5).
type FunctionNode struct {
	Name string // normalised command name, e.g. "frac", "sqrt", "abs"
	Args []Node
	Pos  int
}

func (n *FunctionNode) node() {}

// SubscriptNode represents a subscripted variable x_{i} (spec §5.1).
type SubscriptNode struct {
	Base Node // the symbol being subscripted
	Sub  Node // subscript index expression
}

func (n *SubscriptNode) node() {}

// BigOpNode represents \sum or \prod over a discrete range (spec §5.2).
type BigOpNode struct {
	Op   string // "sum" or "prod"
	Var  string // iteration variable name (e.g. "i")
	From Node   // lower bound expression
	To   Node   // upper bound expression
	Body Node   // body expression (evaluated once per step)
}

func (n *BigOpNode) node() {}

// NormNode represents \lVert expr \rVert — scalar absolute value (spec §5.3).
type NormNode struct {
	Arg Node
}

func (n *NormNode) node() {}
