package gofactory

import (
	"errors"
	"fmt"
	"go/types"
	"slices"
	"strings"

	"github.com/gobwas/glob"
	"golang.org/x/tools/go/analysis"
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

// site is where a candidate bypass was written. fn is the enclosing
// top-level function, so a closure counts as its declaration; nil at
// package scope.
type site struct {
	pkg         *types.Package
	testOnly    bool
	fn          *types.Func
	inTestFile  bool
	inConstDecl bool
}

type blockedStrategy interface {
	IsBlocked(loc site, target *types.Named) bool
}

type nilPkg struct{}

func newNilPkg() nilPkg {
	return nilPkg{}
}

func (nilPkg) IsBlocked(_ site, _ *types.Named) bool {
	return false
}

type anotherPkg struct{}

func newAnotherPkg() anotherPkg {
	return anotherPkg{}
}

func (anotherPkg) IsBlocked(loc site, target *types.Named) bool {
	return !isOwnerPackage(loc, target.Obj().Pkg())
}

// isOwnerPackage reports whether loc.pkg is pkg itself or pkg's external
// test package. The go command gives that external test package the import
// path pkg.Path()+"_test" (TestPackagesAndErrors in
// cmd/go/internal/load/test.go), and go/build requires its package clause to
// be pkg.Name()+"_test", whatever pkg's directory is called. A regular
// package can have that same path and name, and production code could
// import it; only its files tell it apart, as a regular package never has
// a _test.go file, so loc.testOnly must say whether loc.pkg was built
// only from _test.go files.
func isOwnerPackage(loc site, pkg *types.Package) bool {
	return loc.pkg.Path() == pkg.Path() ||
		(loc.testOnly &&
			loc.pkg.Path() == pkg.Path()+"_test" &&
			loc.pkg.Name() == pkg.Name()+"_test")
}

func isTestOnly(pass *analysis.Pass) bool {
	if len(pass.Files) == 0 {
		return false
	}

	for _, file := range pass.Files {
		if !strings.HasSuffix(pass.Fset.File(file.Pos()).Name(), "_test.go") {
			return false
		}
	}

	return true
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

func (s fencedPkgs) IsBlocked(loc site, target *types.Named) bool {
	// foo_test counts as foo, but its path need not lie in foo's fences,
	// so this early return keeps the loop below from blocking it.
	if isOwnerPackage(loc, target.Obj().Pkg()) {
		return false
	}

	identPkgPath := target.Obj().Pkg().Path()

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
