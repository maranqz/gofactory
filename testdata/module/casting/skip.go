package casting

import (
	"io"
	"net/http"
	"unsafe"

	"factory/casting/nested"
)

// A nil argument creates nothing.
func SkipNilArgument() {
	_ = (*nested.Struct)(nil)
}

// A conversion to the argument's own type creates nothing: the no-op
// nested.Struct(nested.NewStruct()).
func SkipNoOpConversion() {
	s := nested.NewStruct()

	_ = nested.Struct(s)
}

// A conversion to an interface creates nothing.
func SkipInterfaceTarget(r io.Reader) {
	_ = io.Reader(r)
}

// A conversion to unsafe.Pointer creates nothing.
func SkipUnsafePointerTarget(x unsafe.Pointer) {
	_ = unsafe.Pointer(x)
}

// A conversion to a func type creates nothing: func types are not protected.
func SkipFuncTypeTarget(f func(http.ResponseWriter, *http.Request)) {
	_ = http.HandlerFunc(f)
}

// A universe type (nil package) is guarded rather than panicking.
func SkipUniverseType(err error) error {
	return error(err)
}
