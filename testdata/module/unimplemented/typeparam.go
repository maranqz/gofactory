package unimplemented

import "factory/unimplemented/nested"

// A type parameter is not a named type, so a literal, conversion, new or
// implicit constant conversion of one is not checked, even when its
// constraint admits a single protected type and every instantiation builds
// that type without its factory.
func TypeParamStruct[T nested.Struct]() []T {
	_ = new(T) // want `Use factory for nested.Struct`

	return []T{{}} // want `Use factory for nested.Struct`
}

func TypeParamLiteral[T nested.Struct]() T {
	return T{} // want `Use factory for nested.Struct`
}

func TypeParamPointer[P *nested.Struct]() []P {
	return []P{{}} // want `Use factory for nested.Struct`
}

func TypeParamConversion[M nested.Mp](m map[bool]bool) M {
	return M(m) // want `Use factory for nested.Mp`
}

func TypeParamConstant[T nested.MyInt]() T {
	return 3 // want `Use factory for nested.MyInt`
}
