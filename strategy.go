package gofactory

import (
	"errors"
	"fmt"
	"go/types"
	"slices"
	"strings"

	"github.com/gobwas/glob"
)

// errEmptyGlobPattern is the configuration error for a glob-flag pattern
// (see compileGlobs) that is empty after TrimSpace: it can never match a
// package path or a qualified name, so it would silently do nothing.
var errEmptyGlobPattern = errors.New("pattern must not be empty")

// errLeadingSlashGlobPattern is the configuration error for a glob-flag
// pattern (see compileGlobs) starting with '/': gitignore gives a leading
// '/' a special "from the root" meaning, but a Go package path, and so a
// qualified name, never starts with '/', so such a pattern would compile
// and silently match nothing.
var errLeadingSlashGlobPattern = errors.New("pattern must not start with '/'")

// site is where a candidate bypass was written: its enclosing package, the
// top-level function enclosing it (nil at package scope), whether it lies
// in a _test.go file, and whether it lies in a const declaration.
type site struct {
	pkg         *types.Package
	fn          *types.Func
	inTestFile  bool
	inConstDecl bool
}

type blockedStrategy interface {
	IsBlocked(loc site, target *types.TypeName) bool
}

type nilPkg struct{}

func newNilPkg() nilPkg {
	return nilPkg{}
}

func (nilPkg) IsBlocked(_ site, _ *types.TypeName) bool {
	return false
}

type anotherPkg struct{}

func newAnotherPkg() anotherPkg {
	return anotherPkg{}
}

func (anotherPkg) IsBlocked(loc site, target *types.TypeName) bool {
	return loc.pkg.Path() != target.Pkg().Path()
}

// fence is one -packageGlobs pattern: the set of packages matching its
// glob. It is compiled with '/' as the separator, so a '*' does not cross
// a package boundary, and is tested against both the package path and the
// path plus "/": the bare-path match is what lets a pattern with no
// wildcard match that exact path, and the path-plus-"/" match is what lets
// `a/**` and `a/*` also match `a` itself.
type fence struct {
	glob glob.Glob
}

func newFences(patterns []string) ([]fence, error) {
	globs, err := compileGlobs(packageGlobsFlag, patterns)
	if err != nil {
		return nil, err
	}

	fences := make([]fence, 0, len(globs))
	for _, g := range globs {
		fences = append(fences, fence{glob: g})
	}

	return fences, nil
}

// compileGlobs compiles the patterns of flag with '/' as the only separator,
// so '*' does not cross a package boundary but does cross the '.' of a
// qualified name (import/path.Name).
func compileGlobs(flag string, patterns []string) ([]glob.Glob, error) {
	globs := make([]glob.Glob, 0, len(patterns))

	for _, pattern := range patterns {
		if pattern == "" {
			return nil, fmt.Errorf("%s %w", flag, errEmptyGlobPattern)
		}

		if strings.HasPrefix(pattern, "/") {
			return nil, fmt.Errorf("%s %w", flag, errLeadingSlashGlobPattern)
		}

		compiled, err := glob.Compile(pattern, '/')
		if err != nil {
			return nil, fmt.Errorf("unable to compile %s pattern %q: %w", flag, pattern, err)
		}

		globs = append(globs, compiled)
	}

	return globs, nil
}

func (f fence) contains(pkgPath string) bool {
	return matchesPackagePath(f.glob, pkgPath)
}

func matchesPackagePath(g glob.Glob, pkgPath string) bool {
	return g.Match(pkgPath) || g.Match(pkgPath+"/")
}

// anyFenceContains reports whether pkgPath lies in any of fences, the
// -packageGlobsOnly reading of "protected".
func anyFenceContains(fences []fence, pkgPath string) bool {
	return slices.ContainsFunc(fences, func(f fence) bool {
		return f.contains(pkgPath)
	})
}

// fencedPkgs applies the intersection rule: a type whose package lies in
// one or more fences may be bypassed only by code inside every one of
// those fences. A type in no fence gains nothing from fences and falls
// back to defaultStrategy.
type fencedPkgs struct {
	fences          []fence
	defaultStrategy blockedStrategy
}

func newFencedPkgs(
	fences []fence,
	defaultStrategy blockedStrategy,
) fencedPkgs {
	return fencedPkgs{
		fences:          fences,
		defaultStrategy: defaultStrategy,
	}
}

func (s fencedPkgs) IsBlocked(loc site, target *types.TypeName) bool {
	identPkgPath := target.Pkg().Path()

	inAnyFence := false

	for _, candidate := range s.fences {
		if !candidate.contains(identPkgPath) {
			continue
		}

		inAnyFence = true

		if !candidate.contains(loc.pkg.Path()) {
			return true
		}
	}

	if !inAnyFence {
		return s.defaultStrategy.IsBlocked(loc, target)
	}

	return false
}
