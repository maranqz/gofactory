package nested1_2

type Nested struct{}

func New() Nested {
	return Nested{}
}

func NewPtr() *Nested {
	return &Nested{}
}

type Mp map[bool]bool

type Slice []bool

type Array [5]bool
