package simple

import "factory/simple/nested"

type DeclStruct nested.Struct
type AliasStruct = nested.Struct
type UnderlyingStruct struct {
	nested.Struct
}

func typeNested() {
	_ = DeclStruct{}  // DeclStruct belongs to the current package: not a bypass.
	_ = AliasStruct{} // want `Use factory for nested.Struct \(nested.NewStruct\)`

	// No diagnostic here: this package does not set -zeroValues, so the
	// embedded nested.Struct field being left at its zero value is a
	// zero-value bypass only under that setting; see
	// zeroValues/fieldpath.go's NamedEmbeddingLiteralIsReported.
	_ = UnderlyingStruct{}
}
