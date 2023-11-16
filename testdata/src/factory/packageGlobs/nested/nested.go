package nested

import (
	"factory/onlyPackageGlobs/blocked"
)

type Struct struct{}

func callNested1() {
	n := blocked.Struct{} // want `Use factory for blocked.Struct`
	_ = n
	_ = blocked.Struct{}       // want `Use factory for blocked.Struct`
	n = blocked.Struct{}.Ret() // want `Use factory for blocked.Struct`
}
