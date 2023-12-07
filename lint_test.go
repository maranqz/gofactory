package factory_test

import (
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/maranqz/go-factory-lint/v2"
)

func TestLinterSuite(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()

	tests := map[string]struct {
		pkgs    []string
		prepare func(t *testing.T, a *analysis.Analyzer) error
	}{
		"simple":  {pkgs: []string{"simple/..."}},
		"casting": {pkgs: []string{"casting/..."}},
		"generic": {pkgs: []string{"generic/..."}},
		"packageGlobs": {
			pkgs: []string{"packageGlobs/..."},
			prepare: func(t *testing.T, a *analysis.Analyzer) error {
				return a.Flags.Set("packageGlobs", "factory/packageGlobs/blocked/**")
			},
		},
		"onlyPackageGlobs": {
			pkgs: []string{"onlyPackageGlobs/main/..."},
			prepare: func(t *testing.T, a *analysis.Analyzer) error {
				if err := a.Flags.Set("packageGlobs", "factory/onlyPackageGlobs/blocked/**"); err != nil {
					return err
				}

				return a.Flags.Set("onlyPackageGlobs", "true")
			},
		},
	}
	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dirs := make([]string, 0, len(tt.pkgs))

			for _, pkg := range tt.pkgs {
				dirs = append(dirs, filepath.Join(testdata, "src", "factory", pkg))
			}

			analyzer := factory.NewAnalyzer()

			if tt.prepare != nil {
				if err := tt.prepare(t, analyzer); err != nil {
					t.Fatal(err)
				}
			}

			analysistest.Run(t, testdata, analyzer, dirs...)
		})
	}
}
