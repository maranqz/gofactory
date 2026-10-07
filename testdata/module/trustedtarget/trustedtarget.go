// Package trustedtarget is the protected type that the trusted-code tests,
// under testdata/module/directive and testdata/module/trusted, bypass
// freely from code that is trusted and still get reported from code that
// is not.
package trustedtarget

type Struct struct{}

func NewStruct() Struct {
	return Struct{}
}
