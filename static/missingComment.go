package static

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var MissingDocsAnalyzer = &analysis.Analyzer{
	Name: "missingdocs",
	Doc:  "reports functions without doc comments in the analyzed code",
	Run:  runMissingDocsAnalyzer,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func runMissingDocsAnalyzer(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok {
			return
		}

		// Check if the function has a doc comment
		if funcDecl.Doc == nil || len(funcDecl.Doc.List) == 0 {
			// No doc comment found
			fmt.Printf("[WARN]: missing doc comment for function: %s at position %s\n", funcDecl.Name.Name, pass.Fset.Position(funcDecl.Pos()))
		}
	})

	return nil, nil
}
