// Package deep sits two path segments below factory/globsyntax, so the
// single-star fence factory/globsyntax/* does not cover it: '*' does not
// cross '/'. deep's type lies in no fence and falls back to plain module
// scope, which protects it regardless of who is bypassing it.
package deep

type Struct struct{}

func New() Struct {
	return Struct{}
}
