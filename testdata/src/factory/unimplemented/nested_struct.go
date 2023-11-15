package unimplemented

type OtherStruct struct{}

type Struct struct {
	Other OtherStruct
}

func NewStruct() Struct {
	return Struct{
		Other: OtherStruct{}, // want `Use factory for nested.Struct`
	}
}
