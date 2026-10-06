// Package domain holds the protected type for the zeroValuesFences case,
// inside the fenced package tree itself: the -packageGlobs entry in
// lint_test.go names "factory/zeroValuesFences/blocked/**", which covers
// this package and its sibling other.
package domain

// Struct is the protected type under test.
type Struct struct {
	Field int
}
