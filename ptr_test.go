package unusual_generics_test

import (
	"fmt"

	"github.com/xakep666/unusual_generics"
)

func ExamplePtr() {
	pInt := unusual_generics.Ptr(10)
	pStr := unusual_generics.Ptr("test")
	pUint64 := unusual_generics.Ptr[uint64](24)
	fmt.Printf("%[1]T %[1]d\n", *pInt)
	fmt.Printf("%[1]T %[1]s\n", *pStr)
	fmt.Printf("%[1]T %[1]d\n", *pUint64)
	// Output:
	// int 10
	// string test
	// uint64 24
}
