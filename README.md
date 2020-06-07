# pushover-cli

Unofficial CLI to send messages with pushover.net on Windows, Linux and MacOS.

## Installation

You can use one of hte following methods:

### Download

Download a package directly from the [release page](https://github.com/adrianrudnik/pushover-cli/releases).

### Snap

```bash
snap install pushover-cli
```

### Build it from sources

Download or clone the sources and build it with:

```bash
go build -mod=vendor -o pushover .
./pushover --version
```

## Configuration

### Environment variables

There are two ways to pass on your pushover credentials: Environment variables or a configuration file.

The setup for environment varialbes is pretty straigt forward, set the following ones and they will be picked up by this tool.

| Variable          | Description |
| ----------------- | ----------- |
| PUSHOVER_CLI_USER | User key    |
| PUSHOVER_CLI_API  | API token   |

### Configuration files

You can also setup configuration files.

There is a small wizard to help you create a valid configuration file nearest to the next configuration folder for the current user profile:

```bash
pushover-cli config setup

# > Enter user key: someuser
# > Enter API token: somekey
# > 2020-06-07T12:36:48+02:00 INF Config saved path=/home/adrian/.config/pushover-cli/config.json
```

The `config.json` can be placed into different locations, to see all available use the following command:

```bash
pushover-cli config paths

# > 2020-06-07T12:39:13+02:00 INF Collecting paths that will be used for config.json lookup
# > 2020-06-07T12:39:13+02:00 INF Folder found path=/home/adrian/.config/pushover-cli
# > 2020-06-07T12:39:13+02:00 INF Folder found path=/etc/xdg/xdg-plasma/pushover-cli
```

To remove the nearest configuration use:

```bash
pushover-cli config clear

# > 2020-06-07T12:43:03+02:00 INF Config cleared file=/home/adrian/.config/pushover-cli/config.json
```

## Usage

### Getting help

There are several ways to get help to specific commands:

```bash
pushover-cli help
pushover-cli --help
pushover-cli config --help
pushover-cli help config
```

### Push messages

To get an overview of all possible flags you can use `pushover-cli help send` or read through the overview here:

```
Flags:
  -e, --api-endpoint string   API endpoint for message submission (default "https://api.pushover.net/1/messages.json")
  -a, --attachment string     path to image attachment (max size 2.5mb)
  -d, --devices strings       devices to limit the push to (comma-separated)
  -h, --help                  help for push
      --link-label string     title for the supplementary URL (max. 100 characters)
      --link-url string       supplementary URL (max. 512 characters)
  -p, --priority string       message priority [none, quiet, normal, high, confirm] (default "normal")
  -s, --sound string          playback sound [see https://pushover.net/api#sounds] (default "pushover")
      --timestamp int         message date and time override as unix timestamp
  -t, --title string          message title (max. 250 characters)

Global Flags:
  -v, --verbose   print debug information
```

All flags are optional, so a simple

```bash
pushover-cli push "hello"

# > 2020-06-07T13:59:47+02:00 INF Message pushed request=a9ee72c0-1e76-476a-bfa5-d421dcd6acca status=1
# > 2020-06-07T13:58:53+02:00 INF Rate limit information requests-per-month=7500 requests-remaining=7468 reset-at=2020-07-01T05:00:00Z
```

will work. Here is a more complex example on how to use it:

```bash
pushover-cli \
    --devices=mobile,workpc \
    --link-url https://localhost/report.html \
    -p quiet \
    -s cashregister \
    "Generated report available now"
```

Sending attachments can be done with:

```bash
pushover-cli --attachment image.jpg "some message"
```

### Rate limits

The following command can be used to request the current rate limits:

```bash
pushover-cli limits

# > 2020-06-07T13:59:47+02:00 INF Message pushed request=a9ee72c0-1e76-476a-bfa5-d421dcd6acca status=1
```

# Related documentations

https://pushover.net/api
