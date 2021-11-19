package unusual_generics_test

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/xakep666/unusual_generics"
)

func ExampleErrorAsPtr() {
	_, err := os.Open("non-existing")
	if pathErr := unusual_generics.ErrorAsPtr[*fs.PathError](err); pathErr != nil {
		fmt.Println("Failed at path:", (*pathErr).Path)
	}

	// Output:
	// Failed at path: non-existing
}

func ExampleErrorAs() {
	_, err := os.Open("non-existing")
	if pathErr, ok := unusual_generics.ErrorAs[*fs.PathError](err); ok {
		fmt.Println("Failed at path:", pathErr.Path)
	}

	// Output:
	// Failed at path: non-existing
}

func ExampleErrorIs() {
	if _, err := os.Open("non-existing"); err != nil {
		if unusual_generics.ErrorIs(err, fs.ErrNotExist) {
			fmt.Println("file does not exist")
		} else {
			fmt.Println(err)
		}
	}

	// Output:
	// file does not exist
}
