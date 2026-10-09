package sibling

type Struct struct{}

func NewStruct() *Struct {
	return &Struct{}
}

//gofactory:ignore
type Ignored struct{} // want Ignored:"gofactory:ignore"

type Declared struct{}

//gofactory:factory
func Restore() Declared { // want Restore:"gofactory:factory"
	return Declared{}
}
