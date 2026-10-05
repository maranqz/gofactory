// Package nested exercises -factoryPatterns.
package nested

// Both has two recognised factories once -factoryPatterns=^Make is set:
// NewBoth matches the default ^New pattern, MakeBoth only the extra one.
type Both struct{}

func NewBoth() Both {
	return Both{}
}

func MakeBoth() Both {
	return Both{}
}
