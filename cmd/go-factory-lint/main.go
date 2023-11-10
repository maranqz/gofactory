package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/maranqz/go-factory-lint"
)

func main() {
	singlechecker.Main(factory.NewAnalyzer())
}
