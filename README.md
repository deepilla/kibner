# Kibner

Kibner is a command-line utility for managing podcasts.

NOTE: Kibner is a work in progress. Commands and flags are subject to change.

## Install

    go install github.com/deepilla/kibner

## Usage

### Common Tasks

#### Subscribe to a feed:

    kibner add [options] <url>

##### Options:

**--itunes**<br/>
Indicate that `url` is an iTunes page rather than an RSS feed.

#### Unsubscribe from a feed:

    kibner remove <feed>

#### Sync subscribed feeds with their RSS feeds:

    kibner sync [feed]

#### List items:

    kibner list [options] [feed]

##### Options:

**--sortby=<pubdate|title|feed|duration|timestamp>**<br/>
Sort items by the given property. The default is pubdate.

**--order=<asc|desc>**<br/>
Sort items in ascending or descending order.

**-u**<br/>
**--unplayed**<br/>
Only show unplayed items.

**-T=<date>**<br/>
**--since=<date>**<br/>
Only show items published on or after the given date. Valid dates are:
- *today* for the current date
- *N*d for *N* days ago
- *N*m for *N* months ago
- *N*y for *N* years ago
- *mon* to *sun* for a particular day of the week
- *jan* to *dec* for a particular month

**--with-title=<title>**<br/>
Only show items with titles that match the given value.

**-N=<limit>**<br/>
**--top=<limit>**<br/>
Set the maximum number of items to display.

**-d**<br/>
**--show-desc**<br/>
Show item descriptions.

**-p**<br/>
**--play**<br/>
Play selected items with the program specified in the `use` flag.

**--mark**<br/>
Mark selected items as played.

**--unmark**<br/>
Mark selected items as unplayed.

**--run**<br/>
Run the program specified in the `use` flag on the selected items.
Unlike the `play` flag, this does not mark items as played.

**--use=<program>**<br/>
Specify a program to use with the `play` or `run` flags.

### More Tasks

`kibner feeds [options]`

`kibner update <options> <feed>`

`kibner import [options] <filename>`

`kibner export [options] <filename>`

`kibner version`

`kibner reset`

## TODO

### Features

#### Provide values for command-line flags in config file.

#### Subscribe to BBC iPlayer audio.

#### Improve terminal output.

#### Windows support.

#### Allow user-defined templates for list/feed commands.

#### Play/download feed content (e.g. audio, video).

#### Pause/Resume, Mute/Unmute feeds.

### Code

#### Refactor database code.

#### Refactor command/flag/config code (use Viper?).

#### Increase test coverage.

#### Handle interrupt signals (e.g. Ctrl-C).

## Licensing

Kibner is provided under an [MIT License](http://choosealicense.com/licenses/mit/). See the [LICENSE](LICENSE) file for details.
