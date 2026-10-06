package gofactory

import (
	"go/types"
	"strings"
)

// currentModule protects the types of the module at path (see
// belongsToModule). An empty path means no module (GOPATH, Bazel nogo):
// every package other than the current one is protected then.
//
// path is the Pass.Module.Path of the package under analysis. Pass.Module.Main
// cannot be used to find it: go vet's unitchecker before Go 1.27 fills
// Module without Main, and in a go.work build every workspace module has
// Main set regardless of which one is the root.
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
	if !newAnotherPkg().IsBlocked(currentPkg, identObj) {
		return false
	}

	if c.path == "" {
		return true
	}

	return belongsToModule(identObj.Pkg().Path(), c.path)
}

func belongsToModule(pkgPath, modulePath string) bool {
	if pkgPath == modulePath {
		return true
	}

	return strings.HasPrefix(pkgPath, modulePath+"/")
}
