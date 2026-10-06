package factoryext

type Struct struct{}

func NewStruct() *Struct {
	return &Struct{}
}
