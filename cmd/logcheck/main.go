package main

import (
	"github.com/ShtrlX00/Selectel-Career-Wave/pkg/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}
