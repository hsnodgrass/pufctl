// Package semver provides a struct object representation of a semantic version, as well
// as several functions for working with both SemVer structs and semantic version
// strings.
package semver

import (
	"fmt"
	"regexp"
	"strconv"
)

// Pattern is a string representation of the regular expression to capture SemVers
const Pattern string = `^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`

// VPattern is a string representation of the regular expressions to capture SemVers with a leading "v"
const VPattern string = `^[vV]?(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`

// MajorGroup is the string name of the capture group for the SemVer major version
const MajorGroup = "major"

// MinorGroup is the string name of the capture group for the SemVer minor version
const MinorGroup = "minor"

// PatchGroup is the string name of the capture group for the SemVer patch version
const PatchGroup = "patch"

// PreReleaseGroup is the string name of the capture group for the SemVer prerelease version
const PreReleaseGroup = "prerelease"

// BuildMetadataGroup is the string name of the capture group for the SemVer build metadata version
const BuildMetadataGroup = "buildmetadata"

var (
	// Regexp is a pointer to a compiled Regexp semver expression
	Regexp = regexp.MustCompile(Pattern)

	// VRegexp is a pointer to a compiled Regexp semver expression that accounts for an optional leading "v"
	VRegexp = regexp.MustCompile(VPattern)
)

// SemVer is a struct representation of a Semantic Version string
type SemVer struct {
	Major         int
	Minor         int
	Patch         int
	PreRelease    string
	BuildMetadata string
	ParsedAsV     bool
}

func (v SemVer) String() string {
	str := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.PreRelease != "" {
		str = fmt.Sprintf("%s-%s", str, v.PreRelease)
	}
	if v.BuildMetadata != "" {
		str = fmt.Sprintf("%s+%s", str, v.BuildMetadata)
	}
	return str
}

// VString returns a string representation of the SemVer with a
// leading "v".
func (v SemVer) VString() string {
	str := fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.PreRelease != "" {
		str = fmt.Sprintf("%s-%s", str, v.PreRelease)
	}
	if v.BuildMetadata != "" {
		str = fmt.Sprintf("%s+%s", str, v.BuildMetadata)
	}
	return str
}

// BumpMajor increments the major version by 1. It also
// sets the minor and patch versions to 0 and clears the
// prelease and buildmetadata versions
func (v *SemVer) BumpMajor() {
	v.Major++
	v.Minor = 0
	v.Patch = 0
	v.PreRelease = ""
	v.BuildMetadata = ""
}

// BumpMinor increments the minor version by 1. It also
// sets the patch version to 0 and clears the prerelease
// and buildmetadata versions.
func (v *SemVer) BumpMinor() {
	v.Minor++
	v.Patch = 0
	v.PreRelease = ""
	v.BuildMetadata = ""
}

// BumpPatch increments the patch version by 1. It also
// clears the prerelease and buildmetadata versions.
func (v *SemVer) BumpPatch() {
	v.Patch++
	v.PreRelease = ""
	v.BuildMetadata = ""
}

// Make returns a SemVer object from a semantic version string
func Make(input string) (SemVer, error) {
	var match map[string]string
	var parsedAsV bool
	match = MapNamedCaptureGroups(Regexp, input)
	if len(match) == 0 {
		match = MapNamedCaptureGroups(VRegexp, input)
		parsedAsV = true
	} else {
		parsedAsV = false
	}
	if len(match) == 0 {
		return SemVer{}, fmt.Errorf("Could not parse input string %s into a SemVer", input)
	}
	major, err := strconv.Atoi(match["major"])
	if err != nil {
		return SemVer{}, fmt.Errorf("Failed to convert mapped major version to int with error: %w", err)
	}
	minor, err := strconv.Atoi(match["minor"])
	if err != nil {
		return SemVer{}, fmt.Errorf("Failed to convert mapped minor version to int with error: %w", err)
	}
	patch, err := strconv.Atoi(match["patch"])
	if err != nil {
		return SemVer{}, fmt.Errorf("Failed to convert mapped patch version to int with error: %w", err)
	}
	prerelease := ""
	buildmetadata := ""
	if val, ok := match["prerelease"]; ok {
		prerelease = val
	}
	if val, ok := match["buildmetadata"]; ok {
		buildmetadata = val
	}
	sv := SemVer{
		Major:         major,
		Minor:         minor,
		Patch:         patch,
		PreRelease:    prerelease,
		BuildMetadata: buildmetadata,
		ParsedAsV:     parsedAsV,
	}
	return sv, nil
}

// Major returns the Major (X.y.z) version of a string SemVer
func Major(input string) (string, error) {
	group := MajorGroup
	var match string
	match, err := GetCaptureGroupMatch(Regexp, group, input)
	if err == nil {
		return match, err
	}
	return GetCaptureGroupMatch(VRegexp, group, input)
}

// Minor returns the Minor (x.Y.z) version of a string SemVer
func Minor(input string) (string, error) {
	group := MinorGroup
	var match string
	match, err := GetCaptureGroupMatch(Regexp, group, input)
	if err == nil {
		return match, err
	}
	return GetCaptureGroupMatch(VRegexp, group, input)
}

// Patch returns the Patch (x.y.Z) version of a string SemVer
func Patch(input string) (string, error) {
	group := PatchGroup
	var match string
	match, err := GetCaptureGroupMatch(Regexp, group, input)
	if err == nil {
		return match, err
	}
	return GetCaptureGroupMatch(VRegexp, group, input)
}

// PreRelease returns the PreRelease (x.y.z-PreRelease) version of a string SemVer
func PreRelease(input string) (string, error) {
	group := PreReleaseGroup
	var match string
	match, err := GetCaptureGroupMatch(Regexp, group, input)
	if err == nil {
		return match, err
	}
	return GetCaptureGroupMatch(VRegexp, group, input)
}

// BuildMetadata returns the BuildMetadata (x.y.z-prereleaseBuildMetadata) version of a string SemVer
func BuildMetadata(input string) (string, error) {
	group := BuildMetadataGroup
	var match string
	match, err := GetCaptureGroupMatch(Regexp, group, input)
	if err == nil {
		return match, err
	}
	return GetCaptureGroupMatch(VRegexp, group, input)
}

// GetCaptureGroupMatch returns the match in string input of a single named capture group.
// A regexp.Regexp pointer to use for the matching is the first argument.
func GetCaptureGroupMatch(re *regexp.Regexp, name, input string) (string, error) {
	match := re.FindStringSubmatch(input)
	if len(match) > 0 {
		for i, n := range re.SubexpNames() {
			if i != 0 && n == name {
				return match[i], nil
			}
		}
	}
	return "", fmt.Errorf("Could not find a match for group %s in text %s", name, input)
}

// MapNamedCaptureGroups returns a map of named regex capture groups where map[<group name>] = match
// Param re: A compiled regular expression
// Param input: A string to evaluate with the regular expression
func MapNamedCaptureGroups(re *regexp.Regexp, input string) map[string]string {
	match := re.FindStringSubmatch(input)
	if len(match) > 0 {
		result := make(map[string]string)
		for i, name := range re.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = match[i]
			}
		}
		return result
	}
	return map[string]string{}
}
