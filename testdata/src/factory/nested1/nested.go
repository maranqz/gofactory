package nested1

import (
	"factory/creatable_anywhere"
	"factory/nested1/nested1_2"
)

type Nested struct{}

func (n Nested) Ret() Nested {
	return n
}

func New() Nested {
	return Nested{}
}

func NewPtr() *Nested {
	return &Nested{}
}

func CallNested2() {
	_ = nested1_2.Nested{}
	_ = &nested1_2.Nested{}

	_ = creatable_anywhere.Struct{}
	_ = &creatable_anywhere.Struct{}
}

func CallMp() {
	_ = nested1_2.Mp{}
}

func CallSlice() {
	_ = nested1_2.Slice{}
}

func CallArray() {
	_ = nested1_2.Array{}
}
