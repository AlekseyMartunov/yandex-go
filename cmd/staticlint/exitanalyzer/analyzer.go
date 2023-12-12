// Package exitanalyzer for checking if call os.exit in main function
package exitanalyzer

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

// ExitAnalyzer can help find os.Exit in main function
var ExitAnalyzer = &analysis.Analyzer{
	Name: "exitCheck",
	Doc:  "this analyzer checks for the presence of a os.Exit call in main function",
	Run:  run,
}

// run the main function in ExitAnalyzer
func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			c, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}

			s, ok := c.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			i, ok := s.X.(*ast.Ident)
			if !ok || i.Name != "os" {
				return true
			}

			if s.Sel.Name == "Exit" {
				pos := s.Sel.NamePos
				pass.Reportf(pos, "calling os.Exit in main is prohibited")
			}
			return true
		})
	}
	return nil, nil
}
