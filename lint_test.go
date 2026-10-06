package gofactory_test

import (
	"net/url"
	"path/filepath"
	"testing"

	"github.com/golangci/plugin-module-register/register"
	"github.com/maranqz/gofactory"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

// pluginName is the name gofactory registers itself under with
// register.Plugin, which golangci-lint would use as the plugin key in its
// own configuration.
const pluginName = "gofactory"

// TestAnalyzerURL checks that both entry points produce an Analyzer.URL
// that parses as an absolute URL: golangci-lint fails to load an analyzer
// whose derived diagnostic URL does not parse.
func TestAnalyzerURL(t *testing.T) {
	t.Parallel()

	assertAbsoluteURL(t, "NewAnalyzer", gofactory.NewAnalyzer().URL)
	assertAbsoluteURL(t, "plugin", pluginAnalyzer(t, caseSettings{}).URL)
}

func assertAbsoluteURL(t *testing.T, entryPoint, raw string) {
	t.Helper()

	parsed, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("%s: Analyzer.URL %q does not parse: %v", entryPoint, raw, err)
	}

	if !parsed.IsAbs() {
		t.Fatalf("%s: Analyzer.URL %q is not an absolute URL", entryPoint, raw)
	}
}

// TestPluginRejectsBadSettings checks that the plugin fails on settings it
// cannot apply instead of silently ignoring them: a key spelled like the
// flag rather than in kebab-case, and a glob that does not compile.
func TestPluginRejectsBadSettings(t *testing.T) {
	t.Parallel()

	tests := map[string]map[string]any{
		"flag spelling": {"packageGlobs": []string{"factory/**"}},
		"invalid glob":  {"package-globs": []string{"["}},
	}
	for name, rawSettings := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			newPlugin, err := register.GetPlugin(pluginName)
			if err != nil {
				t.Fatal(err)
			}

			_, err = newPlugin(rawSettings)
			if err == nil {
				t.Fatalf("settings %v: got no error, want one", rawSettings)
			}
		})
	}
}

// TestLinterSuite runs every case through every entry point that populates
// the shared config: NewAnalyzer configured via Flags.Set, the way a
// command-line user or go vet driver would, and the golangci-lint plugin
// constructor configured via kebab-case settings. The flags analyzer also
// runs with Pass.Module shaped the way go vet passes it (unitcheckerAnalyzer).
func TestLinterSuite(t *testing.T) {
	t.Parallel()

	root := moduleRoot()

	tests := map[string]struct {
		pkgs     []string
		settings caseSettings
	}{
		"simple":    {pkgs: []string{"simple/..."}},
		"casting":   {pkgs: []string{"casting/..."}},
		"generic":   {pkgs: []string{"generic/..."}},
		"factories": {pkgs: []string{"factories/..."}},

		"dotimport": {pkgs: []string{"dotimport/..."}},

		"stdlib": {pkgs: []string{"stdlib/..."}},

		"packageGlobs": {
			pkgs: []string{"packageGlobs/..."},
			settings: caseSettings{
				packageGlobs: []string{"factory/packageGlobs/blocked/**"},
			},
		},
		"packageGlobsOnly": {
			pkgs: []string{"packageGlobsOnly/main/..."},
			settings: caseSettings{
				packageGlobs:     []string{"factory/packageGlobsOnly/blocked/**"},
				packageGlobsOnly: true,
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

			forEachEntryPoint(t, tt.settings,
				func(t *testing.T, analyzer *analysis.Analyzer) {
					analysistest.Run(t, root, analyzer, dirs...)
				})
		})
	}
}

func forEachEntryPoint(
	t *testing.T,
	settings caseSettings,
	check func(t *testing.T, analyzer *analysis.Analyzer),
) {
	t.Helper()

	entryPoints := map[string]func(*testing.T, caseSettings) *analysis.Analyzer{
		"flags":       flagsAnalyzer,
		"plugin":      pluginAnalyzer,
		"unitchecker": unitcheckerAnalyzer,
	}
	for name, build := range entryPoints {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			check(t, build(t, settings))
		})
	}
}

// caseSettings is one case's configuration, applied through either entry
// point: Flags.Set for NewAnalyzer, kebab-case settings for the plugin.
type caseSettings struct {
	packageGlobs     []string
	packageGlobsOnly bool
}

// flagsAnalyzer builds the analyzer through NewAnalyzer, configured via
// Flags.Set the way a command-line user or go vet driver would, mirroring
// pluginAnalyzer's golangci-lint entry point.
func flagsAnalyzer(t *testing.T, s caseSettings) *analysis.Analyzer {
	t.Helper()

	analyzer := gofactory.NewAnalyzer()

	for _, g := range s.packageGlobs {
		err := analyzer.Flags.Set("packageGlobs", g)
		if err != nil {
			t.Fatal(err)
		}
	}

	if s.packageGlobsOnly {
		err := analyzer.Flags.Set("packageGlobsOnly", "true")
		if err != nil {
			t.Fatal(err)
		}
	}

	return analyzer
}

// unitcheckerAnalyzer is flagsAnalyzer handed Pass.Module the way go vet's
// unitchecker before Go 1.27 fills it: Path, Version and GoVersion without
// Main, and nil without a module. analysistest sets Main on every module it
// loads and never passes nil, so only this entry point catches code relying
// on either.
func unitcheckerAnalyzer(t *testing.T, s caseSettings) *analysis.Analyzer {
	t.Helper()

	analyzer := flagsAnalyzer(t, s)
	run := analyzer.Run

	analyzer.Run = func(pass *analysis.Pass) (any, error) {
		vetPass := *pass
		vetPass.Module = nil

		if pass.Module != nil && pass.Module.Path != "" {
			vetPass.Module = &analysis.Module{
				Path:      pass.Module.Path,
				Version:   pass.Module.Version,
				GoVersion: pass.Module.GoVersion,
			}
		}

		return run(&vetPass)
	}

	return analyzer
}

// pluginAnalyzer builds the analyzer through the golangci-lint plugin entry
// point registered by gofactory's init, decoding the same settings a
// golangci-lint YAML/JSON config would supply in kebab-case.
func pluginAnalyzer(t *testing.T, s caseSettings) *analysis.Analyzer {
	t.Helper()

	newPlugin, err := register.GetPlugin(pluginName)
	if err != nil {
		t.Fatal(err)
	}

	rawSettings := map[string]any{
		"package-globs":      s.packageGlobs,
		"package-globs-only": s.packageGlobsOnly,
	}

	linterPlugin, err := newPlugin(rawSettings)
	if err != nil {
		t.Fatal(err)
	}

	analyzers, err := linterPlugin.BuildAnalyzers()
	if err != nil {
		t.Fatal(err)
	}

	if len(analyzers) != 1 {
		t.Fatalf("got %d analyzers, want 1", len(analyzers))
	}

	return analyzers[0]
}

// testdataGoVersion is the go directive of every module-mode testdata module.
const testdataGoVersion = "1.26"

// TestTestdataRoots pins down how analysistest loads the two testdata roots,
// which later cases rely on; it runs through every entry point. Module mode in
// analysistest is undocumented (x/tools v0.50.0, analysistest.loadPackages):
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
//
// The want comments in workspace/ and nomodule/ also pin the current-module
// rule and the no-module fallback; no other test runs those packages.
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

			forEachEntryPoint(t, caseSettings{},
				func(t *testing.T, analyzer *analysis.Analyzer) {
					results := analysistest.Run(t, tt.root, analyzer, tt.pkgs...)
					assertModules(t, results, tt.modules)
				})
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
// with a go.work that also uses the sibling modules "sibling" and
// "factoryext" and the nested module "factory/nestedmodule".
func moduleRoot() string {
	return filepath.Join(analysistest.TestData(), "module")
}

// gopathRoot is the GOPATH-style testdata root, for runs without a module.
func gopathRoot() string {
	return filepath.Join(analysistest.TestData(), "gopath")
}
