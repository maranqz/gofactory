// Package trustedtarget holds the protected type the trusted-code tests
// bypass.
package trustedtarget

type Struct struct{}

func NewStruct() Struct {
	return Struct{}
}
