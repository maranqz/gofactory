package gofactory

import (
	"fmt"
	"strconv"

	"golang.org/x/tools/go/analysis"
)

// crossPackageDirectivesFlag is -crossPackageDirectives. Unlike the other
// flags, which only fill cfg for run to read later, Set here mutates
// analyzer.FactTypes directly: FactTypes is a static field that drivers
// inspect before Run to decide whether to analyse dependencies at all, so
// flipping it inside run would already be too late.
type crossPackageDirectivesFlag struct {
	analyzer *analysis.Analyzer
	enabled  bool
}

func newCrossPackageDirectivesFlag(
	analyzer *analysis.Analyzer,
) *crossPackageDirectivesFlag {
	return &crossPackageDirectivesFlag{analyzer: analyzer, enabled: true}
}

func (f *crossPackageDirectivesFlag) String() string {
	if f == nil {
		return "true"
	}

	return strconv.FormatBool(f.enabled)
}

func (f *crossPackageDirectivesFlag) IsBoolFlag() bool {
	return true
}

func (f *crossPackageDirectivesFlag) Set(s string) error {
	enabled, err := strconv.ParseBool(s)
	if err != nil {
		return fmt.Errorf("crossPackageDirectives: %w", err)
	}

	f.enabled = enabled
	f.analyzer.FactTypes = factTypesFor(enabled)

	return nil
}
