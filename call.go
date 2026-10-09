package gofactory

import (
	"go/ast"
	"go/types"
)

// checkCall reports new(T) of a protected type and explicit conversions to a
// protected type, including parenthesised ones. new(expr) (Go 1.26) and
// calls that are not conversions at all are left alone.
func (d *detector) checkCall(call *ast.CallExpr) {
	if len(call.Args) != 1 {
		return
	}

	if d.checkNew(call) {
		return
	}

	d.checkConversion(call)
}

// checkNew returns true when call is a call to the builtin new, so the
// caller does not also try to treat it as a conversion.
func (d *detector) checkNew(call *ast.CallExpr) bool {
	if d.builtinName(call) != "new" {
		return false
	}

	argTV := d.pass.TypesInfo.Types[call.Args[0]]
	if !argTV.IsType() {
		// new(expr), Go 1.26: allocates from a value, nothing to bypass.
		return true
	}

	// new(*T) allocates a nil *T and builds no T, so argTV.Type is passed
	// as-is rather than through pointee.
	d.reportProtected(call, argTV.Type)

	if d.zeroValues {
		d.reportUnsetFields(call, argTV.Type, nil)
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
// still be reported. Known false negative: go/types also records the target
// type on an untyped non-constant argument, such as a comparison
// (Flag(a == b)) or a shift of an untyped constant (MyInt(1 << n)), so it
// passes as a no-op too.
func (d *detector) checkConversion(call *ast.CallExpr) {
	funTV := d.pass.TypesInfo.Types[call.Fun]
	if !funTV.IsType() {
		return
	}

	argTV := d.pass.TypesInfo.Types[call.Args[0]]
	if argTV.IsNil() {
		return
	}

	isConst := argTV.Value != nil
	if !isConst && types.Identical(funTV.Type, argTV.Type) {
		return
	}

	d.reportProtected(call, pointee(funTV.Type))
}
