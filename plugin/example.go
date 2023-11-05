// This must be package main
package main

import (
	"fmt"

	"golang.org/x/tools/go/analysis"

	"github.com/maranqz/go-factory-lint"
)

func New(conf any) ([]*analysis.Analyzer, error) {
	// TODO: This must be implemented
	fmt.Printf("My configuration (%[1]T): %#[1]v\n", conf)

	// The configuration type will be map[string]any or []interface, it depends on your configuration.
	// You can use https://github.com/mitchellh/mapstructure to convert map to struct.

	return []*analysis.Analyzer{factory.FactoryAnalyzer}, nil
}
