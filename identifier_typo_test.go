package identypo

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
	"testing"
)

func Test_CheckForIdentiferTypos(t *testing.T) {
	// flag/file parsing tests using testdata directory
	type args struct {
		wantLogs []string
		flags    Flags
		cliArgs  []string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "default flags, specifying package and ignoring tests",
			args: args{
				wantLogs: []string{
					"testdata/file.go:6 \"begining\" should be beginning in begining\n",
					"testdata/file.go:9 \"succesful\" should be successful in succesful\n",
					"testdata/file.go:12 \"succesful\" should be successful in succesful\n",
					"testdata/file.go:12 \"begining\" should be beginning in begining\n",
					"testdata/file.go:15 \"Succesful\" should be Successful in constantSuccesful\n",
					"testdata/file.go:19 \"authorithy\" should be authority in authorithyLoop\n",
					"testdata/file.go:22 \"authorithy\" should be authority in authorithyLoop\n",
					"testdata/file.go:26 \"Succesful\" should be Successful in varSuccesful\n",
				},
				flags: Flags{
					Ignores:      "",
					IncludeTests: false,
				},
				cliArgs: []string{
					"testdata",
				},
			},
		},
		{name: "default flags, specifying individual files",
			args: args{
				wantLogs: []string{
					"testdata/file_test.go:8 \"Begining\" should be Beginning in testBegining\n",
					"testdata/file_test.go:11 \"Succesful\" should be Successful in testSuccesful\n",
					"testdata/file_test.go:14 \"Succesful\" should be Successful in testSuccesful\n",
					"testdata/file_test.go:14 \"begining\" should be beginning in begining\n",
					"testdata/file_test.go:17 \"Succesful\" should be Successful in testConstantSuccesful\n",
					"testdata/file_test.go:20 \"Succesful\" should be Successful in TestSuccesful\n",
					"testdata/file_test.go:21 \"authorithy\" should be authority in authorithyLoop\n",
					"testdata/file_test.go:24 \"authorithy\" should be authority in authorithyLoop\n",
					"testdata/file.go:6 \"begining\" should be beginning in begining\n",
					"testdata/file.go:9 \"succesful\" should be successful in succesful\n",
					"testdata/file.go:12 \"succesful\" should be successful in succesful\n",
					"testdata/file.go:12 \"begining\" should be beginning in begining\n",
					"testdata/file.go:15 \"Succesful\" should be Successful in constantSuccesful\n",
					"testdata/file.go:19 \"authorithy\" should be authority in authorithyLoop\n",
					"testdata/file.go:22 \"authorithy\" should be authority in authorithyLoop\n",
					"testdata/file.go:26 \"Succesful\" should be Successful in varSuccesful\n",
				},
				flags: Flags{
					Ignores:      "",
					IncludeTests: true,
				},
				cliArgs: []string{
					"testdata/file_test.go",
					"testdata/file.go",
				},
			},
		},
		{name: "only functions",
			args: args{
				wantLogs: []string{
					"testdata/file.go:6 \"begining\" should be beginning in begining\n",
					"testdata/file_test.go:8 \"Begining\" should be Beginning in testBegining\n",
					"testdata/file_test.go:20 \"Succesful\" should be Successful in TestSuccesful\n",
				},
				flags: Flags{
					Ignores:       "",
					IncludeTests:  true,
					FunctionsOnly: true,
				},
				cliArgs: []string{
					"testdata/file.go",
					"testdata/file_test.go",
				},
			},
		},
		{name: "only constants",
			args: args{
				wantLogs: []string{
					"testdata/file.go:15 \"Succesful\" should be Successful in constantSuccesful\n",
					"testdata/file_test.go:17 \"Succesful\" should be Successful in testConstantSuccesful\n",
				},
				flags: Flags{
					Ignores:       "",
					IncludeTests:  true,
					ConstantsOnly: true,
				},
				cliArgs: []string{
					"testdata/file.go",
					"testdata/file_test.go",
				},
			},
		},
		{name: "only variables",
			args: args{
				wantLogs: []string{
					"testdata/file.go:26 \"Succesful\" should be Successful in varSuccesful\n",
				},
				flags: Flags{
					Ignores:       "",
					IncludeTests:  true,
					VariablesOnly: true,
				},
				cliArgs: []string{
					"testdata/file.go",
					"testdata/file_test.go",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var buf bytes.Buffer
			log.SetFlags(0)
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			err := CheckForIdentiferTypos(tt.args.cliArgs, tt.args.flags)
			if err != nil {
				t.Fatalf("CheckForIdentiferTypos %v", err)
			}

			var testBuffer bytes.Buffer
			testBuffer.Write([]byte(strings.Join(tt.args.wantLogs, "")))

			if buf.String() != testBuffer.String() {
				t.Fatalf("\ngot %v\nexp %v\n", buf.String(), testBuffer.String())
			}
		})
	}
}

func Test_processIdentifiers(t *testing.T) {
	type testFile struct {
		src      string
		name     string
		wantLogs []string
	}
	type args struct {
		testFiles []*testFile
		flags     Flags
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "misspelled function with no receiver",
			args: args{
				testFiles: []*testFile{
					{
						src: `package main
						func Propogate() {
						}`,
						name:     "file.go",
						wantLogs: []string{"file.go:2 \"Propogate\" should be Propagate in Propogate\n"},
					},
				},
				flags: Flags{
					Ignores: "",
				},
			},
		},
		{name: "misspelled function with no receiver which is corrected by hyphenation",
			args: args{
				testFiles: []*testFile{
					{
						src: `package main
						func Alltime() {
						}`,
						name:     "file.go",
						wantLogs: []string{"file.go:2 \"Alltime\" should be AllTime in Alltime\n"},
					},
				},
				flags: Flags{
					Ignores: "",
				},
			},
		},
		{name: "misspelled (unexported) function with no receiver which is corrected by hyphenation",
			args: args{
				testFiles: []*testFile{
					{
						src: `package main
						func alltime() {
						}`,
						name:     "file.go",
						wantLogs: []string{"file.go:2 \"alltime\" should be allTime in alltime\n"},
					},
				},
				flags: Flags{
					Ignores: "",
				},
			},
		},
		{name: "single misspelled function matching ignore",
			args: args{
				testFiles: []*testFile{
					{
						src: `package main
						func Propogate() {
						}`,
						name:     "file.go",
						wantLogs: []string{},
					},
				},
				flags: Flags{
					Ignores: "Propogate",
				},
			},
		},
		{name: "multiple misspelled (unexported) functions matching different ignores",
			args: args{
				testFiles: []*testFile{
					{
						src: `package main
						func nto() {}
						func propogate() {}`,
						name:     "file.go",
						wantLogs: []string{},
					},
				},
				flags: Flags{
					Ignores: "nto,propogate",
				},
			},
		},
		{name: "multiple misspelled functions in different files",
			args: args{
				testFiles: []*testFile{
					{
						src: `package main
					func PropogateMispellings() {
					}`,
						name:     "file1.go",
						wantLogs: []string{"file1.go:2 \"Propogate\" should be Propagate in PropogateMispellings\n"},
					},
					{
						src: `package main
					func AuthorithyFunc() {
					}`,
						name:     "file2.go",
						wantLogs: []string{"file2.go:2 \"Authorithy\" should be Authority in AuthorithyFunc\n"},
					},
				},
				flags: Flags{
					Ignores: "",
				},
			},
		},
		{name: "multiple misspelled function in same file (one with receiver, one without receiver)",
			args: args{
				testFiles: []*testFile{
					{
						src: `package main
						type rec struct{}
						func (r *rec) acheivement() {}
						func creater(){}
						`,
						name: "file.go",
						wantLogs: []string{
							"file.go:3 \"acheivement\" should be achievement in acheivement\n",
							"file.go:4 \"creater\" should be creature in creater\n",
						},
					},
				},
				flags: Flags{
					Ignores: "",
				},
			},
		},
		{name: "multiple misspelled identifiers in same file",
			args: args{
				testFiles: []*testFile{
					{
						src: `package main
						func main() {
							begining := true
							inital := true
						}`,
						name: "file.go",
						wantLogs: []string{
							"file.go:3 \"begining\" should be beginning in begining\n",
							"file.go:4 \"inital\" should be initial in inital\n",
						},
					},
				},
				flags: Flags{
					Ignores: "",
				},
			},
		},
		{name: "multiple misspelled identifiers in different files",
			args: args{
				testFiles: []*testFile{
					{
						src: `package main
							  var begining = true
						`,
						name: "file1.go",
						wantLogs: []string{
							"file1.go:2 \"begining\" should be beginning in begining\n",
						},
					},
					{
						src: `package main
							  var inital = true
						`,
						name: "file1.go",
						wantLogs: []string{
							"file1.go:2 \"inital\" should be initial in inital\n",
						},
					},
				},
				flags: Flags{
					Ignores: "",
				},
			},
		},
		{name: "multiple misspelled constants",
			args: args{
				testFiles: []*testFile{
					{
						src: `package main
							const begining = true
							const inital = true
						`,
						name: "file.go",
						wantLogs: []string{
							"file.go:2 \"begining\" should be beginning in begining\n",
							"file.go:3 \"inital\" should be initial in inital\n",
						},
					},
				},
				flags: Flags{
					Ignores: "",
				},
			},
		},
		{name: "multiple misspelled constants in const block",
			args: args{
				testFiles: []*testFile{
					{
						src: `package main
							const ( 
							   begining = true
							   inital = true
							)
						`,
						name: "file.go",
						wantLogs: []string{
							"file.go:3 \"begining\" should be beginning in begining\n",
							"file.go:4 \"inital\" should be initial in inital\n",
						},
					},
				},
				flags: Flags{
					Ignores: "",
				},
			},
		},
		{name: "misspelled function/variable with no function/variable filter (only constants=true)",
			args: args{
				testFiles: []*testFile{
					{
						src: `package main
						var propogate = false
						func Propogate() {
						}`,
						name:     "file.go",
						wantLogs: []string{},
					},
				},
				flags: Flags{
					Ignores:       "",
					ConstantsOnly: true,
				},
			},
		},
		{name: "misspelled function/variable/constant (with functions/variables/constants=true)",
			args: args{
				testFiles: []*testFile{
					{
						src: `package main
						const begining = true
						var propogate = false
						 func PropogateFunc() {}`,
						name: "file.go",
						wantLogs: []string{
							"file.go:2 \"begining\" should be beginning in begining\n",
							"file.go:3 \"propogate\" should be propagate in propogate\n",
							"file.go:4 \"Propogate\" should be Propagate in PropogateFunc\n",
						},
					},
				},
				flags: Flags{
					Ignores:       "",
					FunctionsOnly: true,
					VariablesOnly: true,
					ConstantsOnly: true,
				},
			},
		},
		{name: "misspelling at beginning, middle, and end of function (and casing permutations)",
			args: args{
				testFiles: []*testFile{
					{
						src: `package main
						const begining = 0
						const Begining = 0
						const fooBegining = 0
						const fooBeginingBar = 0
						const FooBeginingBar = 0
						const FooBegining = 0
						const beginingBar = 0
						const BeginingBar = 0`,
						name: "file.go",
						wantLogs: []string{
							"file.go:2 \"begining\" should be beginning in begining\n",
							"file.go:3 \"Begining\" should be Beginning in Begining\n",
							"file.go:4 \"Begining\" should be Beginning in fooBegining\n",
							"file.go:5 \"Begining\" should be Beginning in fooBeginingBar\n",
							"file.go:6 \"Begining\" should be Beginning in FooBeginingBar\n",
							"file.go:7 \"Begining\" should be Beginning in FooBegining\n",
							"file.go:8 \"begining\" should be beginning in beginingBar\n",
							"file.go:9 \"Begining\" should be Beginning in BeginingBar\n",
						},
					},
				},
				flags: Flags{
					Ignores: "",
				},
			},
		},
		{name: "misspelled package",
			args: args{
				testFiles: []*testFile{
					{
						src:  `package inital`,
						name: "file.go",
						wantLogs: []string{
							"file.go:1 \"inital\" should be initial in inital\n",
						},
					},
				},
				flags: Flags{
					Ignores: "",
				},
			},
		},
		{name: "misspelled label",
			args: args{
				testFiles: []*testFile{
					{
						src: `package main
								func main() {
								initalLabel:
									for i := 0; i < 5; i++ {
										fmt.Println("zoop")
									}
								}`,
						name: "file.go",
						wantLogs: []string{
							"file.go:3 \"inital\" should be initial in initalLabel\n",
						},
					},
				},
				flags: Flags{
					Ignores: "",
				},
			},
		},
		{name: "misspelled type declaration",
			args: args{
				testFiles: []*testFile{
					{
						src: `package main
								type initalType interface{}
								`,
						name: "file.go",
						wantLogs: []string{
							"file.go:2 \"inital\" should be initial in initalType\n",
						},
					},
				},
				flags: Flags{
					Ignores: "",
				},
			},
		},
		{name: "misspelled function inside of interface",
			args: args{
				testFiles: []*testFile{
					{
						src: `package main
								type myInterface interface{
									inital() error
								}
								`,
						name: "file.go",
						wantLogs: []string{
							"file.go:3 \"inital\" should be initial in inital\n",
						},
					},
				},
				flags: Flags{
					Ignores: "",
				},
			},
		},
		{name: "misspelled function call",
			args: args{
				testFiles: []*testFile{
					{
						src: `package main
								func main() {
									err := a.inital()
								}
								`,
						name: "file.go",
						wantLogs: []string{
							"file.go:3 \"inital\" should be initial in inital\n",
						},
					},
				},
				flags: Flags{
					Ignores: "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fset := token.NewFileSet() // positions are relative to fset
			files := make([]*ast.File, len(tt.args.testFiles))
			for _, testFile := range tt.args.testFiles {
				f, err := parser.ParseFile(fset, testFile.name, testFile.src, 0)
				if err != nil {
					t.Fatalf("Did not expect error parsing file, %v", err)
				}
				files = append(files, f)
			}

			var buf bytes.Buffer
			log.SetFlags(0)
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			err := processIdentifiers(fset, files, tt.args.flags)
			if err != nil {
				t.Fatalf("processIdentifiers %v", err)
			}

			var testBuffer bytes.Buffer
			for _, testFile := range tt.args.testFiles {
				testBuffer.Write([]byte(strings.Join(testFile.wantLogs, "")))
			}

			if buf.String() != testBuffer.String() {
				t.Fatalf("\ngot %v\nexp %v\n", buf.String(), testBuffer.String())
			}
		})
	}
}
