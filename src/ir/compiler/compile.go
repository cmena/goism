package compiler

import (
	"ir/instr"
	"lisp"
	"lisp/function"
	"sexp"
)

func compileStmtList(cl *Compiler, forms []sexp.Form) {
	for _, form := range forms {
		compileStmt(cl, form)
	}
}

func compileExprList(cl *Compiler, forms []sexp.Form) {
	for _, form := range forms {
		compileExpr(cl, form)
	}
}

func compileBlock(cl *Compiler, form *sexp.Block) {
	compileStmtList(cl, form.Forms)
	cl.st.Discard(uint16(form.Scope.Len()))
}

func compileReturn(cl *Compiler, form *sexp.Return) {
	if len(form.Results) == 0 {
		// Any function in Emacs Lisp must return a value.
		// To avoid Emacs crash, we always return "nil" for void functions.
		emit(cl, instr.ConstRef(cl.cvec.InsertSym("nil")))
		emit(cl, instr.Return)
	} else {
		compileExpr(cl, form.Results[0])
		for i := 1; i < len(form.Results); i++ {
			compileExpr(cl, form.Results[i])
			sym := lisp.RetVars[i]
			emit(cl, instr.VarSet(cl.cvec.InsertSym(sym)))
		}
		emit(cl, instr.Return)
	}
}

func compileIf(cl *Compiler, form *sexp.If) {
	if form.Else == nil {
		elseLabel := labelCreate(cl, "else")
		compileExpr(cl, form.Cond)
		emitJmpNil(cl, elseLabel)
		compileBlock(cl, form.Then)
		labelBind(cl, elseLabel)
	} else {
		elseLabel := labelCreate(cl, "else")
		endifLabel := labelCreate(cl, "endif")
		compileExpr(cl, form.Cond)
		emitJmpNil(cl, elseLabel)
		compileBlock(cl, form.Then)
		emitJmp(cl, endifLabel)
		labelBind(cl, elseLabel)
		compileStmt(cl, form.Else)
		labelBind(cl, endifLabel)
	}
}

func compileWhile(cl *Compiler, form *sexp.While) {
	condLabel := labelCreate(cl, "loop-cond")
	bodyLabel := labelCreate(cl, "loop-body")
	emitJmp(cl, condLabel)
	labelBind(cl, bodyLabel)
	compileBlock(cl, form.Body)
	labelBind(cl, condLabel)
	compileExpr(cl, form.Cond)
	emitJmpNotNil(cl, bodyLabel)
}

func compileBind(cl *Compiler, form *sexp.Bind) {
	compileExpr(cl, form.Init)
	cl.st.Bind(form.Name)
}

func compileRebind(cl *Compiler, form *sexp.Rebind) {
	compileExpr(cl, form.Expr)
	if stIndex := cl.st.Find(form.Name); stIndex != -1 {
		emit(cl, instr.StackSet(stIndex))
	} else {
		emit(cl, instr.VarSet(cl.cvec.InsertSym(form.Name)))
	}
}

func compileCallStmt(cl *Compiler, form sexp.CallStmt) {
	compileCall(cl, form.Call)

	if !form.Fn.IsPanic() && !form.Fn.IsVoid() {
		emit(cl, instr.Discard(1))
	}
}

func compileBinOp(cl *Compiler, ins instr.Instr, args [2]sexp.Form) {
	compileExpr(cl, args[0])
	compileExpr(cl, args[1])
	emit(cl, ins)
}

func compileVariadicOp(cl *Compiler, ins instr.Instr, args []sexp.Form) {
	compileExprList(cl, args)
	emit(cl, ins)
}

func compileUnaryOp(cl *Compiler, ins instr.Instr, arg sexp.Form) {
	compileExpr(cl, arg)
	emit(cl, ins)
}

func compileUnaryOps(cl *Compiler, instr instr.Instr, arg sexp.Form, n int64) {
	compileExpr(cl, arg)
	for i := int64(0); i < n; i++ {
		emit(cl, instr)
	}
}

func compileBool(cl *Compiler, form sexp.Bool) {
	if form.Val {
		emit(cl, instr.ConstRef(cl.cvec.InsertSym("t")))
	} else {
		emit(cl, instr.ConstRef(cl.cvec.InsertSym("nil")))
	}
}

func compileVar(cl *Compiler, form sexp.Var) {
	if stIndex := cl.st.Find(form.Name); stIndex != -1 {
		emit(cl, instr.StackRef(stIndex))
	} else {
		emit(cl, instr.VarRef(cl.cvec.InsertSym(form.Name)))
	}
	cl.st.Bind(form.Name)
}

func compileSparseArrayLit(cl *Compiler, form *sexp.SparseArrayLit) {
	compileCall(cl, form.Ctor)
	for _, val := range form.Vals {
		emit(cl, instr.StackRef(0)) // Array
		emit(cl, instr.ConstRef(cl.cvec.InsertInt(val.Index)))
		compileExpr(cl, val.Expr)
		emit(cl, instr.ArraySet)
	}
}

func compileArrayCopy(cl *Compiler, form *sexp.ArrayCopy) {
	emit(cl, instr.ConstRef(cl.cvec.InsertSym("copy-sequence")))
	compileExpr(cl, form.Array)
	emit(cl, instr.Call(1))
}

func compilePanic(cl *Compiler, form *sexp.Panic) {
	call(cl, function.Panic, form.ErrorData)
}

func compileCall(cl *Compiler, form *sexp.Call) {
	emit(cl, instr.ConstRef(cl.cvec.InsertSym(form.Fn.Name())))
	compileExprList(cl, form.Args)
	emit(cl, instr.Call(len(form.Args)))

	if form.Fn.IsPanic() {
		cl.st.Discard(1)
	} else if form.Fn.IsVoid() {
		emit(cl, instr.Discard(1))
	}
}

func compileMultiValueRef(cl *Compiler, form *sexp.MultiValueRef) {
	if form.Index+1 > len(lisp.RetVars) {
		panic("too many return values")
	}
	sym := lisp.RetVars[form.Index]
	emit(cl, instr.VarRef(cl.cvec.InsertSym(sym)))
}

func compileArrayIndex(cl *Compiler, form *sexp.ArrayIndex) {
	compileExpr(cl, form.Array)
	compileExpr(cl, form.Index)
	emit(cl, instr.ArrayRef)
}

func compileArrayUpdate(cl *Compiler, form *sexp.ArrayUpdate) {
	compileExpr(cl, form.Array)
	compileExpr(cl, form.Index)
	compileExpr(cl, form.Expr)
	emit(cl, instr.ArraySet)
}
