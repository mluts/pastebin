# pastebin
[![GoDoc](https://godoc.org/github.com/prologic/pastebin?status.svg)](https://godoc.org/github.com/prologic/pastebin)
[![Go Report Card](https://goreportcard.com/badge/github.com/prologic/pastebin)](https://goreportcard.com/report/github.com/prologic/pastebin)

pastebin is a self-hosted pastebin web app that lets you create and share
"ephemeral" data between devices and users. There is a configurable expiry
(TTL) afterwhich the paste expires and is purged. There is also a handy
CLI for interacting with the service in a easy way or you can also use curl!

### Source

```#!bash
$ go install github.com/prologic/pastebin/...
```

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

There is also an included command line utility for convenience:

```#!bash
echo hello | pb
```

## Configuration

When running the `pastebin` server there are a few default options you might
want to tweak:

```
$ ./pastebin --help
  ...
  -expiry duration
        expiry time for pastes (default 5m0s)
  -fqdn string
        FQDN for public access (default "localhost")
```

Setting a custom `-expiry` lets you change when pastes are automatically
expired (*the purge time is 2x this value*). The ``-fqdn` option is used as
a namespace for generating the UUID(s) for pastes, change this to be your
domain name.

The command-line utility by default talk to http://localhost:8000 which can be
changed via the `-url` option or by creating a `$HOME/.pastebin.conf`
configuration file with contents similar to:

```
$ cat ~/.pastebin.conf
url=https://paste.mydomain.com/
```
## License


MIT
