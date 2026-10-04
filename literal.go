package gofactory

import (
	"go/ast"
	"go/types"
)

// checkLiteral reports a composite literal of a protected type. The literal's
// type comes straight from the type checker, which already resolves elided
// element literals at every nesting and keyed form, so a named container
// literal and its elements are each visited and checked independently.
func (d *detector) checkLiteral(lit *ast.CompositeLit) {
	litType := d.pass.TypesInfo.TypeOf(lit)
	if litType == nil {
		return
	}

	// An elided &T{} records the element type: *T, an alias of it, or a
	// defined pointer type such as type P *T.
	if ptr, ok := litType.Underlying().(*types.Pointer); ok {
		litType = ptr.Elem()
	}

	d.reportProtected(lit, litType)
}
