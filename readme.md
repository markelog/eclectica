<p align="center">
	<img alt="Eclectica" src="./assets/logo.svg" width="300">
</p>

<p align="center">
  Cool and eclectic version manager for any language
</p>

<p align="center">
  <a href="https://travis-ci.org/markelog/eclectica">
		<img alt="Build Status" src="https://travis-ci.org/markelog/eclectica.svg?branch=master">
	</a><a href="https://godoc.org/github.com/markelog/eclectica">
		<img alt="GoDoc" src="https://godoc.org/github.com/markelog/eclectica?status.svg">
	</a><a href="https://goreportcard.com/report/github.com/markelog/eclectica">
		<img alt="Go Report" src="https://goreportcard.com/badge/github.com/markelog/eclectica">
	</a>
</p>

Eclectica unifies management of any language under one cohesive and minimalistic interface.

Like [pyenv](https://github.com/pyenv/pyenv) for Python,
[rbenv](https://github.com/rbenv/rbenv) for Ruby, [nvm](https://github.com/creationix/nvm) Node.js and etc. Managing multiple languages and doing it in a little more enjoyable fashion

```
Usage:
  ec [command] [flags] [<language>@<version>]

Examples:
  Install specifc version
  $ ec node@6.4.0

  Choose local version with interactive list
  $ ec go

  Choose remote version with interactive list
  $ ec -r rust

Available Commands:
  ls                list installed language versions
  remove-everything removes everything related to the eclectica
  rm                remove language version
  version           print version of Eclectica

Flags:
  -h, --help           help for ec
  -l, --local          install to the current folder only
  -r, --remote         ask for remote versions
  -w, --with-modules   reinstall global modules from the previous version (currently works only for node.js)

Use "ec [command] --help" for more information about a command
```
# Install

- [go get](#go-get)
- [npm](#npm)
- [pip](#pip)
- [gem](#gem)
- [cargo](#cargo)
- [curl](#curl)
- [wget](#wget)

## go get

```sh
go get github.com/markelog/eclectica/bin/{ec,ec-proxy}
```

## npm

```sh
[sudo] npm install -g eclectica
```

## pip

```sh
sudo -H pip install -v eclectica
```

## gem

```sh
sudo gem install eclectica
```

## cargo

```sh
cargo install eclectica
```

## curl

```sh
curl -s https://raw.githubusercontent.com/markelog/ec-install/master/scripts/install.sh | sh
```

Default installation folder is `/usr/local/bin`, so you might need to execute `sh` with `sudo` like this –

```sh
curl -s https://raw.githubusercontent.com/markelog/ec-install/master/scripts/install.sh | sudo sh
```

if you need to install it to your `$HOME` for example, do this

```sh
curl -s https://raw.githubusercontent.com/markelog/ec-install/master/scripts/install.sh | EC_DEST=~/bin sh
```

## wget

```sh
wget -qO - https://raw.githubusercontent.com/markelog/ec-install/master/scripts/install.sh | sh
```

Default installation folder is `/usr/local/bin`, so you might need to execute `sh` with `sudo` like this –

```sh
wget -qO - https://raw.githubusercontent.com/markelog/ec-install/master/scripts/install.sh | sudo sh
```

if you need to install it to your `$HOME` for example, do this

```sh
wget -qO - https://raw.githubusercontent.com/markelog/ec-install/master/scripts/install.sh | EC_DEST=~/bin sh
```
