package local

// not implemented

type OtherStruct struct{}

type Struct struct {
	Other OtherStruct
}

func NewOtherStruct() OtherStruct {
	return OtherStruct{}
}

func NewStruct() Struct {
	return Struct{
		Other: NewOtherStruct(),
	}

	return Struct{
		Other: OtherStruct{}, // want `Use factory for OtherStruct`
	}
}

func NewStructWithoutFields() Struct {
	return Struct{NewOtherStruct()}

	return Struct{
		OtherStruct{}, // want `Use factory for OtherStruct`
	}
}
