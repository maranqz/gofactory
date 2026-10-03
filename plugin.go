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
// module-plugin configuration. It is today's only, default, settings group;
// per-glob settings groups can be added later without breaking this shape.
// golangci-lint plugin settings are kebab-case by convention.
type settings struct {
	PackageGlobs     []string `json:"package-globs"`      //nolint:tagliatelle
	PackageGlobsOnly bool     `json:"package-globs-only"` //nolint:tagliatelle
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

	cfg := &config{}

	for _, g := range decoded.PackageGlobs {
		err = cfg.pkgGlobs.Set(g)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
	}

	cfg.onlyPkgGlobs = decoded.PackageGlobsOnly

	return &plugin{analyzer: newAnalyzer(cfg)}, nil
}

func (p *plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{p.analyzer}, nil
}

func (p *plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
