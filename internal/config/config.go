package config

import (
	"fmt"
	"os/user"

	"github.com/hsnodgrass/pufctl/internal/logging"
)

// ConfigFile is the default config file name (not file extension)
const ConfigFile string = ".pufctl"

// ConfigFileExtension is the default config file extension
const ConfigFileExtension string = "yaml"

// AltConfPath is the default alternate config path
const AltConfPath string = "."

// EnvVarPrefix is the default prefix for config environment variables
const EnvVarPrefix string = "PUF"

// Puppetfile is the default Puppetfile path
const Puppetfile string = "./Puppetfile"

// PuppetfileBranch is the branch to use with the configured Puppetfile
// is you are using a Puppetfile from a Git source.
const PuppetfileBranch string = "production"

// AlwaysVerbose is the default setting for the verbose global flag
const AlwaysVerbose bool = false

// AlwaysConfirm is the default setting for the auto-confirm flag
const AlwaysConfirm bool = false

// ForgeAPI is the default URL for the forge-url flag
const ForgeAPI string = "https://forgeapi-cdn.puppet.com"

// AlwaysWriteInPlace is the default setting for the write-in-place flag
const AlwaysWriteInPlace bool = false

// AlwaysShow is the default setting for the show flag
const AlwaysShow bool = false

// AlwaysPreferGit is the default setting for the prefer-git flag
const AlwaysPreferGit bool = false

// DocGenPath is the default directory path where all docs will be generated
const DocGenPath string = "./doc/"

var (
	userName, homeDir string = GetUsernameAndHomedir()

	// ConfPath is the default path to directory that houses the config file
	ConfPath string = homeDir

	// FullConfPath is the full path to (and including) the config file
	FullConfPath string = fmt.Sprintf("%s/%s.%s", homeDir, ConfigFile, ConfigFileExtension)

	// SSHKeyPath is the default path for an SSH key
	SSHKeyPath string = fmt.Sprintf("%s/.ssh/id_rsa", homeDir)

	// ForgeStrDefaults is a map of default values under the "forge" config key
	// used in setting Viper defaults. The key user_agent is set at runtime.
	ForgeStrDefaults = map[string]string{
		"user_agent": "",
		"api_url":    ForgeAPI,
	}

	//AlwaysBoolDefaults is a map of default values under the "always" config key
	// used in setting Viper defaults.
	AlwaysBoolDefaults = map[string]bool{
		"verbose":        AlwaysVerbose,
		"show":           AlwaysShow,
		"prefer_git":     AlwaysPreferGit,
		"write_in_place": AlwaysWriteInPlace,
	}

	// AuthStrDefaults is a map of default values under the "auth" config key
	// used in setting Viper defaults.
	AuthStrDefaults = map[string]string{
		"ssh_key":  SSHKeyPath,
		"username": userName,
		"password": "",
		"token":    "",
	}

	// GenoptsStrDefaults is a map of default values under the "genopts" config key
	// used in setting Viper defaults.
	GenoptsStrDefaults = map[string]string{
		"doc_gen_path": DocGenPath,
	}
)

// GetUsernameAndHomedir returns the current user's username and home directory or,
// if an error occurs in determining the current user, these anonymous defaults:
// username: ""
// homedir: "."
func GetUsernameAndHomedir() (string, string) {
	_user, err := user.Current()
	if err != nil {
		logging.Warnln("Could not identify the current user, using anonymous defaults")
		return "", "."
	}
	return _user.Name, _user.HomeDir
}
