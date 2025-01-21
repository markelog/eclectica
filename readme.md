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
<br/><br/>

Eclectica unifies management of any language under one cohesive and minimalistic interface. Like [pyenv](https://github.com/pyenv/pyenv) for Python,
[rbenv](https://github.com/rbenv/rbenv) for Ruby, [nvm](https://github.com/creationix/nvm) Node.js and etc.

But instead of having all of those, you have only one binary

## Usage

After you [install](#install) eclectica, `ec` program will be available in your terminal, I used to have a nice site with fancy animation explaning how to used it, but help output will do too -

```
$ ec --help

Usage:
  ec [command] [flags] [<language>@<version>]

Examples:
  Install specifc, say, node version
  $ ec node@6.4.0

  Or choose from already installed Go versions
  $ ec go

  Same way to choose, plus install available Rust versions
  $ ec -r rust

Available Commands:
  completion        generate the autocompletion script for the specified shell
  install           same as "ec [<language>@<version>]"
  ls                list installed language versions
  remove-everything removes everything related to eclectica
  rm                remove language version
  version           print version of eclectica

Flags:
  -h, --help   help for ec

Use "ec [command] --help" for more information about a command

```

## Install

Since eclectica is language manager for any language, it should be installed through any package manager :-)

- [go get](#go-get)
- [npm](#npm)
- [pip](#pip)
- [gem](#gem)
- [cargo](#cargo)
- [curl](#curl)
- [wget](#wget)

## go get

```sh
go install github.com/markelog/eclectica/bin/{ec,ec-proxy}@latest
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
wget -qO - https://raw.githubusercontent.com/markelog/ec-install/master/scripts/wget-install.sh | sudo sh
```

if you need to install it to your `$HOME` for example, do this

```sh
wget -qO - https://raw.githubusercontent.com/markelog/ec-install/master/scripts/wget-install.sh | EC_DEST=~/bin sh
```
