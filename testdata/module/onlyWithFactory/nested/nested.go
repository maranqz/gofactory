// Package nested exercises -onlyWithFactory.
package nested

// WithFactory has an accessible factory.
type WithFactory struct{}

func NewWithFactory() WithFactory {
	return WithFactory{}
}

// WithoutFactory has none.
type WithoutFactory struct{}

// Hidden has a recognised factory, but it is a method of builder, an
// unexported type: outside this package nobody can write nested.builder,
// so the factory is never accessible from another package's diagnostic.
type Hidden struct{}

type builder struct{}

func (builder) NewHidden() Hidden { return Hidden{} }
