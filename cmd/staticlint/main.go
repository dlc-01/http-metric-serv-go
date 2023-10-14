package main

import (
	"github.com/dlc-01/http-metric-serv-go/internal/analytics"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(analytics.AllAnalyzers()...)
}
