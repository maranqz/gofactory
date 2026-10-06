package replaceFactoryPattern

import "factory/replaceFactoryPattern/nested"

// ReplacedAndRepeated exercises -useDefaultFactoryPattern=false with
// -factoryPatterns=^Make -factoryPatterns=^Restore: NewX no longer counts,
// since the default pattern was dropped, and both MakeX and RestoreX count,
// since repeating -factoryPatterns appends rather than replaces.
func ReplacedAndRepeated() {
	_ = nested.X{} // want `Use factory for nested.X \(nested.MakeX, nested.RestoreX\)$`
}
