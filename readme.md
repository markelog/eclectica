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
Eclectica unifies management of any language under one cohesive and minimalistic interface.

Like [pyenv](https://github.com/pyenv/pyenv) for Python,
[rbenv](https://github.com/rbenv/rbenv) for Ruby, [nvm](https://github.com/creationix/nvm) Node.js and etc. Managing multiple languages and doing it in a little more enjoyable fashion

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
