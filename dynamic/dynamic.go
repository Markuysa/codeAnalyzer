package dynmic

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func AnalyzeCode(projectPath string) error {
	// Находим все файлы с расширением .go в проекте
	fileSet := token.NewFileSet()
	pkgs, err := parser.ParseDir(fileSet, projectPath, nil, parser.AllErrors)
	if err != nil {
		return err
	}

	// range through pkgs
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			// range through the file syntax tree
			ast.Inspect(file, func(n ast.Node) bool {
				// check "*t"
				if expr, ok := n.(*ast.StarExpr); ok {
					checkNilPointer(expr)
				}
				// check assignments to nil
				if assign, ok := n.(*ast.AssignStmt); ok {
					checkNilAssignment(assign)
				}

				return true
			})
		}
	}

	return nil
}

func checkNilPointer(expr *ast.StarExpr) {
	// Проверяем, что это указатель
	if _, ok := expr.X.(*ast.Ident); ok {
		// Выводим предупреждение в случае использования nil-указателя
		if isNil(expr) {
			fmt.Printf("Предупреждение: Использование nil-указателя в строке %d\n", expr.Pos())
		}
	}
}

func checkNilAssignment(assign *ast.AssignStmt) {
	// Проверяем, что присваивается nil
	for _, rhs := range assign.Rhs {
		if isNil(rhs) {
			fmt.Printf("Предупреждение: Присваивание nil в строке %d\n", assign.Pos())
		}
	}
}

func isNil(expr ast.Expr) bool {
	// Проверяем, что выражение - это nil
	// Используем switch для обработки разных типов выражений
	switch expr := expr.(type) {
	case *ast.Ident:
		// Если выражение - это идентификатор, проверяем, что он равен "nil"
		return expr.Name == "nil"
	case *ast.BasicLit:
		// Если выражение - это базовая литералка, проверяем, что она равна "nil"
		return expr.Value == "nil"
	default:
		return false
	}
}
