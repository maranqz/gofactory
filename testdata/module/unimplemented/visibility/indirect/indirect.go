package indirect

// not implemented: a declared factory's fact is exported by the package that
// applies the directive, and x/tools attaches a function fact only to
// packages that import that package directly (checker.exportedFrom,
// facts.Encode); a -factories glob is matched against direct imports only.
// indirect imports wrapper only through mid, so wrapper.New's fact for
// owner.T never reaches it.

import (
	"factory/unimplemented/visibility/mid"
	"factory/unimplemented/visibility/owner"
)

func Indirect() owner.T {
	mid.Mid()

	return owner.T{} // want `Use factory for owner.T \(wrapper.New\)`
}
