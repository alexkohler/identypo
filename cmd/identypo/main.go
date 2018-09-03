package main

import (
	"flag"
	"go/build"
	"log"
	"os"

	"github.com/alexkohler/identypo"
)

func init() {
	build.Default.UseAllFiles = false
}

func usage() {
	log.Printf("Usage of %s:\n", os.Args[0])
	log.Printf("\nidentypo[flags] # runs on package in current directory\n")
	log.Printf("\nidentypo [flags] [packages]\n")
	log.Printf("Flags:\n")
	flag.PrintDefaults()
	log.Printf("\nNOTE: by default, identypo will check for typos in every identifier (functions, function calls, variables, constants, type declarations, packages, labels). In this case, no flag needs specified.\n")
}

func main() {

	// Remove log timestamp
	log.SetFlags(0)

	ignores := flag.String("i", "", "ignore the following words requiring correction, comma separated (e.g. -i=\"nto,creater\")")
	includeTests := flag.Bool("tests", true, "include test (*_test.go) files")
	functionsOnly := flag.Bool("functions", false, "find typos in function declarations only")
	constantsOnly := flag.Bool("constants", false, "find typos in constants only")
	variablesOnly := flag.Bool("variables", false, "find typos in variables only")
	setExitStatus := flag.Bool("set_exit_status", false, "Set exit status to 1 if any issues are found")
	flag.Usage = usage
	flag.Parse()

	flags := identypo.Flags{
		Ignores:       *ignores,
		IncludeTests:  *includeTests,
		FunctionsOnly: *functionsOnly,
		ConstantsOnly: *constantsOnly,
		VariablesOnly: *variablesOnly,
		SetExitStatus: *setExitStatus,
	}

	if err := identypo.CheckForIdentiferTypos(flag.Args(), flags); err != nil {
		log.Println(err)
	}
}
