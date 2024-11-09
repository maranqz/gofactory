package main

import (
	"github.com/maranqz/gofactory"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(gofactory.NewAnalyzer())
}
