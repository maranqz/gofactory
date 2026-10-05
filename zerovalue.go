package gofactory

import (
	"go/ast"
	"go/token"
	"go/types"
)

const zeroValueSuffix = ": zero value"

// A package-level var has no function to decide a first interaction in, so
// it is always a candidate.
func (d *detector) checkPackageVars(file *ast.File) {
	if !d.zeroValues {
		return
	}

	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			continue
		}

		for _, name := range zeroValueSpecNames(genDecl) {
			d.reportProtectedSuffix(name, d.pass.TypesInfo.TypeOf(name), zeroValueSuffix)
		}
	}
}

type zeroVar struct {
	typ      types.Type
	isResult bool
}

func (d *detector) checkFuncZeroValues(
	fnType *ast.FuncType, body *ast.BlockStmt,
) {
	if !d.zeroValues || body == nil {
		return
	}

	tracked := d.collectZeroVars(fnType, body)
	if len(tracked) == 0 {
		return
	}

	walk := firstInteractionWalk{
		uses:    d.pass.TypesInfo.Uses,
		tracked: tracked,
		first:   map[types.Object]interaction{},
	}
	ast.Inspect(body, walk.visit)

	for obj, first := range walk.first {
		if !first.safe {
			d.reportProtectedSuffix(first.node, tracked[obj].typ, zeroValueSuffix)
		}
	}
}

func (d *detector) collectZeroVars(
	fnType *ast.FuncType, body *ast.BlockStmt,
) map[types.Object]zeroVar {
	tracked := map[types.Object]zeroVar{}

	// A named result called `_` (func F() (_ T, err error)) still has an
	// object, and only a naked return can reach it.
	if fnType.Results != nil {
		for _, field := range fnType.Results.List {
			for _, name := range field.Names {
				if obj := d.pass.TypesInfo.ObjectOf(name); obj != nil {
					tracked[obj] = zeroVar{typ: obj.Type(), isResult: true}
				}
			}
		}
	}

	ast.Inspect(body, func(node ast.Node) bool {
		// A function literal is checked as its own function, with its own
		// locals and results.
		if _, ok := node.(*ast.FuncLit); ok {
			return false
		}

		genDecl, ok := node.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			return true
		}

		for _, name := range zeroValueSpecNames(genDecl) {
			if obj := d.pass.TypesInfo.ObjectOf(name); obj != nil {
				tracked[obj] = zeroVar{typ: obj.Type()}
			}
		}

		return true
	})

	return tracked
}

// A var with a value is not a zero value; a literal or conversion in it is
// checked by its own route.
func zeroValueSpecNames(genDecl *ast.GenDecl) []*ast.Ident {
	var names []*ast.Ident

	for _, spec := range genDecl.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		if !ok || valueSpec.Values != nil {
			continue
		}

		for _, name := range valueSpec.Names {
			if name.Name != "_" {
				names = append(names, name)
			}
		}
	}

	return names
}

type interaction struct {
	node ast.Node
	safe bool
}

// firstInteractionWalk records each tracked object's first interaction, in
// evaluation order rather than source order: an assignment's right-hand side
// before its targets, a range expression before its key and value, and a for
// loop's body before its post statement.
type firstInteractionWalk struct {
	uses    map[*ast.Ident]types.Object
	tracked map[types.Object]zeroVar
	first   map[types.Object]interaction
	inLit   bool
}

func (w *firstInteractionWalk) visit(node ast.Node) bool {
	switch node := node.(type) {
	case *ast.FuncLit:
		w.visitFuncLit(node)
	case *ast.AssignStmt:
		w.visitAssign(node)
	case *ast.RangeStmt:
		w.inspect(node.X)
		w.target(node.Key, true)
		w.target(node.Value, true)
		w.inspect(node.Body)
	case *ast.ForStmt:
		w.inspect(node.Init)
		w.inspect(node.Cond)
		w.inspect(node.Body)
		w.inspect(node.Post)
	case *ast.CallExpr:
		w.visitCall(node)
	case *ast.ReturnStmt:
		w.recordNakedReturn(node)

		return true
	case *ast.Ident:
		w.record(node, false)

		return true
	default:
		return true
	}

	return false
}

func (w *firstInteractionWalk) visitFuncLit(lit *ast.FuncLit) {
	inLit := w.inLit
	w.inLit = true
	ast.Inspect(lit.Body, w.visit)
	w.inLit = inLit
}

func (w *firstInteractionWalk) visitAssign(assign *ast.AssignStmt) {
	wholeValue := assign.Tok == token.ASSIGN || assign.Tok == token.DEFINE

	for _, rhs := range assign.Rhs {
		w.inspect(rhs)
	}

	for _, lhs := range assign.Lhs {
		w.target(lhs, wholeValue)
	}
}

func (w *firstInteractionWalk) visitCall(call *ast.CallExpr) {
	w.inspect(call.Fun)

	for _, arg := range call.Args {
		if ident, ok := addressedIdent(arg); ok {
			w.record(ident, true)
		} else {
			w.inspect(arg)
		}
	}
}

func (w *firstInteractionWalk) inspect(node ast.Node) {
	if node != nil {
		ast.Inspect(node, w.visit)
	}
}

func (w *firstInteractionWalk) target(expr ast.Expr, safe bool) {
	if ident, ok := expr.(*ast.Ident); ok {
		w.record(ident, safe)

		return
	}

	w.inspect(expr)
}

func (w *firstInteractionWalk) record(ident *ast.Ident, safe bool) {
	obj := w.uses[ident]
	if _, ok := w.tracked[obj]; !ok {
		return
	}

	if _, ok := w.first[obj]; !ok {
		w.first[obj] = interaction{node: ident, safe: safe}
	}
}

func (w *firstInteractionWalk) recordNakedReturn(ret *ast.ReturnStmt) {
	// A naked return in a function literal returns from the literal, not
	// from the function whose results are tracked.
	if len(ret.Results) != 0 || w.inLit {
		return
	}

	for obj, v := range w.tracked {
		if _, ok := w.first[obj]; v.isResult && !ok {
			w.first[obj] = interaction{node: ret}
		}
	}
}

func addressedIdent(arg ast.Expr) (*ast.Ident, bool) {
	unary, ok := ast.Unparen(arg).(*ast.UnaryExpr)
	if !ok || unary.Op != token.AND {
		return nil, false
	}

	ident, ok := unary.X.(*ast.Ident)

	return ident, ok
}
