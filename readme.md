# Eclectica [![Build Status](https://travis-ci.org/markelog/eclectica.svg?branch=master)](https://travis-ci.org/markelog/eclectica)

## Install

- [go get](#go-get)
- [npm](#npm)
- [gem](#gem)
- [cargo](#cargo)
- [curl](#curl)
- [wget](#wget)

## go get

```sh
go get github.com/markelog/eclectica/ec
```

## npm

```sh
[sudo] npm install -g eclectica
```

## gem

```sh
gem install eclectica
```

## cargo

```sh
cargo install eclectica
```

## curl

```sh
curl -s https://raw.githubusercontent.com/markelog/ec-install/master/install.sh | sh
```

Default installation folder is `/usr/local/bin`, so you might need to execute `sh` with `sudo` like this –

```sh
curl -s https://raw.githubusercontent.com/markelog/ec-install/master/install.sh | sudo sh
```

if you need to install it to your `$HOME` for example, do this

```sh
curl -s https://raw.githubusercontent.com/markelog/ec-install/master/install.sh | EC_DEST=~/bin sh
```

## wget

```sh
wget -qO- https://raw.githubusercontent.com/markelog/ec-install/master/install.sh | sh
```

Default installation folder is `/usr/local/bin`, so you might need to execute `sh` with `sudo` like this –

```sh
wget -qO- https://raw.githubusercontent.com/markelog/ec-install/master/install.sh | sudo sh
```

if you need to install it to your `$HOME` for example, do this

```sh
wget -qO- https://raw.githubusercontent.com/markelog/ec-install/master/install.sh | EC_DEST=~/bin sh
```

