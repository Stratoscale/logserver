# logserver

Logserver is a web log viewer that combines logs from several sources.

## Usage

### Config File

The config file has the following specifications:

- `sources` (list of source dicts): Logs sources, from which the logs are merged ans served.
- `parsers` (list of parser dicts): Which parsers to apply to the log files.

### Source Dict

- `name` (string): Name of source, the name that this source will be shown as
- `url` (URL string): URL of source, see [supported sources](#Supported URL schemes) schemes below.
- `open_tar_files` (bool): Wether to treat tar files as directories, used for logs that are packed
                           into a tar file.

#### Supported URL Schemes

- `file://` (URL string): Address in local file system.
- `sftp://` (URL string): Address of sftp server (or SSH server).
- `ssh://` (URL string): Address of ssh server.

### Parser Dict

- `glob` (string): File pattern to apply this parser on.
- `time_formats` (list of strings): Parse timestamp string according to those time formats.
                                    The given format should be in Go style time formats, or
                                    `unix_int` or `unix_float`.
- `json_mapping` (dict): Parse each log line as a json, and map keys from that json to the UI expected keys.
                         The keys are values that the UI expect, the values are keys from the file json.
- `regexp` (Go style regular expression string): Parse each line in the long with this regular expression.
                                                 the given regular expression should have named groups with
                                                 the keys that the UI expects.

#### UI Keys

The UI expects the following keys in each parsed log:

- `msg`: Log message.
- `time`: Time stamp of log.
- `level`: Log level.
- `args`: If args are given, they will be injected into the log msg. Args value can be `[]interface{}`
          Or `map[string]interface{}`, According to the log message.