// Package ext holds a protected type bypassed from the other packages
// under generatedFiles, to show that a bypass in a generated file is
// skipped and the same bypass in a plain file of the same package is not.
package ext

type Struct struct{}

func NewStruct() Struct {
	return Struct{}
}
