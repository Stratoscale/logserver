# logserver

[![docker](https://img.shields.io/docker/build/stratoscale/logserver.svg)](https://hub.docker.com/r/stratoscale/logserver)

Logserver is a web log viewer that combines logs from several sources.

![diagram](./doc/diagram.png)

### Screenshots

![screenshot general view](./doc/screenshot-1.png)

![screenshot log view](./doc/screenshot-2.png)

![screenshot filter files](./doc/screenshot-3.png)

## Usage

### Install

* Get the binary with go get:

`go get -u github.com/Stratoscale/logserver`

### Run a Docker Container

Assuming:

* A logserver config is in `/etc/logserver/logserver.json`.
* Logserver should listen to port 80

The command line:

```
docker run -d --restart always --name logserver --net host \
    -v /etc/logserver/logserver.json:/logserver:json \
    stratoscale/logserver -addr :80
```

### Configuration

Logserver is configured with a json configuration file. See [example](./example/logserver.json).

The json should be a dict with the following keys:

- `sources` (list of [source dicts](./README.md#source-dict)): Logs sources, from which the logs are merged ans served.
- `parsers` (list of [parser dicts](./README.md#parser-dict)): Which parsers to apply to the log files.
- `global` (dict of [attributes](./README.md#global-dict)): General configuration
- `cache` (dict of [attributes](./README.md#cache-dict)): Cache configuration
- `route` (dict of [attributes](./README.md#route-dict)): Route configuration

#### Source Dict

- `name` (string): Name of source, the name that this source will be shown as
- `url` (URL string with [supported schemes](./README.md#supported-url-schemes)): URL of source.
- `open_tar` (bool): Weather to treat tar files as directories, used for logs that are packed
                           into a tar file.
- `open_journal` (string): Open a journalctl directory as a log file. The value
                           should be the journalctl directory from the source root.

#### Supported URL Schemes

Logserver supports different type of log sources, each has a different scheme:

- `file://` (URL string): Address in local file system.
    The file location can be absolute with `file:///var/log` or relative to the directory
    from which the command line was executed: `file://./log`.
- `sftp://` (URL string): Address of sftp server (or SSH server).
    The URL can contain a user name and password and path from the system root.
    For example: `sftp://user:password@example.com:22/var/log`
- `ssh://` (URL string): Address of ssh server. Obey the same rules as sftp server.
- `nginx+http://`, `nginx+https://` (URL string): Address of an nginx configured to serve files with `autoindex on;`
    directive. It supports both HTML and JSON `autoindex_format`.

#### Parser Dict

Logserver can parse each log line according to a list of configured parsers.
Parsers can handle json logs (each line is a json dict), or regular logs with regular expression rules.

Each parser can be defined with the following keys:

- `glob` (string): File pattern to apply this parser on.
- `time_formats` (list of strings): Parse timestamp string according to those time formats.
                                    The given format should be in Go style time formats, or
                                    `unix_int` or `unix_float`.
- `json_mapping` (dict): Parse each log line as a json, and map keys from that json to the
                         [UI expected keys](./README.md#ui-keys). The keys are values that the
                         UI expect, the values are keys from the file json.
- `regexp` (Go style regular expression string): Parse each line in the long with this regular expression.
                                                 the given regular expression should have named groups with
                                                 the keys that the UI expects.
- `append_args` (bool): (for json log) Add to msg all remaining json keys in format: key=value.

#### UI Keys

The UI expects the following keys in each parsed log:

- `msg`: Log message.
- `time`: Time stamp of log.
- `level`: Log level.
- `args`: If args are given, they will be injected into the log msg. Args value can be `[]interface{}`
          Or `map[string]interface{}`, According to the log message.

#### Global Dict

- `content_batch_size`
- `content_batch_time`
- `search_max_size`

#### Cache Dict

- `size`
- `expiration`

#### Route Dict

- `base_path`
- `root_path`
