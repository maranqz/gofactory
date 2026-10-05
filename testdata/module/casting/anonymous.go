package casting

import "factory/casting/nested"

// Converting an anonymous struct literal to a protected type creates
// something new and is reported, unlike the no-op in SkipNoOpConversion.
func ToNestedStructFromAnonymousLiteral() {
	var l nested.Struct

	l = nested.Struct(struct{ Field int }{}) // want `Use factory for nested.Struct \(nested.NewStruct\)`

	_ = l
}
