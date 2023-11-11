package blocked_nested

type Struct struct{}

func New() Struct {
	return Struct{}
}

func NewPtr() *Struct {
	return &Struct{}
}
