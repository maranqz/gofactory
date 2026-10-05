// Package nested exercises -onlyWithFactory together with
// -factoryPatterns=^Make -useDefaultFactoryPattern=false.
package nested

// OnlyMake has an accessible factory under the replaced pattern.
type OnlyMake struct{}

func MakeOnlyMake() OnlyMake { return OnlyMake{} }

// OnlyNew's only factory matches the dropped default pattern, so it is not
// recognised and -onlyWithFactory keeps it silent.
type OnlyNew struct{}

func NewOnlyNew() OnlyNew { return OnlyNew{} }
