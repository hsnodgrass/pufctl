package uitext

// StarSep provides a 75 character star separator for terminal output
const StarSep = "*****************************************************************"

// DashSep provides a 75 character dash separator for terminal output
const DashSep = "-----------------------------------------------------------------"

// ConfirmNewWritePrompt provides a basic text prompt asking for confirmation of writing new file to disk
const ConfirmNewWritePrompt = "Are you sure you want to write the new file to disk? (y/n)"

// ConfirmOverwritePrompt provides a basic text prompt asking for confirmation of overwriting a files on disk
const ConfirmOverwritePrompt = "Are you sure you want to overwrite the file with new content? This cannot be reversed. (y/n)"

// RootUse is the usage description for the pufctl command
const RootUse = "pufctl"

// RootShort is the short description for the pufctl command
const RootShort = "pufctl is a multitool for Puppetfiles"

// RootLong is the long description for the pufctl command
const RootLong = `
A tool for doing everything you need to with Puppetfiles right from the command line.`

// ShowUse is the usage description for the pufctl show command
const ShowUse = "show"

// ShowShort is the short description for the pufctl show command
const ShowShort = "show prints a sorted and organized Puppetfile to screen"

// ShowLong is the long description for the pufctl show command
const ShowLong = `
The pufctl show command prints a sorted Puppetfile to screen. As pufctl is
an opinionated tool, pufctl show gives you a look at how all pufctl commands
will organize your Puppetfile, should you decide to save any results.
`

// AddUse is the usage description for the pufctl add command
const AddUse = "add [type]"

// AddShort is the short description for the pufctl add command
const AddShort = "add new content to the Puppetfile"

// AddLong is the long description for the pufctl add command
const AddLong = `
The pufctl add command adds new content such as modules and metadata to the Puppetfile`

// AddModuleUse is the usage description for the pufctl add module command
const AddModuleUse = "module [modulename]"

// AddModuleShort is the short description for the pufctl add module command
const AddModuleShort = "module adds a new module to the Puppetfile"

// AddModuleLong is the long description for the pufctl add module command
const AddModuleLong = `
The pufctl add module command adds a new module to the Puppetfile. By just
specifying a name, the module will be added with the format "mod '<name>'".

In order to resolve the modules dependencies, the module name must be in slug
format (<namespace>-<modulename>) if it is a Forge module. You can also specify
a URL to a git repository to add non-Forge modules. 

Flags can be passed to add properties and metadata to the module entry.
`

// AddMetaUse is the usage description for the pufctl add meta command
const AddMetaUse = "meta [<modulename>|top|bottom]"

// AddMetaShort is the short description for the pufctl add meta command
const AddMetaShort = "meta adds new metadata to a module or the top/bottom comment block"

// AddMetaLong is the long description for the pufctl add meta command
const AddMetaLong = `
The pufctl add meta command allows you to add metadata to a Puppetfile programatically.
This metadata can be a regular comment or a meta tag with optional data.

To add a module comment of metadata to a module, the first argument should be the module
name as it appears in the Puppetfile, such as "puppetlabs-apache".

To add a top-block comment or metadata to the top block, the first argument should be "top".

To add a bottom-block comment or metadata to the bottom block, the first argument should be "bottom".

Use the --metadata (-m) flag to specify the metadata you want to add. You can specify multiple metadata
statements / comments at a time by putting them in a comma-separated list. Metadata and comments should
be single-quoted.

Examples:

$ pufctl add meta puppetlabs-apache -m '# @maintainer: puppetlabs'

$ pufctl add meta top -m '# Production Puppetfile'
`

// DiffUse is the usage description for the pufctl diff command
const DiffUse = "diff [Puppetfile] [<optional> Puppetfile]"

// DiffShort is the short description for the pufctl diff command
const DiffShort = "diff finds the difference between two Puppetfiles"

// DiffLong is the long description for the pufctl diff command
const DiffLong = `
The pufctl diff command compares Puppetfiles at the parsed object level to 
find differences in them. You must supply the path to at least one Puppetfile.
If only one Puppetfile is specified, the configured default Puppetfile will
be used as the first Puppetfile in the comparison.

Meaningful exit codes have been added to the diff command to assist with
programatic implementations of this command (such as use in CI/CD systems).

Exit Codes:

0: No difference between the Puppetfiles

3: Only the first Puppetfile has differences

4: Only the second Puppetfile has differences

5: Differences in both Puppetfiles
`

// SearchUse is the usage description for the pufctl search command
const SearchUse = "search [subcommand]"

// SearchShort is the short description for the pufctl search command
const SearchShort = "search operations for all things Puppetfile related"

// SearchLong is the long description of the pufctl search command
const SearchLong = `
The pufctl search command allows you to search for all things Puppetfile
realted. Read the descriptions of the subcommands for more information.`

// ForgeSearchUse is the usage description of the pufctl search forge command
const ForgeSearchUse = "forge [query]"

// ForgeSearchShort is the short description of the pufctl search forge command
const ForgeSearchShort = "forge allows you to search for modules on the Puppet Forge"

// ForgeSearchLong is the long description of the pufctl search forge command
const ForgeSearchLong = `
The pufctl search forge command allows you to search for modules on the Puppet Forge.

Search queries are passed as command args, and the search results can be fine-tuned
with flags.
`

// BumpUse is the usage description of the pufctl bump command
const BumpUse = "bump [module]..."

// BumpShort is the short desccription of the pufctl bump command
const BumpShort = "bump the semantic version of a module"

// BumpLong is the long description of the pufctl bump command
const BumpLong = `
The pufctl bump command allows you "bump" (increment by 1) a module's
semantic version, if it has one. You must use valid module name slugs,
<org>-<module name>, as command arguments.

The command will look for a semver in the modules properties,
specifically the :ref and :tag symbols.

The command works with both regular semver strings, as well as semver
strings that have a leading "v".

By default, pufctl bump increments the Patch portion of the module's 
semver (x.y.Z). Which portion of the semver gets bumped can be
changed with flags.
`

// EditUse is the usage description of the pufctl edit command
const EditUse = "edit [subcommand]"

//EditShort is the short description of the pufctl edit command
const EditShort = "edit objects in a Puppetfile"

//EditLong is the long description of the pufctl edit command
const EditLong = `
The pufctl edit command allows you to edit objects in a Puppetfile such as
modules, metadata, and the Forge declaration.

Read the descriptions of the subcommands for more details.`

// EditModuleUse is the usage description of the pufctl edit module command
const EditModuleUse = "module [name]"

// EditModuleShort is the short description of the pufctl edit module command
const EditModuleShort = "edit module name/properties"

// EditModuleLong is the long description of the pufctl edit module command
const EditModuleLong = `
The pufctl edit module command allows you to edit a specific module in the
Puppetfile.

You can use the --name (-n) flag to edit a module's name.

You can use the --key (-k) flag to specify the Module property you would like to edit.
When using the --key (-k) flag, you should specify Module properties as strings
with the form ':key=>value'. The --key (-k) flag accepts multiple values as
a comma-separated string, ex. --key ':key=>value',':key=>value',':key=>value'.

If you attempt to edit a Module property that doesn't exist, that property
will be created.

If you want to edit a "bare" property (i.e. :latest or single version string),
your --key (-k) flag string would look like this: --key 'bare=>value'. Since
bare values are mutually exclusive, a bare value will overwrite all other
properties of the Module.
`

// CompletionUse is the usage description for the pufctl completion command
const CompletionUse = "completion [bash|zsh|powershell]"

// CompletionShort is the short description of the pufctl completion command
const CompletionShort = "completion generates a completion script"

// CompletionLong is the long description of the pufctl completion command
const CompletionLong = `
The completion command generates scripts that provide tab completions for
all Pufctl commands.

To load completions:

Bash:

$ source <(pufctl completion bash)

To load completions for each session, execute once:

Linux:

$ pufctl completion bash > /etc/bash_completion.d/pufctl

MacOS:

$ pufctl completion bash > /usr/local/etc/bash_completion.d/pufctl

Zsh:

If shell completion is not already enabled in your environment you will need
to enable it. You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions for each session, execute once:

$ pufctl completion zsh > "${fpath[1]}/_pufctl"

You will need to start a new shell for this setup to take effect.

PowerShell:

PS > pufctl.exe completion powershell | Out-File pufctl-completion.ps1

Add the following line to your $PROFILE:

. C:\path\to\where\you\saved\pufctl-completion.ps1

You will need to close and reopen PowerShell for it to take effect.
`

// DocGenUse is the usage description of the pufctl docgen command
const DocGenUse = "docgen"

// DocGenShort is the short description of the pufctl docgen command
const DocGenShort = "docgen generates pufctl documentation"

// DocGenLong is the long description of the pufctl docgen command
const DocGenLong = `
Generate pufctl documentation in markdown. The documentation 
will be saved to a directory specified by the docgen_path config option.
`

// ConfGenUse is the usage description of the pufctl confgen command
const ConfGenUse = "confgen"

// ConfGenShort is the short description of the pufctl confgen command
const ConfGenShort = "confgen generates a default config file"

// ConfGenLong is the long description of the pufctl confgen command
const ConfGenLong = `
Generate a default pufctl config file. 

You can use the --out-file flag to specify where the config file will be generated.

If --out-file is not used, the config file will be generated at $HOME/.pufctl.yaml.
`

// NewUse is the usage description of the pufctl new command
const NewUse = "new [subcommand]"

// NewShort is the the short description of the pufctl new command
const NewShort = "new allows you to create new Puppetfiles or .fixtures.yml files"

// NewLong is the long description of the pufctl new command
const NewLong = `
Create a new Puppetfile or .fixtures.yml file.

NOT CURRENTLY IMPLEMENTED
`

// NewFixturesUse is the usage description of the pufctl new fixtures command
const NewFixturesUse = "fixtures [nameslug|URL]"

// NewFixturesShort is the short description of the pufctl new fixtures command
const NewFixturesShort = "fixtures allows you to create a new .fixtures.yml file"

// NewFixturesLong is the long description of the pufctl new fixtures command
const NewFixturesLong = `
Create a new .fixtures.yml file from a module name slug (<org>-<module>) or a Git URL.

If you pass in a Git URL, the module must have a metadata.json file so that the dependencies
can be properly parsed.

NOT CURRENTLY IMPLEMENTED
`
