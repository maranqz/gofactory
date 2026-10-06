package nomodule

import (
	"strings"

	"nomodule/ext"
)

func NoModule() {
	_ = ext.Struct{}  // want `Use factory for ext.Struct`
	_ = &ext.Struct{} // want `Use factory for ext.Struct`
	_ = ext.NewStruct()
}

// Stdlib is protected too without a module: the fallback protects every
// package other than the current one, standard library included.
func Stdlib() {
	_ = strings.Builder{} // want `Use factory for strings.Builder`
}
