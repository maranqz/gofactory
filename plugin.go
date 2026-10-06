package gofactory

import (
	"fmt"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

// register.Plugin must run on package load to make the plugin visible to
// golangci-lint's module-plugin loader.
//
//nolint:gochecknoinits
func init() {
	register.Plugin(name, newPlugin)
}

// settings mirrors the command-line flags in kebab-case, for golangci-lint
// module-plugin configuration (golangci-lint plugin settings are kebab-case
// by convention). These flat fields are today's only settings; a later
// ticket adding per-glob settings groups will need to introduce a grouping
// type, not just extend this struct.
//
// UseDefaultFactoryPattern is a pointer so that an absent key keeps the
// flags entry point's true default, while an explicit false drops it;
// register.DecodeSettings leaves it nil rather than false when the key is
// missing from JSON.
//
//nolint:tagliatelle
type settings struct {
	PackageGlobs     []string `json:"package-globs"`
	PackageGlobsOnly bool     `json:"package-globs-only"`
	ZeroValues       bool     `json:"zero-values"`

	FactoryPatterns          []string `json:"factory-patterns"`
	UseDefaultFactoryPattern *bool    `json:"use-default-factory-pattern"`
	OnlyWithFactory          bool     `json:"only-with-factory"`
}

type plugin struct {
	analyzer *analysis.Analyzer
}

//nolint:ireturn // register.NewPlugin dictates this exact signature.
func newPlugin(rawSettings any) (register.LinterPlugin, error) {
	decoded, err := register.DecodeSettings[settings](rawSettings)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}

	cfg := &config{useDefaultFactoryPattern: true}

	for _, g := range decoded.PackageGlobs {
		err = cfg.pkgGlobs.Set(g)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
	}

	cfg.onlyPkgGlobs = decoded.PackageGlobsOnly
	cfg.zeroValues = decoded.ZeroValues

	for _, p := range decoded.FactoryPatterns {
		err = cfg.extraFactoryPatterns.Set(p)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
	}

	if decoded.UseDefaultFactoryPattern != nil {
		cfg.useDefaultFactoryPattern = *decoded.UseDefaultFactoryPattern
	}

	cfg.onlyWithFactory = decoded.OnlyWithFactory

	return &plugin{analyzer: newAnalyzer(cfg)}, nil
}

func (p *plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{p.analyzer}, nil
}

func (p *plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
