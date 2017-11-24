package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
	"unicode"

	kibner "github.com/deepilla/kibner/internal/types"

	"github.com/mmcdole/gofeed"
)

func initDB(db *sql.DB) error {

	sql := []string{

		`DROP TABLE IF EXISTS items`,

		`DROP TABLE IF EXISTS feeds`,

		// The feeds table has an AUTOINCREMENT id which means
		// that the ids of deleted rows are not reused. This
		// gives us an immutable id which we can use to uniquely
		// identify a feed. So we can select a feed, store its
		// id, and operate on it some time later without fear
		// of modifying a different row with the same id.

		`CREATE TABLE feeds (
			id				INTEGER PRIMARY KEY AUTOINCREMENT,
			title			TEXT NOT NULL,
			author			TEXT NOT NULL,
			desc			TEXT,
			type			TEXT NOT NULL,
			url				TEXT NOT NULL,
			image			TEXT,
			link			TEXT,
			timestamp		DATETIME NOT NULL
		)`,

		// We could create the following index by adding a UNIQUE
		// constraint to the url field. But creating it explicitly
		// allows us to specify the index name (which makes unit
		// tests more robust).

		`CREATE UNIQUE INDEX unique_feed_url ON feeds(url)`,

		`CREATE TABLE items (
			feedid			INTEGER NOT NULL REFERENCES feeds(id),
			title			TEXT NOT NULL,
			desc			TEXT,
			pubdate			DATETIME NOT NULL,
			url				TEXT NOT NULL,
			filesize		INTEGER DEFAULT 0,
			duration		INTEGER DEFAULT 0,
			guid			TEXT NOT NULL,
			unplayed		BOOLEAN DEFAULT 0,
			timestamp		DATETIME NOT NULL
		)`,

		// TODO: Do we need to rely on the GUID given that an item
		// is uniquely identified by its URL?

		`CREATE UNIQUE INDEX unique_item_guid ON items(feedid, guid)`,
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, q := range sql {
		if _, err := tx.Exec(q); err != nil {
			return rollback(tx, err)
		}
	}

	return tx.Commit()
}

type syncResult struct {
	ID    int64
	URL   string
	Title string
	Items int
	Err   error
}

func addFeed(db *sql.DB, url string) (*syncResult, error) {

	feed, err := fetchAndParse(url)
	if err != nil {
		return nil, errors.New("could not fetch feed: " + err.Error())
	}

	id, err := saveFeed(db, feed, time.Now())
	if err != nil {
		return nil, errors.New("could not save feed: " + err.Error())
	}

	return &syncResult{
		ID:    id,
		URL:   feed.URL,
		Title: feed.Title,
		Items: len(feed.Items),
	}, nil
}

func addFeedMultiple(db *sql.DB, urls []string) []*syncResult {

	feeds, errs := fetchAndParseMultiple(urls, defaults.MaxWorkers)

	i := 0
	now := time.Now()
	results := make([]*syncResult, 0, len(urls))

	for url, feed := range feeds {

		i++
		fmt.Printf("Adding %d of %d feeds\r", i, len(feeds))

		id, err := saveFeed(db, feed, now)
		if err != nil {
			errs[url] = err
			continue
		}

		results = append(results, &syncResult{
			ID:    id,
			URL:   feed.URL,
			Title: feed.Title,
			Items: len(feed.Items),
		})
	}

	for url, err := range errs {
		results = append(results, &syncResult{
			URL: url,
			Err: err,
		})
	}

	return results
}

func removeFeed(db *sql.DB, id int64) error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = deleteItems(tx, id)
	if err != nil {
		return rollback(tx, err)
	}

	err = deleteFeed(tx, id)
	if err != nil {
		return rollback(tx, err)
	}

	return tx.Commit()
}

func syncOne(db *sql.DB, id int64) (*syncResult, error) {

	infos, err := loadSyncInfo(db, id)
	if err != nil {
		return nil, err
	}

	if n := len(infos); n != 1 {
		if n == 0 {
			return nil, errNoFeedFound
		}
		return nil, errors.New("could not retrieve feed")
	}

	info := &infos[0]
	feed, err := fetchAndParse(info.URL)
	if err != nil {
		return nil, err
	}

	items, err := syncItems(db, info, feed.Items)
	if err != nil {
		return nil, err
	}

	err = syncFeed(db, info, feed)
	if err != nil {
		// Swallow this error
		log.Println(err)
	}

	return &syncResult{
		ID:    info.ID,
		URL:   info.URL,
		Title: info.Title,
		Items: items,
	}, nil
}

func syncAll(db *sql.DB) ([]*syncResult, error) {

	infos, err := loadSyncInfo(db, 0)
	if err != nil {
		return nil, err
	}

	if len(infos) == 0 {
		return nil, errors.New("no feeds to sync")
	}

	urls := make([]string, len(infos))
	mURLToInfo := make(map[string]*syncInfo, len(infos))
	results := make([]*syncResult, 0, len(infos))

	for i := range infos {
		urls[i] = infos[i].URL
		mURLToInfo[infos[i].URL] = &infos[i]
	}

	feeds, errs := fetchAndParseMultiple(urls, defaults.MaxWorkers)

	i := 0
	for url, feed := range feeds {

		i++
		fmt.Printf("Syncing %d of %d feeds\r", i, len(feeds))

		info := mURLToInfo[url]

		if err := syncFeed(db, info, feed); err != nil {
			// Swallow this error.
			log.Println(err)
		}

		items, err := syncItems(db, info, feed.Items)
		if err != nil {
			errs[url] = err
			continue
		}

		results = append(results, &syncResult{
			ID:    info.ID,
			URL:   info.URL,
			Title: info.Title,
			Items: items,
		})
	}

	for url, err := range errs {

		info := mURLToInfo[url]

		results = append(results, &syncResult{
			ID:    info.ID,
			URL:   info.URL,
			Title: info.Title,
			Err:   err,
		})
	}

	return results, nil
}

type listFeedOptions struct {
	SortBy    sortFeedsBy
	SortOrder sortOrder
	Limit     uint
	Title     string
	Author    string
	ShowDesc  bool
}

func listFeeds(db *sql.DB, w io.Writer, tmpl *template.Template, opts listFeedOptions) error {

	feeds, err := loadFeedViews(db, opts)
	if err != nil {
		return err
	}

	return tmpl.Execute(w, map[string]interface{}{
		"Feeds":    feeds,
		"ShowDesc": opts.ShowDesc,
	})
}

type feedView struct {
	Title         string
	Author        string
	Desc          string
	Items         int64
	UnplayedItems int64
	LastPubdate   time.Time
}

func loadFeedViews(db *sql.DB, opts listFeedOptions) ([]feedView, error) {

	// ORDER BY clause

	var order string

	switch opts.SortOrder {
	case sortOrderAsc:
		order = "ASC"
	case sortOrderDesc:
		order = "DESC"
	case sortOrderDefault:
		switch opts.SortBy {
		case sortFeedsByTitle:
			order = "ASC"
		default:
			order = "DESC"
		}
	}

	secondaryFields := "f.title COLLATE NOCASE ASC, f.ROWID DESC"

	var sortFields string

	switch opts.SortBy {
	case sortFeedsByPubdate:
		sortFields = fmt.Sprintf("max_pubdate %s, %s", order, secondaryFields)
	case sortFeedsByTitle:
		sortFields = fmt.Sprintf("f.title COLLATE NOCASE %s, f.ROWID DESC", order)
	case sortFeedsByItemCount:
		sortFields = fmt.Sprintf("item_count %s, %s", order, secondaryFields)
	case sortFeedsByUnplayedCount:
		sortFields = fmt.Sprintf("unplayed_count %s, %s", order, secondaryFields)
	case sortFeedsByTimestamp:
		sortFields = fmt.Sprintf("f.timestamp %s, %s", order, secondaryFields)
	default:
		return nil, errors.New("unsupported sort order")
	}

	// WHERE clause

	var params []interface{}
	var conditions []string

	params = append(params, time.Time{}.Unix()) // Default value for NULL pubdates

	if title := opts.Title; title != "" {
		conditions = append(conditions, "f.title LIKE ?")
		params = append(params, "%"+title+"%")
	}

	if author := opts.Author; author != "" {
		conditions = append(conditions, "f.author LIKE ?")
		params = append(params, "%"+author+"%")
	}

	whereClause := "1=1"
	if len(conditions) > 0 {
		whereClause = strings.Join(conditions, " AND ")
	}

	// LIMIT

	limit := int(opts.Limit)
	if limit <= 0 {
		limit = -1
	}

	// Build and execute query

	// A few points about this query.
	//
	// 1. The join between the feeds and items tables is an
	// outer join, just in case we have any feeds with zero
	// items. I haven't seen any empty feeds in the wild but
	// the RSS spec allows them and there's no reason not to
	// support them.
	//
	// 2. COUNT(i.ROWID) gives the number of items for each
	// feed. ROWID is a convenient value to count because it
	// can never be NULL. Using COUNT(*) would not work as
	// that would include the feed row in the count, reporting
	// 1 instead of 0 for feeds with no items.
	//
	// 3. i.unplayed is a BOOLEAN field with values 1 (True)
	// or 0 (False). So SUM(i.unplayed) gives the number of
	// unplayed items for each feed. SUM will return NULL
	// if a feed has zero items (hence the IFNULL).
	//
	// 4. i.pubdate is a DATETIME field with integer values
	// representing the Unix time in seconds. MAX(i.pubdate)
	// gives the date of the most recent feed item. Like SUM,
	// MAX will return NULL if a feed has zero items (hence
	// the IFNULL).
	//
	// TODO: Can we rewrite this query so that MAX(i.pubdate)
	// can be scanned directly into a date?

	q :=
		`SELECT
			f.title,
			f.author,
			f.desc,
			COUNT(i.ROWID)					AS item_count,
			IFNULL(SUM(i.unplayed), 0)		AS unplayed_count,
			IFNULL(MAX(i.pubdate), ?)		AS max_pubdate
		FROM
			feeds f
		LEFT OUTER JOIN
			items i ON i.feedid = f.id
		WHERE
			` + whereClause + `
		GROUP BY
			f.id
		ORDER BY
			` + sortFields + `
		LIMIT ?`

	var rows []struct {
		Title           string
		Author          string
		Desc            string
		Items           int64
		UnplayedItems   int64
		LastPubdateUnix int64
	}

	err := queryRows(&rows, db, q, append(params, limit)...)
	if err != nil {
		return nil, err
	}

	feeds := make([]feedView, len(rows))

	for i, r := range rows {
		feeds[i] = feedView{
			Title:         r.Title,
			Author:        r.Author,
			Desc:          r.Desc,
			Items:         r.Items,
			UnplayedItems: r.UnplayedItems,
			LastPubdate:   time.Unix(r.LastPubdateUnix, 0),
		}
	}

	return feeds, nil
}

func defaultFeedTemplate(now time.Time) *template.Template {

	layout := `
{{- with len .Feeds -}}
Showing {{.}} feed{{if ne . 1}}s{{end}}:
{{range $index, $feed := $.Feeds}}{{println -}}
         {{$index | plus 1 | printf "%*d" 7}}. {{$feed.Title}}
         {{$feed.Author}}
         {{if $.ShowDesc}}{{range $feed.Desc | lines 70}}{{.}}
         {{end}}{{end -}}
         {{$feed.Items}} item{{if ne $feed.Items 1}}s{{end}}{{with $feed.UnplayedItems}}, {{.}} unplayed{{end}}{{with $feed.LastPubdate}}{{if not .IsZero}} (updated {{. | ago}}){{end}}{{end}}
{{end -}}
{{else -}}
No feeds found
{{end}}`

	funcs := template.FuncMap{
		"lines": formatLines,
		"ago": func(t time.Time) string {
			return timeRelativeTo(now, t)
		},
		"plus": func(i, j int) int {
			return i + j
		},
	}

	return template.Must(
		template.New("feed").Funcs(funcs).Parse(layout),
	)
}

type listItemAction uint

const (
	actionNone listItemAction = iota
	actionPlay
	actionRun
	actionMark
	actionUnmark
)

type listItemOptions struct {
	SortBy    sortItemsBy
	SortOrder sortOrder
	Limit     uint
	Unplayed  bool
	StartDate time.Time
	Title     string
	FeedID    int64
	ShowDesc  bool
	Action    listItemAction
	Use       string
}

func listItems(db *sql.DB, w io.Writer, tmpl *template.Template, opts listItemOptions) error {

	app := opts.Use

	if opts.Action == actionPlay || opts.Action == actionRun {
		_, err := parseCommand(app)
		if err != nil {
			return err
		}
	}

	items, err := loadItemViews(db, opts)
	if err != nil {
		return err
	}

	t := tmpl
	var urls []string

	if opts.Action != actionNone {
		t, err = addCallbackTemplate(t, "prompt", listItemCallback(opts.Action, app, items, &urls))
		if err != nil {
			return err
		}
	}

	err = t.Execute(w, map[string]interface{}{
		"Items":      items,
		"SingleFeed": opts.FeedID != 0,
		"ShowDesc":   opts.ShowDesc,
		"ShowPrompt": opts.Action != actionNone,
	})

	if err != nil {
		switch {
		case strings.Contains(err.Error(), errListItemsDone.Error()):
		case strings.Contains(err.Error(), errListItemsAborted.Error()):
			return nil
		default:
			return err
		}
	}

	if opts.Action == actionNone || len(urls) == 0 {
		return nil
	}

	if opts.Action == actionMark {
		return updatePlayedStatus(db, true, urls...)
	}

	if opts.Action == actionUnmark {
		return updatePlayedStatus(db, false, urls...)
	}

	for _, url := range urls {

		cmd, err := parseCommand(app, url)
		if err != nil {
			return err
		}

		if err := cmd.Run(); err != nil {
			return err
		}

		if opts.Action == actionPlay {
			if err := updatePlayedStatus(db, true, url); err != nil {
				return err
			}
		}
	}

	return nil
}

var errListItemsDone = errors.New("__list_items_exit_loop__")
var errListItemsAborted = errors.New("__list_items_exit_function__")

func listItemCallback(action listItemAction, app string, items []itemView, urls *[]string) func(i int) error {

	var prompt string

	switch action {
	case actionPlay:
		prompt = "Play"
	case actionMark:
		prompt = "Mark as played"
	case actionUnmark:
		prompt = "Unmark as played"
	case actionRun:
		_, name := filepath.Split(strings.Fields(app)[0])
		prompt = "Run " + name
	default:
		return func(int) error {
			return nil
		}
	}

	prompt += "? Yes, No, All, Done, Quit"

	return func(i int) error {

		c, err := ask(prompt, "ynadq")
		if err != nil {
			return err
		}

		switch c {
		case 'y':
			*urls = append(*urls, items[i].url)
		case 'a':
			for i := range items[i:] {
				*urls = append(*urls, items[i].url)
			}
			return errListItemsDone
		case 'd':
			return errListItemsDone
		case 'q':
			return errListItemsAborted
		}

		return nil
	}
}

func addCallbackTemplate(t *template.Template, name string, callback func(int) error) (*template.Template, error) {

	layout := "{{with __callback__ .}}{{end}}"

	funcs := template.FuncMap{
		"__callback__": func(i int) (bool, error) {
			return false, callback(i)
		},
	}

	t, err := t.Clone()
	if err != nil {
		return nil, err
	}

	_, err = t.New(name).Funcs(funcs).Parse(layout)
	if err != nil {
		return nil, err
	}

	return t, nil
}

type itemView struct {
	Title      string
	Desc       string
	Duration   int64
	Pubdate    time.Time
	IsUnplayed bool
	FeedTitle  string
	feedID     int64
	url        string
}

func loadItemViews(db *sql.DB, opts listItemOptions) ([]itemView, error) {

	// ORDER BY clause

	var order string

	switch opts.SortOrder {
	case sortOrderAsc:
		order = "ASC"
	case sortOrderDesc:
		order = "DESC"
	case sortOrderDefault:
		switch opts.SortBy {
		case sortItemsByTitle, sortItemsByFeed:
			order = "ASC"
		default:
			order = "DESC"
		}
	}

	secondaryFields := "i.pubdate DESC, i.ROWID DESC"

	var sortFields string

	switch opts.SortBy {
	case sortItemsByPubdate:
		sortFields = fmt.Sprintf("i.pubdate %s, i.ROWID DESC", order)
	case sortItemsByTitle:
		sortFields = fmt.Sprintf("i.title COLLATE NOCASE %s, %s", order, secondaryFields)
	case sortItemsByFeed:
		sortFields = fmt.Sprintf("f.title COLLATE NOCASE %s, %s", order, secondaryFields)
	case sortItemsByDuration:
		sortFields = fmt.Sprintf("i.duration %s, %s", order, secondaryFields)
	case sortItemsByTimestamp:
		sortFields = fmt.Sprintf("i.timestamp %s, %s", order, secondaryFields)
	default:
		return nil, errors.New("unsupported sort type")
	}

	// WHERE clause

	var params []interface{}
	var conditions []string

	if opts.Unplayed {
		conditions = append(conditions, "i.unplayed = 1")
	}

	if t := opts.StartDate; !t.IsZero() {
		conditions = append(conditions, "i.pubdate >= ?")
		params = append(params, t.Unix())
	}

	if title := opts.Title; title != "" {
		conditions = append(conditions, "i.title LIKE ?")
		params = append(params, "%"+title+"%")
	}

	if id := opts.FeedID; id != 0 {
		conditions = append(conditions, "i.feedid = ?")
		params = append(params, id)
	}

	whereClause := "1=1"
	if len(conditions) > 0 {
		whereClause = strings.Join(conditions, " AND ")
	}

	// Limit

	limit := int(opts.Limit)
	if limit <= 0 {
		limit = -1
	}

	// Build and execute query

	q :=
		`SELECT
			i.title,
			i.desc,
			i.url,
			i.duration,
			i.pubdate,
			i.unplayed,
			f.id,
			f.title
		FROM
			items i
		INNER JOIN
			feeds f ON f.id = i.feedid
		WHERE
			` + whereClause + `
		ORDER BY
			` + sortFields + `
		LIMIT ?`

	var rows []struct {
		Title     string
		Desc      string
		URL       string
		Duration  int64
		Pubdate   time.Time
		Unplayed  bool
		FeedID    int64
		FeedTitle string
	}

	err := queryRows(&rows, db, q, append(params, limit)...)
	if err != nil {
		return nil, err
	}

	items := make([]itemView, len(rows))

	for i, r := range rows {
		items[i] = itemView{
			Title:      r.Title,
			FeedTitle:  r.FeedTitle,
			Desc:       r.Desc,
			Duration:   r.Duration,
			Pubdate:    r.Pubdate,
			IsUnplayed: r.Unplayed,
			url:        r.URL,
			feedID:     r.FeedID,
		}
	}

	return items, nil
}

func defaultItemTemplate(now time.Time) *template.Template {

	layout := `
{{- with len .Items -}}
Showing {{.}} item{{if ne . 1}}s{{end}}:
{{range $index, $item := $.Items}}{{println -}}
         {{if $item.IsUnplayed}}{{$index | plus 1 | printf "* %d" | printf "%*s" 7}}{{else}}{{$index | plus 1 | printf "%*d" 7}}{{end}}. {{$item.Title}}
         {{if $.SingleFeed}}Released{{else}}From {{$item.FeedTitle}},{{end}} {{if $item.Pubdate.IsZero}}Date unknown{{else}}{{$item.Pubdate.Local | ago}}{{end}}
         {{if $.ShowDesc}}{{range $item.Desc | lines 70}}{{.}}
         {{end}}{{end -}}
         Duration: {{with $item.Duration}}{{. | duration}}{{else}}Unknown{{end}}{{if $.ShowPrompt}}
         {{template "prompt" $index}}{{else}}{{println}}{{end -}}
{{end -}}
{{else -}}
No items found
{{end}}`

	funcs := template.FuncMap{
		"lines": formatLines,
		"ago": func(t time.Time) string {
			return timeRelativeTo(now, t)
		},
		"plus": func(i, j int) int {
			return i + j
		},
		"duration": formatSeconds,
	}

	return template.Must(
		template.New("item").Funcs(funcs).Parse(layout),
	)
}

func exportList(db *sql.DB, w io.Writer) error {

	q := `SELECT url FROM feeds ORDER BY url COLLATE NOCASE`

	var rows []struct {
		URL string
	}

	err := queryRows(&rows, db, q)
	if err != nil {
		return err
	}

	if len(rows) == 0 {
		return errors.New("no feeds to export")
	}

	for _, r := range rows {
		fmt.Fprintln(w, r.URL)
	}

	return nil
}

func exportOPML(db *sql.DB, w io.Writer, timestamp time.Time) error {

	q := `SELECT type, title, desc, url FROM feeds ORDER BY title COLLATE NOCASE, url COLLATE NOCASE`

	var rows []struct {
		Type  string
		Title string
		Desc  string
		URL   string
	}

	err := queryRows(&rows, db, q)
	if err != nil {
		return err
	}

	if len(rows) == 0 {
		return errors.New("no feeds to export")
	}

	data := newOPML("Kibner subscriptions")
	data.Pubdate = opmltime(timestamp)
	data.Entries = make([]*entry, len(rows))

	for i, r := range rows {
		data.Entries[i] = &entry{
			Type:  r.Type,
			Text:  r.Title,
			Title: r.Title,
			Desc:  r.Desc,
			URL:   r.URL,
		}
	}

	return data.writeTo(w)
}

var errInvalidTarget = errors.New("invalid target")

func loadFeedURL(db *sql.DB, feedID int64, target target) (string, error) {

	var url string
	var image, link sql.NullString

	q := "SELECT url, link, image FROM feeds WHERE id = ?"

	err := db.QueryRow(q, feedID).Scan(&url, &link, &image)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errNoFeedFound
		}
		return "", err
	}

	switch target {
	case targetFeed:
		return url, nil
	case targetLink:
		return link.String, nil
	case targetImage:
		return image.String, nil
	default:
		return "", errInvalidTarget
	}
}

func syncItems(db *sql.DB, info *syncInfo, items []*kibner.Item) (int, error) {

	var newItems []*kibner.Item

	for _, item := range items {
		if !info.guids[item.GUID] {
			newItems = append(newItems, item)
		}
	}

	if len(newItems) == 0 {
		return 0, nil
	}

	return len(newItems), saveNewItems(db, info.ID, newItems)
}

func syncFeed(db *sql.DB, info *syncInfo, feed *kibner.Feed) error {

	if info.URL == feed.URL {
		return nil
	}

	return updateFeed(db, info.ID, map[string]interface{}{
		"url": feed.URL,
	})
}

type syncInfo struct {
	ID    int64
	Title string
	URL   string
	guids map[string]bool
}

func loadSyncInfo(db *sql.DB, id int64) ([]syncInfo, error) {

	infos, err := loadSyncInfoFeeds(db, id)
	if err != nil {
		return nil, err
	}

	guids, err := loadSyncInfoGUIDs(db, id)
	if err != nil {
		return nil, err
	}

	for i := range infos {
		info := &infos[i]
		for _, g := range guids[info.ID] {
			if info.guids == nil {
				info.guids = make(map[string]bool)
			}
			info.guids[g] = true
		}
	}

	return infos, nil
}

func loadSyncInfoFeeds(db *sql.DB, id int64) ([]syncInfo, error) {

	var rows []syncInfo
	var params []interface{}

	q := "SELECT id, title, url FROM feeds"

	if id > 0 {
		q += " WHERE id = ?"
		params = append(params, id)
	}

	if err := queryRows(&rows, db, q, params...); err != nil {
		return nil, err
	}

	return rows, nil
}

func loadSyncInfoGUIDs(db *sql.DB, id int64) (map[int64][]string, error) {

	var rows []struct {
		ID   int64
		GUID string
	}
	var params []interface{}

	q := "SELECT feedid, guid FROM items"

	if id > 0 {
		q += " WHERE feedid = ?"
		params = append(params, id)
	}

	if err := queryRows(&rows, db, q, params...); err != nil {
		return nil, err
	}

	var guids map[int64][]string

	for _, r := range rows {
		if guids == nil {
			guids = make(map[int64][]string)
		}
		guids[r.ID] = append(guids[r.ID], r.GUID)
	}

	return guids, nil
}

func extractURLsList(r io.Reader) ([]string, error) {

	var urls []string

	seen := map[string]bool{}
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {

		u := strings.TrimSpace(scanner.Text())

		if u == "" || u[0] == '#' || seen[u] {
			continue
		}

		if _, err := url.Parse(u); err != nil {
			continue
		}

		seen[u] = true
		urls = append(urls, u)
	}

	return urls, scanner.Err()
}

func extractURLsOPML(r io.Reader) ([]string, error) {

	var data opml
	if err := data.readFrom(r); err != nil {
		return nil, err
	}

	seen := map[string]bool{}
	urls := make([]string, 0, len(data.Entries))

	for _, entry := range data.Entries {

		u := strings.TrimSpace(entry.URL)

		if u == "" || seen[u] {
			continue
		}

		if _, err := url.Parse(u); err != nil {
			continue
		}

		seen[u] = true
		urls = append(urls, u)
	}

	return urls, nil
}

func rollback(tx *sql.Tx, err error) error {
	// TODO: Handle rollback errors
	tx.Rollback()
	return err
}

func newRequest(url string) (*http.Request, error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "kibner/"+version)
	return req, nil
}

func fetchAndParse(feedURL string) (*kibner.Feed, error) {

	req, err := newRequest(feedURL)
	if err != nil {
		return nil, errors.New("bad request: " + err.Error())
	}

	resp, err := defaultClient.Do(req)
	if err != nil {
		return nil, errors.New("fetch error: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("bad status: " + resp.Status)
	}

	parser := gofeed.NewParser()
	parser.RSSTranslator = NewRSSTranslator()
	f, err := parser.Parse(resp.Body)
	if err != nil {
		return nil, errors.New("parse error: " + err.Error())
	}

	var newFeedURL string
	if f.ITunesExt != nil {
		newFeedURL = f.ITunesExt.NewFeedURL
	}
	if newFeedURL != "" && newFeedURL != feedURL {
		return fetchAndParse(newFeedURL)
	}

	feed := translateFeed(f)
	if feed.Title == "" {
		return nil, errors.New("bad feed: no title")
	}

	feed.URL = resp.Request.URL.String()
	if feed.URL == "" {
		feed.URL = feedURL
	}

	return feed, nil
}

func translateFeed(f *gofeed.Feed) *kibner.Feed {

	var author string
	if f.Author != nil {
		author = f.Author.Name
	}
	if author == "" {
		author = "Unknown Author"
	}

	var image string
	if f.Image != nil {
		image = f.Image.URL
	}

	return &kibner.Feed{
		Title:  f.Title,
		Author: author,
		Desc:   f.Description,
		Type:   f.FeedType,
		Link:   f.Link,
		Image:  image,
		Items:  translateItems(f.Items),
	}
}

func translateItems(items []*gofeed.Item) []*kibner.Item {

	feedItems := make([]*kibner.Item, 0, len(items))

	for _, item := range items {
		newItem := translateItem(item)
		if newItem != nil {
			feedItems = append(feedItems, newItem)
		}
	}

	return feedItems
}

func translateItem(item *gofeed.Item) *kibner.Item {

	url, filesize := translateItemEnclosures(item.Enclosures)
	if url == "" {
		// Ignore items with no download URL.
		// TODO: Log/output this error.
		return nil
	}

	pubdate := translateItemPubdate(item)

	title := item.Title
	if title == "" {
		title = "Untitled"
		if !pubdate.IsZero() {
			title += ": " + pubdate.Format("Jan 2, 2006")
		}
	}

	guid := item.GUID
	if guid == "" {
		guid = url
	}

	return &kibner.Item{
		Title:    title,
		Pubdate:  pubdate,
		Desc:     item.Description,
		URL:      url,
		Filesize: filesize,
		Duration: translateItemDuration(item),
		GUID:     guid,
	}
}

func translateItemEnclosures(encs []*gofeed.Enclosure) (string, int64) {

	if len(encs) == 0 {
		return "", 0
	}

	// TODO: Handle multiple enclosures?
	enc := encs[0]

	filesize, err := strconv.ParseInt(enc.Length, 10, 64)
	if err != nil {
		// Swallow this error.
		// No big deal if we don't have the filesize.
		// TODO: Log/output this error.
		filesize = 0
	}

	return enc.URL, filesize
}

func translateItemPubdate(item *gofeed.Item) time.Time {

	if item.PublishedParsed == nil {
		// Date was in an invalid format and could not be parsed.
		// Swallow this error.
		// TODO: Log/output this error.
		// TODO: Check for common typos? e.g. Tues, Thur, Sept
		return time.Time{}
	}

	return *item.PublishedParsed
}

func translateItemDuration(item *gofeed.Item) time.Duration {

	if item.ITunesExt == nil {
		return 0
	}

	duration, err := parseDuration(item.ITunesExt.Duration)
	if err != nil {
		// Swallow this error.
		// No big deal if we don't have the duration.
		// TODO: Log/output this error.
		return 0
	}

	return duration
}

func fetchAndParseMultiple(urls []string, maxWorkers int) (map[string]*kibner.Feed, map[string]error) {

	feeds := make(chan struct {
		string
		*kibner.Feed
	})
	errors := make(chan struct {
		string
		error
	})
	workers := make(chan struct{}, maxWorkers)

	for i := range urls {
		go func(url string) {

			workers <- struct{}{}

			feed, err := fetchAndParse(url)
			switch {
			case err != nil:
				errors <- struct {
					string
					error
				}{url, err}
			default:
				feeds <- struct {
					string
					*kibner.Feed
				}{url, feed}
			}
		}(urls[i])
	}

	oks := map[string]*kibner.Feed{}
	errs := map[string]error{}

	for i := range urls {

		select {
		case f := <-feeds:
			oks[f.string] = f.Feed
		case e := <-errors:
			errs[e.string] = e.error
		}

		<-workers
		fmt.Printf("Fetched %d of %d feeds...\r", i+1, len(urls))
	}

	return oks, errs
}

func deleteItems(tx *sql.Tx, feedID int64) error {

	_, err := tx.Exec("DELETE FROM items WHERE feedid = ?", feedID)
	return err
}

func deleteFeed(tx *sql.Tx, feedID int64) error {

	_, err := tx.Exec("DELETE FROM feeds WHERE id = ?", feedID)
	return err
}

func saveFeed(db *sql.DB, feed *kibner.Feed, timestamp time.Time) (int64, error) {

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	feedID, err := insertFeed(tx, feed, timestamp)
	if err != nil {
		return 0, rollback(tx, err)
	}

	err = insertItems(tx, feedID, feed.Items, false, timestamp)
	if err != nil {
		return 0, rollback(tx, err)
	}

	return feedID, tx.Commit()
}

func saveNewItems(db *sql.DB, feedID int64, items []*kibner.Item) error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = insertItems(tx, feedID, items, true, time.Now())
	if err != nil {
		return rollback(tx, err)
	}

	return tx.Commit()
}

func updatePlayedStatus(db *sql.DB, played bool, urls ...string) error {

	params := make([]interface{}, len(urls)+1)
	params[0] = !played
	placeholders := make([]string, len(urls))

	for i := range urls {
		params[i+1] = urls[i]
		placeholders[i] = "?"
	}

	q := fmt.Sprintf("UPDATE items SET unplayed = ? WHERE url IN(%s)", strings.Join(placeholders, ", "))

	_, err := db.Exec(q, params...)
	return err
}

func updateFeed(db *sql.DB, feedID int64, values map[string]interface{}) error {

	if len(values) == 0 {
		return errors.New("no values provided")
	}

	setters := make([]string, 0, len(values))
	params := make([]interface{}, 0, len(values)+1)

	for field, value := range values {
		setters = append(setters, fmt.Sprintf("%s = ?", field))
		params = append(params, value)
	}

	q := fmt.Sprintf("UPDATE feeds SET %s WHERE id = ?", strings.Join(setters, ", "))

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	res, err := tx.Exec(q, append(params, feedID)...)
	if err != nil {
		return rollback(tx, err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return rollback(tx, err)
	}

	switch n {
	case 1:
		return tx.Commit()
	case 0:
		return rollback(tx, fmt.Errorf("feed not found"))
	default:
		return rollback(tx, fmt.Errorf("%d rows affected, expected 1", n))
	}
}

func insertFeed(tx *sql.Tx, feed *kibner.Feed, timestamp time.Time) (int64, error) {

	sql :=
		`INSERT INTO feeds(
			title,
			author,
			desc,
			url,
			type,
			link,
			image,
			timestamp
		) VALUES(?, ?, ?, ?, ?, ?, ?, ?)`

	res, err := tx.Exec(sql, normaliseText(feed.Title), normaliseText(feed.Author), normaliseText(feed.Desc), feed.URL, feed.Type, feed.Link, feed.Image, timestamp.Unix())
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func insertItems(tx *sql.Tx, feedid int64, items []*kibner.Item, unplayed bool, timestamp time.Time) error {

	sql :=
		`INSERT INTO items(
			feedid,
			unplayed,
			title,
			desc,
			pubdate,
			url,
			filesize,
			duration,
			guid,
			timestamp
		) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	stmt, err := tx.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, item := range items {

		_, err := stmt.Exec(feedid, unplayed, normaliseText(item.Title), normaliseText(item.Desc), item.Pubdate.Unix(), item.URL, item.Filesize, item.Duration.Seconds(), item.GUID, timestamp.Unix())
		if err != nil {
			return err
		}
	}

	return nil
}

func normaliseText(desc string) string {
	if desc == "" {
		return ""
	}
	return replaceWhitespace(replaceTags(desc, ""), " ")
}

var rxTags = []*regexp.Regexp{

	// Opening HTML tag: <p> or <a href="...">
	// <                    Opening bracket
	// [0-9A-Za-z:]+        One or more alphanumeric characters
	// (\s+(?s:.+?))?       Optional text (non-greedy, including line breaks) preceded by whitespace
	// >                    Closing bracket
	regexp.MustCompile(`<[0-9A-Za-z:]+(\s+(?s:.+?))?>`),

	// Closing HTML tag: </p>
	// </                   Opening bracket + slash
	// [0-9A-Za-z:]+        One or more alphanumeric characters
	// >                    Closing bracket
	regexp.MustCompile(`</[0-9A-Za-z:]+>`),

	// Self-closing HTML tag: <br />
	// <                    Opening bracket
	// [0-9A-Za-z:]+        One or more alphanumeric characters
	// \s?                  Optional whitespace
	// />                   Closing slash + bracket
	regexp.MustCompile(`<[0-9A-Za-z:]+\s?/>`),
}

func replaceTags(s, repl string) string {
	for _, rx := range rxTags {
		s = rx.ReplaceAllString(s, repl)
	}
	return s
}

var rxSpaces = []*regexp.Regexp{
	regexp.MustCompile(`\x{00a0}`), // \u00a0 is the unicode code point for nbsp
	regexp.MustCompile(`&nbsp;`),
	regexp.MustCompile(`\s+`),
}

func replaceWhitespace(s, repl string) string {
	for _, rx := range rxSpaces {
		s = rx.ReplaceAllString(s, repl)
	}
	return strings.TrimSpace(s)
}

func formatSeconds(secs int64) string {

	if secs < 60 {
		return fmt.Sprintf("%ds", secs)
	}

	div := func(i, j int64) (int64, int64) {
		return i / j, i % j
	}

	s := ""
	h, rem := div(secs, 60*60)
	if h > 0 {
		s += fmt.Sprintf("%dh", h)
	}

	m, rem := div(rem, 60)
	if rem >= 30 {
		m++
	}
	if m > 0 {
		s += fmt.Sprintf("%dm", m)
	}

	return s
}

var errNoFeedChosen = errors.New("No feed selected")
var errNoFeedFound = errors.New("No such feed")

func chooseFeed(db *sql.DB, feedTitle string, prompt string, alwaysPrompt bool) (int64, error) {

	q := `SELECT id, title FROM feeds WHERE title LIKE ? ORDER BY title`

	var rows []struct {
		ID    int64
		Title string
	}

	err := queryRows(&rows, db, q, "%"+feedTitle+"%")
	if err != nil {
		return 0, err
	}

	if len(rows) == 0 {
		return 0, errNoFeedFound
	}

	if len(rows) == 1 && !alwaysPrompt {
		return rows[0].ID, nil
	}

	for i, r := range rows {

		s := fmt.Sprintf("[%d/%d] %s? Yes, No, Quit", i+1, len(rows), fmt.Sprintf(prompt, r.Title))
		c, err := ask(s, "ynq")
		if err != nil {
			return 0, err
		}

		if c == 'q' {
			break
		}

		if c == 'y' {
			return r.ID, nil
		}
	}

	return 0, errNoFeedChosen
}

func ask(prompt string, responses string) (rune, error) {

	responses = strings.ToLower(responses)

	for {
		fmt.Printf("%s: ", prompt)
		c, _, err := bufio.NewReader(os.Stdin).ReadRune()
		if err != nil {
			return 0, err
		}

		if c = unicode.ToLower(c); strings.ContainsRune(responses, c) {
			return c, nil
		}
	}
}

func timeRelativeTo(t0, t time.Time) string {

	days := daysFrom(t0, t)
	if days < 0 {
		days = 0 - days
	}

	if days == 0 {
		return "Today"
	}

	if days == 1 {
		return "Yesterday"
	}

	for month := 1; month < 12; month++ {

		interval := daysFrom(t0.AddDate(0, -month, 0), t0)

		if days < interval {
			switch month {
			case 1:
				return fmt.Sprintf("%d days ago", days)
			case 2:
				return "over a month ago"
			default:
				return fmt.Sprintf("over %d months ago", month-1)
			}
		}

		if days == interval {
			if month == 1 {
				return "a month ago"
			}
			return fmt.Sprintf("%d months ago", month)
		}
	}

	aYearAgo := t0.AddDate(-1, 0, 0)
	interval := daysFrom(aYearAgo, t0)

	switch {
	case days < interval:
		return "over 11 months ago"
	case days == interval:
		return "a year ago"
	case t.Year() == aYearAgo.Year():
		return "over a year ago"
	default:
		return fmt.Sprintf("in %d", t.Year())
	}
}

func startOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func daysFrom(refdate, t time.Time) int {
	t = t.In(refdate.Location())
	return int(startOfDay(t).Sub(startOfDay(refdate)).Seconds() / (60 * 60 * 24))
}

func parseDuration(hhmmss string) (time.Duration, error) {

	var units []time.Duration
	var duration time.Duration

	parts := strings.Split(hhmmss, ":")

	switch len(parts) {
	case 1:
		units = []time.Duration{time.Second}
	case 2:
		units = []time.Duration{time.Minute, time.Second}
	case 3:
		units = []time.Duration{time.Hour, time.Minute, time.Second}
	default:
		return 0, errors.New("invalid duration format: " + hhmmss)
	}

	for i := range parts {
		n, err := strconv.Atoi(parts[i])
		if err != nil {
			return 0, errors.New("could not parse duration " + hhmmss + ": " + err.Error())
		}
		duration += time.Duration(n) * units[i]
	}

	return duration, nil
}

func formatLines(linelen int, s string) []string {

	var lines []string

	for {
		if len(s) <= linelen {
			if len(s) > 0 {
				lines = append(lines, s)
			}
			break
		}

		pos := linelen
		for pos > 0 && s[pos] > ' ' {
			pos--
		}

		if pos == 0 {
			pos = linelen
		}

		lines = append(lines, s[:pos])

		for pos < len(s) && s[pos] <= ' ' {
			pos++
		}

		if pos >= len(s) {
			break
		}

		s = s[pos:]
	}

	return lines
}

func parseCommand(command string, args ...string) (*exec.Cmd, error) {

	parts := strings.Fields(command)
	if len(parts) < 1 {
		return nil, errors.New("no command given")
	}

	exe, err := exec.LookPath(parts[0])
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(exe, append(parts[1:], args...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd, nil
}
