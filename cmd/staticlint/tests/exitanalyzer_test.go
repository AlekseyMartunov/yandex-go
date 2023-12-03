// Package tests for testing ExitAnalyzer
package tests

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/AlekseyMartunov/yandex-go.git/cmd/staticlint/exitanalyzer"
)

// TestMyAnalyzer for testing ExitAnalyzer
func TestMyAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), exitanalyzer.ExitAnalyzer, "./...")
}
