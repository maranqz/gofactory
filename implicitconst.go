package gofactory

import (
	"go/ast"
	"go/token"
	"go/types"
	"slices"
)

func (d *detector) checkVarConstants(decl *ast.GenDecl) {
	if decl.Tok != token.VAR {
		return
	}

	for _, spec := range decl.Specs {
		if valueSpec, ok := spec.(*ast.ValueSpec); ok {
			d.checkStoredConstants(valueSpec.Values)
		}
	}
}

// checkAssignedConstants leaves op-assignments such as st += 1 to
// arithmetic, which is not a storing position.
func (d *detector) checkAssignedConstants(assign *ast.AssignStmt) {
	if assign.Tok == token.ASSIGN || assign.Tok == token.DEFINE {
		d.checkStoredConstants(assign.Rhs)
	}
}

// checkArgConstants skips conversions, which the conversion route reports,
// and every builtin but append: the others (min, max, delete, ...) do not
// store their arguments.
func (d *detector) checkArgConstants(call *ast.CallExpr) {
	fun := ast.Unparen(call.Fun)
	if d.pass.TypesInfo.Types[fun].IsType() {
		return
	}

	if name := d.builtinName(call); name != "" && name != "append" {
		return
	}

	d.checkStoredConstants(call.Args)
}

func (d *detector) checkElementConstants(lit *ast.CompositeLit) {
	for _, elt := range lit.Elts {
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			d.checkStoredConstants([]ast.Expr{kv.Key, kv.Value})

			continue
		}

		d.checkStoredConstants([]ast.Expr{elt})
	}
}

// go/types records an implicitly converted constant with its target type:
// in st = 3, 3 has the type of st.
func (d *detector) checkStoredConstants(exprs []ast.Expr) {
	for _, expr := range exprs {
		if d.isUntypedConstant(expr) {
			d.reportProtected(expr, d.pass.TypesInfo.TypeOf(expr))
		}
	}
}

// isUntypedConstant tells an untyped constant from a typed one by its
// operands, since go/types records both with their final, typed type.
func (d *detector) isUntypedConstant(expr ast.Expr) bool {
	if d.pass.TypesInfo.Types[expr].Value == nil {
		return false
	}

	switch expr := ast.Unparen(expr).(type) {
	case *ast.BasicLit:
		return true
	case *ast.Ident:
		return isUntypedConstObject(d.pass.TypesInfo.Uses[expr])
	case *ast.SelectorExpr:
		return isUntypedConstObject(d.pass.TypesInfo.Uses[expr.Sel])
	case *ast.UnaryExpr:
		return d.isUntypedConstant(expr.X)
	case *ast.BinaryExpr:
		return d.isUntypedBinary(expr)
	case *ast.CallExpr:
		return d.isUntypedBuiltinCall(expr)
	}

	return false
}

// A comparison is untyped even when its operands are typed; a shift takes
// the type of its left operand.
func (d *detector) isUntypedBinary(expr *ast.BinaryExpr) bool {
	comparisons := []token.Token{
		token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ,
	}
	if slices.Contains(comparisons, expr.Op) {
		return true
	}

	if expr.Op == token.SHL || expr.Op == token.SHR {
		return d.isUntypedConstant(expr.X)
	}

	return d.isUntypedConstant(expr.X) && d.isUntypedConstant(expr.Y)
}

// len, cap and unsafe.Sizeof give typed constants.
func (d *detector) isUntypedBuiltinCall(call *ast.CallExpr) bool {
	switch d.builtinName(call) {
	case "min", "max", "complex", "real", "imag":
	default:
		return false
	}

	for _, arg := range call.Args {
		if !d.isUntypedConstant(arg) {
			return false
		}
	}

	return true
}

func (d *detector) builtinName(call *ast.CallExpr) string {
	ident, ok := ast.Unparen(call.Fun).(*ast.Ident)
	if !ok {
		return ""
	}

	builtin, ok := d.pass.TypesInfo.Uses[ident].(*types.Builtin)
	if !ok {
		return ""
	}

	return builtin.Name()
}

func isUntypedConstObject(obj types.Object) bool {
	constant, ok := obj.(*types.Const)
	if !ok {
		return false
	}

	basic, ok := constant.Type().(*types.Basic)

	return ok && basic.Info()&types.IsUntyped != 0
}
