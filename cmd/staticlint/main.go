// Пакет предназначенный для статического анализа кода
// Запускается командой go run ./cmd/staticlint/main.go ./cmd/staticlint/no_os_exit.go ./...
package main

import (
	"github.com/kisielk/errcheck/errcheck"
	"go.uber.org/nilaway"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	var mychecks []*analysis.Analyzer

	mychecks = append(
		mychecks,
		// Стандартные статические анализаторы
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		// Публичные анализаторы
		nilaway.Analyzer,
		errcheck.Analyzer,
	)

	// Анализаторы из пакета staticcheck.io
	for _, v := range staticcheck.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}

	// Добавление моего анализатора
	mychecks = append(mychecks, Analyzer)

	multichecker.Main(
		mychecks...,
	)
}
