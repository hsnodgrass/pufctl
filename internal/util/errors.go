package util

import (
	"fmt"
	"os"
)

// ExitErr prints the error to the console and then exits with a code of 1
func ExitErr(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}
