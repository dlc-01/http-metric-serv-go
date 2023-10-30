package main

import (
	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/dlc-01/http-metric-serv-go/internal/analytics"
)

func main() {
	multichecker.Main(analytics.AllAnalyzers()...)
}
