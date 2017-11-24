package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/deepilla/itunes"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/go-homedir"
)

const version = "0.1"

const (
	flagITunes     = "itunes"
	flagTitle      = "title"
	flagAuthor     = "author"
	flagDesc       = "desc"
	flagLink       = "link"
	flagImage      = "image"
	flagSortBy     = "sortby"
	flagSortOrder  = "order"
	flagLimit      = "top"
	flagStartDate  = "since"
	flagUnplayed   = "unplayed"
	flagWithTitle  = "with-title"
	flagWithAuthor = "with-author"
	flagShowDesc   = "show-desc"
	flagPlay       = "play"
	flagMark       = "mark"
	flagUnmark     = "unmark"
	flagRun        = "run"
	flagUse        = "use"
	flagFormat     = "format"
	flagTarget     = "target"
)

// TODO: Make these settings configurable.
var defaults = struct {
	Timeout    time.Duration
	MaxWorkers int
}{
	Timeout:    10 * time.Second,
	MaxWorkers: 10,
}

var defaultClient = &http.Client{
	Timeout: defaults.Timeout,
}

func main() {

	cs := NewCommandSet("kibner", getCommands()...)

	if len(os.Args) < 2 {
		cs.Usage()
		os.Exit(1)
	}

	err := cs.Run(os.Args[1], os.Args[2:])
	if err != nil {
		fmt.Println("Whoops:", err)
		os.Exit(2)
	}
}

func getCommands() []*Command {

	var sortItemsOpt uintFlag
	sortItemsOpt.AddValue("pubdate", sortItemsByPubdate, "Sort by pubdate")
	sortItemsOpt.AddValue("title", sortItemsByTitle, "Sort by title")
	sortItemsOpt.AddValue("feed", sortItemsByFeed, "Sort by feed")
	sortItemsOpt.AddValue("duration", sortItemsByDuration, "Sort by duration")
	sortItemsOpt.AddValue("timestamp", sortItemsByTimestamp, "Sort by timestamp")
	sortItemsOpt.MustSet("pubdate")

	var sortFeedsOpt uintFlag
	sortFeedsOpt.AddValue("pubdate", sortFeedsByPubdate, "Sort by last pubdate")
	sortFeedsOpt.AddValue("title", sortFeedsByTitle, "Sort by title")
	sortFeedsOpt.AddValue("items", sortFeedsByItemCount, "Sort by item count")
	sortFeedsOpt.AddValue("unplayed", sortFeedsByUnplayedCount, "Sort by unplayed count")
	sortFeedsOpt.AddValue("timestamp", sortFeedsByTimestamp, "Sort by timestamp")
	sortFeedsOpt.MustSet("pubdate")

	var sortOrderOpt uintFlag
	sortOrderOpt.AddValue("asc", sortOrderAsc, "Ascending order")
	sortOrderOpt.AddValue("desc", sortOrderAsc, "Descending order")

	var fileFormatOpt uintFlag
	fileFormatOpt.AddValue("list", fileFormatList, "Plain text")
	fileFormatOpt.AddValue("opml", fileFormatOPML, "OPML file")
	fileFormatOpt.MustSet("list")

	var targetOpt uintFlag
	targetOpt.AddValue("link", targetLink, "Website")
	targetOpt.AddValue("feed", targetFeed, "RSS Feed")
	targetOpt.AddValue("image", targetImage, "Image")
	targetOpt.MustSet("feed")

	return []*Command{

		NewCommand("add",
			runAdd,
			WithSyntax("kibner add [options] <url>"),
			WithDescription("Subscribe to a feed"),
			WithOption(flagITunes, "add an iTunes page instead of an RSS feed", false),
		),

		NewCommand("remove",
			runRemove,
			WithAlias("rm"),
			WithSyntax("kibner remove <name>"),
			WithDescription("Unsubscribe from a feed"),
		),

		NewCommand("update",
			runUpdate,
			WithSyntax("kibner update <options> <name>"),
			WithDescription("Edit feed details"),
			WithOption(flagTitle, "set feed title to the given value", ""),
			WithOption(flagAuthor, "set feed author to the given value", ""),
			WithOption(flagDesc, "set feed description to the given value", ""),
			WithOption(flagLink, "set feed website to the given value", ""),
			WithOption(flagImage, "set feed image to the given value", ""),
		),

		NewCommand("sync",
			runSync,
			WithSyntax("kibner sync [name]"),
			WithDescription("Check for new items"),
		),

		NewCommand("feeds",
			runFeeds,
			WithSyntax("kibner feeds [options]"),
			WithDescription("List feeds"),
			WithOption(flagSortBy, "sort feeds by the given property", sortFeedsOpt),
			WithOption(flagSortOrder, "sort in ascending or descending order", sortOrderOpt),
			WithOptionAlias(flagLimit, "N", "the maximum `number` of feeds to display", uint(0)),
			WithOption(flagWithTitle, "show feeds that match the given title", ""),
			WithOption(flagWithAuthor, "show feeds that match the given author", ""),
			WithOptionAlias(flagShowDesc, "d", "show feed descriptions", false),
		),

		NewCommand("list",
			runList,
			WithAlias("ls"),
			WithSyntax("kibner list [options] [name]"),
			WithDescription("List items"),
			WithOption(flagSortBy, "sort items by the given property", sortItemsOpt),
			WithOption(flagSortOrder, "sort in ascending or descending order", sortOrderOpt),
			WithOptionAlias(flagLimit, "N", "the maximum `number` of items to display", uint(0)),
			WithOptionAlias(flagStartDate, "T", "show items released on or after the given `date`", reldate{}),
			WithOptionAlias(flagUnplayed, "u", "show unplayed items", false),
			WithOption(flagWithTitle, "show items that match the given title", ""),
			WithOptionAlias(flagPlay, "p", "play selected items", false),
			WithOption(flagMark, "mark selected items as played", false),
			WithOption(flagUnmark, "mark selected items as unplayed", false),
			WithOption(flagRun, "run the specified program on selected items", false),
			WithOption(flagUse, "a `program` to play or run items", ""),
			WithOptionAlias(flagShowDesc, "d", "show item descriptions", false),
		),

		NewCommand("import",
			runImport,
			WithSyntax("kibner import [options] <filename>"),
			WithDescription("Import subscriptions from a file"),
			WithOption(flagFormat, "the `type` of file being imported", fileFormatOpt),
		),

		NewCommand("export",
			runExport,
			WithSyntax("kibner export [options] <filename>"),
			WithDescription("Export subscriptions to a file"),
			WithOption(flagFormat, "the `type` of file to export", fileFormatOpt),
		),

		NewCommand("open",
			runOpen,
			WithSyntax("kibner open [options] <name>"),
			WithDescription("View a feed's website, RSS feed, or image"),
			WithOption(flagTarget, "which `property` of the feed to view", targetOpt),
			WithOption(flagUse, "a `program` to use as the viewer", "xdg-open"),
		),

		NewCommand("reset",
			runReset,
			WithSyntax("kibner nuke"),
			WithDescription("Wipe all existing data"),
		),

		NewCommand("version",
			runVersion,
			WithSyntax("kibner version"),
			WithDescription("Print kibner's version number"),
		),
	}
}

func runDB(fn func(*sql.DB) error) error {

	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	return fn(db)
}

// Commands

func runAdd(opts Options, args []string, env *Env) error {

	if len(args) != 1 {
		return ErrBadArgs
	}

	url := args[0]

	var err error
	if opts.Get(flagITunes).Bool() {
		url, err = itunes.ToRSSClient(url, defaultClient)
		if err != nil {
			return errors.New("could not extract feed from iTunes page: " + err.Error())
		}
	}

	return runDB(func(db *sql.DB) error {

		res, err := addFeed(db, url)
		if err != nil {
			return err
		}

		fmt.Printf("Added %s - %d items\n", res.Title, res.Items)
		return nil
	})
}

func runRemove(opts Options, args []string, env *Env) error {

	if len(args) != 1 {
		return ErrBadArgs
	}

	return runDB(func(db *sql.DB) error {

		id, err := chooseFeed(db, args[0], "Remove %s", true)
		if err != nil {
			return err
		}

		return removeFeed(db, id)
	})
}

func runUpdate(opts Options, args []string, env *Env) error {

	if len(args) != 1 {
		return ErrBadArgs
	}

	fieldsByFlag := map[string]string{
		flagTitle:  "title",
		flagAuthor: "author",
		flagDesc:   "desc",
		flagLink:   "link",
		flagImage:  "image",
	}

	values := map[string]interface{}{}

	for _, opt := range opts {

		s := strings.TrimSpace(opt.String())
		if s == "" {
			continue
		}

		fieldName, ok := fieldsByFlag[opt.Name]
		if !ok {
			return errors.New("unsupported update " + opt.Name)
		}

		values[fieldName] = s
	}

	if len(values) == 0 {
		return errors.New("no values to update")
	}

	return runDB(func(db *sql.DB) error {

		id, err := chooseFeed(db, args[0], "Update %s", true)
		if err != nil {
			return err
		}

		return updateFeed(db, id, values)
	})
}

func runSync(opts Options, args []string, env *Env) error {

	switch len(args) {
	case 0:
		return runDB(runSyncAll)
	case 1:
		return runDB(func(db *sql.DB) error {
			return runSyncOne(db, args[0])
		})
	default:
		return ErrBadArgs
	}
}

func runSyncAll(db *sql.DB) error {

	results, err := syncAll(db)
	if err != nil {
		return err
	}

	resetOutput()
	printSyncResults(results)
	return nil
}

func runSyncOne(db *sql.DB, feedName string) error {

	id, err := chooseFeed(db, feedName, "Sync %s", false)
	if err != nil {
		return err
	}

	res, err := syncOne(db, id)
	if err != nil {
		return err
	}

	fmt.Printf("%s: ", res.Title)

	resetOutput()
	printSyncResults([]*syncResult{res})
	return nil
}

func printSyncResults(results []*syncResult) {

	items, errs := 0, 0

	for _, res := range results {

		items += res.Items
		if res.Err != nil {
			errs++
		}
	}

	switch items {
	case 0:
		fmt.Print("No new items")
	case 1:
		fmt.Print("1 new item")
	default:
		fmt.Printf("%d new items", items)
	}

	switch errs {
	case 0:
	case 1:
		fmt.Print(", 1 error")
	default:
		fmt.Printf(", %d errors", errs)
	}

	fmt.Println()

	i := 0
	for _, res := range results {

		if res.Err == nil {
			continue
		}

		fmt.Println(i+1, res.Title, res.Err)
		i++
	}
}

func runFeeds(opts Options, args []string, env *Env) error {

	if len(args) != 0 {
		return ErrBadArgs
	}

	listOpts := listFeedOptions{
		SortBy:    opts.Get(flagSortBy).Value().(sortFeedsBy),
		SortOrder: opts.Get(flagSortOrder).Value().(sortOrder),
		Limit:     opts.Get(flagLimit).Uint(),
		Title:     opts.Get(flagWithTitle).String(),
		Author:    opts.Get(flagWithAuthor).String(),
		ShowDesc:  opts.Get(flagShowDesc).Bool(),
	}

	return runDB(func(db *sql.DB) error {
		return listFeeds(db, env.Stdout, defaultFeedTemplate(time.Now()), listOpts)
	})
}

func runList(opts Options, args []string, env *Env) error {

	nArgs := len(args)
	if nArgs != 0 && nArgs != 1 {
		return ErrBadArgs
	}

	action := actionNone

	switch {
	case opts.Get(flagPlay).Bool():
		action = actionPlay
	case opts.Get(flagRun).Bool():
		action = actionRun
	case opts.Get(flagMark).Bool():
		action = actionMark
	case opts.Get(flagUnmark).Bool():
		action = actionUnmark
	}

	listOpts := listItemOptions{
		SortBy:    opts.Get(flagSortBy).Value().(sortItemsBy),
		SortOrder: opts.Get(flagSortOrder).Value().(sortOrder),
		Limit:     opts.Get(flagLimit).Uint(),
		Unplayed:  opts.Get(flagUnplayed).Bool(),
		StartDate: opts.Get(flagStartDate).Value().(time.Time),
		Title:     opts.Get(flagWithTitle).String(),
		ShowDesc:  opts.Get(flagShowDesc).Bool(),
		Action:    action,
		Use:       opts.Get(flagUse).String(),
	}

	return runDB(func(db *sql.DB) error {

		if nArgs == 1 {
			feedID, err := chooseFeed(db, args[0], "List %s", false)
			if err != nil {
				return err
			}
			listOpts.FeedID = feedID
		}

		return listItems(db, env.Stdout, defaultItemTemplate(time.Now()), listOpts)
	})
}

func runImport(opts Options, args []string, env *Env) error {

	if len(args) != 1 {
		return ErrBadArgs
	}

	var extractURLs func(io.Reader) ([]string, error)

	switch opts.Get(flagFormat).Value().(fileFormat) {
	case fileFormatList:
		extractURLs = extractURLsList
	case fileFormatOPML:
		extractURLs = extractURLsOPML
	default:
		return errors.New("unsupported file format")
	}

	f, err := os.Open(args[0])
	if err != nil {
		return err
	}
	defer f.Close()

	urls, err := extractURLs(f)
	if err != nil {
		return err
	}

	return runDB(func(db *sql.DB) error {
		resetOutput()
		printImportResults(addFeedMultiple(db, urls))
		return nil
	})
}

func printImportResults(results []*syncResult) {

	var oks, errs int

	for _, res := range results {
		if res.Err != nil {
			errs++
		} else {
			oks++
		}
	}

	fmt.Printf("Added %d of %d feeds, %d errors\n", oks, len(results), errs)

	for _, res := range results {
		if res.Err != nil {
			fmt.Println(res.Err)
		}
	}
}

func runExport(opts Options, args []string, env *Env) error {

	if len(args) != 1 {
		return ErrBadArgs
	}

	var export func(*sql.DB, io.Writer) error

	switch opts.Get(flagFormat).Value().(fileFormat) {
	case fileFormatList:
		export = exportList
	case fileFormatOPML:
		export = func(db *sql.DB, w io.Writer) error {
			return exportOPML(db, w, time.Now())
		}
	default:
		return errors.New("unsupported file format")
	}

	f, err := os.Create(args[0])
	if err != nil {
		return err
	}
	defer f.Close()

	return runDB(func(db *sql.DB) error {
		return export(db, f)
	})
}

func runOpen(opts Options, args []string, env *Env) error {

	if len(args) != 1 {
		return ErrBadArgs
	}

	var name string
	target := opts.Get(flagTarget).Value().(target)

	switch target {
	case targetFeed:
		name = "RSS feed"
	case targetImage:
		name = "image"
	case targetLink:
		name = "website"
	default:
		return errors.New("unsupported target")
	}

	cmd, err := parseCommand(opts.Get(flagUse).String())
	if err != nil {
		return err
	}

	return runDB(func(db *sql.DB) error {

		id, err := chooseFeed(db, args[0], "Open "+name+" for %s", false)
		if err != nil {
			return err
		}

		url, err := loadFeedURL(db, id, target)
		if err != nil {
			return err
		}

		if url == "" {
			return errors.New("no " + name + " found for this feed")
		}

		cmd.Args = append(cmd.Args, url)
		return cmd.Start()
	})
}

func runReset(opts Options, args []string, env *Env) error {

	if len(args) != 0 {
		return ErrBadArgs
	}

	c, err := ask("Reset Kibner (this will wipe all existing data!)? Yes, No", "yn")
	if err != nil {
		return err
	}

	if c == 'n' {
		return nil
	}

	return runDB(initDB)
}

func runVersion(opts Options, args []string, env *Env) error {

	if len(args) != 0 {
		return ErrBadArgs
	}

	fmt.Printf("Kibner v%s\n", version)
	return nil
}

// Custom flags

type sortItemsBy uint

const (
	sortItemsByPubdate sortItemsBy = iota
	sortItemsByTitle
	sortItemsByFeed
	sortItemsByDuration
	sortItemsByTimestamp
)

type sortFeedsBy uint

const (
	sortFeedsByPubdate sortFeedsBy = iota
	sortFeedsByTitle
	sortFeedsByItemCount
	sortFeedsByUnplayedCount
	sortFeedsByTimestamp
)

type sortOrder uint

const (
	sortOrderDefault sortOrder = iota
	sortOrderAsc
	sortOrderDesc
)

type fileFormat uint

const (
	fileFormatList fileFormat = iota
	fileFormatOPML
)

type target uint

const (
	targetFeed target = iota
	targetLink
	targetImage
)

type uintFlag struct {
	typ     reflect.Type
	value   interface{}
	inputs  map[string]interface{}
	strings map[interface{}]string
	usage   string
}

func (f *uintFlag) Get() interface{} {
	return f.value
}

func (f *uintFlag) Set(val string) error {

	v, ok := f.inputs[val]
	if !ok {
		return errors.New(f.usage)
	}

	f.value = v
	return nil
}

func (f *uintFlag) MustSet(val string) {
	if err := f.Set(val); err != nil {
		panic(err)
	}
}

func (f *uintFlag) String() string {
	return f.strings[f.value]
}

func (f *uintFlag) AddValue(name string, value interface{}, desc string) {

	typ := reflect.TypeOf(value)

	if f.typ == nil {
		f.typ = typ
		f.value = reflect.Zero(typ).Interface()
	}

	if f.typ != typ {
		panic(fmt.Sprintf("invalid uintFlag: expected Type %v, got %v", f.typ, typ))
	}

	if f.inputs == nil {
		f.inputs = make(map[string]interface{})
	}
	f.inputs[name] = value

	if f.strings == nil {
		f.strings = make(map[interface{}]string)
	}
	f.strings[value] = desc

	if f.usage == "" {
		f.usage = "valid values -- " + name
	} else {
		f.usage += " | " + name
	}
}

type reldate time.Time

func (v *reldate) Set(val string) error {

	weekdays := map[string]time.Weekday{
		"mon": time.Monday,
		"tue": time.Tuesday,
		"wed": time.Wednesday,
		"thu": time.Thursday,
		"fri": time.Friday,
		"sat": time.Saturday,
		"sun": time.Sunday,
	}

	months := map[string]time.Month{
		"jan": time.January,
		"feb": time.February,
		"mar": time.March,
		"apr": time.April,
		"may": time.May,
		"jun": time.June,
		"jul": time.July,
		"aug": time.August,
		"sep": time.September,
		"oct": time.October,
		"nov": time.November,
		"dec": time.December,
	}

	now := startOfDay(time.Now())

	val = strings.ToLower(val)

	if val == "today" || val == "0d" {
		*v = reldate(now)
		return nil
	}

	if weekday, ok := weekdays[val]; ok {
		d := now
		for d.Weekday() != weekday {
			d = d.AddDate(0, 0, -1)
		}
		*v = reldate(d)
		return nil
	}

	if month, ok := months[val]; ok {
		d := time.Date(now.Year(), month, 1, 0, 0, 0, 0, time.Local)
		if d.After(now) {
			d = d.AddDate(-1, 0, 0)
		}
		*v = reldate(d)
		return nil
	}

	if y, m, d, err := parseYMD(val); err == nil {
		*v = reldate(now.AddDate(-y, -m, -d))
		return nil
	}

	return errors.New("valid values: today | mon...sun | jan...dec | Nd | Nm | Ny")
}

func (v *reldate) String() string {

	t := time.Time(*v)
	if t.IsZero() {
		return ""
	}

	return t.Format("02 Jan, 2006")
}

func (v *reldate) Get() interface{} {
	return time.Time(*v)
}

func parseYMD(val string) (int, int, int, error) {

	if len(val) < 2 {
		return 0, 0, 0, errors.New("could not parse YMD " + val)
	}

	last := len(val) - 1
	n, err := strconv.Atoi(val[:last])
	if err != nil {
		return 0, 0, 0, errors.New("could not parse YMD " + val + ": " + err.Error())
	}
	if n <= 0 {
		return 0, 0, 0, errors.New("could not parse YMD " + val + ": number must be greater than zero")
	}

	switch val[last:] {
	case "y":
		return n, 0, 0, nil
	case "m":
		return 0, n, 0, nil
	case "d":
		return 0, 0, n, nil
	default:
		return 0, 0, 0, errors.New("could not parse YMD " + val + ": unit must be y, m or d")
	}
}

// Database functions

func openDB() (*sql.DB, error) {

	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(home, ".config", "kibner", "db", "kibner.db")
	return dbAtPath(path)
}

func dbAtPath(path string) (*sql.DB, error) {

	isNew := false
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		isNew = true
		dir, _ := filepath.Split(path)
		err = os.MkdirAll(dir, 0755)
	}
	if err != nil {
		return nil, err
	}

	db, err := connect(path, map[string]string{
		"_foreign_keys": "1",
	})
	if err != nil {
		return nil, err
	}

	if isNew {
		if err := initDB(db); err != nil {
			db.Close()
			return nil, err
		}
	}

	return db, nil
}

func connect(path string, options map[string]string) (*sql.DB, error) {

	vals := url.Values{}
	for k, v := range options {
		vals.Set(k, v)
	}

	if s := vals.Encode(); s != "" {
		path += "?" + s
	}

	return sql.Open("sqlite3", "file:"+path)
}

func resetOutput() {
	// TODO: Implement on Windows
	fmt.Printf("%c[2K\r", 27)
}
