<p align="center">
	<img alt="Eclectica" src="https://cdn.rawgit.com/markelog/eclectica/81e8e049f825f6ceaf7a2a8e402e16cf59799bb1/assets/logo.png" width="300">
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


## Install

- [go get](#go-get)
- [npm](#npm)
- [pip](#pip)
- [gem](#gem)
- [cargo](#cargo)
- [curl](#curl)
<!-- - [wget](#wget) -->

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

<!-- ## wget

```sh
wget -qO- https://raw.githubusercontent.com/markelog/ec-install/master/scripts/install.sh | sh
```

Default installation folder is `/usr/local/bin`, so you might need to execute `sh` with `sudo` like this –

```sh
wget -qO- https://raw.githubusercontent.com/markelog/ec-install/master/scripts/install.sh | sudo sh
```

if you need to install it to your `$HOME` for example, do this

```sh
wget -qO- https://raw.githubusercontent.com/markelog/ec-install/master/scripts/install.sh | EC_DEST=~/bin sh
``` -->
