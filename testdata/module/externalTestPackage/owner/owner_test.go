// owner_test counts as owner's own code, so the bypass below is silent.
package owner_test

import (
	"testing"

	"factory/externalTestPackage/owner"
)

func TestBypass(t *testing.T) {
	t.Helper()

	_ = owner.T{}
}
