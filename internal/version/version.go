package version

import "fmt"

var (
	// Ver is the semantic version of this package
	Ver string
	// Time is the build time of this package
	Time string
	//OS is which operating system this package was built for
	OS string
	//Arch is the architecture this package was built for
	Arch string
)

// Version returns the semver, os, arch, and build time of this package
func Version() string {
	return fmt.Sprintf("%s-%s-%s-%s", Ver, OS, Arch, Time)
}

// UserAgent returns a string used in User-Agent HTTP request headers
func UserAgent() string {
	if Ver == "" {
		return fmt.Sprintf("pufctl/0.0.0")
	}
	return fmt.Sprintf("pufctl/%s", Ver)
}

// License returns the License statement
func License() string {
	return fmt.Sprintln(licenseShorthand)
}
