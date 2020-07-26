package validators

import (
	"regexp"
	"strings"

	"github.com/hsnodgrass/pufctl/internal/logging"
)

// IsConfirmed validates that input is a confirmation and returns a bool where "yes" = true
func IsConfirmed(input string) bool {
	confirm := strings.TrimSpace(input)
	switch strings.ToLower(confirm) {
	case "y", "ye", "yes":
		return true
	default:
		return false
	}
}

// IsSlug validates that the input is a "slug", a string valid in URLs, and returns a boolean
func IsSlug(input string) bool {
	matched, err := regexp.MatchString(`^[a-zA-Z-]*$`, input)
	if err != nil {
		return false
	}
	logging.Debugf("Input %s is a valid slug", input)
	return matched
}

// IsModuleString validates that an input string follows the Puppet
// module identifier string convention of <org>-<module_name>
func IsModuleString(input string) (bool, bool) {
	slug := reModSlug.MatchString(input)
	if slug {
		return true, false
	}
	ref := reModString.MatchString(input)
	if ref {
		return true, true
	}
	return false, false
}

// IsModuleSlug validates that an input string is a slug
// identifier of a module: <org>-<module>
func IsModuleSlug(input string) bool {
	slug := reModSlug.MatchString(input)
	if slug {
		return true
	}
	return false
}

// IsModuleRef validates that an input string is a ref
// identifier of a module: <org>/<module>
func IsModuleRef(input string) bool {
	s := reModString.MatchString(input)
	if s {
		return true
	}
	return false
}

// IsGitURL validates that an input string is a git URL.
func IsGitURL(input string) bool {
	url := reGitURL.MatchString(input)
	if url {
		return true
	}
	return false
}

// IsGitURLNoSuffix validates that an input is a git URL
// without the ".git" suffix
func IsGitURLNoSuffix(input string) bool {
	url := reGitURLNoSuffix.MatchString(input)
	if url {
		return true
	}
	return false
}
