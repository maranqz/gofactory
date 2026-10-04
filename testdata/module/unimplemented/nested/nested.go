package nested

type Flag bool

type MyInt int

type Struct struct {
	Field int
}

func NewStruct() Struct {
	return Struct{}
}

type Mp map[bool]bool

type PU *struct{ Field int }
