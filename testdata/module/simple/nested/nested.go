package nested

type Struct struct{}

func NewStruct() Struct {
	return Struct{}
}

type Mp map[bool]bool

type Slice []bool

type Array [5]bool

type Structs []Struct

type Err struct{}

func (*Err) Error() string { return "" }
