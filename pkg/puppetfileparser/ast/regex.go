// Package ast provides the abstract syntax tree, parser, and everything else in puppetfileparser
package ast

import "regexp"

// ReKeyword is a compiled regular expression for the "mod" keyword
var ReKeyword = regexp.MustCompile(`^mod`)

// ReString is a compiled regular expression for Puppetfile strings
var ReString = regexp.MustCompile(`'([^']*)'`)

// ReIdent is a compiled regular expression for Puppetfile symbols
var ReIdent = regexp.MustCompile(`:[A-Za-z0-9_]+`)

// ReComment is a compiled regular expression for Puppetfile comments
var ReComment = regexp.MustCompile(`#.*`)

// ReForgeURL is a compiled regular expression for extracting a Forge
// URL from a forge line in a Puppetfile.
var ReForgeURL = regexp.MustCompile(`forge ["']([^"']*)["']`)

// ReMetaTag is a compiled regular expression for Puppetfile metadata
var ReMetaTag = regexp.MustCompile(`#\s?@(?P<Tag>[\w_-]+):.*`)

// ReMetaData is a compiled regular expression for Puppetfile meta data
var ReMetaData = regexp.MustCompile(`#\s?@[\w_-]+:\s*(?P<Data>.*)`)

// ReMeta is a compiled regular expression for Puppetfile meta tags and data.
// ReMeta matches should return two named capture groups, Tag and Data.
var ReMeta = regexp.MustCompile(`#\s?@(?P<Tag>[\w_-]+):\s*(?P<Data>.*)`)

//ReAssignment is a compiled regular expression for Module properties with assignements.
// ReAssignment matches should return two named groups, Key and Value.
var ReAssignment = regexp.MustCompile(`(?P<Key>:[A-Za-z0-9_]+)\s?=?>?\s?'?(?P<Value>[^\s=>']*)'?`)

// ReValidIdent is a compiled regular expression used to validate Idents in the Puppetfile.
var ReValidIdent = regexp.MustCompile(`(:latest|:git|:install_path|:ref|:branch|:tag|:commit|:default_branch|:install_path)`)

// ReGitData is a compiled regular expression for parsing Modules' :git property value.
// ReGitData matches should return a named group Proto. Proto is the protocol used by
// git. If Proto is "git" or "ssh://git", this means SSH is to be used. The full match is
// the URI. ReGitData does not validate URLs.
var ReGitData = regexp.MustCompile(`(?P<Proto>https?|(?:ssh:\/\/)?git)(?::\/\/.*|@.*)`)

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

// SingleMatch returns the first matching capture group. Useful for when you only
// want to return one value from a match.
// Param re: A compiled regular expression
// Param inputL A string to evaluate with the regular expression
func SingleMatch(re *regexp.Regexp, input string) string {
	match := re.FindStringSubmatch(input)
	if len(match) > 0 {
		return match[1]
	}
	return ""
}
