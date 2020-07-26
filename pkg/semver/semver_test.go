package semver

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

const simpleRegexpStr = `(?P<first>[a-zA-Z]*)\.(?P<second>[a-zA-Z]*)`
const basicSemverStr = "1.2.3"
const vBasicSemverStr = "v1.2.3"
const preSemverStr = "1.2.3-beta"
const vPreSemverStr = "v1.2.3-beta"
const metaSemverStr = "1.2.3-beta+exp.sha.5114f85"
const vMetaSemverStr = "v1.2.3-beta+exp.sha.5114f85"
const invalidLeadingZeros = "001.002.003-beta+exp.sha.5114f85"
const invalidNoY = "1-beta+exp.sha.5114f85"
const invalidNoZ = "1.2-beta+exp.sha.5114f85"
const invalidAlphaX = "a.2.3-beta+exp.sha.5114f85"
const invalidAlphaY = "1.a.3-beta+exp.sha.5114f85"
const invalidAlphaZ = "1.2.a-beta+exp.sha.5114f85"
const invalidCharPR = "1.2.3-****+exp.sha.5114f85"

var matchTestCases = []string{"test.case", "TEST.CASE"}
var mapTestResultOne = map[string]string{"first": "test", "second": "case"}
var mapTestResultTwo = map[string]string{"first": "TEST", "second": "CASE"}
var strTestCases = []string{basicSemverStr, vBasicSemverStr, preSemverStr, vPreSemverStr, metaSemverStr, vMetaSemverStr}
var strPreTestCases = []string{preSemverStr, vPreSemverStr, metaSemverStr, vMetaSemverStr}
var strMetaTestCases = []string{metaSemverStr, vMetaSemverStr}
var invalidTestCases = []string{invalidLeadingZeros, invalidNoY, invalidNoZ, invalidAlphaX, invalidAlphaY, invalidAlphaZ, invalidCharPR}

func newTestStruct() SemVer {
	return SemVer{
		Major:         1,
		Minor:         2,
		Patch:         3,
		PreRelease:    "beta",
		BuildMetadata: "exp.sha.5114f85",
	}
}

func TestGetCaptureGroupMatch(t *testing.T) {
	simpleRegexp := regexp.MustCompile(simpleRegexpStr)
	for _, s := range matchTestCases {
		matchone, err := GetCaptureGroupMatch(simpleRegexp, "first", s)
		if err != nil {
			t.Errorf("Failed to find capture group first with error: %#w", err)
			t.FailNow()
		} else {
			if strings.ToLower(matchone) != "test" {
				t.Errorf("Failed to properly parse match for group first. Expected: test or TEST, Got: %s", matchone)
				t.FailNow()
			}
		}
		matchtwo, err := GetCaptureGroupMatch(simpleRegexp, "second", s)
		if err != nil {
			t.Errorf("Failed to find capture group second with error: %#w", err)
			t.FailNow()
		} else {
			if strings.ToLower(matchtwo) != "case" {
				t.Errorf("Failed to properly parse match for group second. Expected: case or CASE, Got: %s", matchtwo)
				t.FailNow()
			}
		}
	}
}

func TestMapNamedCaptureGroups(t *testing.T) {
	simpleRegexp := regexp.MustCompile(simpleRegexpStr)
	mapOne := MapNamedCaptureGroups(simpleRegexp, matchTestCases[0])
	if !reflect.DeepEqual(mapOne, mapTestResultOne) {
		t.Fatalf("Failed to map capture groups for test case one")
	}
	mapTwo := MapNamedCaptureGroups(simpleRegexp, matchTestCases[1])
	if !reflect.DeepEqual(mapTwo, mapTestResultTwo) {
		t.Fatalf("Failed to map capture groups for test case two")
	}
}

func TestSemVerString(t *testing.T) {
	testSemVerStruct := newTestStruct()
	svStr := fmt.Sprintf("%s", testSemVerStruct)
	if svStr != metaSemverStr {
		t.Errorf("Failed to format SemVer as string")
	}
}

func TestSemVerVString(t *testing.T) {
	testSemVerStruct := newTestStruct()
	svStr := testSemVerStruct.VString()
	if svStr != vMetaSemverStr {
		t.Errorf("Failed to format SemVer as v-string")
	}
}

func TestMake(t *testing.T) {
	for _, v := range strTestCases {
		_, err := Make(v)
		if err != nil {
			t.Errorf("Make failed test case %s with error: %w", v, err)
		}
	}
	for _, v := range invalidTestCases {
		_, err := Make(v)
		if err == nil {
			t.Errorf("Make parsed the invalid semver string %s", v)
		}
	}
}

func TestBumpMajor(t *testing.T) {
	testSemVerStruct := newTestStruct()
	testSemVerStruct.BumpMajor()
	if testSemVerStruct.Major != 2 {
		t.Errorf("Failed to increase Major version number by exactly 1. Major version: %d", testSemVerStruct.Major)
	}
	if testSemVerStruct.Minor != 0 || testSemVerStruct.Patch != 0 {
		t.Errorf("Failed to zero Minor or Patch")
	}
	if testSemVerStruct.PreRelease != "" || testSemVerStruct.BuildMetadata != "" {
		t.Errorf("Failed to wipe PreRelease or BuildMetadata")
	}
}

func TestBumpMinor(t *testing.T) {
	testSemVerStruct := newTestStruct()
	testSemVerStruct.BumpMinor()
	if testSemVerStruct.Major != 1 {
		t.Errorf("Major version differs from expected. Expected: 1, Got: %d", testSemVerStruct.Major)
	}
	if testSemVerStruct.Minor != 3 {
		t.Errorf("Failed to increase Minor version number by exactly 1. Minor version: %d", testSemVerStruct.Minor)
	}
	if testSemVerStruct.Patch != 0 {
		t.Errorf("Failed to zero Patch")
	}
	if testSemVerStruct.PreRelease != "" || testSemVerStruct.BuildMetadata != "" {
		t.Errorf("Failed to wipe PreRelease or BuildMetadata")
	}
}

func TestBumpPatch(t *testing.T) {
	testSemVerStruct := newTestStruct()
	testSemVerStruct.BumpPatch()
	if testSemVerStruct.Major != 1 {
		t.Errorf("Major version differs from expected. Expected: 1, Got: %d", testSemVerStruct.Major)
	}
	if testSemVerStruct.Minor != 2 {
		t.Errorf("Minor version differs from expected. Expected: 1, Got: %d", testSemVerStruct.Minor)
	}
	if testSemVerStruct.Patch != 4 {
		t.Errorf("Failed to increase Patch version number by exactly 1. Patch version: %d", testSemVerStruct.Patch)
	}
	if testSemVerStruct.PreRelease != "" || testSemVerStruct.BuildMetadata != "" {
		t.Errorf("Failed to wipe PreRelease or BuildMetadata")
	}
}

func TestMajor(t *testing.T) {
	for _, v := range strTestCases {
		ver, err := Major(v)
		if err != nil {
			t.Errorf("Function failed. Input: %s, error: %w", v, err)
		}
		if ver != "1" {
			t.Errorf("Unexpected function return. Expected: 1, Got: %s", ver)
		}
	}
}

func TestMinor(t *testing.T) {
	for _, v := range strTestCases {
		ver, err := Minor(v)
		if err != nil {
			t.Errorf("Function failed. Input: %s, error: %w", v, err)
		}
		if ver != "2" {
			t.Errorf("Unexpected function return. Expected: 0, Got: %s", ver)
		}
	}
}

func TestPatch(t *testing.T) {
	for _, v := range strTestCases {
		ver, err := Patch(v)
		if err != nil {
			t.Errorf("Function failed. Input: %s, error: %w", v, err)
		}
		if ver != "3" {
			t.Errorf("Unexpected function return. Expected: 0, Got: %s", ver)
		}
	}
}

func TestPreRelease(t *testing.T) {
	for _, v := range strPreTestCases {
		ver, err := PreRelease(v)
		if err != nil {
			t.Errorf("Function failed. Input: %s, error: %w", v, err)
		}
		if ver != "beta" {
			t.Errorf("Unexpected function return. Expected: beta, Got: %s", ver)
		}
	}
}

func TestBuildMetadata(t *testing.T) {
	for _, v := range strMetaTestCases {
		ver, err := BuildMetadata(v)
		if err != nil {
			t.Errorf("Function failed. Input: %s, error: %w", v, err)
		}
		if ver != "exp.sha.5114f85" {
			t.Errorf("Unexpected function return. Expected: exp.sha.5114f85, Got: %s", ver)
		}
	}
}
