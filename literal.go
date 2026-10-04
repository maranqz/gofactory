package gofactory

import "go/ast"

// checkLiteral reports a composite literal of a protected type. The literal's
// type comes straight from the type checker, which already resolves elided
// element literals at every nesting and keyed form, so a named container
// literal and its elements are each visited and checked independently.
func (d *detector) checkLiteral(lit *ast.CompositeLit) {
	t := d.pass.TypesInfo.TypeOf(lit)
	if t == nil {
		return
	}

	d.reportProtected(lit, t)
}
