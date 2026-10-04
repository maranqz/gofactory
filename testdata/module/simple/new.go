package simple

import (
	"errors"

	"factory/simple/nested"
)

func NewBuiltin() {
	_ = new(nested.Struct)   // want `Use factory for nested.Struct`
	_ = *new(nested.Struct)  // want `Use factory for nested.Struct`
	_ = (new)(nested.Struct) // want `Use factory for nested.Struct`

	// new(expr), Go 1.26: allocates from an already-produced value, so there
	// is nothing left to bypass.
	_ = new(nested.NewStruct())

	// new(*T) allocates a nil *T and builds no T.
	_ = new(*nested.Struct)

	// make is a bypass too, but it is deferred until fill analysis exists.
	_ = make(nested.Mp)
}

func NewPointerErrorsAs(err error) bool {
	return errors.As(err, new(*nested.Err))
}
