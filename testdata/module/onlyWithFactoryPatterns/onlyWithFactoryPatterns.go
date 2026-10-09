package onlyWithFactoryPatterns

import "factory/onlyWithFactoryPatterns/nested"

// Reported exercises -onlyWithFactory gated on the replaced pattern set:
// OnlyMake's factory matches -factoryPatterns=^Make, so it is reported.
func Reported() {
	_ = nested.OnlyMake{} // want `Use factory for nested.OnlyMake \(nested.MakeOnlyMake\)$`
}

// Silent exercises -onlyWithFactory gated on the replaced pattern set:
// OnlyNew's candidate factory matches only the default pattern, dropped by
// -useDefaultFactoryPattern=false, so -onlyWithFactory keeps it silent.
func Silent() {
	_ = nested.OnlyNew{}
}
