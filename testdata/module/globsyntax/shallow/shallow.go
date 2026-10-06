// Package shallow sits one path segment below factory/globsyntax, so the
// single-star fence factory/globsyntax/* covers it.
package shallow

type Struct struct{}

func New() Struct {
	return Struct{}
}
