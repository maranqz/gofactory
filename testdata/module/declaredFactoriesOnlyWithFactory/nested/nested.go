// Package nested exercises -onlyWithFactory together with a declared
// factory, including with -useDefaultFactoryPattern=false.
package nested

// WithDeclaredFactory has no recognised factory in its own package; its
// only factory, other.Build, is declared elsewhere.
type WithDeclaredFactory struct{}

// WithoutFactory has none at all.
type WithoutFactory struct{}
