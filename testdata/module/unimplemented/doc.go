// Package unimplemented holds false negatives the linter does not catch yet.
// It is not run by lint_test.go, so its // want comments are not asserted;
// they record the diagnostic each case should get once it is handled.
//
// Every case added here must carry a // want comment on the line that should
// be reported, so that moving it into a tested package is all it takes to
// turn it into a test.
package unimplemented
