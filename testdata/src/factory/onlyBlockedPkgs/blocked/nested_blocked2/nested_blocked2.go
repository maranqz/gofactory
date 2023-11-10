package nested_blocked2

type Struct struct{}

func New() Struct {
	return Struct{}
}

func NewPtr() *Struct {
	return &Struct{}
}
