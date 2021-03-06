## pufctl

pufctl is a multitool for Puppetfiles

### Synopsis


A tool for doing everything you need to with Puppetfiles right from the command line.

```
pufctl [flags]
```

### Options

```
      --config string              path to config file (default "/home/sharkdeth/.pufctl.yaml")
  -y, --confirm                    skip all confirmation checks
      --forge-api string           Puppet Forge API URL (default "https://forgeapi-cdn.puppet.com")
  -h, --help                       help for pufctl
  -L, --license                    show the license shorthand statement
  -o, --out-file string            Write command output or changed Puppetfile to specified file
      --pass string                Forge / Git authentication password
  -p, --puppetfile string          path to the Puppetfile to parse (default "./Puppetfile")
      --puppetfile-branch string   The branch to use for a Puppetfile from Git (default "production")
  -s, --show                       Show Puppetfile after each command
      --ssh-key string             Path to your SSH key (default "/home/sharkdeth/.ssh/id_rsa")
      --token string               Forge / Git authentication token
      --user string                Forge / Git authentication username
  -v, --verbose                    verbose logging
```

### SEE ALSO

* [pufctl add](pufctl_add.md)	 - add new content to the Puppetfile
* [pufctl bump](pufctl_bump.md)	 - bump the semantic version of a module
* [pufctl completion](pufctl_completion.md)	 - completion generates a completion script
* [pufctl confgen](pufctl_confgen.md)	 - confgen generates a default config file
* [pufctl diff](pufctl_diff.md)	 - diff finds the difference between two Puppetfiles
* [pufctl docgen](pufctl_docgen.md)	 - docgen generates pufctl documentation
* [pufctl edit](pufctl_edit.md)	 - edit objects in a Puppetfile
* [pufctl search](pufctl_search.md)	 - search operations for all things Puppetfile related
* [pufctl show](pufctl_show.md)	 - show prints a sorted and organized Puppetfile to screen

###### Auto generated by spf13/cobra on 24-Aug-2020
