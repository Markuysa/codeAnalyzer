package main

import (
	"flag"
	"fmt"
	"log"

	dynmic "github.com/Markuysa/codeAnalyzer/dynamic"
	"github.com/Markuysa/codeAnalyzer/static"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
)

func Execute() {
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
