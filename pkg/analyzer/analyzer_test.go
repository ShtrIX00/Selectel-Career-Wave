package analyzer_test

import (
	"testing"

	"github.com/retro/vas/pkg/analyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analyzer.SetConfig(analyzer.DefaultConfig())
	analysistest.Run(t, testdata, analyzer.Analyzer, "a")
}
