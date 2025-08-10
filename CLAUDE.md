# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`nyan` is a Go CLI tool that provides syntax-highlighted output for files, similar to `cat` but with colorization. It's named after Nyan Cat and uses the Chroma syntax highlighting library.

## Architecture

- **Entry point**: `main.go` - Simple wrapper that calls `cmd.Execute()`
- **CLI logic**: `cmd/root.go` - Contains all command logic using Cobra framework
- **Syntax highlighting**: Uses `github.com/alecthomas/chroma/v2` for tokenization and formatting
- **Themes**: Custom theme registry in `styles/` directory with individual theme files
- **Core functionality**: 
  - File reading with automatic language detection via Chroma's lexers
  - Terminal detection for conditional colorization
  - Line numbering support via custom `numberWriter`
  - Theme switching and listing capabilities

## Development Commands

### Build and Test
```bash
go build -v
go test ./...
```

### Run locally
```bash
go run main.go [flags] FILE
```

### Test output
```bash
# Normal output test
go run main.go -- main.go

# Test with different themes
go run main.go -t dracula main.go
go run main.go -l go main.go
```

## Key Dependencies

- `github.com/alecthomas/chroma/v2` - Syntax highlighting engine
- `github.com/spf13/cobra` - CLI framework
- `github.com/mattn/go-colorable` - Cross-platform terminal color support
- `github.com/mattn/go-isatty` - Terminal detection

## Theme System

The `styles/` package maintains a registry of color themes. Each theme is implemented as a separate Go file that registers itself with the `Registry` map. The `api.go` file provides the theme management interface.

## Testing Strategy

Tests are located alongside source files (`*_test.go`). The CI runs tests on multiple platforms (Ubuntu, macOS, Windows) and includes both unit tests and integration tests that verify actual command output.

## Project Structure & Module Organization

- `main.go`: CLI entry; calls `cmd.Execute()`.
- `cmd/`: Cobra commands, flags, and tests (e.g., `cmd/root.go`, `cmd/root_test.go`).
- `styles/`: Built-in Chroma themes and small API.
- `testdata/`: Fixtures used by tests.
- `.github/workflows/`: CI for build, tests, coverage, release.
- `dist/`: GoReleaser artifacts (ignored in commits).
