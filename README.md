# pastebin
[![Build Status](https://travis-ci.org/prologic/pastebin.svg)](https://travis-ci.org/prologic/pastebin)
[![GoDoc](https://godoc.org/github.com/prologic/pastebin?status.svg)](https://godoc.org/github.com/prologic/pastebin)
[![Wiki](https://img.shields.io/badge/docs-wiki-blue.svg)](https://github.com/prologic/pastebin/wiki)
[![Go Report Card](https://goreportcard.com/badge/github.com/prologic/pastebin)](https://goreportcard.com/report/github.com/prologic/pastebin)
[![Coverage](https://coveralls.io/repos/prologic/pastebin/badge.svg)](https://coveralls.io/r/prologic/pastebin)

pastebin is a self-hosted pastebin web app that lets you create and share
"ephemeral" data between devices and users. There is a configurable expiry
(TTL) afterwhich the paste expires and is purged. There is also a handy
CLI for interacting with the service in a easy way or you can also use curl!

### Source

```#!bash
$ go install github.com/prologic/pastebin/...
```

### OS X Homebrew

**Coming**

There is a formula provided that you can tap and install from
[prologic/homebrew-pastebin](https://github.com/prologic/homebrew-pastebin):

```#!bash
$ brew tap prologic/pastebin
$ brew install pastebin
```

**NB:** This installs the latest released binary; so if you want a more
recent unreleased version from master you'll have to clone the repository
and build yourself.

pastebin is still early days so contributions, ideas and expertise are
much appreciated and highly welcome!

## Usage

Run pastebin:

```#!bash
$ pastebin
```

Create a paste:

```#!bash
$ echo "Hello World" | pb
http://localhost:8000/92sHUeGPfoFctazBxdEhae
```

Or use the Web UI: http://localhost:8000/

Or curl:

```#bash
$ echo "hello World" | curl -q -L -d @- -o - http://localhost:8000/
...
```

## License

MIT
