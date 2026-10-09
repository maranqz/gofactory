package declaredFactoriesOnlyWithFactory

import (
	"factory/declaredFactoriesOnlyWithFactory/nested"
	"factory/declaredFactoriesOnlyWithFactory/other"
)

// Reported exercises -onlyWithFactory: WithDeclaredFactory's only factory
// is declared in another package rather than recognised by name pattern,
// and Restored's is declared in its own package, the README's
// reconstitution recipe; both still count, even with
// -useDefaultFactoryPattern=false.
func Reported() {
	_ = other.Build()
	_ = nested.Restore()

	_ = nested.WithDeclaredFactory{} // want `Use factory for nested.WithDeclaredFactory \(other.Build\)$`
	_ = nested.Restored{}            // want `Use factory for nested.Restored \(nested.Restore\)$`
}

// Silent exercises -onlyWithFactory: WithoutFactory has no factory at
// all.
func Silent() {
	_ = nested.WithoutFactory{}
}
