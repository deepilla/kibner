# Kibner

Kibner is a command-line utility for managing podcasts.

NOTE: Kibner is a work in progress. Commands and flags are subject to change.

## Install

    go install github.com/deepilla/kibner

## Usage

### Common Tasks

Subscribe to a feed:

    kibner add [options] <url>

Options:

**--itunes**<br/>
Indicate that `url` is an iTunes page rather than an RSS feed.

Unsubscribe from a feed:

    kibner remove <feed>

Sync subscriptions with their RSS feeds:

    kibner sync [feed]

List items:

    kibner list [options] [feed]

Options:

**-u**<br/>
**--unplayed**<br/>
Only show unplayed items.

**-T=<start>**<br/>
**--since=<start>**<br/>
Only show items published on or after the given date. Valid dates are:
- *today* for the current date
- *N*d for *N* days ago
- *N*m for *N* months ago
- *N*y for *N* years ago
- *mon* to *sun* for a particular day of the week
- *jan* to *dec* for a particular month

**--with-title=<title>**<br/>
Only show items with titles that match the given value.

**-d**<br/>
**--show-desc**<br/>
Show item descriptions.

**-N=<number>**<br/>
**--top=<number>**<br/>
Set the maximum number of items to display.

**--sortby=<pubdate|title|feed|duration|timestamp>**<br/>
Sort items by the given property. Default sort is by pubdate -- the most recently published items appear first.

**--order=<asc|desc>**<br/>
Sort items in ascending or descending order. Default is ascending if sorting by title or feed, otherwise descending.

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

List feeds:

    kibner feeds [options]

Options:

**--with-title=<title>**<br/>
Only show feeds with titles that match the given value.

**--with-author=<author>**<br/>
Only show feeds with authors that match the given value.

**-d**<br/>
**--show-desc**<br/>
Show feed descriptions.

**-N=<number>**<br/>
**--top=<number>**<br/>
Set the maximum number of feeds to display.

**--sortby=<pubdate|title|items|unplayed|timestamp>**<br/>
Sort items by the given property. Default sort is by pubdate -- the most recently updated feeds appear first.

**--order=<asc|desc>**<br/>
Sort feeds in ascending or descending order. Default is ascending if sorting by title, otherwise descending.

### Other Tasks

Update feed details:

    kibner update <options> <feed>

Options:

**--title**<br/>
Set the feed title to the given value.

**--author**<br/>
Set the feed author to the given value.

**--desc**<br/>
Set the feed description to the given value.

**--link**<br/>
Set the feed website to the given URL.

**--image**<br/>
Set the feed artwork to the given URL.

Import feeds from a file:

    kibner import [options] <filename>

Options:

**--format=<list|opml>**<br/>
Specify the file format. Format 'list' is a plain text file with one feed URL per line. Format 'opml' is an [OPML file](https://en.wikipedia.org/wiki/OPML).

Export feeds to a file:

    kibner export [options] <filename>

Options:

**--format=<list|opml>**<br/>
Specify the file format. Format 'list' is a plain text file with one feed URL per line. Format 'opml' is an [OPML file](https://en.wikipedia.org/wiki/OPML).

Print Kibner's version number.

    kibner version

Reset Kibner (wipes existing data!)

    kibner reset

## TODO

### Features

- Provide values for command-line flags in config file.
- Subscribe to BBC iPlayer audio.
- Improve terminal output.
- Windows support.
- Allow user-defined templates for list/feed commands.
- Play/download feed content (e.g. audio, video).
- Pause/Resume, Mute/Unmute feeds.

### Code

- Refactor database code.
- Refactor command/flag/config code (use Viper?).
- Increase test coverage.
- Handle interrupt signals (e.g. Ctrl-C).

## Licensing

Kibner is provided under an [MIT License](http://choosealicense.com/licenses/mit/). See the [LICENSE](LICENSE) file for details.
