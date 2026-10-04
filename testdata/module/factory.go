// Package factory is the current module's own root package: its import
// path, "factory", equals the module path exactly, the other half of the
// membership rule besides the "+ /" prefix nested modules rely on.
package factory

type Struct struct{}

func NewStruct() *Struct {
	return &Struct{}
}
