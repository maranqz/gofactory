package nested

type Struct struct{}

func NewStruct() Struct {
	return Struct{}
}

type Mp map[bool]bool

type Slice []bool

type Array [5]bool
