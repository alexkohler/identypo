package testdata

import "fmt"

// misspelled function
func begining() {}

// misspelled type declaration
type succesful int

// misspelled function with receiver
func (s *succesful) begining() {}

// misspelled constant
const constantSuccesful = 0

// misspelled label
func main() {
authorithyLoop:
	for i := 0; i < 5; i++ {
		fmt.Println("loooooooool")
		continue authorithyLoop
	}
}

var varSuccesful = 0
