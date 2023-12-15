package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/maranqz/gofactory"
)

func main() {
	singlechecker.Main(gofactory.NewAnalyzer())
}
