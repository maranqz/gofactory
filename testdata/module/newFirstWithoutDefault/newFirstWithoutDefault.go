package newFirstWithoutDefault

import "factory/newFirstWithoutDefault/nested"

// NewFirst exercises -factoryPatterns=Both$ -useDefaultFactoryPattern=false:
// NewBoth is recognised through Both$, not the dropped default pattern, and
// is still suggested before MakeBoth.
func NewFirst() {
	_ = nested.Both{} // want `Use factory for nested.Both \(nested.NewBoth, nested.MakeBoth\)$`
}
