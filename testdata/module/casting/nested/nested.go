package nested

type MyInt int

type Struct struct {
	Field int
}

func NewStruct() Struct {
	return Struct{}
}
