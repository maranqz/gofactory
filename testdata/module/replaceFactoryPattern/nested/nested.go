// Package nested exercises -useDefaultFactoryPattern=false together with a
// repeated -factoryPatterns.
package nested

// X has three candidate factories: NewX matches only the dropped default
// pattern, MakeX and RestoreX each match one of two repeated
// -factoryPatterns.
type X struct{}

func NewX() X {
	return X{}
}

func MakeX() X {
	return X{}
}

func RestoreX() X {
	return X{}
}
