// owner_test counts as owner's own code even once owner is fenced by
// -packageGlobs, so the bypass below stays silent.
package owner_test

import (
	"testing"

	"factory/externalTestPackageFence/owner"
)

func TestBypass(t *testing.T) {
	t.Helper()

	_ = owner.T{}
}
