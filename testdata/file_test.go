package testdata

import "testing"

import "fmt"

// misspelled function
func testBegining() {}

// misspelled type declaration
type testSuccesful int

// misspelled function with receiver
func (s *testSuccesful) begining() {}

// misspelled constant
const testConstantSuccesful = 0

// misspelled test and label
func TestSuccesful(t *testing.T) {
authorithyLoop:
	for i := 0; i < 5; i++ {
		fmt.Println("loooooooool")
		continue authorithyLoop
	}
}
