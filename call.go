package gofactory

import (
	"go/ast"
	"go/types"
)

// checkCall reports new(T) of a protected type and explicit conversions to a
// protected type, including parenthesised ones. new(expr) (Go 1.26) and
// calls that are not conversions at all are left alone.
func (v *detector) checkCall(call *ast.CallExpr) {
	if len(call.Args) != 1 {
		return
	}

	if v.checkNew(call) {
		return
	}

	v.checkConversion(call)
}

// checkNew reports new(T) and returns true when call is a call to the
// builtin new, so the caller does not also try to treat it as a conversion.
func (v *detector) checkNew(call *ast.CallExpr) bool {
	ident, ok := call.Fun.(*ast.Ident)
	if !ok {
		return false
	}

	builtin, ok := v.pass.TypesInfo.ObjectOf(ident).(*types.Builtin)
	if !ok || builtin.Name() != "new" {
		return false
	}

	argTV := v.pass.TypesInfo.Types[call.Args[0]]
	if !argTV.IsType() {
		// new(expr), Go 1.26: allocates from a value, nothing to bypass.
		return true
	}

	named, ok := protectedNamed(argTV.Type)
	if ok && v.isProtected(named) {
		v.report(call, named)
	}

	return true
}

// checkConversion reports a type conversion to a protected type: Go's T(x)
// syntax, where T names a type rather than a function or value, including
// parenthesised forms like (*T)(x). go/types marks this by recording
// call.Fun itself as denoting a type (IsType() below) instead of resolving
// it to a function signature, which is also what tells it apart from an
// ordinary call or from new(T) (checkNew, handled separately before this
// is reached). A conversion whose argument is nil, or whose argument
// already has the target's exact type and is not a constant, creates
// nothing and stays silent; an untyped constant looks identical but must
// still be reported.
func (v *detector) checkConversion(call *ast.CallExpr) {
	funTV := v.pass.TypesInfo.Types[call.Fun]
	if !funTV.IsType() {
		return
	}

	named, ok := protectedNamed(funTV.Type)
	if !ok {
		return
	}

	argTV := v.pass.TypesInfo.Types[call.Args[0]]
	if argTV.IsNil() {
		return
	}

	isConst := argTV.Value != nil
	if !isConst && types.Identical(funTV.Type, argTV.Type) {
		return
	}

	if !v.isProtected(named) {
		return
	}

	v.report(call, named)
}
