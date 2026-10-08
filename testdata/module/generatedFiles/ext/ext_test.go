package ext

import "testing"

// TestEmpty's only purpose is to make `go test` synthesize a test main for
// this package: see TestLinterSuite's "generatedFiles" case, fencing
// testing.InternalTest and testdeps.TestDeps so that without the
// generated-file skip, the synthesized main would be reported.
func TestEmpty(t *testing.T) {
	t.Helper()
}
