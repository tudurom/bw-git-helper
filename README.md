# bw-git-helper

A git credential helper for using BitWarden as a backend for credential storage.
It supports only credential retrieval. Passwords are matched based on a mapping
defined in the config file.

## Dependencies

* go 1.13+
* [BitWarden CLI](https://github.com/bitwarden/cli)

## Installation

``` sh
# install to $GOPATH/bin. Defaults to ~/go/bin.
$ go get github.com/tudurom/bw-git-helper
```

## Usage

First, you need to tell Git to use this credential helper:

``` sh
$ git config --global credential.helper '!bw-git-helper $@'
```

If you want to match entries based not only on the host, but also on the path,
set `credential.useHttpPath` to `true`:

``` sh
$ git config --global credential.useHttpPath true
```

Create the file `~/.config/bw-git-helper/config.ini`. Each section has a host
pattern as its name and a single property named `target` that specified the
BitWarden vault entry (either by UUID, or by a string to search for).

Example:

``` ini
[github.com*]
target=GitHub

[*.fooo-bar.*]
target=5d80865a-dc01-4dc1-b376-7a00087d6214
```

If you enable `useHttpPath`, you can use it for example to write mappings for
different accounts:

``` ini
[github.com/user1/*]
target=GitHub user1

[github.com/user2/*]
target=GitHub user2

; GitHub catch-all
[github.com/*]
target=My GitHub

[gitlab.com/*]
target=GitLab
```

This helper asks the user for their password through the means of `pinentry`.
You can choose the pinentry implementation by adding a special `[config]`
section:

``` ini
[config]
pinentry=pinentry-gnome3
```

If you don't want to have `bw-git-helper` handle your username, you can disable
that functionality:
``` ini
[config]
use_username=false
```

You can also use the `-c` flag to load the config file from an arbitrary path:

``` sh
git config --global credential.helper 'bw-git-helper -c ~/some/config.ini $@'
```
