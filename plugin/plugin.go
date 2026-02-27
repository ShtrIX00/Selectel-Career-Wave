package main

import (
	"github.com/ShtrlX00/Selectel-Career-Wave/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

// New is the entry point for golangci-lint v2 module plugins.
func New(conf any) ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{analyzer.Analyzer}, nil
}
