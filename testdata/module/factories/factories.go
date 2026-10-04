package factories

import "factory/factories/nested"

// NoFactory exercises every rule that disqualifies a New-named candidate at
// once: a wither (method of the type itself), a function taking the type,
// a function taking a pointer to it, an unexported function and a function
// whose name does not match the default pattern. None of them count, so
// the type keeps the bare prefix.
func NoFactory() {
	_ = nested.Struct{} // want `Use factory for nested.Struct$`
}

// MethodFactory exercises a factory method of another type, rendered as
// pkg.Type.Method in the suffix.
func MethodFactory() {
	_ = nested.Widget{} // want `Use factory for nested.Widget \(nested.Builder.NewWidget\)`
}

// FactoryWithError exercises a factory whose error result is ignored.
func FactoryWithError() {
	_ = nested.WithErr{} // want `Use factory for nested.WithErr \(nested.NewWithErr\)`
}

// CappedFactories exercises the three-factory cap and the deterministic,
// alphabetical tie-break among four equally-ranked candidates: NewD is
// dropped.
func CappedFactories() {
	_ = nested.Multi{} // want `Use factory for nested.Multi \(nested.NewA, nested.NewB, nested.NewC\)`
}
