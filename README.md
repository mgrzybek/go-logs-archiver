# go-logs-archiver

A tool to consume json-based messages from a data source, sort and send them to a persistant storage.

## Architecture

![Architecture](./docs/Architecture.svg)

## Usage

```bash
$ ./go-logs-archiver --help
Reads the incoming messages from the configured consumer driver
and send them to the backend.

For example: get messages from a kafka topic and send them to a S3 storage.

Usage:
  go-logs-archiver [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  run         Start processing
  validate    Validate the given configuration

Flags:
      --config string      config file (default is $HOME/.go-logs-archiver.yaml)
  -h, --help               help for go-logs-archiver
  -l, --log-level string   Level of verbosity (development, production).
  -t, --toggle             Help message for toggle

Use "go-logs-archiver [command] --help" for more information about a command.
```
