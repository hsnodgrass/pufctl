package validators

import "testing"

var (
	testModSlug          = "puppetlabs-apache"
	testModRef           = "puppetlabs/apache"
	testGitHTTP          = "http://github.com/puppetlabs/apache.git"
	testGitHTTPS         = "https://github.com/puppetlabs/apache.git"
	testGitSSH           = "git@github.com:puppetlabs/apache.git"
	testGitSSHPrefixed   = "ssh://git@github.com:puppetlabs/apache.git"
	testMissingURLSuffix = "https://github.com/puppetlabs/apache"
	testSSHWrongSep      = "git@github.com/puppetlabs/apache.git"
)

var (
	testGitURLs        = []string{testGitHTTP, testGitHTTPS, testGitSSH, testGitSSHPrefixed}
	testInvalidGitURLs = []string{testMissingURLSuffix, testSSHWrongSep}
	testConfirm        = []string{"y", "ye", "yes", "Y", "YE", "YES", " y ", " ye ", " yes "}
	testInvalidConfirm = []string{"n", "N", "yeah", "yup", "you", "asjdhfasjdhfkajsdfk"}
)

func TestIsConfirmed(t *testing.T) {
	for _, i := range testConfirm {
		if !IsConfirmed(i) {
			t.Errorf("Failed to validation confirmation string: %s", i)
		}
	}
	for _, i := range testInvalidConfirm {
		if IsConfirmed(i) {
			t.Errorf("Validated invalid confirmation string: %s", i)
		}
	}
}

func TestIsModuleSlug(t *testing.T) {
	if !IsModuleSlug(testModSlug) {
		t.Errorf("Failed to validate module slug: %s", testModSlug)
	}
	if IsModuleSlug(testModRef) {
		t.Errorf("Validated invalid module slug: %s", testModRef)
	}
}

func TestIsModuleRef(t *testing.T) {
	if !IsModuleRef(testModRef) {
		t.Errorf("Failed to validate module ref: %s", testModRef)
	}
	if IsModuleRef(testModSlug) {
		t.Errorf("Validated invalid module ref: %s", testModSlug)
	}
}

func TestIsGitURL(t *testing.T) {
	for _, i := range testGitURLs {
		if !IsGitURL(i) {
			t.Errorf("Failed to validate Git URL: %s", i)
		}
	}
	for _, i := range testInvalidGitURLs {
		if IsGitURL(i) {
			t.Errorf("Falsely validated invalid Git URL: %s", i)
		}
	}
}
