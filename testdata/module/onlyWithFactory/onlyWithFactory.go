package onlyWithFactory

import "factory/onlyWithFactory/nested"

// Reported exercises -onlyWithFactory: WithFactory has an accessible
// factory, so it is still reported.
func Reported() {
	_ = nested.WithFactory{} // want `Use factory for nested.WithFactory \(nested.NewWithFactory\)$`
}

// Silent exercises -onlyWithFactory: WithoutFactory has no accessible
// factory, so -onlyWithFactory keeps it silent.
func Silent() {
	_ = nested.WithoutFactory{}
}

// SilentInaccessible exercises -onlyWithFactory: Hidden has a recognised
// factory, but it belongs to an unexported type and so is not accessible
// from here, keeping -onlyWithFactory silent just as with no factory at
// all.
func SilentInaccessible() {
	_ = nested.Hidden{}
}
