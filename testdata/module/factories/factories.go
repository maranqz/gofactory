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

// FactoryWithError exercises a factory with an extra error result, which
// doesn't disqualify it: the rule only asks for target among the results.
func FactoryWithError() {
	_ = nested.WithErr{} // want `Use factory for nested.WithErr \(nested.NewWithErr\)`
}

// SortedFactories exercises the alphabetical tie-break among methods
// declared out of alphabetical order: NewSortedA must come first even though
// NewSortedB is declared first.
func SortedFactories() {
	_ = nested.Sorted{} // want `Use factory for nested.Sorted \(nested.Maker.NewSortedA, nested.Maker.NewSortedB\)$`
}

// InaccessibleFactory exercises a factory method of an unexported type:
// outside nested, nothing can write nested.hidden, so the factory is never
// suggested here even though it is a recognised factory of Secret.
func InaccessibleFactory() {
	_ = nested.Secret{} // want `Use factory for nested.Secret$`
}

// CappedFactories exercises the three-factory cap: NewD is dropped.
func CappedFactories() {
	_ = nested.Multi{} // want `Use factory for nested.Multi \(nested.NewA, nested.NewB, nested.NewC\)`
}

// NamedPointerFactory exercises a protected type that is itself a defined
// pointer type: BoxPtr's own factory, NewBoxPtr, must be suggested for
// BoxPtr, not just for the Box it points to.
func NamedPointerFactory() {
	_ = new(nested.BoxPtr) // want `Use factory for nested.BoxPtr \(nested.NewBoxPtr\)$`
	_ = nested.Box{}       // want `Use factory for nested.Box \(nested.NewBoxPtr\)$`
}
