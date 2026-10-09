// owner_test is the external test package of a package named owner, but
// not of externalTestPackage/owner, so the bypass below is reported.
package owner_test

import (
	"testing"

	"factory/externalTestPackage/owner"
)

func TestBypass(t *testing.T) {
	t.Helper()

	_ = owner.T{} // want `Use factory for owner.T \(owner.NewT\)`
}
