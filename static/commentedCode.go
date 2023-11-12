package static

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var CommentedCodeAnalyzer = &analysis.Analyzer{
	Name: "commentedcode",
	Doc:  "reports potentially unnecessary commented code in the analyzed code",
	Run:  runCommentedCodeAnalyzer,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func runCommentedCodeAnalyzer(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CommentGroup)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		commentGroup, ok := n.(*ast.CommentGroup)
		if !ok {
			return
		}

		for _, comment := range commentGroup.List {
			// Проверяем, начинается ли комментарий с "//" (однострочный) или "/*" (многострочный)
			if comment.Text[0:2] == "//" || (len(comment.Text) > 1 && comment.Text[0:2] == "/*") {
				fmt.Printf("[WARN]: potentially unnecessary commented code at position %s\n", pass.Fset.Position(comment.Pos()))
			}
		}
	})

	return nil, nil
}
