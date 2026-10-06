// Package nested exercises -useDefaultFactoryPattern=false.
package nested

// OnlyNew's only candidate factory matches just the default ^New pattern.
type OnlyNew struct{}

func NewOnlyNew() OnlyNew {
	return OnlyNew{}
}
