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
