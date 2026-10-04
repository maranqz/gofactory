package invariant

// Struct's factory, NewStruct, is the only thing that is supposed to be able
// to set secret. Code outside this package cannot name secret, but it can
// still leave it at its zero value by reusing Struct's underlying layout
// under a different, locally-owned type (see ../underlying.go).
type Struct struct {
	Field  int
	secret int
}

func NewStruct() Struct {
	return Struct{Field: 1, secret: 42}
}
