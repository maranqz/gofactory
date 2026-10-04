package casting

import "factory/casting/nested"

// Parenthesising the target type does not hide a conversion.
func ToNestedStructParenthesis() {
	l := Struct{Field: 1}

	_ = (nested.Struct)(l) // want `Use factory for nested.Struct \(nested.NewStruct\)`
}

func ToNestedStructPtrParenthesis() {
	l := Struct{Field: 1}

	_ = (*nested.Struct)(&l) // want `Use factory for nested.Struct \(nested.NewStruct\)`
}
