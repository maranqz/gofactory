// Package nested exercises the New… first ordering of suggested factories
// once -useDefaultFactoryPattern=false drops the default pattern.
package nested

// Both has two candidate factories, NewBoth and MakeBoth, and both match
// -factoryPatterns=Both$.
type Both struct{}

func MakeBoth() Both { return Both{} }

func NewBoth() Both { return Both{} }
