package factory_test

// Tests for linters.

import (
	"path/filepath"
	"runtime"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/maranqz/go-factory-lint"
)

func TestLinterSuite(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()

	tests := map[string]struct {
		pkgs []string
	}{
		"creatable_anywhere": {pkgs: []string{"creatable_anywhere/..."}},
		"nested1":            {pkgs: []string{"nested1/..."}},
		"main":               {pkgs: []string{"..."}},
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
			if err := analyzer.Flags.Set("b", "factory/nested1"); err != nil {
				t.Fatal(err)
			}
			if err := analyzer.Flags.Set("blockedPkgs", "factory/nested2"); err != nil {
				t.Fatal(err)
			}

			analysistest.Run(t, TestdataDir(),
				analyzer, dirs...)
		})
	}
}

func TestdataDir() string {
	_, testFilename, _, ok := runtime.Caller(1)
	if !ok {
		panic("unable to get current test filename")
	}

	return filepath.Join(filepath.Dir(testFilename), "testdata")
}
