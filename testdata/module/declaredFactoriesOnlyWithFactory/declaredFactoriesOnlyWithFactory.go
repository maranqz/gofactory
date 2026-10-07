package declaredFactoriesOnlyWithFactory

import (
	"factory/declaredFactoriesOnlyWithFactory/nested"
	"factory/declaredFactoriesOnlyWithFactory/other"
)

// Reported exercises -onlyWithFactory: WithDeclaredFactory's only factory
// is declared in another package rather than recognised by name pattern,
// and it still counts, even with -useDefaultFactoryPattern=false.
func Reported() {
	_ = other.Build()

	_ = nested.WithDeclaredFactory{} // want `Use factory for nested.WithDeclaredFactory \(other.Build\)$`
}

// Silent exercises -onlyWithFactory: WithoutFactory has no factory at
// all.
func Silent() {
	_ = nested.WithoutFactory{}
}
