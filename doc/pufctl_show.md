## pufctl show

show prints a sorted and organized Puppetfile to screen

### Synopsis


The pufctl show command prints a sorted Puppetfile to screen. As pufctl is
an opinionated tool, pufctl show gives you a look at how all pufctl commands
will organize your Puppetfile, should you decide to save any results.


```
pufctl show [flags]
```

### Options

```
  -h, --help            help for show
      --versions-only   Show a truncated for of Puppetfile with only <module name>: <version>
```

### Options inherited from parent commands

```
      --config string              path to config file (default "/home/sharkdeth/.pufctl.yaml")
  -y, --confirm                    skip all confirmation checks
      --forge-api string           Puppet Forge API URL (default "https://forgeapi-cdn.puppet.com")
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

* [pufctl](pufctl.md)	 - pufctl is a multitool for Puppetfiles

###### Auto generated by spf13/cobra on 24-Aug-2020
