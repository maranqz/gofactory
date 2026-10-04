package gofactory

import (
	"go/types"
	"strings"
)

// currentModule protects only the current module's own packages, so stdlib
// and dependency types stay silent by default while a developer's domain
// types stay guarded without any configuration (strict mode still reports
// every bypass of a type this strategy does protect).
//
// path is the Pass.Module.Path of the package under analysis. Pass.Module.Main
// cannot be used to find it: go vet's unitchecker before Go 1.27 fills
// Module without Main, and in a go.work build every workspace module has
// Main set regardless of which one is the root (confirmed for
// golang.org/x/tools v0.50.0 by TestTestdataRoots in lint_test.go). The
// caller is expected to fall back to anotherPkg when path is empty, the
// signal that pass.Module carries no module (GOPATH, Bazel nogo).
type currentModule struct {
	path string
}

func newCurrentModule(path string) currentModule {
	return currentModule{path: path}
}

func (c currentModule) IsBlocked(
	currentPkg *types.Package,
	identObj types.Object,
) bool {
	identPkg := identObj.Pkg()
	if currentPkg.Path() == identPkg.Path() {
		return false
	}

	return belongsToModule(identPkg.Path(), c.path)
}

// belongsToModule reports whether pkgPath is the module at modulePath, or
// nested under it, so that a nested module in the same repository counts
// as the current module too. A go.work sibling module, whose path does not
// start with modulePath, does not belong, and is silent until a fence
// names it.
func belongsToModule(pkgPath, modulePath string) bool {
	if pkgPath == modulePath {
		return true
	}

	return strings.HasPrefix(pkgPath, modulePath+"/")
}
