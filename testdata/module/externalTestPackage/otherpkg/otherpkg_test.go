// otherpkg_test is an external test package too, but not owner's, so the
// bypass below is still reported.
package otherpkg_test

import (
	"testing"

	"factory/externalTestPackage/owner"
)

func TestBypass(t *testing.T) {
	t.Helper()

	_ = owner.T{} // want `Use factory for owner.T \(owner.NewT\)`
}
