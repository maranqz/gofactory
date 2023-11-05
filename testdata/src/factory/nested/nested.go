package nested

import (
	"factory/creatable_anywhere"
	"factory/nested/nested2"
)

type Nested struct{}

func New() Nested {
	return Nested{}
}

func NewPtr() *Nested {
	return &Nested{}
}

func CallNested2() {
	_ = nested2.Nested{}
	_ = &nested2.Nested{}

	_ = creatable_anywhere.Struct{}
	_ = &creatable_anywhere.Struct{}
}

func CallMp() {
	_ = nested2.Mp{}
}

func CallSlice() {
	_ = nested2.Slice{}
}

func CallArray() {
	_ = nested2.Array{}
}
