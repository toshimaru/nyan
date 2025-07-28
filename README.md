[![release](https://img.shields.io/github/release/toshimaru/nyan.svg)](https://github.com/toshimaru/nyan/releases/latest)
[![Go Build & Test](https://github.com/toshimaru/nyan/actions/workflows/ci.yml/badge.svg)](https://github.com/toshimaru/nyan/actions/workflows/ci.yml)
[![Release](https://github.com/toshimaru/nyan/actions/workflows/release.yml/badge.svg)](https://github.com/toshimaru/nyan/actions/workflows/release.yml)
[![Maintainability](https://qlty.sh/gh/toshimaru/projects/nyan/maintainability.svg)](https://qlty.sh/gh/toshimaru/projects/nyan)
[![Code Coverage](https://qlty.sh/gh/toshimaru/projects/nyan/coverage.svg)](https://qlty.sh/gh/toshimaru/projects/nyan)

# nyan

Colorizing `cat` command with syntax highlighting.

![OG image for nyan command](https://repository-images.githubusercontent.com/195893425/0a7e7dfc-3a80-49d5-8193-5482fe2e7848)

## Installation

### Homebrew

```console
$ brew install nyan
```

<details>
<summary>Homebrew Tap</summary>

```console
$ brew install --cask toshimaru/nyan/nyan
```

</details>

### go get

```console
$ go get github.com/toshimaru/nyan
```

### go install (requires Go 1.16+)

```console
$ go install github.com/toshimaru/nyan@latest
```

## Usage

```console
$ nyan FILE
```

### Available Options

| Option | Description |
| --- | --- |
| `-h`, `--help` | Show help |
| `-l`, `--language` lang | Specify language for syntax highlighting |
| `-T`, `--list-themes` | List available color themes |
| `-n`, `--number` | Output with line numbers |
| `-t`, `--theme` theme | Set color theme for syntax highlighting |

## Available Color Themes

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

You can list and preview available color themes with the command:

```console
$ nyan --list-themes
```

![Available Themes](https://user-images.githubusercontent.com/803398/67260792-42a91000-f4d8-11e9-9b92-19c0072987e3.png)

## What is nyan?

`nyan` originates from [Nyan Cat](https://www.nyan.cat/) (Music by [daniwell](https://aidn.jp/about/)).

![nyancat](https://giphygifs.s3.amazonaws.com/media/sIIhZliB2McAo/giphy.gif)
