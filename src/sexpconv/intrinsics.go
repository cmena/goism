package sexpconv

import (
	"assert"
	"go/ast"
	"go/constant"
	"magic_pkg/emacs/lisp"
	"magic_pkg/emacs/rt"
	"sexp"
)

func (conv *converter) lispObjectMethod(fn string, recv ast.Expr, args []ast.Expr) sexp.Form {
	switch fn {
	case "Bool":
		return conv.call(rt.FnCoerceBool, recv)
	case "Int":
		return conv.call(rt.FnCoerceInt, recv)
	case "Float":
		return conv.call(rt.FnCoerceFloat, recv)
	case "String":
		return conv.call(rt.FnCoerceString, recv)
	case "Symbol":
		return conv.call(rt.FnCoerceSymbol, recv)
	}

	assert.Unreachable()
	return nil
}

func (conv *converter) intrinFuncCall(sym string, args []ast.Expr) sexp.Form {
	switch sym {
	case "Int", "Float", "Str", "Symbol", "Bool":
		// These types can be constructed only from
		// typed values that does not require any
		// conversion, so we ignore them.
		return conv.Expr(args[0])

	case "DynCall":
		return &sexp.DynCall{
			Callable: conv.Expr(args[0]),
			Args:     conv.exprList(args[1:]),
			Typ:      lisp.TypObject,
		}

	case "Call":
		cv := conv.valueOf(args[0])
		if cv != nil {
			name := constant.StringVal(cv)
			return conv.lispApply(
				lisp.InternFunc(name),
				conv.exprList(args[1:]),
			)
		}
		return &sexp.DynCall{
			Callable: conv.Expr(args[0]),
			Args:     conv.exprList(args[1:]),
			Typ:      lisp.TypObject,
		}

	case "Intern":
		return conv.intrinIntern(args[0])

	default:
		fn := lisp.FFI[sym]
		args := conv.exprList(args)
		return &sexp.LispCall{Fn: fn, Args: args}
	}
}

func (conv *converter) intrinIntern(arg ast.Expr) sexp.Form {
	if cv := conv.valueOf(arg); cv != nil {
		s := constant.StringVal(cv)
		if s == "" {
			return sexp.Symbol{Val: "##"}
		}
		return sexp.Symbol{Val: s}
	}

	return conv.lispCall(lisp.FnIntern, arg)
}
