package factory

// Tests for linters.

import (
	"path/filepath"
	"runtime"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestLinterSuite(t *testing.T) {
	testdata := analysistest.TestData()

	tests := map[string]struct {
		pkgs []string
	}{
		"creatable_anywhere": {pkgs: []string{"creatable_anywhere/..."}},
		"nested":             {pkgs: []string{"nested/..."}},
		"main":               {pkgs: []string{"..."}},
	}
	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			dirs := make([]string, 0, len(tt.pkgs))

			for _, pkg := range tt.pkgs {
				dirs = append(dirs, filepath.Join(testdata, "src", "factory", pkg))
			}

			analysistest.Run(t, TestdataDir(),
				FactoryAnalyzer, dirs...)
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
