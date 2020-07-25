package util

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hsnodgrass/pufctl/pkg/ast"
)

// ParseFile reads a file and parses it using the Puppetfile AST
func ParseFile(path string) (*ast.Puppetfile, error) {
	text, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	file, err := ast.Parse(string(text))
	if err != nil {
		return nil, err
	}
	return file, nil
}

// ValidatePath checks that the file at the given path is a regular file
func ValidatePath(path string) error {
	pathStat, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !pathStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", path)
	}
	return nil
}
