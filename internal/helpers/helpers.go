package helpers

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/hsnodgrass/pufctl/internal/sources/gitsource"

	"github.com/hsnodgrass/pufctl/internal/auth"

	"github.com/hsnodgrass/pufctl/internal/validators"

	"github.com/hsnodgrass/pufctl/internal/logging"
	"github.com/hsnodgrass/pufctl/internal/uitext"
	"github.com/hsnodgrass/pufctl/pkg/puppetfileparser/ast"
)

// ParseOptions holds options related to parsing Puppetfiles
// Specifically, these options are used for parsing Puppetfiles
// from Git sources.
type ParseOptions struct {
	Verbose  bool
	GitRef   string
	SSHKey   string
	Username string
	Password string
	Token    string
}

// PromptForInput prompts a user to enter input in the terminal and
// returns their input as a string
func PromptForInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s  ", prompt)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	return text
}

// PromptConfirmFile is used with the out-file or write-in-place
// flags and wraps appropriate user confirmation prompts. If confirm
// is true, the user prompts are bypassed.
func PromptConfirmFile(path, data string, confirm bool) error {
	err := ValidatePath(path)
	if err != nil {
		logging.Infof("New file path: %s\n", path)
		return PromptNew(path, data, confirm)
	}
	logging.Infof("Existing file path: %s\n", path)
	return PromptOverwrite(path, data, confirm)
}

// PromptOverwrite is a wrapper for the WriteInPlace function that
// gives the user a confirmation prompt before writing.
func PromptOverwrite(path, data string, confirm bool) error {
	if confirm {
		return OverwriteFile(path, data)
	}
	text := PromptForInput(uitext.ConfirmOverwritePrompt)
	if strings.HasPrefix(text, "y") {
		return OverwriteFile(path, data)
	}
	logging.Infof("Not overwriting file %s\n", path)
	return nil
}

// PromptNew is a wrapper for the OutFile function that
// gives the user a confirmation prompt before writing.
func PromptNew(path, data string, confirm bool) error {
	if confirm {
		return NewFile(&path, data)
	}
	text := PromptForInput(uitext.ConfirmNewWritePrompt)
	if strings.HasPrefix(text, "y") {
		return NewFile(&path, data)
	}
	logging.Infof("Not creating file %s\n", path)
	return nil
}

// MaxBools takes two bools and returns the max "truthiness".
// I.E. if bool1 is false and bool2 is true, return true.
// The only time this function returns false is when both bools are false.
func MaxBools(one, two bool) bool {
	if one || two {
		return true
	}
	return false
}

// Parse parses a Puppetfile from either a Git source or a file on disk
func Parse(target string, opts ParseOptions) (*ast.Puppetfile, error) {
	isGit := validators.IsGitURL(target)
	isGitNoSuf := validators.IsGitURLNoSuffix(target)
	if isGit || isGitNoSuf {
		return parseFromGit(target, opts)
	}
	return ParseFile(target)
}

func parseFromGit(url string, opts ParseOptions) (*ast.Puppetfile, error) {
	var normalized string
	if validators.IsGitURLNoSuffix(url) {
		logging.Warnln("URL does not have the suffix \".git\"")
		logging.Warnln("Appending the \".git\" suffix to the URL")
		normalized = fmt.Sprintf("%s.git", url)
	} else {
		normalized = url
	}
	auth, err := auth.GitAuth(normalized, opts.Username, opts.Password, opts.Token, opts.SSHKey)
	if err != nil {
		return nil, fmt.Errorf("Failed to retreive Git auth object: %w", err)
	}
	pfile, branch, err := gitsource.CloneFile(normalized, "Puppetfile", opts.GitRef, auth)
	if err != nil {
		return nil, fmt.Errorf("Failed to clone Puppetfile from %s#%s: %w", normalized, opts.GitRef, err)
	}
	logging.Debugf("Cloned Puppetfile from Git repo %s#%s\n", normalized, branch)
	textbytes, err := ioutil.ReadAll(pfile)
	if err != nil {
		return nil, fmt.Errorf("Failed to read Puppetfile: %w", err)
	}
	puppetfile, err := ast.Parse(string(textbytes))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse Puppetfile text: %w", err)
	}
	logging.Debugf("Using Puppetfile from %s#%s\n", normalized, branch)
	return puppetfile, nil
}

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
	logging.Debugln("Using Puppetfile ", path)
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

// NewFile writes given data to a new file at the given path.
// The function should only be used for creating new files as
// it does not care if there is a file at the target path already
// and will not attempt to recover the original file in the event
// that an error occurs during the write process.
func NewFile(path *string, content string) error {
	logging.Infof("Preparing to write Puppetfile to %s\n", *path)
	f, err := os.Create(*path)
	if err != nil {
		return err
	}
	defer f.Close()
	logging.Infoln("Writing sorted Puppetfile to file")
	_, err = f.WriteString(content)
	if err != nil {
		return err
	}
	logging.Infof("Puppetfile successfully written to %s\n", *path)
	return nil
}

// OverwriteFile writes given data to the given file that exists already.
// The function attempts to recover the original file if an error occurs
// during the write process.
func OverwriteFile(path string, data string) error {
	logging.Infoln("Preparing to overwrite file", path)
	original, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, []byte(data), 0755)
	if err != nil {
		e := ioutil.WriteFile(path, []byte(original), 0755)
		if e != nil {
			logging.Errorln("Error occured and original file could not be restored")
		}
		logging.Errorln("Error occured, but original file was restored")
		return err
	}
	logging.Infof("Successfully saved data to %s\n", path)
	return nil
}

// ModStringToSlug converts a module string "<org>/<module> to
// a slug "<org>-<module>".
func ModStringToSlug(mod string) (string, error) {
	parts := strings.Split(mod, "/")
	if len(parts) == 2 {
		return fmt.Sprintf("%s-%s", parts[0], parts[1]), nil
	}
	return "", fmt.Errorf("Failed to convert module string \"%s\" to slug", mod)
}

// ModSlugToString converts a module slug "<org>-<module>" to
// a string "<org>/<module>"
func ModSlugToString(mod string) (string, error) {
	parts := strings.Split(mod, "-")
	if len(parts) == 2 {
		return fmt.Sprintf("%s/%s", parts[0], parts[1]), nil
	}
	return "", fmt.Errorf("Failed to convert slug \"%s\" to module string", mod)
}
