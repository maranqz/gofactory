package gofactory_test

import (
	"path/filepath"
	"testing"

	"github.com/maranqz/gofactory"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestLinterSuite(t *testing.T) {
	t.Parallel()

	root := moduleRoot()

	tests := map[string]struct {
		pkgs    []string
		prepare func(t *testing.T, a *analysis.Analyzer) error
	}{
		"simple":  {pkgs: []string{"simple/..."}},
		"casting": {pkgs: []string{"casting/..."}},
		"generic": {pkgs: []string{"generic/..."}},
		"packageGlobs": {
			pkgs: []string{"packageGlobs/..."},
			prepare: func(_ *testing.T, a *analysis.Analyzer) error {
				return a.Flags.Set("packageGlobs", "factory/packageGlobs/blocked/**")
			},
		},
		"packageGlobsOnly": {
			pkgs: []string{"packageGlobsOnly/main/..."},
			prepare: func(_ *testing.T, a *analysis.Analyzer) error {
				err := a.Flags.Set("packageGlobs", "factory/packageGlobsOnly/blocked/**")
				if err != nil {
					return err
				}

				return a.Flags.Set("packageGlobsOnly", "true")
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dirs := make([]string, 0, len(tt.pkgs))

			for _, pkg := range tt.pkgs {
				dirs = append(dirs, filepath.Join(root, pkg))
			}

			analyzer := gofactory.NewAnalyzer()

			if tt.prepare != nil {
				err := tt.prepare(t, analyzer)
				if err != nil {
					t.Fatal(err)
				}
			}

			analysistest.Run(t, root, analyzer, dirs...)
		})
	}
}

// testdataGoVersion is the go directive of every module-mode testdata module.
const testdataGoVersion = "1.26"

// TestTestdataRoots pins down how analysistest loads the two testdata roots,
// which later cases rely on. Module mode in analysistest is undocumented
// (x/tools v0.50.0, analysistest.loadPackages):
//
//   - A root holding a go.mod is loaded with GO111MODULE=on, GOPROXY=off and,
//     when the root also holds a go.work, GOWORK=<root>/go.work. Packages of
//     every module used by that go.work resolve without require directives,
//     including a nested module whose path lies under the main module path.
//   - Pass.Module is the module of the analysed package, not of the root:
//     factory/workspace sees "factory", the sibling package sees "sibling",
//     and the nested module package sees "factory/nestedmodule". Every
//     workspace module has Main set, so Main cannot single out the root.
//   - A root without a go.mod is a GOPATH-style tree (GO111MODULE=off,
//     packages under src/). Pass.Module is then non-nil but empty.
//
// Other drivers differ: go vet's unitchecker before Go 1.27 fills only Path,
// Version and GoVersion, and leaves Pass.Module nil without a module.
func TestTestdataRoots(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		root    string
		pkgs    []string
		modules map[string]analysis.Module
	}{
		"module": {
			root: moduleRoot(),
			pkgs: []string{
				filepath.Join(moduleRoot(), "workspace"),
				filepath.Join(moduleRoot(), "sibling"),
				filepath.Join(moduleRoot(), "nestedmodule"),
			},
			modules: map[string]analysis.Module{
				"factory/workspace": {
					Path:      "factory",
					Main:      true,
					GoVersion: testdataGoVersion,
				},
				"sibling": {
					Path:      "sibling",
					Main:      true,
					GoVersion: testdataGoVersion,
				},
				"factory/nestedmodule": {
					Path:      "factory/nestedmodule",
					Main:      true,
					GoVersion: testdataGoVersion,
				},
			},
		},
		"gopath": {
			root: gopathRoot(),
			pkgs: []string{"nomodule/..."},
			modules: map[string]analysis.Module{
				"nomodule":     {},
				"nomodule/ext": {},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			results := analysistest.Run(
				t, tt.root, gofactory.NewAnalyzer(), tt.pkgs...,
			)

			assertModules(t, results, tt.modules)
		})
	}
}

// assertModules checks that results cover exactly the packages in want and
// that each Pass.Module has the wanted Path, Version, Main and GoVersion.
func assertModules(
	t *testing.T,
	results []*analysistest.Result,
	want map[string]analysis.Module,
) {
	t.Helper()

	if len(results) != len(want) {
		t.Fatalf("got %d packages, want %d", len(results), len(want))
	}

	for _, res := range results {
		path := res.Pass.Pkg.Path()

		mod, ok := want[path]
		if !ok {
			t.Errorf("unexpected package %s", path)

			continue
		}

		got := res.Pass.Module
		if got == nil {
			t.Errorf("%s: Pass.Module is nil", path)

			continue
		}

		if got.Path != mod.Path || got.Version != mod.Version ||
			got.Main != mod.Main || got.GoVersion != mod.GoVersion {
			t.Errorf("%s: Pass.Module is "+
				"{Path: %q, Version: %q, Main: %t, GoVersion: %q}, "+
				"want {Path: %q, Version: %q, Main: %t, GoVersion: %q}",
				path, got.Path, got.Version, got.Main, got.GoVersion,
				mod.Path, mod.Version, mod.Main, mod.GoVersion)
		}
	}
}

// moduleRoot is the module-mode testdata root: module "factory" (go 1.26)
// with a go.work that also uses the sibling module "sibling" and the nested
// module "factory/nestedmodule".
func moduleRoot() string {
	return filepath.Join(analysistest.TestData(), "module")
}

// gopathRoot is the GOPATH-style testdata root, for runs without a module.
func gopathRoot() string {
	return filepath.Join(analysistest.TestData(), "gopath")
}
