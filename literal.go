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
	d.checkZeroValueFields(lit, litType)
}

func (d *detector) checkZeroValueFields(
	lit *ast.CompositeLit, litType types.Type,
) {
	if !d.zeroValues || isFullyPositional(lit) {
		return
	}

	d.reportUnsetFields(lit, litType, keyedFieldNames(lit))
}

// isFullyPositional reports whether every field of lit's struct type is
// set: Go requires a positional struct literal to supply every field, so
// a non-empty, unkeyed literal leaves nothing unset.
func isFullyPositional(lit *ast.CompositeLit) bool {
	if len(lit.Elts) == 0 {
		return false
	}

	_, keyed := lit.Elts[0].(*ast.KeyValueExpr)

	return !keyed
}

func keyedFieldNames(lit *ast.CompositeLit) map[string]bool {
	names := map[string]bool{}

	for _, elt := range lit.Elts {
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			if ident, ok := kv.Key.(*ast.Ident); ok {
				names[ident.Name] = true
			}
		}
	}

	return names
}
