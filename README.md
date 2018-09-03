# identypo

identypo is a Go static analysis tool to find typos in identifiers (functions, function calls, variables, constants, type declarations, packages, labels). It is built on top of [client9's misspell package](https://github.com/client9/misspell).

## Installation

    go get -u github.com/alexkohler/identypo/cmd/identypo

## Usage

Similar to other Go static analysis tools (such as golint, go vet), identypo can be invoked with one or more filenames, directories, or packages named by its import path. Identypo also supports the `...` wildcard. By default, it will search for typos in every identifier (functions, function calls, variables, constants, type declarations, packages, labels).

    identypo [flags] files/directories/packages

### Flags
- **-tests** (default true) - Include test files in analysis
- **-i** - Comma separated list of corrections to be ignored (for example, to stop corrections on "nto" and "creater", pass `-i="nto,creater"). This is a direct passthrough to the misspell package.
- **-functions** - Find typos in function declarations only.
- **-constants** - Find typos in constants only.
- **-variables** - Find typos in variables only.
- **-set_exit_status** (default false) - Set exit status to 1 if any issues are found.

NOTE: by default, identypo will check for typos in every identifier (functions, function calls, variables, constants, type declarations, packages, labels). In this case, no flag needs specified. Due to a lack of frequency, there are currently no flags to find only type declarations, packages, or labels.

## Example uses in popular Go repos

Some selected examples from [Kubernetes](https://github.com/kubernetes/kubernetes):
```Bash
$ identypo ./...
cmd/kubeadm/app/cmd/phases/kubeconfig_test.go:325 "Authorithy" should be Authority in SetupPkiDirWithCertificateAuthorithy
cmd/kubeadm/app/util/apiclient/wait.go:51 "inital" should be initial in initalTimeout
pkg/apis/certificates/types.go:125 "Committment" should be Commitment in UsageContentCommittment
controller/nodeipam/ipam/cidrset/cidr_set.go:158 "Begining" should be Beginning in getBeginingAndEndIndices
staging/src/k8s.io/apimachinery/pkg/conversion/converter_test.go:358 "Overriden" should be Overridden in TestConverter_WithConversionOverriden
```

```Go
// cmd/kubeadm/app/cmd/phases/kubeconfig_test.go:325 "Authorithy" should be Authority in SetupPkiDirWithCertificateAuthorithy
pkidir := testutil.SetupPkiDirWithCertificateAuthorithy(t, tmpdir)

// cmd/kubeadm/app/util/apiclient/wait.go:51 "inital" should be initial in initalTimeout
WaitForHealthyKubelet(initalTimeout time.Duration, healthzEndpoint string) error

// pkg/apis/certificates/types.go:125 "Committment" should be Commitment in UsageContentCommittment
UsageContentCommittment KeyUsage = "content commitment"

// controller/nodeipam/ipam/cidrset/cidr_set.go:158 "Begining" should be Beginning in getBeginingAndEndIndices
func (s *CidrSet) getBeginingAndEndIndices(cidr *net.IPNet) (begin, end int, err error) {

// staging/src/k8s.io/apimachinery/pkg/conversion/converter_test.go:358 "Overriden" should be Overridden in TestConverter_WithConversionOverriden
func TestConverter_WithConversionOverriden(t *testing.T) {
```


Some examples from the [Go standard library](https://github.com/golang/go) (utilizing the `-i` flag to suppress some non-isses):

```Bash
$ identypo -i="rela,nto,onot,alltime" ./...
cmd/trace/goroutines.go:169 "dividened" should be dividend in dividened
cmd/trace/goroutines.go:173 "dividened" should be dividend in dividened
cmd/trace/goroutines.go:175 "dividened" should be dividend in dividened
cmd/trace/goroutines.go:179 "dividened" should be dividend in dividened
cmd/trace/annotations.go:1162 "dividened" should be dividend in dividened
cmd/trace/annotations.go:1166 "dividened" should be dividend in dividened
cmd/trace/annotations.go:1168 "dividened" should be dividend in dividened
cmd/trace/annotations.go:1172 "dividened" should be dividend in dividened
crypto/x509/verify.go:208 "Comparisions" should be Comparisons in MaxConstraintComparisions
crypto/x509/verify.go:585 "Comparisions" should be Comparisons in MaxConstraintComparisions
```

```Go
// cmd/trace/annotations.go:1162-1172 dividened" should be dividend in dividened
"percent": func(dividened, divisor int64) template.HTML {
	if divisor == 0 {
		return ""
	}
	return template.HTML(fmt.Sprintf("(%.1f%%)", float6(dividened)/float64(divisor)*100))
},
"barLen": func(dividened, divisor int64) template.HTML {
	if divisor == 0 {
		return "0"
	}
	return template.HTML(fmt.Sprintf("%.2f%%", float6(dividened)/float64(divisor)*100))
},

// crypto/x509/verify.go:208 "Comparisions" should be Comparisons in MaxConstraintComparisions
type VerifyOptions struct {
	...
	Roots         *CertPool // if nil, the system roots are used
	CurrentTime   time.Time // if zero, the current time is used
	...
	MaxConstraintComparisions int
}
```


## Packages used
- https://github.com/client9/misspell
- https://github.com/fatih/camelcase



## Contributing

Please open an issue and/or a PR for any features/bugs. 


## Other static analysis tools

If you've enjoyed identypo, take a look at my other static anaylsis tools!
- [prealloc](https://github.com/alexkohler/prealloc) - Finds slice declarations that could potentially be preallocated.
- [nakedret](https://github.com/alexkohler/nakedret) - Finds naked returns.
- [unimport](https://github.com/alexkohler/unimport) - Finds unnecessary import aliases.