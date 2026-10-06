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

// Many has four recognised factories, one more than a message lists, and
// its three Make… factories sort before NewMany by name.
type Many struct{}

func MakeManyA() Many {
	return Many{}
}

func MakeManyB() Many {
	return Many{}
}

func MakeManyC() Many {
	return Many{}
}

func NewMany() Many {
	return Many{}
}
