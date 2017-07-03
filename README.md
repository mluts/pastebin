# pastebin
[![Build Status](https://travis-ci.org/prologic/pastebin.svg)](https://travis-ci.org/prologic/pastebin)
[![GoDoc](https://godoc.org/github.com/prologic/pastebin?status.svg)](https://godoc.org/github.com/prologic/pastebin)
[![Wiki](https://img.shields.io/badge/docs-wiki-blue.svg)](https://github.com/prologic/pastebin/wiki)
[![Go Report Card](https://goreportcard.com/badge/github.com/prologic/pastebin)](https://goreportcard.com/report/github.com/prologic/pastebin)
[![Coverage](https://coveralls.io/repos/prologic/pastebin/badge.svg)](https://coveralls.io/r/prologic/pastebin)

pastebin is a web app that allows you to create smart bookmarks, commands and aliases by pointing your web browser's default search engine at a running instance. Similar to bunny1 or yubnub.

## Installation

### Source

```#!bash
$ go install github.com/prologic/pastebin/...
```

### OS X Homebrew

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
$ pastebin -bind 127.0.0.1:8000
```

Set your browser's default pastebin engine to http://localhost:8000/?q=%s

Then type `help` to view the main help page, `g foo bar` to perform a [Google](https://google.com) search for "foo bar" or `list` to list all available commands.

## License

MIT
