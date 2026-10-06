package useDefaultFactoryPattern

import "factory/useDefaultFactoryPattern/nested"

// Disabled exercises -useDefaultFactoryPattern=false with no extra
// patterns: NewOnlyNew no longer matches anything, so no factory is
// recognised at all and the message keeps the bare prefix.
func Disabled() {
	_ = nested.OnlyNew{} // want `Use factory for nested.OnlyNew$`
}
