# aidoku-cli
Aidoku development tools in a single program

# Usage
```sh
Aidoku development toolkit

Usage:
  aidoku [command]

Available Commands:
  build       Build a source list from packages
  help        Help about any command
  logcat      Log streaming
  serve       Build a source list and serve it on the local network
  version     Print version

Flags:
  -h, --help      help for aidoku
  -v, --verbose   verbose output

Use "aidoku [command] --help" for more information about a command.
```

# Commands
## `aidoku init [rust-template|rust|as|c] [DIR]`
```sh
Create initial code for an Aidoku source

Usage:
  aidoku init [rust-template|rust|as|c] [DIR] [flags]

Flags:
  -h, --help              help for init
  -p, --homepage string   Source homepage
  -l, --language string   Source language
  -n, --name string       Source name
      --nsfw int          Source NSFW level (default -1)
      --version           version for init

Global Flags:
  -v, --verbose   verbose output
```

## `aidoku build <FILES>`
```sh
Build a source list from packages

Usage:
  aidoku build <FILES> [flags]

Flags:
  -h, --help            help for build
  -o, --output string   Output folder (default "public")

Global Flags:
  -v, --verbose   verbose output
```

## `aidoku serve <FILES>`
```sh
Build a source list and serve it on the local network

Usage:
  aidoku serve <FILES> [flags]

Flags:
  -h, --help            help for serve
  -o, --output string   The source list folder (default "public")
  -p, --port string     The port to broadcast the source list on (default "8080")

Global Flags:
  -v, --verbose   verbose output
```

## `aidoku logcat`
```sh
Log streaming

Usage:
  aidoku logcat [flags]

Flags:
  -h, --help          help for logcat
  -p, --port string   The port to listen to logs on (default "9000")

Global Flags:
  -v, --verbose   verbose output
```
