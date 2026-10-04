package sibling

type Struct struct{}

func NewStruct() *Struct {
	return &Struct{}
}

//gofactory:ignore
type Ignored struct{} // want Ignored:"gofactory:ignore"
