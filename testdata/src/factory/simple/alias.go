package simple

import alias "factory/simple/aliasnested"

func Alias() {
	_ = alias.Struct{}     // want `Use factory for aliasnested.Struct`
	_ = &alias.Struct{}    // want `Use factory for aliasnested.Struct`
	_ = []alias.Struct{{}} // want `Use factory for aliasnested.Struct`
}
