package indirect

// not implemented: a declared factory's fact is exported by the package that
// applies the directive or the -factories glob, and x/tools attaches a
// function fact only to packages that import that package directly
// (checker.exportedFrom, facts.Encode). indirect imports owner but never
// wrapper, so wrapper.New's fact for owner.T never reaches it, even though
// wrapper.New is just as reachable here as any recognised factory would be.

import "factory/unimplemented/visibility/owner"

func Indirect() owner.T {
	return owner.T{} // want `Use factory for owner.T \(wrapper.New\)`
}
