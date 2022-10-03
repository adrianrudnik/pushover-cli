# pushover-cli

- [Installation](#installation)
  - [Download](#download)
  - [Snap](#snap)
  - [Build it from sources](#build-it-from-sources)
- [Configuration](#configuration)
  - [Environment variables](#environment-variables)
  - [Configuration files](#configuration-files)
- [Usage](#usage)
  - [Getting help](#getting-help)
  - [Push messages](#push-messages)
  - [Rate limits](#rate-limits)
  - [Sending text via images](#sending-text-via-images)
- [Related documentations](#related-documentations)

Unofficial CLI to send messages with [pushover.net](https://pushover.net/) on Windows, Linux and macOS.

## Installation

The following installation channels are available:

### Download

Download a package directly from the [release page](https://github.com/adrianrudnik/pushover-cli/releases).

### Snap

```bash
snap install pushover-cli
```

There are some  small limitations to using file configurations with snaps.:

- You can still use `pushover-cli config setup` to create a configuration file, but you cannot use any other location as a fallback, as this snap is operating in `strict` mode.
- Attaching files is not easily done as the snap does not have access to paths inside home.

If unsure, I would recommend to use environment variables instead when operating outside the confined snap home.

### Build it from sources

Download or clone the sources and build it with:

```bash
go build -mod=vendor .
./pushover --version
```

## Configuration

There are two ways to pass on your pushover credentials: Environment variables or a configuration file.

### Environment variables

The setup for environment variables is pretty straightforward, set the following ones, and they will be picked up by this tool.

| Variable          | Description |
| ----------------- | ----------- |
| PUSHOVER_CLI_USER | User key    |
| PUSHOVER_CLI_API  | API token   |

### Configuration files

A configuration file can be easily be configured.

There is a small wizard to help you create a valid configuration file in your current user profile folder.

```bash
pushover-cli config setup

# > Enter user key: someuser
# > Enter API token: somekey
# > 2020-06-07T12:36:48+02:00 INF Config saved path=/home/adrian/.config/pushover-cli/config.json
```

The `config.json` will be picked up in different locations, to get a list of all recognized, see all available use the following command:

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

To get an overview of all possible flags, see `pushover-cli help send` or read through the overview here:

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

# > 2020-06-07T17:28:51+02:00 INF Rate limit information requests-per-month=7500 requests-remaining=7461 reset-at=2020-07-01T05:00:00Z
```

### Sending text via images

Push messages are limited to 1024 characters. A way to transport a longer or preformatted text is to send it as an image.

This example is based on ImageMagick on Linux, but should be easily adaptable to other operating systems.

```bash
First, ensure that you have ImageMagick installed via `convert --version` and that your desired files are allowed to convert by checking a policy *like* this exists in your `/etc/ImageMagick-{version}/policy.xml`:

```xml
<policy domain="path" rights="all" pattern="@*.log" />
```

Please understand that this might be considered a security risk, as anyone could use the same technique to do simple text to image conversions, that then might leave the system.

After that we can try to convert a text-snipped into an image:

```bash
convert \
  -size 4000x4000 xc:white \
  -font "FreeMono" \
  -pointsize 12 \
  -fill black \
  -annotate +15+30 "@your.log" \
  -trim \
  -bordercolor "#FFF" \
  -border 10 \
  +repage \
  result.png

pushover-cli push --attachment result.png "daily report"
```

Please be graceful with the maximum size of the image, the given example of `4000x4000` is just the maximum that could be reached, the final `+repage` will reduce the image size to the actual content while still respecting the maximum size.

## Related documentations

https://pushover.net/api
