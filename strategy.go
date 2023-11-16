package factory

import (
	"go/types"

	"github.com/gobwas/glob"
)

type blockedStrategy interface {
	IsBlocked(currentPkg *types.Package, identObj types.Object) bool
}

type nilPkg struct{}

func newNilPkg() nilPkg {
	return nilPkg{}
}

func (nilPkg) IsBlocked(_ *types.Package, _ types.Object) bool {
	return false
}

type anotherPkg struct{}

func newAnotherPkg() anotherPkg {
	return anotherPkg{}
}

func (anotherPkg) IsBlocked(
	currentPkg *types.Package,
	identObj types.Object,
) bool {
	return currentPkg.Path() != identObj.Pkg().Path()
}

type blockedPkgs struct {
	pkgs            []glob.Glob
	defaultStrategy blockedStrategy
}

func newBlockedPkgs(
	pkgs []glob.Glob,
	defaultStrategy blockedStrategy,
) blockedPkgs {
	return blockedPkgs{
		pkgs:            pkgs,
		defaultStrategy: defaultStrategy,
	}
}

func (b blockedPkgs) IsBlocked(
	currentPkg *types.Package,
	identObj types.Object,
) bool {
	currentPkgPath := currentPkg.Path() + "/"
	isIncludedInBlocked := containsMatchGlob(b.pkgs, currentPkgPath)

	if isIncludedInBlocked {
		return false
	}

	identPkgPath := identObj.Pkg().Path() + "/"
	isBlocked := containsMatchGlob(b.pkgs, identPkgPath)

	if isBlocked {
		return true
	}

	if b.defaultStrategy.IsBlocked(currentPkg, identObj) {
		return true
	}

	return false
}
