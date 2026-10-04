// Package zeroValues: this file exercises the three closure special cases
// the detector handles, which no other testdata file reaches: a closure's
// own locals, a closure's own naked return, and an outer local or result
// captured and used inside a closure.
package zeroValues

import "factory/zeroValues/nested"

// ClosureLocalVarIsReported shows that a function literal is its own
// function entity: checkFuncZeroValues runs on it separately from the
// enclosing function, and tracks its own `var` independently.
func ClosureLocalVarIsReported() {
	fn := func() int {
		var x nested.Struct

		return x.Field // want `Use factory for nested.Struct: zero value`
	}

	_ = fn()
}

// ClosureNamedReturnIsReported shows that a naked return inside a function
// literal is an interaction with that literal's own named result.
func ClosureNamedReturnIsReported() {
	fn := func() (hack nested.Struct) {
		return // want `Use factory for nested.Struct: zero value`
	}

	_ = fn()
}

// OuterVarUsedInsideClosureIsReported shows that firstIdentUses, unlike
// collectZeroVars and firstNakedReturn, does descend into a nested function
// literal: a capture of an outer local counts as that local's first
// interaction.
func OuterVarUsedInsideClosureIsReported() {
	var x nested.Struct

	fn := func() int {
		return x.Field // want `Use factory for nested.Struct: zero value`
	}

	_ = fn()
}

// OuterResultSilentDespiteInnerNakedReturnIsSilent shows that
// firstNakedReturn does not descend into a nested function literal: the
// closure's naked return is an interaction with its own result, inner, not
// with the enclosing function's hack, which is assigned before either
// return runs.
func OuterResultSilentDespiteInnerNakedReturnIsSilent() (hack nested.Struct) {
	fn := func() (inner nested.Struct) {
		return // want `Use factory for nested.Struct: zero value`
	}

	hack = fn()

	return hack
}
