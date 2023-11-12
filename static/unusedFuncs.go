package static

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var UnusedFuncAnalyzer = &analysis.Analyzer{
	Name: "unusedfuncs",
	Doc:  "reports unused functions in the analyzed code",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	// get the func names
	funcs := make(map[string]*ast.FuncDecl)
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			if funcDecl, ok := decl.(*ast.FuncDecl); ok {
				if funcDecl.Name.Name == "main" {
					continue
				}

				funcs[funcDecl.Name.Name] = funcDecl
			}
		}
	}

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	usedFuncs := make(map[string]bool)

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return
		}

		ast.Inspect(callExpr, func(node ast.Node) bool {
			if ident, ok := node.(*ast.Ident); ok {
				if funcDecl, found := funcs[ident.Name]; found {
					usedFuncs[funcDecl.Name.Name] = true
				}
			}
			return true
		})

		return
	})

	// check for unused functions
	for name, funcDecl := range funcs {
		if !usedFuncs[name] {
			fmt.Printf("[WARN]: unused function %s at position %s\n", name, pass.Fset.Position(funcDecl.Pos()))
		}
	}

	return nil, nil
}
