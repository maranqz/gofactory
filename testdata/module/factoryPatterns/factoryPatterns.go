package factoryPatterns

import "factory/factoryPatterns/nested"

// ExtraPattern exercises -factoryPatterns=^Make: MakeBoth is recognised
// alongside the default NewBoth, which keeps sorting first.
func ExtraPattern() {
	_ = nested.Both{} // want `Use factory for nested.Both \(nested.NewBoth, nested.MakeBoth\)$`
}

// CapKeepsNewFirst checks that NewMany survives the three-factory cap.
func CapKeepsNewFirst() {
	_ = nested.Many{} // want `Use factory for nested.Many \(nested.NewMany, nested.MakeManyA, nested.MakeManyB\)$`
}
