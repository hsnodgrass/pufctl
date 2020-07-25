package config

import "fmt"

var (
	// SemVer is the semver of this package
	SemVer string
	// Build is the build time of this package
	Build string
	//OS is which operating system this package was built for
	OS string
	//Arch is the architecture this package was built for
	Arch string
)

// Version returns the semver, os, arch, and build time of this package
func Version() string {
	return fmt.Sprintf("pufctl %s, OS: %s, Arch:%s, Build Time: %s\n", SemVer, OS, Arch, Build)
}
