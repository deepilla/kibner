# Kibner

[![Build Status](https://travis-ci.org/deepilla/kibner.svg?branch=master)](https://travis-ci.org/deepilla/kibner)
[![Go Report Card](https://goreportcard.com/badge/github.com/deepilla/kibner)](https://goreportcard.com/report/github.com/deepilla/kibner)

Kibner is a command-line utility for managing podcasts.

Subscribe to podcasts, keep them in sync, and play episodes with
a few simple commands. Kibner is designed to be minimal. It keeps
track of your subscriptions but relies on other programs for media
playback and downloads.

**NOTE**: This is a work in progress. Commands and flags are likely
to change.

## Install

    go install github.com/deepilla/kibner

## Usage

Everyday usage requires just a few commands:

- `add` to subscribe to feeds
- `sync` to check for new feed items
- `list` to display, play and download items

Type `kibner` with no arguments to see all available commands.

### Subscribe to a feed

    kibner add [options] <url>

Add a feed to your subscriptions. The url can either be an RSS
feed or, for podcasts that don't provide a feed, an iTunes page.
To add multiple feeds, see the `import` command.

Options:

**--itunes**<br/>
Indicate that `url` is an iTunes page rather than an RSS feed.

### Unsubscribe from a feed

    kibner remove <feed>

Remove a feed from your subscriptions. If the provided feed name
matches more than one feed, you will be prompted to choose between
them. For example, `kibner remove this` would match both *This
American Life* and *Answer Me This* and ask you to confirm which
you wanted to remove. To unsubscribe from all feeds, see the
`reset` command.

### Synchronise feeds

    kibner sync [feed]

Check feeds for new items and update your subscriptions.
Synchronise an individual feed by specifying a feed name.

### List/Play items

    kibner list [options] [feed]

Display and interact with items from your subscribed feeds.
Use the options to sort and filter the results, as well as
play selected items. Playback is via the media player of your
choice. You can also run arbitrary programs on items (such as
[curl](https://curl.haxx.se/) for downloading).

Options:

**-d**, **--show-desc**<br/>
Show item descriptions.

**-N**, **--top**=*number*<br/>
Set the maximum number of items to display.

**-u**, **--unplayed**<br/>
Only show unplayed items.

**-T**, **--since**=*date*<br/>
Only show items published on or after the given date. Valid date
values are:

- *today* for the current date
- *Nd* for *N* days ago (0d is the current date)
- *Nm* for *N* months ago (0m is the start of the current month)
- *Ny* for *N* years ago (0y is the start of the current year)
- *mon*, *tue*, *wed*, *thu*, *fri*, *sat* or *sun* for a specific
day of the week
- *jan*, *feb*, *mar*, *apr*, *may*, *jun*, *jul*, *aug*, *sep*,
*oct*, *nov* or *dec* for a specific month

**--with-title**=*title*<br/>
Only show items with titles that match the given value.

**--sortby**=*property*<br/>
Sort items by the given property. The available properties are:

- *pubdate* to sort by publish date, most recent first
- *title* to sort alphabetically by title
- *feed* to sort alphabetically by feed title
- *duration* to sort by duration, longest first
- *timestamp* to sort by local creation time, most recent first

The default is pubdate.

**--order**=*order*<br/>
Display items in ascending (*asc*) or descending (*desc*) order.
The default is ascending if sorting by title or feed, otherwise
descending.

**-p**, **--play**<br/>
Play the selected items using the program specified by the `--use`
option. If playback successfully completes, the item is marked as
played (and therefore no longer appears in the output when the
`--unplayed` flag is specified).

**--mark**<br/>
Mark selected items as played.

**--unmark**<br/>
Mark selected items as unplayed.

**--run**<br/>
Run the program specified by the `--use` option on the selected
items. Unlike the `--play` option, this does not mark the items
as played.

**--use**=*program*<br/>
Specify a program to use with the `--play` or `--run` options.

### List feeds

    kibner feeds [options]

Display subscribed feeds. Use the options to sort and filter the
results.

Options:

**-d**, **--show-desc**<br/>
Show feed descriptions.

**-N**, **--top**=*number*<br/>
Set the maximum number of feeds to display.

**--with-title**=*title*<br/>
Only show feeds with titles that match the given value.

**--with-author**=*author*<br/>
Only show feeds with authors that match the given value.

**--sortby**=*property*<br/>
Sort feeds by the given property. The available properties are:

- *pubdate* sorts by last published date, most recently updated
first
- *title* sorts alphabetically by title
- *items* sorts by number of items, highest first
- *unplayed* sorts by number of unplayed items, highest first
- *timestamp* sorts by local creation time, most recent first

The default is pubdate.

**--order**=*order*<br/>
Display feeds in ascending (*asc*) or descending (*desc*) order.
The default is ascending if sorting by title, otherwise descending.

### Other tasks

The following commands are less commonly used.

#### Import feeds

    kibner import [options] <filename>

Import feeds from a file.

Options:

**--format**=*format*<br/>
Specify a file format. Valid values are:

- *list* for a plain text file with one feed URL per line
- *opml* for an [OPML file](https://en.wikipedia.org/wiki/OPML)

The default is list.

#### Export feeds

    kibner export [options] <filename>

Export feeds to a file.

Options:

**--format**=*format*<br/>
Specify a file format. Valid values are:

- *list* for a plain text file with one feed URL per line
- *opml* for an [OPML file](https://en.wikipedia.org/wiki/OPML)

The default is list.

#### Open a feed URL

    kibner open [options] <feed>

Open a URL associated with the given feed.

Options:

**--target**=*type*<br/>
Specify the type of URL to open. Valid values are:

- *feed* for the feed's RSS feed
- *link* for the website associated with the feed
- *image* for the feed's artwork

The default is feed. Note that link and image URLs are not
guaranteed to exist.

**--use**=*program*<br/>
Specify a program to open the URL.

#### Update feed details

    kibner update <options> <feed>

If for some reason you don't like a feed's values as specified
in its RSS feed, you can change them for your local subscription.

Options:

**--title**=*title*<br/>
Set the feed title to the given value.

**--author**=*author*<br/>
Set the feed author to the given value.

**--desc**=*desc*<br/>
Set the feed description to the given value.

**--link**=*link*<br/>
Set the feed website to the given URL.

**--image**=*image*<br/>
Set the feed artwork to the given URL.

#### Nuke your data

    kibner reset

Wipe all of your existing data and start fresh with a clean database.
It's probably a good idea to export your feeds to a file before doing
this!

#### Version number

    kibner version

## TODO

### Features

- Specify command-line options via config file.
- Make sure everything works in Windows.
- Subscribe to BBC iPlayer audio.
- Allow user-defined templates for list/feed commands.
- Pause/Resume, Mute/Unmute feeds.

### Code

- Come up with an SQLite vacuum strategy.
- Handle interrupt signals (e.g. Ctrl-C).
- Increase test coverage.
- Refactor command/flag/config code (use Viper?).

## Licensing

Kibner is provided under an [MIT License](http://choosealicense.com/licenses/mit/). See the [LICENSE](LICENSE) file for details.
