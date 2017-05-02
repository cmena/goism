// Package sexp provides a high level intermediate representation
// that contains both Go and Emacs Lisp traits.
package sexp

type Form interface {
	// IsAtom returns true if underlying implementations
	// belongs to atom category.
	//
	// #FIXME: this method can be obsolete.
	// Its currently unused and probably should be
	// as free function if ever needed.
	IsAtom() bool
}

// Var - reference to lexical variable.
// #FIXME: what category does Var belong to?
type Var struct{ Name string }

/* Atoms */

type Bool struct{ Val bool }
type Char struct{ Val rune }
type Int struct{ Val int64 }
type Float struct{ Val float64 }
type String struct{ Val string }

/* Composite literals */

type ArrayLit struct {
	Vals []Form
}

type QuotedArray struct {
	Vals []Form
}

/* Special forms */

// Block can be described as "let-like form",
// but unlike classical Lisp let, it is a statement,
// not expression.
type Block struct {
	// Forms form block body.
	Forms []Form
	// Bindings that extend block lexical environment.
	Locals []Binding
}

// If statement evaluates test expression and,
// depending on the result, one of the branches gets
// executed. Else branch is optional.
type If struct {
	Test Form
	Then Form
	Else Form
}

// Return statement exits the function and returns
// one or more values to the caller.
type Return struct {
	Results []Form
}

/* Builtin ops */

type (
	IntAdd       struct{ Args []Form }
	IntSub       struct{ Args []Form }
	IntMul       struct{ Args []Form }
	IntDiv       struct{ Args []Form }
	IntBitOr     struct{ Args []Form }
	IntBitAnd    struct{ Args []Form }
	IntBitXor    struct{ Args []Form }
	IntRem       struct{ Args []Form }
	IntEq        struct{ Args []Form }
	IntNotEq     struct{ Args []Form }
	IntLess      struct{ Args []Form }
	IntLessEq    struct{ Args []Form }
	IntGreater   struct{ Args []Form }
	IntGreaterEq struct{ Args []Form }

	FloatAdd       struct{ Args []Form }
	FloatSub       struct{ Args []Form }
	FloatMul       struct{ Args []Form }
	FloatDiv       struct{ Args []Form }
	FloatEq        struct{ Args []Form }
	FloatNotEq     struct{ Args []Form }
	FloatLess      struct{ Args []Form }
	FloatLessEq    struct{ Args []Form }
	FloatGreater   struct{ Args []Form }
	FloatGreaterEq struct{ Args []Form }

	Concat          struct{ Args []Form }
	StringEq        struct{ Args []Form }
	StringNotEq     struct{ Args []Form }
	StringLess      struct{ Args []Form }
	StringLessEq    struct{ Args []Form }
	StringGreater   struct{ Args []Form }
	StringGreaterEq struct{ Args []Form }
)

/* Call expressions */

// Call expression is normal (direct) function invocation.
type Call struct {
	Fn   string
	Args []Form
}

/* Helper types (not forms themself) */

// Binding represents named value.
type Binding struct {
	Name string
	Init Form
}
