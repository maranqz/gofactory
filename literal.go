package gofactory

import "go/ast"

// checkLiteral reports a composite literal of a protected type. The literal's
// type comes straight from the type checker, which already resolves elided
// element literals at every nesting and keyed form, so a named container
// literal and its elements are each visited and checked independently.
func (v *detector) checkLiteral(lit *ast.CompositeLit) {
	t := v.pass.TypesInfo.TypeOf(lit)
	if t == nil {
		return
	}

	named, ok := protectedNamed(t)
	if !ok {
		return
	}

	if !v.isProtected(named) {
		return
	}

	v.report(lit, named)
}
