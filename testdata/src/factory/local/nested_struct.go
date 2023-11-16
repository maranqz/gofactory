package local

type OtherStruct struct{}

type Struct struct {
	Other OtherStruct
}

func NewOtherStruct() OtherStruct {
	return OtherStruct{}
}

func NewStruct() {
	_ = Struct{Other: NewOtherStruct()}

	_ = Struct{
		Other: OtherStruct{}, // want `Use factory for nested.Struct`
	}
}

func NewStructWithoutFields() {
	_ = Struct{NewOtherStruct()}

	_ = Struct{OtherStruct{}} // want `Use factory for nested.Struct`
}
