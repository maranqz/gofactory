// Package ownPackage exercises -ownPackage: inside this package, Loan and
// Status's factories may be bypassed only in a producer — a top-level
// function or method whose results include the type, *the type, a named
// interface it implements, or a container holding it at any depth.
package ownPackage

type Loan struct {
	Amount int
}

// Status is this package's exported protected int-based type, for the
// conversion and implicit-constant routes a struct can't exercise.
type Status int
