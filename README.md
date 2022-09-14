[![Release](https://img.shields.io/github/release/ddrugeon/terrafactor.svg)](https://github.com/ddrugeon/terrafactor/releases/latest)
![Go Build & Test](https://github.com/ddrugeon/terrafactor/workflows/Go%20Build%20&%20Test/badge.svg)
![Release with goreleaser](https://github.com/ddrugeon/terrafactor/workflows/Release%20with%20goreleaser/badge.svg)
![Issues](https://img.shields.io/github/issues/ddrugeon/terrafactor)
![License](https://img.shields.io/github/license/ddrugeon/terrafactor)

# terrafactor

A simple cli to generate move instruction for terraform when refactoring files.

## Installation

### go get

```console
$ go get github.com/ddrugeon/terrafactor
```

### go install (requires Go 1.16+)

```console
$ go install github.com/ddrugeon/terrafactor@latest
```

## Usage

```console
$ terrafactor COMMAND [FLAGS]
```
### Available Commands

| Command  | Description                             |
|----------|-----------------------------------------|
| help     | Help about any command                  |
| resources| command related to terraform resources  |
| version  | Print the version number of terrafactor |

### Available Subcommands
```console
$ terrafactor resources SUBCOMMAND [FLAGS]
```

| Command  | Description                           |
|----------|---------------------------------------|
| list     | list resources found in given tfstate |
| refactor | generate terraform moved directives   |


### Available Options

| Option                  | Description                                                                                    |
|-------------------------|------------------------------------------------------------------------------------------------|
| `-h`, `--help`          | Show help                                                                                      |
| `-f`, `--filter` string | (optional) Filter string to apply - Example: resource.datadog_synthetics_private_location.main |
| `-t`, `--tfstate` path  | (required) Path of the terraform state in json                                                 |


