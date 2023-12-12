// Package main for start analysis application
package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"strings"

	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"

	"github.com/AlekseyMartunov/yandex-go.git/cmd/staticlint/exitanalyzer"
)

func main() {
	var myChecks []*analysis.Analyzer

	checks := map[string]bool{
		"ST1003": true,
		"ST1006": true,
	}

	for _, v := range stylecheck.Analyzers {
		if checks[v.Analyzer.Name] {
			myChecks = append(myChecks, v.Analyzer)
		}
	}

	for _, v := range staticcheck.Analyzers {
		if strings.HasPrefix(v.Analyzer.Name, "SA") {
			myChecks = append(myChecks, v.Analyzer)
		}
	}

	myChecks = append(myChecks,
		printf.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,

		exitanalyzer.ExitAnalyzer,
		sortslice.Analyzer,
		unmarshal.Analyzer,
	)

	multichecker.Main(
		myChecks...,
	)
}
