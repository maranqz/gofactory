package casting

import "factory/casting/nested"

// Parenthesising the target type does not hide a conversion.
func ToNestedStructParen() {
	l := Struct{Field: 1}

	_ = (nested.Struct)(l) // want `Use factory for nested.Struct`
}

func ToNestedStructPtrParen() {
	l := Struct{Field: 1}

	_ = (*nested.Struct)(&l) // want `Use factory for nested.Struct`
}
