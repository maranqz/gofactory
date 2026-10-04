package gofactory

import (
	"go/ast"
	"go/token"
	"go/types"
)

// zeroValueSuffix marks the zero-value route in a diagnostic, after the
// "Use factory for pkg.T" prefix every route shares.
const zeroValueSuffix = ": zero value"

// checkPackageVars reports every package-level `var x T` of a protected
// type. Unlike a local var or a named result, a package-level var has no
// function scope to decide a first interaction in, so it is always a
// candidate bypass; the usual permission policy (owner package, fences)
// still decides whether it is actually reported.
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

// zeroVar is a local var or named result tracked for its first interaction.
type zeroVar struct {
	typ      types.Type
	isResult bool
}

// checkFuncZeroValues reports a local `var x T` or named result of fnType
// whose first interaction anywhere in the function is not a whole-value
// assignment or &x passed to a call. It is called once per function
// entity (a *ast.FuncDecl or a *ast.FuncLit); a nested function literal is
// visited separately, with its own locals and results.
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

	safeIdents := collectSafeIdents(body)
	firstIdent := firstIdentUses(d.pass.TypesInfo.Uses, body, tracked)
	nakedReturn := firstNakedReturn(body)

	for obj, v := range tracked {
		node, safe := firstInteraction(v, firstIdent[obj], nakedReturn, safeIdents)
		if node == nil || safe {
			continue
		}

		d.reportProtectedSuffix(node, v.typ, zeroValueSuffix)
	}
}

// collectZeroVars gathers fnType's named results and every local `var x T`
// declared directly in body, i.e. not inside a nested function literal:
// that literal is its own function entity, visited separately by the
// inspector.
func (d *detector) collectZeroVars(
	fnType *ast.FuncType, body *ast.BlockStmt,
) map[types.Object]zeroVar {
	tracked := map[types.Object]zeroVar{}

	if fnType.Results != nil {
		for _, field := range fnType.Results.List {
			for _, name := range field.Names {
				if obj := d.objectOf(name); obj != nil {
					tracked[obj] = zeroVar{typ: obj.Type(), isResult: true}
				}
			}
		}
	}

	ast.Inspect(body, func(n ast.Node) bool {
		if _, ok := n.(*ast.FuncLit); ok {
			return false
		}

		genDecl, ok := n.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			return true
		}

		for _, name := range zeroValueSpecNames(genDecl) {
			if obj := d.objectOf(name); obj != nil {
				tracked[obj] = zeroVar{typ: obj.Type()}
			}
		}

		return true
	})

	return tracked
}

// objectOf resolves name's object, blank named results included: a named
// result called `_` (func F() (_ T, err error)) still has an object in
// TypesInfo.Defs, and can only be reached through a naked return, which
// firstInteraction already treats as the one interaction with every named
// result. A blank local var is filtered earlier, by zeroValueSpecNames.
func (d *detector) objectOf(name *ast.Ident) types.Object {
	return d.pass.TypesInfo.ObjectOf(name)
}

// zeroValueSpecNames returns the declared names of every `var x T` (or
// `var x, y T`) spec in genDecl that has no initializer: a var with a value
// is a literal, conversion or other route, not a zero value.
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

// collectSafeIdents marks every *ast.Ident that occurs as the whole target
// of a whole-value assignment or as &x passed to any call: the two OK
// interactions for CONTEXT.md's First interaction. Every other mention of
// a tracked variable's identifier is left unmarked, so it is reported if
// it is the first interaction.
//
// A whole-value assignment is `x = …`, `x, err = …`, a `:=` that
// redeclares x rather than shadowing it (go/types records that x in Uses
// as the same object as the earlier declaration, so collectZeroVars never
// tracks the ident on its own), or a range clause's key or value with
// `=` rather than `:=`.
func collectSafeIdents(body *ast.BlockStmt) map[*ast.Ident]bool {
	safe := map[*ast.Ident]bool{}

	ast.Inspect(body, func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.AssignStmt:
			markSafeAssignTargets(safe, node)
		case *ast.RangeStmt:
			markSafeRangeTargets(safe, node)
		case *ast.CallExpr:
			markSafeAddressArgs(safe, node)
		}

		return true
	})

	return safe
}

// markSafeAssignTargets marks node's left-hand targets safe when node is a
// whole-value assignment: plain `=`, or a `:=` that redeclares its targets
// rather than shadowing them.
func markSafeAssignTargets(safe map[*ast.Ident]bool, node *ast.AssignStmt) {
	if node.Tok != token.ASSIGN && node.Tok != token.DEFINE {
		return
	}

	for _, lhs := range node.Lhs {
		if ident, ok := lhs.(*ast.Ident); ok {
			safe[ident] = true
		}
	}
}

// markSafeRangeTargets marks node's key and value safe when node assigns
// them with `=` rather than declaring them with `:=`: that overwrites the
// whole value on every iteration, like a plain assignment.
func markSafeRangeTargets(safe map[*ast.Ident]bool, node *ast.RangeStmt) {
	if node.Tok != token.ASSIGN {
		return
	}

	if ident, ok := node.Key.(*ast.Ident); ok {
		safe[ident] = true
	}

	if ident, ok := node.Value.(*ast.Ident); ok {
		safe[ident] = true
	}
}

// markSafeAddressArgs marks every &x argument of node safe: passing a
// variable's address to any call is the other OK interaction.
func markSafeAddressArgs(safe map[*ast.Ident]bool, node *ast.CallExpr) {
	for _, arg := range node.Args {
		unary, ok := ast.Unparen(arg).(*ast.UnaryExpr)
		if !ok || unary.Op != token.AND {
			continue
		}

		if ident, ok := unary.X.(*ast.Ident); ok {
			safe[ident] = true
		}
	}
}

// firstIdentUses returns, for every object in tracked, the earliest
// *ast.Ident in body (including inside a nested function literal, which
// may capture an outer local or result) that uses records as that object.
// The declaring identifier itself is never recorded as a use.
//
// "Earliest" is evaluation order, not source position: an *ast.AssignStmt's
// right-hand side runs before its left-hand targets are written, and an
// *ast.RangeStmt's range expression runs once before any key or value is
// assigned, so each is visited first here even though it comes later in
// the text. Once an object's first use is recorded it is kept, so the
// visit order alone decides it.
func firstIdentUses(
	uses map[*ast.Ident]types.Object,
	body *ast.BlockStmt,
	tracked map[types.Object]zeroVar,
) map[types.Object]*ast.Ident {
	first := map[types.Object]*ast.Ident{}

	var visit func(node ast.Node) bool

	visit = func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.AssignStmt:
			inspectInOrder(visit, assignEvalOrder(node))

			return false
		case *ast.RangeStmt:
			inspectInOrder(visit, []ast.Node{node.X, node.Key, node.Value, node.Body})

			return false
		}

		recordIdentUse(first, uses, tracked, node)

		return true
	}

	ast.Inspect(body, visit)

	return first
}

// assignEvalOrder returns assign's operands in the order they run: every
// right-hand side expression (the values), then every left-hand target
// (the assignment itself).
func assignEvalOrder(assign *ast.AssignStmt) []ast.Node {
	nodes := make([]ast.Node, 0, len(assign.Rhs)+len(assign.Lhs))

	for _, e := range assign.Rhs {
		nodes = append(nodes, e)
	}

	for _, e := range assign.Lhs {
		nodes = append(nodes, e)
	}

	return nodes
}

// inspectInOrder runs visit over each of nodes in turn, skipping a nil
// entry (an *ast.RangeStmt's Key or Value is nil when the clause omits it).
func inspectInOrder(visit func(ast.Node) bool, nodes []ast.Node) {
	for _, n := range nodes {
		if n != nil {
			ast.Inspect(n, visit)
		}
	}
}

// recordIdentUse records node as the first use of its object when node is
// an *ast.Ident that uses records as one of tracked, and no earlier use was
// already recorded for that object.
func recordIdentUse(
	first map[types.Object]*ast.Ident,
	uses map[*ast.Ident]types.Object,
	tracked map[types.Object]zeroVar,
	node ast.Node,
) {
	ident, ok := node.(*ast.Ident)
	if !ok {
		return
	}

	obj := uses[ident]
	if obj == nil {
		return
	}

	if _, ok := tracked[obj]; !ok {
		return
	}

	if _, ok := first[obj]; !ok {
		first[obj] = ident
	}
}

// firstNakedReturn returns the earliest bare `return` statement directly in
// body, not descending into a nested function literal: such a return exits
// the literal, not the enclosing function, so it is not an interaction with
// the enclosing function's named results.
func firstNakedReturn(body *ast.BlockStmt) *ast.ReturnStmt {
	var first *ast.ReturnStmt

	ast.Inspect(body, func(n ast.Node) bool {
		if _, ok := n.(*ast.FuncLit); ok {
			return false
		}

		ret, ok := n.(*ast.ReturnStmt)
		if !ok || len(ret.Results) != 0 {
			return true
		}

		if first == nil || ret.Pos() < first.Pos() {
			first = ret
		}

		return true
	})

	return first
}

// firstInteraction picks a variable's first mention, textually, between its
// earliest identifier use (ident, nil if none) and, for a named result, the
// function's earliest naked return, and reports whether that mention is
// safe (an OK interaction). A nil node means the variable is never
// mentioned, so there is nothing to decide: its zero value never leaks.
func firstInteraction(
	v zeroVar,
	ident *ast.Ident,
	nakedReturn *ast.ReturnStmt,
	safeIdents map[*ast.Ident]bool,
) (ast.Node, bool) {
	if v.isResult && nakedReturn != nil {
		if ident == nil || nakedReturn.Pos() < ident.Pos() {
			return nakedReturn, false
		}
	}

	if ident == nil {
		return nil, false
	}

	return ident, safeIdents[ident]
}
