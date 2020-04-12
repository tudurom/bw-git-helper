# bw-git-helper

A git credential helper for using BitWarden as a backend for credential storage.
It supports only credential retrieval. Passwords are matched based on a mapping
defined in the config file.

## Dependencies

* go 1.13+
* [BitWarden CLI](https://github.com/bitwarden/cli)

## Installation

``` sh
$ go get github.com/tudurom/bw-git-helper
```

## Usage

Create the file `~/.config/bw-git-helper/config.ini`. Each section has a host
pattern as its name and a signle property named `target` that specified the
BitWarden vault entry (either by UUID, or by a string to search for).

Example:

``` ini
[github.com*]
target=GitHub

[*.fooo-bar.*]
target=5d80865a-dc01-4dc1-b376-7a00087d6214
```

This helper asks the user for their password through the means of `pinentry`.
You can choose the pinentry implementation by adding a special `[config]`
section:

``` ini
[config]
pinentry=pinentry-gnome3
```
