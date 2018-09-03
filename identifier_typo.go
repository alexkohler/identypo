package identypo

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"os"
	"strings"

	"github.com/client9/misspell"
	"github.com/fatih/camelcase"
)

// Flags contains configuration specific to identypo.
// * Ignores - comma separated list of corrections to be ignored (for example, to stop corrections on "nto" and "creater", pass `-i="nto,creater"). This is a direct passthrough to the misspell package.
// * IncludeTests - include test files in analysis
// * FunctionsOnly - Find typos in function declarations only.
// * ConstantsOnly - Find typos in constants only.
// * VariablesOnly - Find typos in variables only.
// * SetExitStatus - Set exit status to 1 if any issues are found.
// Note: If FunctionsOnly, ConstantsOnly, and VariablesOnly are all false, every identifier will be searched for typos.
// (functions, function calls, variables, constants, type declarations, packages, labels).
type Flags struct {
	Ignores                                     string
	IncludeTests                                bool
	FunctionsOnly, ConstantsOnly, VariablesOnly bool
	SetExitStatus                               bool
}

// CheckForIdentiferTypos takes a slice of file arguments (this could be file names, directories, or packages (with or without the ... wildcard).
// Further configuration (such as words to ignore, whether or not to include tests, etc.) can be specified with the flags argument. Output is written
// using the log.Printf function. This is currently not configurable. For redirection to a file/buffer, see the log.SetOutput() method.
func CheckForIdentiferTypos(args []string, flags Flags) error {

	fset := token.NewFileSet()

	files, err := parseInput(args, fset, flags.IncludeTests)
	if err != nil {
		return fmt.Errorf("could not parse input %v", err)
	}

	return processIdentifiers(fset, files, flags)
}

func processIdentifiers(fset *token.FileSet, files []*ast.File, flags Flags) error {
	all := !flags.FunctionsOnly && !flags.ConstantsOnly && !flags.VariablesOnly

	retVis := &returnsVisitor{
		f:        fset,
		replacer: misspell.New(),
	}

	if len(flags.Ignores) > 0 {
		lci := strings.ToLower(flags.Ignores)
		retVis.replacer.RemoveRule(strings.Split(lci, ","))
	}

	retVis.replacer.Compile()

	for _, f := range files {
		if f == nil {
			continue
		}
		ast.Walk(retVis, f)
	}

	exitStatus := 0

	for _, ident := range retVis.identifiers {
		for _, word := range camelcase.Split(ident.Name) {
			v, d := retVis.replacer.Replace(word)
			if len(d) > 0 {
				exitStatus = 1
				file := retVis.f.File(ident.Pos())
				fileName := file.Name()
				line := file.Position(ident.Pos()).Line

				if all {
					// if we're including everything, no need to look at the kind of identifier we have
					log.Printf("%v:%v %q should be %v in %v\n", fileName, line, word, v, ident.Name)
				} else if ident.Obj != nil {
					switch ident.Obj.Kind {
					case ast.Fun:
						if !flags.FunctionsOnly {
							continue
						}
					case ast.Var:
						if !flags.VariablesOnly {
							continue
						}
					case ast.Con:
						if !flags.ConstantsOnly {
							continue
						}
					default:
						// labels, packages, etc. currently do not have individual flags and will be skipped
						continue
					}
					log.Printf("%v:%v %q should be %v in %v\n", fileName, line, word, v, ident.Name)
				}
			}
		}
	}

	if flags.SetExitStatus {
		os.Exit(exitStatus)
	}
	return nil
}

type returnsVisitor struct {
	f           *token.FileSet
	identifiers []*ast.Ident
	replacer    *misspell.Replacer
}

func (v *returnsVisitor) Visit(node ast.Node) ast.Visitor {
	funcDecl, ok := node.(*ast.Ident)
	if !ok {
		return v
	}

	v.identifiers = append(v.identifiers, funcDecl)

	return v
}
