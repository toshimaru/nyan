[![Release](https://img.shields.io/github/release/toshimaru/nyan.svg)](https://github.com/toshimaru/nyan/releases/latest)
[![Actions Status](https://github.com/toshimaru/nyan/workflows/Go/badge.svg)](https://github.com/toshimaru/nyan/actions)
[![Actions Status](https://github.com/toshimaru/nyan/workflows/Release/badge.svg)](https://github.com/toshimaru/nyan/actions)
[![Maintainability](https://api.codeclimate.com/v1/badges/f5063da42c2e2b00e625/maintainability)](https://codeclimate.com/github/toshimaru/nyan/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/f5063da42c2e2b00e625/test_coverage)](https://codeclimate.com/github/toshimaru/nyan/test_coverage)

# nyan

Colored `cat` command which supports syntax highlighting.

![Screen Capture](https://user-images.githubusercontent.com/803398/63024853-00b18b80-bee3-11e9-853a-eea7e790a575.png)

## Installation

### Homebrew

```console
$ brew install toshimaru/nyan/nyan
```

### go get

```console
$ go get github.com/toshimaru/nyan
```

## Usage

```
$ nyan FILE
```

## Available Themes

- abap
- dracula
- emacs
- monokai (default)
- monokailight
- pygments
- solarized-dark
- solarized-light
- swapoff
- vim

You can see available color themes with the command:

```
$ nyan --list-themes
```

![Available Themes](https://user-images.githubusercontent.com/803398/67260792-42a91000-f4d8-11e9-9b92-19c0072987e3.png)

## What is nyan?

`nyan` originates from [nyan-cat](http://www.nyan.cat/).

![nyancat](https://giphygifs.s3.amazonaws.com/media/sIIhZliB2McAo/giphy.gif)
