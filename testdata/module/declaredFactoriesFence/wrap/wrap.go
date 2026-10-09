// Package wrap checks that a declared factory may bypass owner.T even
// though a -packageGlobs fence covers owner but not this package, while an
// undeclared bypass in the same package is still reported.
package wrap

import "factory/declaredFactoriesFence/owner"

//gofactory:factory
func Build() owner.T { // want Build:"gofactory:factory"
	return owner.T{}
}

func Other() owner.T {
	return owner.T{} // want `Use factory for owner.T \(wrap.Build\)$`
}
