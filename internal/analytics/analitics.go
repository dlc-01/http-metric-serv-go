package analytics

import (
	"github.com/charithe/durationcheck"
	"github.com/kisielk/errcheck/errcheck"
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"honnef.co/go/tools/staticcheck"
)

var Analyzers = []*analysis.Analyzer{}

// StandardPasses uses the standard static analyser from the package golang.org/x/tools/go/analysis/passes
func StandardPasses() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		atomicalign.Analyzer,
		bools.Analyzer,
		buildssa.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		ctrlflow.Analyzer,
		deepequalerrors.Analyzer,
		directive.Analyzer,
		errorsas.Analyzer,
		fieldalignment.Analyzer,
		findcall.Analyzer,
		framepointer.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		inspect.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		pkgfact.Analyzer,
		printf.Analyzer,
		reflectvaluecompare.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		sigchanyzer.Analyzer,
		stdmethods.Analyzer,
		structtag.Analyzer,
		tests.Analyzer,
		timeformat.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
		unusedwrite.Analyzer,
		usesgenerics.Analyzer,
	}

}

// StaticChecks uses the analyser from the staticcheck package
func StaticChecks() []*analysis.Analyzer {
	checks := map[string]bool{
		"SA":  true,
		"SA1": true,
		"SA2": true,
		"SA3": true,
		"SA4": true,
		"SA5": true,
		"SA6": true,
		"SA9": true,
		"S1":  true,
		"ST1": true,
		"QF1": true,
	}
	var staticchecks []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {

		if checks[v.Analyzer.Name] {
			staticchecks = append(staticchecks, v.Analyzer)
		}
	}
	return staticchecks
}

// PublicAnalysisTool use public analysers, namely errcheck (to check for unhandled errors),
// durationcheck check the project for two time.Duration values which are multiplied.
func PublicAnalysisTool() []*analysis.Analyzer {
	return []*analysis.Analyzer{errcheck.Analyzer, durationcheck.Analyzer}
}

func init() {
	Analyzers = []*analysis.Analyzer{
		ExitCheckAnalyzer,
	}
}

// ExitCheckAnalyzer - checking os.Exit() calls in main function main packages
var ExitCheckAnalyzer = &analysis.Analyzer{
	Name: "exitCheck",
	Doc:  "Checking os.Exit() calls in main function main packages",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if file.Name.Name == "main" {
			continue
		}
		inMain := false
		ast.Inspect(file, func(node ast.Node) bool {
			if isNameFunc(node, "main") {
				inMain = true
				return true
			}
			if inMain && isNameFunc(node, "Exit") {
				pass.Reportf(node.Pos(), "don't use exit in main")
			}

			return true
		})
	}
	return nil, nil
}

func isNameFunc(node ast.Node, name string) bool {
	if id, ok := node.(*ast.Ident); ok {
		if id.Name == name {
			return true
		}
	}
	return false
}

func AllAnalyzers() []*analysis.Analyzer {
	return GroupingAnalyzers(StaticChecks(), StandardPasses(), StaticChecks(), Analyzers, PublicAnalysisTool())
}

// GroupingAnalyzers - Grouping of analysers
func GroupingAnalyzers(a ...[]*analysis.Analyzer) []*analysis.Analyzer {
	joined := make([]*analysis.Analyzer, 0, 100)
	for _, l := range a {
		joined = append(joined, l...)
	}
	return joined
}
