# Pufctl

An opinionated, easy to use multitool for Puppetfiles. Pufctl aims to extract the difficulty of working with
Puppetfiles and, also, provide some helpful programatic tooling for Puppetfiles for command-line warriors
and CI/CD systems alike.

Pufctl allows you to quickly add modules to a Puppetfile from a git source or the Puppet Forge (along with their dependencies, if desired), get a comprehensive diff of two Puppetfiles, sort and organize your Puppetfile, and more.

**Pull Requests and Issues are encouraged and appreciated!**

## Quickstart

* Download the appropriate zip file for your operating system [here](https://github.com/hsnodgrass/pufctl/releases)
* Unzip the file (`unzip` on *nix, 7zip or WinRar on Windows)
    * (optional) Validate the checksum of the binary (`pufctl` or `pufctl.exe`) with `sha256sum`. The correct checksum is located in the file `sha256checksum.txt`
* Move the binary to somewhere in your `$PATH`
* Run `pufctl --help` (`pufctl.exe --help` on Windows)

## Table of Contents

* [Reference](#reference)
    * [Command Overview](#command-overview)
    * [Working With Puppetfiles](#working-with-puppetfiles)
    * [Git and Authentication](#git-and-authentication)
    * [Using a Config File](#using-a-config-file)
    * [Proxy Settings](#proxy-settings)
* [Features](#features)
* [Usage](#usage)
    * [Basic Examples](#basic-examples)
* [Roadmap](#roadmap)
* [Development](#development)
    * [Building From Source](#building-from-source)
* [License](#license)

## Reference

### Command Overview

See [the docs](https://github.com/hsnodgrass/pufctl/blob/main/doc/pufctl.md) for full documentation.

* `pufctl help` - Shows the help for any command. You can also use the `--help` (`-h`) flag.
* `pufctl add` - Adds new content to the specified Puppetfile
  * `pufctl add meta` - Add metadata (comments, etc.) to the Puppetfile
  * `pufctl add module` - Add new module statements to the Puppetfile. Use with the `-D` flag to add the module's dependencies as well.
* `pufctl bump` - "Bump" (increment by one) a module's semver in the Puppetfile. Flags determine which part of the semver is bumped.
* `pufctl completion` - Generate completion script for Pufctl. These can be used with your profile to provide tab completion for Pufctl. Supports `bash`, `zsh`, and `powershell`).
* `pufctl confgen` - Generate a default config file for Pufctl.
* `pufctl diff` - Diff two Puppetfiles at the object level.
* `pufctl docgen` - Generate markdown documentation for Pufctl.
* `pufctl edit module` - Edit a module's properties in the Puppetfile.
* `pufctl search forge` - Search the Puppet Forge for modules with a simple string query.
* `pufctl show` - Prints a sorted and organized version of your Puppetfile to screen

### Working with Puppetfiles

Most Pufctl commands require you to specify a Puppetfile to work with using the flag `--puppetfile`(`-p`). You can pass either a path to a Puppetfile on disk, or you can pass a Git repository.

When you pass a Git repository, you can specify a branch to get the Puppetfile from as well using the `--puppetfile_branch` flag. Pufctl uses the `production` branch by default.

Example:

```sh
# Using a path on disk
pufctl show -p /path/to/Puppetfile

# Using a git repo
pufctl show -p git@github.com:fakeorg/control-repo.git
```

### Git and Authentication

There are several Pufctl commands that use git and, most of the time, these commands will require authentication.

#### SSH

Using SSH is fully supported by these commands, and Pufctl defaults to using an SSH key located at `$HOME/.ssh/id_rsa`. This path can be configured using the `--ssh-key` global flag.

If you use an `ssh-agent` on Linux or MacOS, Pufctl will automatically use that connection. If you don't and your ssh key is password-protected, Pufctl will prompt you for your password at the appropriate time.

#### HTTP(S)

Pufctl also supports HTTP(S) authentication. You can pass your username to a command using the `--user` global flag and your password using the `--pass` global flag. 

If you have two-factor authentication enabled, in Github at least, you will need to use an access token instead of a password. Your access token can be passed using the `--token` global flag.

### Using a Config File

Nearly every option provided by a command line flag can be configured using a config file. Pufctl can generate a default config file for you using the `pufctl confgen` command.

The default location for the config file is `$HOME/.pufctl.yaml`. You can use a config file at a different path by passing the path to the `--config` global flag.

### Proxy Settings

If you are behind a proxy, so can configure Pufctl to use the proxy by settings the following environment variables:

* `HTTP_PROXY`
* `HTTPS_PROXY`
* `USE_PROXIES` (optional)

## Features

* [Automatic Puppetfile Sorting](#automatic-puppetfile-sorting)
* [Object Level Diffs of Puppetfiles](#object-level-diffs-of-puppetfiles)
* [Organized Comments](#organized-comments)
* [Puppetfile Metadata](#puppetfile-metadata)

### Automatic Puppetfile Sorting

By default, Pufctl sorts your Puppetfile alphabetically by module name. 

Before using Pufctl:
```ruby
mod 'zanyorg-module1',
    :git => 'https://github.com/zanyorg/module1.git',
    :branch => 'production'

mod 'puppetlabs-apache',
    :latest

```

After using Pufctl:
```ruby
mod 'puppetlabs-apache',
    :latest

mod 'zanyorg-module1',
    :git => 'https://github.com/zanyorg/module1.git',
    :branch => 'production'

```

### Object Level Diffs of Puppetfiles

Have you ever tried to diff two Puppetfiles and realized that it's kind of a pain?
If the modules are different or the Puppetfiles are oranized differently, the diff
can be incredibly hard to parse. Fortunately, Pufctl is a full Puppetfile parser
and this allows you to diff Puppetfiles by their abstract syntax trees! In simpler
terms, you only see the relevant stuff that's different instead of comments, whitespace,
etc.

Please see the [diff command docs](doc/pufctl_diff.md) for more information.

### Organized Comments

Pufctl introduces the concept of organized comments. What does this mean? It means that depending
on where you put comments carries significance now!

Pufctl has three "types" of organized comments:

#### Top-Block Comments

Comments at the top of a Puppetfile will stay at the top of a Puppetfile, always.
To separate top-block comment from module comments, ensure there is a blank line inbetween your
top-block comments and the first module comment (or bare module).

```ruby
# My Puppetfile at /path/to/control-repo/Puppetfile
# The line above, as well as this line, are considered top-block comments.
# This line is, too. Top-block comments will always be at the top of the
# Puppetfile no matter what you add, remove, or change with pufctl.

# I'm a module comment! We're talking about me next!
mod 'puppetlabs-apache',
    :latest

```

#### Module Comments

Comments directly above modules are considered module comments and will follow their
associated module around no matter how your Puppetfile changes.

```ruby
# My Puppetfile at /path/to/control-repo/Puppetfile
# The line above, as well as this line, are considered top-block comments.

# I'm a module comment! Told you we would talk about me!
# I'm also a module comment. Both me and the line above will
# follow the "puppetlabs-apache" module declaration wherever it may go.
mod 'puppetlabs-apache',
    :latest

```

#### Bottom-Block Comments

Comments at the bottom of the Puppetfile are considered bottom-block comments.
Bottom-block comments will always stay at the bottom of the Puppetfile.

```ruby
# My Puppetfile at /path/to/control-repo/Puppetfile
# The line above, as well as this line, are considered top-block comments.

# I'm a module comment!
# I'm also a module comment.
mod 'puppetlabs-apache',
    :latest

# I'm a bottom-block comment, always bringing up the rear.
```

#### All Other Comments

All other comments in a Puppetfile that don't fall under top-block, module, or bottom-block
comments stay exactly where they were at relative to the file itself. This means that you may
have some weird comment behavior when you first start using Pufctl.

Before using `pufctl`:
```ruby
mod 'zanyorg-module1',
    :git => 'https://github.com/zanyorg/module1.git',
    :branch => 'production'

# Modules from Puppet Forge

mod 'puppetlabs-apache',
    :latest

```

After using `pufctl`:
```ruby
mod 'puppetlabs-apache',
    :latest

# Modules from Puppet Forge

mod 'zanyorg-module1',
    :git => 'https://github.com/zanyorg/module1.git',
    :branch => 'production'

```

This behavior is due to the fact that Pufctl alphabetizes your Puppetfile by module name
by default. This can be a bit of a pain, but, in my opinion, having a uniformly organized
Puppetfile is worth the little bit of initial work fixing things like this. 

### Puppetfile Metadata

Another cool feature of Pufctl is the addition of metadata to your Puppetfile. Metadata
is simple to add:

```ruby
# @tag: data
```

You can place comments like this anywhere to give semantic significance to things in your
Puppetfile.

```ruby
# @maintainer: IaCTeam@zanyorg.com
mod 'zanyorg-module1',
    :git => 'https://github.com/zanyorg/module1.git',
    :branch => 'production'
```

## Usage Overview

One of the goals of Pufctl is to make it easy to use. To see a help message, use the command `pufctl --help`.
For more comprehensive documentation, see the [included docs](doc/pufctl.md).

### Basic Examples

#### Add a new Forge module and it's dependencies to your Puppetfile

```sh
pufctl add -p /path/to/your/Puppetfile 'puppetlabs-apache' -D -w -s
```

The -p flag allows you to specify the Puppetfile you want to use.

The -D flag means resolve dependencies. Don't worry, Pufctl won't
add duplicate modules to your Puppetfile.

The -w flag means write-in-place, or overwrite your current Puppetfile
with the changes. Don't worry, it will prompt you for confirmation.

The -s flag means show, or print your Puppetfile, with changes, to the
terminal.

#### Edit a module entry in a Puppetfile

Before:
```ruby
mod 'fakemod',
  :git => 'https://fake.com/fakeorg/fakeorg-fakemod',
  :tag => 'v1.6.5'
```

Command:
```sh
pufctl edit -p /path/to/your/Puppetfile fakemod --name 'fakeorg-fakemod' --key ':tag=>1.8.5' -w
```

After:
```ruby
mod 'fakeorg-fakemod',
  :git => 'https://fake.com/fakeorg/fakeorg-fakemod',
  :tag => 'v1.8.5'
```

#### Bump the semantic version of a module in a Puppetfile

Before:
```ruby
mod 'fakeorg-fakemod',
  :git => 'https://fake.com/fakeorg/fakeorg-fakemod',
  :tag => 'v1.8.5'
```

Command:
```sh
pufctl bump -p /path/to/your/Puppetfile fakeorg-fakemod -w
```

After:
```ruby
mod 'fakeorg-fakemod',
  :git => 'https://fake.com/fakeorg/fakeorg-fakemod',
  :tag => 'v1.8.6'
```

## Roadmap

Here are some features I'd like to implement in the future, as well as some housekeeping work I'd like to get done:

* Check for updates for all modules in a Puppetfile.
* Search through Puppetfiles themselves.
* Quickly update a module to latest.
* Resolve dependency conflicts. Right now, Pufctl won't add or update a module if it exists in the Puppetfile already.
* Disk-based caching of remote resources.
* Generation of `.fixtures.yaml` files.
* Ruby bindings via C-Go and `ffi`
* More tests!
* Some refactoring to make the public API more sensical.
* Possibly split `forgeapi` and `puppetfileparser` into their own, separate modules.

Feel free to make [feature suggestions](https://github.com/hsnodgrass/pufctl/issues) or even submit PRs!

## Development

### Build From Source

Building Pufctl from source requires you have Go v1.14.1 or higher installed on
your machine.

The default way to build Pufctl is by using the included Makefile. This has only
been tested on Linux (specifically Ubuntu 20.04 LTS using Windows Subsystem for Linux), however it should
also work just fine on MacOS. If you're using Windows, either use WSL like I do or
build it using Go's provided tooling. You  can reference the included Makefile for
the proper commands.

* Clone this repository and `cd` to the `pufctl` dir
* Use the `make` command to compile Pufctl
    * All platforms:
        * `make build`
    * Target platform:
        * `make [nix|dar|win]build`
* Binaries will now be located in `bin/[linux|darwin|windows]/`

### Code Conventions

Try to use idomatic Go when possible. Also, try to write tests for new features. That being said, I'm still learning Go myself and don't always adhear to these rules.

## License

Pufctl is licensed under the Apache 2.0 license. See [LICENSE.txt](LICENSE.txt)