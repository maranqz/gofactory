package nested

type MyInt int

type Flag bool

type Struct struct {
	Field int
}

func NewStruct() Struct {
	return Struct{}
}
