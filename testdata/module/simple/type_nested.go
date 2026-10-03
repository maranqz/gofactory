package simple

import "factory/simple/nested"

type DeclStruct nested.Struct
type AliasStruct = nested.Struct
type UnderlyingStruct struct {
	nested.Struct
}

func typeNested() {
	_ = DeclStruct{}  // DeclStruct belongs to the current package: not a bypass.
	_ = AliasStruct{} // want `Use factory for nested.Struct`

	// No diagnostic, known false negative: the embedded nested.Struct field
	// is left at its zero value, which is a field-path bypass (its own ticket).
	_ = UnderlyingStruct{}
}
