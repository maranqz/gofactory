package casting

import (
	"bytes"
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

// No diagnostic, known false negative: an untyped non-constant argument is
// recorded with the target's type, so it looks like a no-op conversion.
func SkipUntypedNonConstant(a, b int, n uint) {
	_ = nested.Flag(a == b)
	_ = nested.MyInt(1 << n)
}

// A conversion to an interface creates nothing, nor does new of one.
func SkipInterfaceTarget(buf *bytes.Buffer) {
	_ = io.Reader(buf)
	_ = new(io.Reader)
}

// A conversion to unsafe.Pointer creates nothing.
func SkipUnsafePointerTarget(s *nested.Struct) {
	_ = unsafe.Pointer(s)
}

// A conversion to a func type creates nothing: func types are not protected.
func SkipFuncTypeTarget(f func(http.ResponseWriter, *http.Request)) {
	_ = http.HandlerFunc(f)
}

type myErr struct{}

func (*myErr) Error() string { return "" }

// A universe type (nil package) does not panic. Every universe named type
// (error, comparable) is an interface, so the interface rule silences it
// before protectedNamed's nil-package guard matters; that guard is defensive.
func SkipUniverseType() error {
	return error(&myErr{})
}
