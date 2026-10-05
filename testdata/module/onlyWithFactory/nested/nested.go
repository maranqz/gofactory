// Package nested exercises -onlyWithFactory.
package nested

// WithFactory has an accessible factory.
type WithFactory struct{}

func NewWithFactory() WithFactory {
	return WithFactory{}
}

// WithoutFactory has none.
type WithoutFactory struct{}
