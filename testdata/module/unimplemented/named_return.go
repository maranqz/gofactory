package unimplemented

import "factory/unimplemented/nested"

func Call() {
	st := HackFactory()

	_ = st.Field
}

func HackFactory() (hack nested.Struct) {
	return // want `Use factory for nested.Struct`
}
