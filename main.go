package main

import (
	dynmic "code_analyzer/dynamic"
	"code_analyzer/static"
	"flag"
	"fmt"
	"log"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	analyzers := []*analysis.Analyzer{
		static.MissingDocsAnalyzer,
		static.UnusedFuncAnalyzer,
		static.CommentedCodeAnalyzer,
	}

	multichecker.Main(analyzers...)

	err := dynmic.AnalyzeCode(".")
	if err != nil {
		fmt.Println("Ошибка при анализе кода:", err)
	}
}
