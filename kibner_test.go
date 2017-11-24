package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/deepilla/kibner/internal/testdata"
	kibner "github.com/deepilla/kibner/internal/types"

	"github.com/deepilla/sqlitemeta"
)

type testCase struct {
	Filename string
	NewFeed  func() *kibner.Feed
}

var allTestCases = map[string]*testCase{

	"Constitutional": {
		Filename: "constitutional.xml",
		NewFeed:  testdata.Constitutional,
	},
	"Crooked Conversations": {
		Filename: "crooked-conversations.xml",
		NewFeed:  testdata.Crooked,
	},
	"Ear Hustle": {
		Filename: "ear-hustle.xml",
		NewFeed:  testdata.EarHustle,
	},
	"Gadget Lab": {
		Filename: "gadget-lab.xml",
		NewFeed:  testdata.GadgetLab,
	},
	"Homecoming": {
		Filename: "homecoming.xml",
		NewFeed:  testdata.Homecoming,
	},
	"Mogul": {
		Filename: "mogul.xml",
		NewFeed:  testdata.Mogul,
	},
	"No Such Thing As A Fish": {
		Filename: "fish.xml",
		NewFeed:  testdata.Fish,
	},
	"Rabbits": {
		Filename: "rabbits.xml",
		NewFeed:  testdata.Rabbits,
	},
	"Revisionist History": {
		Filename: "revisionist-history.xml",
		NewFeed:  testdata.RevisionistHistory,
	},
	"Rework": {
		Filename: "rework.xml",
		NewFeed:  testdata.Rework,
	},
	"S-Town": {
		Filename: "s-town.xml",
		NewFeed:  testdata.STown,
	},
	"Serial": {
		Filename: "serial.xml",
		NewFeed:  testdata.Serial,
	},
	"The Turnaround": {
		Filename: "turnaround.xml",
		NewFeed:  testdata.Turnaround,
	},
	"Uncivil": {
		Filename: "uncivil.xml",
		NewFeed:  testdata.Uncivil,
	},
	"What's Good": {
		Filename: "whats-good.xml",
		NewFeed:  testdata.WhatsGood,
	},
	"Why We Eat": {
		Filename: "why-we-eat.xml",
		NewFeed:  testdata.WhyWeEat,
	},
	"Why We Eat 2": {
		Filename: "why-we-eat-2.xml",
		NewFeed:  testdata.WhyWeEat2,
	},

	"Columbo": {
		Filename: "columbo.xml",
		NewFeed:  testdata.Columbo,
	},

	"Empty Feed": {
		Filename: "test-feed-empty.xml",
		NewFeed:  testdata.EmptyFeed,
	},
	"Minimal Feed": {
		Filename: "test-feed-minimal.xml",
		NewFeed:  testdata.MinimalFeed,
	},
}

func newFileServer() *httptest.Server {
	return httptest.NewServer(http.FileServer(http.Dir("internal/testdata/rss")))
}

func newNotFoundServer() *httptest.Server {
	return httptest.NewServer(http.NotFoundHandler())
}

func newTestDB(t *testing.T) *sql.DB {

	opts := map[string]string{
		"cache":         "shared",
		"_foreign_keys": "1",
	}

	db, err := connect(":memory:", opts)
	if err != nil {
		t.Fatalf("sql.Open returned error %q", err)
	}

	verifyTables(t, db, nil)
	verifyForeignKeysEnabled(t, db, true)
	return db
}

func testWithDB(t *testing.T, fn func(t *testing.T, db *sql.DB)) {
	db := newTestDB(t)
	defer db.Close()
	fn(t, db)
}

func testWithInitDB(t *testing.T, fn func(t *testing.T, db *sql.DB)) {
	testWithDB(t, func(t *testing.T, db *sql.DB) {
		testInitDB(t, db)
		fn(t, db)
	})
}

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())
	os.Exit(m.Run())
}

func TestInitDB(t *testing.T) {
	testWithDB(t, testInitDB)
}

func testInitDB(t *testing.T, db *sql.DB) {

	if err := initDB(db); err != nil {
		t.Fatalf("initDB returned error %q", err)
	}

	verifyTables(t, db, map[string]int{
		"feeds":           0,
		"items":           0,
		"sqlite_sequence": 0,
	})

	verifyColumns(t, db, "feeds", []sqlitemeta.Column{
		{
			ID:         0,
			Name:       "id",
			Type:       "INTEGER",
			PrimaryKey: 1,
		},
		{
			ID:      1,
			Name:    "title",
			Type:    "TEXT",
			NotNull: true,
		},
		{
			ID:      2,
			Name:    "author",
			Type:    "TEXT",
			NotNull: true,
		},
		{
			ID:   3,
			Name: "desc",
			Type: "TEXT",
		},
		{
			ID:      4,
			Name:    "type",
			Type:    "TEXT",
			NotNull: true,
		},
		{
			ID:      5,
			Name:    "url",
			Type:    "TEXT",
			NotNull: true,
		},
		{
			ID:   6,
			Name: "image",
			Type: "TEXT",
		},
		{
			ID:   7,
			Name: "link",
			Type: "TEXT",
		},
		{
			ID:      8,
			Name:    "timestamp",
			Type:    "DATETIME",
			NotNull: true,
		},
	})

	verifyIndexes(t, db, "feeds", []sqlitemeta.Index{
		{
			Name:        "unique_feed_url",
			Type:        sqlitemeta.IndexTypeNormal,
			IsUnique:    true,
			ColumnNames: nullStrings("url"),
		},
	})

	verifyForeignKeys(t, db, "feeds", nil)

	verifyColumns(t, db, "items", []sqlitemeta.Column{
		{
			ID:      0,
			Name:    "feedid",
			Type:    "INTEGER",
			NotNull: true,
		},
		{
			ID:      1,
			Name:    "title",
			Type:    "TEXT",
			NotNull: true,
		},
		{
			ID:   2,
			Name: "desc",
			Type: "TEXT",
		},
		{
			ID:      3,
			Name:    "pubdate",
			Type:    "DATETIME",
			NotNull: true,
		},
		{
			ID:      4,
			Name:    "url",
			Type:    "TEXT",
			NotNull: true,
		},
		{
			ID:      5,
			Name:    "filesize",
			Type:    "INTEGER",
			Default: []byte("0"),
		},
		{
			ID:      6,
			Name:    "duration",
			Type:    "INTEGER",
			Default: []byte("0"),
		},
		{
			ID:      7,
			Name:    "guid",
			Type:    "TEXT",
			NotNull: true,
		},
		{
			ID:      8,
			Name:    "unplayed",
			Type:    "BOOLEAN",
			Default: []byte("0"),
		},
		{
			ID:      9,
			Name:    "timestamp",
			Type:    "DATETIME",
			NotNull: true,
		},
	})

	verifyIndexes(t, db, "items", []sqlitemeta.Index{
		{
			Name:        "unique_item_guid",
			Type:        sqlitemeta.IndexTypeNormal,
			IsUnique:    true,
			ColumnNames: nullStrings("feedid", "guid"),
		},
	})

	verifyForeignKeys(t, db, "items", []sqlitemeta.ForeignKey{
		{
			ID:         0,
			ChildTable: "items",
			ChildKey: []string{
				"feedid",
			},
			ParentTable: "feeds",
			ParentKey:   nullStrings("id"),
			OnUpdate:    sqlitemeta.ForeignKeyActionNone,
			OnDelete:    sqlitemeta.ForeignKeyActionNone,
		},
	})
}

func TestAddFeed(t *testing.T) {
	testWithInitDB(t, testAddFeed)
}

func testAddFeed(t *testing.T, db *sql.DB) {

	ts := newFileServer()
	defer ts.Close()

	items, feeds := 0, 0

	for name, test := range allTestCases {

		url := serverURL(ts, test.Filename)

		feed := test.NewFeed()
		feed.URL = url

		res, err := addFeed(db, url)
		if err != nil {
			t.Fatalf("%s: addFeed returned error %q", name, err)
			continue
		}

		verifyAddResult(t, res, feed)
		verifyFeed(t, db, res.ID, feed)
		verifyUnplayedItemCount(t, db, res.ID, 0)

		feeds++
		items += res.Items

		verifyTables(t, db, map[string]int{
			"feeds":           feeds,
			"items":           items,
			"sqlite_sequence": 1,
		})

		verifySequenceNumbers(t, db, map[string]int64{
			"feeds": res.ID,
		})
	}
}

func TestAddFeedNotFound(t *testing.T) {
	testWithInitDB(t, testAddFeedNotFound)
}

func testAddFeedNotFound(t *testing.T, db *sql.DB) {

	ts := newNotFoundServer()
	defer ts.Close()

	exp := errors.New("could not fetch feed: bad status: 404 Not Found")

	for name, test := range allTestCases {

		url := serverURL(ts, test.Filename)

		if _, err := addFeed(db, url); !equalErrors(exp, err) {
			t.Errorf("%s: expected addFeed to return error %q, got %v", name, exp, err)
		}
	}
}

func TestAddFeedDuplicate(t *testing.T) {
	testWithInitDB(t, testAddFeedDuplicate)
}

func testAddFeedDuplicate(t *testing.T, db *sql.DB) {

	ts := newFileServer()
	defer ts.Close()

	exps := []error{
		nil,
		errors.New("could not save feed: UNIQUE constraint failed: feeds.url"),
	}

	for name, test := range allTestCases {
		for _, exp := range exps {

			url := serverURL(ts, test.Filename)

			if _, err := addFeed(db, url); !equalErrors(exp, err) {
				t.Errorf("%s: expected addFeed to return error %v, got %v", name, exp, err)
			}
		}
	}
}

func TestAddFeedDuplicateGUID(t *testing.T) {
	testWithInitDB(t, testAddFeedDuplicateGUID)
}

func testAddFeedDuplicateGUID(t *testing.T, db *sql.DB) {

	ts := newFileServer()
	defer ts.Close()

	url := serverURL(ts, "errors/duplicate-guid.xml")
	exp := errors.New("could not save feed: UNIQUE constraint failed: items.feedid, items.guid")

	if _, err := addFeed(db, url); !equalErrors(exp, err) {
		t.Errorf("expected addFeed to return %v, got %v", exp, err)
	}
}

func TestAddFeedNoTitle(t *testing.T) {
	testWithInitDB(t, testAddFeedNoTitle)
}

func testAddFeedNoTitle(t *testing.T, db *sql.DB) {

	ts := newFileServer()
	defer ts.Close()

	url := serverURL(ts, "errors/no-title.xml")
	exp := errors.New("could not fetch feed: bad feed: no title")

	if _, err := addFeed(db, url); !equalErrors(exp, err) {
		t.Errorf("expected addFeed to return %q, got %v", exp, err)
	}
}

func TestAddFeedMultiple(t *testing.T) {
	testWithInitDB(t, testAddFeedMultiple)
}

func testAddFeedMultiple(t *testing.T, db *sql.DB) {

	ts := newFileServer()
	defer ts.Close()

	urls := make([]string, 0, len(allTestCases))
	feedsByURL := map[string]*kibner.Feed{}

	for _, test := range allTestCases {

		url := serverURL(ts, test.Filename)
		urls = append(urls, url)

		feed := test.NewFeed()
		feed.URL = url
		feedsByURL[url] = feed
	}

	results := addFeedMultiple(db, urls)

	if len(results) != len(allTestCases) {
		t.Fatalf("expected %d results from addFeedMultiple, got %d", len(allTestCases), len(results))
	}

	maxID := int64(0)
	feeds, items := 0, 0

	for _, res := range results {

		feed, ok := feedsByURL[res.URL]
		if !ok {
			t.Fatalf("addFeedMultiple returned unexpected URL %q", res.URL)
		}

		if res.Err != nil {
			t.Fatalf("%s: addFeedMultiple returned error %q", feed.Title, res.Err)
			continue
		}

		verifyAddResult(t, res, feed)
		verifyFeed(t, db, res.ID, feed)
		verifyUnplayedItemCount(t, db, res.ID, 0)

		feeds++
		items += res.Items
		if res.ID > maxID {
			maxID = res.ID
		}
	}

	verifyTables(t, db, map[string]int{
		"feeds":           feeds,
		"items":           items,
		"sqlite_sequence": 1,
	})

	verifySequenceNumbers(t, db, map[string]int64{
		"feeds": maxID,
	})
}

func TestAddFeedMultipleNotFound(t *testing.T) {
	testWithInitDB(t, testAddFeedMultipleNotFound)
}

func testAddFeedMultipleNotFound(t *testing.T, db *sql.DB) {

	ts := newNotFoundServer()
	defer ts.Close()

	urls := make([]string, 0, len(allTestCases))
	namesByURL := map[string]string{}

	for name, test := range allTestCases {

		url := serverURL(ts, test.Filename)

		urls = append(urls, url)
		namesByURL[url] = name
	}

	results := addFeedMultiple(db, urls)

	if len(results) != len(allTestCases) {
		t.Fatalf("expected %d results from addFeedMultiple, got %d", len(allTestCases), len(results))
	}

	exp := errors.New("bad status: 404 Not Found")

	for _, res := range results {

		name, ok := namesByURL[res.URL]
		if !ok {
			t.Fatalf("addFeedMultiple returned unexpected URL: %s", res.URL)
		}

		if !equalErrors(exp, res.Err) {
			t.Errorf("%s: expected addFeedMultiple to return error %v, got %v", name, exp, res.Err)
		}
	}
}

func TestRemoveFeed(t *testing.T) {
	testWithInitDB(t, testRemoveFeed)
}

func testRemoveFeed(t *testing.T, db *sql.DB) {

	ts := newFileServer()
	defer ts.Close()

	maxID := int64(0)
	feeds, items := 0, 0

	results := make([]*syncResult, 0, len(allTestCases))

	for name, test := range allTestCases {

		res, err := addFeed(db, serverURL(ts, test.Filename))
		if err != nil {
			t.Fatalf("%s: addFeed returned error %q", name, err)
		}

		feeds++
		items += res.Items
		maxID = res.ID
		results = append(results, res)
	}

	for _, res := range results {

		err := removeFeed(db, res.ID)
		if err != nil {
			t.Fatalf("%s: removeFeed returned error %q", res.Title, err)
		}

		feeds--
		items -= res.Items

		verifyTables(t, db, map[string]int{
			"feeds":           feeds,
			"items":           items,
			"sqlite_sequence": 1,
		})

		verifySequenceNumbers(t, db, map[string]int64{
			"feeds": maxID,
		})
	}
}

func TestSyncFeed(t *testing.T) {
	testWithDB(t, testSyncFeed)
}

func testSyncFeed(t *testing.T, db *sql.DB) {

	ts := newFileServer()
	defer ts.Close()

	now := time.Now()

	for name, test := range allTestCases {

		feed := test.NewFeed()
		feed.URL = serverURL(ts, test.Filename)

		for unplayed := 0; unplayed <= len(feed.Items); unplayed++ {

			testInitDB(t, db)

			feed2 := &kibner.Feed{
				Title: feed.Title,
				URL:   feed.URL,
				Items: feed.Items[unplayed:],
			}

			id, err := saveFeed(db, feed2, now)
			if err != nil {
				t.Fatalf("%s: saveFeed returned error %q", name, err)
			}

			verifyTables(t, db, map[string]int{
				"feeds":           1,
				"items":           len(feed2.Items),
				"sqlite_sequence": 1,
			})
			verifyUnplayedItemCount(t, db, id, 0)
			verifyFeed(t, db, id, feed2)

			res, err := syncOne(db, id)
			if err != nil {
				t.Fatalf("%s: syncOne returned error %s", name, err)
			}

			if res.ID != id {
				t.Errorf("%s: expected syncOne to return id %d, got %d", name, id, res.ID)
			}

			if res.Items != unplayed {
				t.Errorf("%s: expected syncOne to return %d unplayed items, got %d", name, unplayed, res.Items)
			}

			verifyTables(t, db, map[string]int{
				"feeds":           1,
				"items":           len(feed.Items),
				"sqlite_sequence": 1,
			})
			verifyUnplayedItemCount(t, db, id, unplayed)

			feed2.Items = feed.Items
			verifyFeed(t, db, id, feed2)
		}
	}
}

func TestSyncAll(t *testing.T) {
	testWithInitDB(t, testSyncAll)
}

func testSyncAll(t *testing.T, db *sql.DB) {

	ts := newFileServer()
	defer ts.Close()

	now := time.Now()

	feeds, items := 0, 0
	unplayedByID := map[int64]int{}

	for name, test := range allTestCases {

		feed := test.NewFeed()
		feed.URL = serverURL(ts, test.Filename)

		unplayed := 0
		if n := len(feed.Items); n > 0 {
			unplayed = rand.Intn(n)
		}

		feed2 := &kibner.Feed{
			URL:   feed.URL,
			Items: feed.Items[unplayed:],
		}

		id, err := saveFeed(db, feed2, now)
		if err != nil {
			t.Fatalf("%s: saveFeed returned error %q", name, err)
		}

		feeds++
		items += len(feed2.Items)
		unplayedByID[id] = unplayed

		verifyUnplayedItemCount(t, db, id, 0)
		verifyFeed(t, db, id, feed2)
	}

	verifyTables(t, db, map[string]int{
		"feeds":           feeds,
		"items":           items,
		"sqlite_sequence": 1,
	})

	results, err := syncAll(db)
	if err != nil {
		t.Fatalf("syncAll returned error %q", err)
	}

	if len(results) != len(allTestCases) {
		t.Fatalf("expected syncAll to return %d results, got %d", len(allTestCases), len(results))
	}

	for _, res := range results {

		unplayed, ok := unplayedByID[res.ID]
		if !ok {
			t.Fatalf("syncAll returned unexpected ID %d", res.ID)
		}

		if res.Err != nil {
			t.Fatalf("%s: syncAll returned error %s", res.Title, res.Err)
		}

		if res.Items != unplayed {
			t.Errorf("%s: expected %d unplayed items, got %d", res.Title, unplayed, res.Items)
			continue
		}

		verifyUnplayedItemCount(t, db, res.ID, unplayed)
		items += unplayed
	}

	verifyTables(t, db, map[string]int{
		"feeds":           feeds,
		"items":           items,
		"sqlite_sequence": 1,
	})
}

type feedSortInfo struct {
	Title         string
	lcTitle       string
	lcAuthor      string
	lastPubdate   time.Time
	itemCount     int
	unplayedCount int
	timestamp     time.Time
	id            int64
}

func newFeedSortInfo(feed *kibner.Feed, unplayed int, timestamp time.Time, id int64) *feedSortInfo {

	var pubdate time.Time
	if len(feed.Items) > 0 {
		pubdate = feed.Items[0].Pubdate
	}

	return &feedSortInfo{
		Title:         feed.Title,
		lcTitle:       strings.ToLower(feed.Title),
		lcAuthor:      strings.ToLower(feed.Author),
		itemCount:     len(feed.Items),
		unplayedCount: unplayed,
		lastPubdate:   pubdate,
		timestamp:     timestamp,
		id:            id,
	}
}

func filterFeedSortInfo(feeds []*feedSortInfo, test func(*feedSortInfo) bool) []*feedSortInfo {

	if len(feeds) == 0 {
		return feeds
	}

	results := make([]*feedSortInfo, 0, len(feeds))
	for _, f := range feeds {
		if test(f) {
			results = append(results, f)
		}
	}

	return results
}

func filterFeedsAuthor(author string) func([]*feedSortInfo) []*feedSortInfo {
	return func(feeds []*feedSortInfo) []*feedSortInfo {
		return filterFeedSortInfo(feeds, func(f *feedSortInfo) bool {
			return strings.Contains(f.lcAuthor, author)
		})
	}
}

func filterFeedsTitle(title string) func([]*feedSortInfo) []*feedSortInfo {
	return func(feeds []*feedSortInfo) []*feedSortInfo {
		return filterFeedSortInfo(feeds, func(f *feedSortInfo) bool {
			return strings.Contains(f.lcTitle, title)
		})
	}
}

func limitFeedSortInfo(feeds []*feedSortInfo, n int) []*feedSortInfo {

	if len(feeds) <= n {
		return feeds
	}

	return feeds[:n]
}

func limitFeeds(n int) func([]*feedSortInfo) []*feedSortInfo {
	return func(feeds []*feedSortInfo) []*feedSortInfo {
		return limitFeedSortInfo(feeds, n)
	}
}

func lessFeedID(feeds []*feedSortInfo, i, j int) bool {
	return !(feeds[i].id < feeds[j].id)
}

func lessFeedSecondary(feeds []*feedSortInfo, i, j int) bool {
	if feeds[i].lcTitle == feeds[j].lcTitle {
		return lessFeedID(feeds, i, j)
	}
	return feeds[i].lcTitle < feeds[j].lcTitle
}

func sortFeedSortInfo(feeds []*feedSortInfo, less func(i, j int) bool) []*feedSortInfo {

	if len(feeds) < 2 {
		return feeds
	}

	sort.Slice(feeds, less)
	return feeds
}

func sortFeedsPubdate(ascending bool) func([]*feedSortInfo) []*feedSortInfo {
	return func(feeds []*feedSortInfo) []*feedSortInfo {
		return sortFeedSortInfo(feeds, func(i, j int) bool {
			if feeds[i].lastPubdate.Equal(feeds[j].lastPubdate) {
				return lessFeedSecondary(feeds, i, j)
			}
			return ascending == (feeds[i].lastPubdate.Before(feeds[j].lastPubdate))
		})
	}
}

func sortFeedsTitle(ascending bool) func([]*feedSortInfo) []*feedSortInfo {
	return func(feeds []*feedSortInfo) []*feedSortInfo {
		return sortFeedSortInfo(feeds, func(i, j int) bool {
			if feeds[i].lcTitle == feeds[j].lcTitle {
				return lessFeedID(feeds, i, j)
			}
			return ascending == (feeds[i].lcTitle < feeds[j].lcTitle)
		})
	}
}

func sortFeedsItemCount(ascending bool) func([]*feedSortInfo) []*feedSortInfo {
	return func(feeds []*feedSortInfo) []*feedSortInfo {
		return sortFeedSortInfo(feeds, func(i, j int) bool {
			if feeds[i].itemCount == feeds[j].itemCount {
				return lessFeedSecondary(feeds, i, j)
			}
			return ascending == (feeds[i].itemCount < feeds[j].itemCount)
		})
	}
}

func sortFeedsUnplayedCount(ascending bool) func([]*feedSortInfo) []*feedSortInfo {
	return func(feeds []*feedSortInfo) []*feedSortInfo {
		return sortFeedSortInfo(feeds, func(i, j int) bool {
			if feeds[i].unplayedCount == feeds[j].unplayedCount {
				return lessFeedSecondary(feeds, i, j)
			}
			return ascending == (feeds[i].unplayedCount < feeds[j].unplayedCount)
		})
	}
}

func sortFeedsTimestamp(ascending bool) func([]*feedSortInfo) []*feedSortInfo {
	return func(feeds []*feedSortInfo) []*feedSortInfo {
		return sortFeedSortInfo(feeds, func(i, j int) bool {
			if feeds[i].timestamp.Equal(feeds[j].timestamp) {
				return lessFeedSecondary(feeds, i, j)
			}
			return ascending == (feeds[i].timestamp.Before(feeds[j].timestamp))
		})
	}
}

func TestListFeeds(t *testing.T) {
	testWithInitDB(t, testListFeeds)
}

func testListFeeds(t *testing.T, db *sql.DB) {

	timestamp := time.Now()
	allFeeds := make([]*feedSortInfo, 0, len(allTestCases))

	for name, test := range allTestCases {

		feed := test.NewFeed()

		id, err := saveFeed(db, feed, timestamp)
		if err != nil {
			t.Fatalf("%s: saveFeed returned error %q", name, err)
		}

		var unplayed int
		if n := len(feed.Items); n > 0 {
			unplayed = rand.Intn(n)
			markUnplayedItems(t, db, id, unplayed)
		}

		allFeeds = append(allFeeds, newFeedSortInfo(feed, unplayed, timestamp, id))
		timestamp = timestamp.AddDate(0, 0, -1)
	}

	data := []struct {
		Name    string
		Opts    []listFeedOptions
		Actions []func([]*feedSortInfo) []*feedSortInfo
	}{
		{
			Name: "Feeds By Pubdate DESC",
			Opts: []listFeedOptions{
				{},
				{
					SortOrder: sortOrderDefault,
				},
				{
					SortOrder: sortOrderDesc,
				},
				{
					SortBy: sortFeedsByPubdate,
				},
				{
					SortBy:    sortFeedsByPubdate,
					SortOrder: sortOrderDefault,
				},
				{
					SortBy:    sortFeedsByPubdate,
					SortOrder: sortOrderDesc,
				},
			},
			Actions: []func([]*feedSortInfo) []*feedSortInfo{
				sortFeedsPubdate(false),
			},
		},
		{
			Name: "Feeds By Pubdate ASC",
			Opts: []listFeedOptions{
				{
					SortOrder: sortOrderAsc,
				},
				{
					SortBy:    sortFeedsByPubdate,
					SortOrder: sortOrderAsc,
				},
			},
			Actions: []func([]*feedSortInfo) []*feedSortInfo{
				sortFeedsPubdate(true),
			},
		},
		{
			Name: "Top 5 Feeds By Pubdate",
			Opts: []listFeedOptions{
				{
					SortBy: sortFeedsByPubdate,
					Limit:  5,
				},
			},
			Actions: []func([]*feedSortInfo) []*feedSortInfo{
				sortFeedsPubdate(false),
				limitFeeds(5),
			},
		},
		{
			Name: "Bottom 5 Feeds By Pubdate",
			Opts: []listFeedOptions{
				{
					SortBy:    sortFeedsByPubdate,
					SortOrder: sortOrderAsc,
					Limit:     5,
				},
			},
			Actions: []func([]*feedSortInfo) []*feedSortInfo{
				sortFeedsPubdate(true),
				limitFeeds(5),
			},
		},
		{
			Name: "Feeds By Title ASC",
			Opts: []listFeedOptions{
				{
					SortBy: sortFeedsByTitle,
				},
				{
					SortBy:    sortFeedsByTitle,
					SortOrder: sortOrderDefault,
				},
				{
					SortBy:    sortFeedsByTitle,
					SortOrder: sortOrderAsc,
				},
			},
			Actions: []func([]*feedSortInfo) []*feedSortInfo{
				sortFeedsTitle(true),
			},
		},
		{
			Name: "Feeds By Title DESC",
			Opts: []listFeedOptions{
				{
					SortBy:    sortFeedsByTitle,
					SortOrder: sortOrderDesc,
				},
			},
			Actions: []func([]*feedSortInfo) []*feedSortInfo{
				sortFeedsTitle(false),
			},
		},
		{
			Name: "Feeds By Item Count DESC",
			Opts: []listFeedOptions{
				{
					SortBy: sortFeedsByItemCount,
				},
				{
					SortBy:    sortFeedsByItemCount,
					SortOrder: sortOrderDefault,
				},
				{
					SortBy:    sortFeedsByItemCount,
					SortOrder: sortOrderDesc,
				},
			},
			Actions: []func([]*feedSortInfo) []*feedSortInfo{
				sortFeedsItemCount(false),
			},
		},
		{
			Name: "Feeds By Item Count ASC",
			Opts: []listFeedOptions{
				{
					SortBy:    sortFeedsByItemCount,
					SortOrder: sortOrderAsc,
				},
			},
			Actions: []func([]*feedSortInfo) []*feedSortInfo{
				sortFeedsItemCount(true),
			},
		},
		{
			Name: "Feeds By Unplayed Count DESC",
			Opts: []listFeedOptions{
				{
					SortBy: sortFeedsByUnplayedCount,
				},
				{
					SortBy:    sortFeedsByUnplayedCount,
					SortOrder: sortOrderDefault,
				},
				{
					SortBy:    sortFeedsByUnplayedCount,
					SortOrder: sortOrderDesc,
				},
			},
			Actions: []func([]*feedSortInfo) []*feedSortInfo{
				sortFeedsUnplayedCount(false),
			},
		},
		{
			Name: "Feeds By Unplayed Count ASC",
			Opts: []listFeedOptions{
				{
					SortBy:    sortFeedsByUnplayedCount,
					SortOrder: sortOrderAsc,
				},
			},
			Actions: []func([]*feedSortInfo) []*feedSortInfo{
				sortFeedsUnplayedCount(true),
			},
		},
		{
			Name: "Feeds By Timestamp DESC",
			Opts: []listFeedOptions{
				{
					SortBy: sortFeedsByTimestamp,
				},
				{
					SortBy:    sortFeedsByTimestamp,
					SortOrder: sortOrderDefault,
				},
				{
					SortBy:    sortFeedsByTimestamp,
					SortOrder: sortOrderDesc,
				},
			},
			Actions: []func([]*feedSortInfo) []*feedSortInfo{
				sortFeedsTimestamp(false),
			},
		},
		{
			Name: "Feeds By Timestamp ASC",
			Opts: []listFeedOptions{
				{
					SortBy:    sortFeedsByTimestamp,
					SortOrder: sortOrderAsc,
				},
			},
			Actions: []func([]*feedSortInfo) []*feedSortInfo{
				sortFeedsTimestamp(true),
			},
		},
		{
			Name: "Feeds With Author 'Gimlet'",
			Opts: []listFeedOptions{
				{
					SortBy: sortFeedsByTitle,
					Author: "Gimlet",
				},
			},
			Actions: []func([]*feedSortInfo) []*feedSortInfo{
				filterFeedsAuthor("gimlet"),
				sortFeedsTitle(true),
			},
		},
		{
			Name: "Feeds With Author 'Gladwell'",
			Opts: []listFeedOptions{
				{
					SortBy: sortFeedsByTitle,
					Author: "Gladwell",
				},
			},
			Actions: []func([]*feedSortInfo) []*feedSortInfo{
				filterFeedsAuthor("gladwell"),
				sortFeedsTitle(true),
			},
		},
		{
			Name: "Feeds With Nonsense Author",
			Opts: []listFeedOptions{
				{
					SortBy: sortFeedsByTitle,
					Author: "xxxxxxxxxxxxx",
				},
			},
			Actions: []func([]*feedSortInfo) []*feedSortInfo{
				filterFeedsAuthor("xxxxxxxxxxxxx"),
				sortFeedsTitle(true),
			},
		},
		{
			Name: "Feeds With Title 'The'",
			Opts: []listFeedOptions{
				{
					SortBy: sortFeedsByTitle,
					Title:  "The",
				},
			},
			Actions: []func([]*feedSortInfo) []*feedSortInfo{
				filterFeedsTitle("the"),
				sortFeedsTitle(true),
			},
		},
		{
			Name: "Feeds With Title 'Fish'",
			Opts: []listFeedOptions{
				{
					SortBy: sortFeedsByTitle,
					Title:  "Fish",
				},
			},
			Actions: []func([]*feedSortInfo) []*feedSortInfo{
				filterFeedsTitle("fish"),
				sortFeedsTitle(true),
			},
		},
		{
			Name: "Feeds With Nonsense Title",
			Opts: []listFeedOptions{
				{
					SortBy: sortFeedsByTitle,
					Title:  "xxxxxxxxxxxxxxxx",
				},
			},
			Actions: []func([]*feedSortInfo) []*feedSortInfo{
				filterFeedsTitle("xxxxxxxxxxxxxxxx"),
				sortFeedsTitle(true),
			},
		},
		{
			Name: "Top 3 Feeds With Author 'Gimlet' And Title 'Serial",
			Opts: []listFeedOptions{
				{
					SortBy: sortFeedsByTitle,
					Author: "Gimlet",
					Title:  "Serial",
					Limit:  3,
				},
			},
			Actions: []func([]*feedSortInfo) []*feedSortInfo{
				filterFeedsAuthor("gimlet"),
				filterFeedsTitle("serial"),
				sortFeedsTitle(true),
				limitFeeds(3),
			},
		},
	}

	tmpl := template.Must(
		template.New("test").Parse(`{{range $i, $f := .Feeds}}{{with $i}}#{{end}}{{$f.Title}}{{end}}`),
	)

	for _, test := range data {

		feeds := allFeeds
		for _, apply := range test.Actions {
			feeds = apply(feeds)
		}

		titles := make([]string, len(feeds))
		for i, f := range feeds {
			titles[i] = f.Title
		}

		exp := strings.Join(titles, "#")

		for i := range test.Opts {

			w := &bytes.Buffer{}

			err := listFeeds(db, w, tmpl, test.Opts[i])
			if err != nil {
				t.Fatalf("%s: listFeeds returned error %q", test.Name, err)
			}

			if result := w.String(); result != exp {
				t.Errorf("%s (%d/%d): Expected titles %s, got %s", test.Name, i+1, len(test.Opts), jsonify(titles), jsonify(strings.Split(result, "#")))
			}
		}
	}
}

func TestLoadFeedViews(t *testing.T) {
	testWithDB(t, testLoadFeedViews)
}

func testLoadFeedViews(t *testing.T, db *sql.DB) {

	for name, test := range allTestCases {

		testInitDB(t, db)

		feed := test.NewFeed()

		_, err := saveFeed(db, feed, time.Now())
		if err != nil {
			t.Fatalf("%s: saveFeed returned error %q", name, err)
		}

		data, err := loadFeedViews(db, listFeedOptions{})
		if err != nil {
			t.Fatalf("%s: loadFeedViews returned error %q", name, err)
		}

		if len(data) != 1 {
			t.Fatalf("%s: expected loadFeedViews to return 1 item, got %d", name, len(data))
		}

		verifyFeedView(t, &data[0], feed)
	}
}

func TestDefaultFeedTemplate(t *testing.T) {
	testWithInitDB(t, testDefaultFeedTemplate)
}

func testDefaultFeedTemplate(t *testing.T, db *sql.DB) {

	feeds := map[string]int{
		"Crooked Conversations": 1,
		"Rework":                2,
		"Homecoming":            0,
		"Why We Eat":            0,
	}

	for name, unplayed := range feeds {

		id, err := saveFeed(db, allTestCases[name].NewFeed(), time.Now())
		if err != nil {
			t.Fatalf("%s: saveFeed returned error %q", name, err)
		}

		if unplayed > 0 {
			markUnplayedItems(t, db, id, unplayed)
		}
	}

	data := []struct {
		Opts   listFeedOptions
		Output string
	}{
		{
			Opts: listFeedOptions{
				Title: "xxxxxxx",
			},
			Output: `No feeds found`,
		},
		{
			Opts: listFeedOptions{
				Title: "Crooked",
			},
			Output: `Showing 1 feed:

      1. Crooked Conversations
         Crooked Media
         1 item, 1 unplayed (updated Today)`,
		},
		{
			Opts: listFeedOptions{
				Author: "Gimlet",
			},
			Output: `Showing 2 feeds:

      1. Homecoming
         Gimlet
         13 items (updated over a month ago)

      2. Why We Eat What We Eat
         Blue Apron / Gimlet Creative
         0 items`,
		},
		{
			Opts: listFeedOptions{
				ShowDesc: true,
			},
			Output: `Showing 4 feeds:

      1. Crooked Conversations
         Crooked Media
         One side effect of our national addiction to Trump’s tweets and
         other news cycle garbage is that fascinating issues, brilliant books
         and important debates aren't getting the attention they deserve. With
         a rotating crew of your favorite Crooked Media hosts, contributors,
         and special guests, Crooked Conversations brings Pod Save America's
         no-b.s., conversational style to topics in politics, media, culture,
         sports, and technology that aren’t making headlines but still have a
         major impact on our world.
         1 item, 1 unplayed (updated Today)

      2. REWORK
         Basecamp
         A podcast by Basecamp about a better way to work and run your
         business. We bring you stories and unconventional wisdom from
         Basecamp’s co-founders and other business owners.
         5 items, 2 unplayed (updated 8 days ago)

      3. Homecoming
         Gimlet
         A new psychological thriller from Gimlet Media, starring Catherine
         Keener, Oscar Isaac, and David Schwimmer.
         13 items (updated over a month ago)

      4. Why We Eat What We Eat
         Blue Apron / Gimlet Creative
         A podcast from Blue Apron and Gimlet Creative for anyone who has ever
         eaten.
         0 items`,
		},
	}

	now := time.Date(2017, time.October, 4, 12, 30, 0, 0, time.Local)
	tmpl := defaultFeedTemplate(now)

	for _, test := range data {

		wexp := &bytes.Buffer{}
		fmt.Fprintln(wexp, test.Output)

		wgot := &bytes.Buffer{}
		err := listFeeds(db, wgot, tmpl, test.Opts)
		if err != nil {
			t.Fatalf("listFeeds returned error %q", err)
		}

		if exp, got := wexp.String(), wgot.String(); got != exp {
			t.Errorf("Expected listFeeds to output %q, got %q", exp, got)
		}
	}
}

type itemSortInfo struct {
	Title       string
	FeedTitle   string
	lcTitle     string
	lcFeedTitle string
	pubdate     time.Time
	duration    float64
	unplayed    bool
	feedID      int64
	timestamp   time.Time
	id          int64
}

func newItemSortInfo(item *kibner.Item, feedID int64, feedTitle string, unplayed bool, timestamp time.Time, id int64) *itemSortInfo {
	return &itemSortInfo{
		Title:       item.Title,
		FeedTitle:   feedTitle,
		lcTitle:     strings.ToLower(item.Title),
		lcFeedTitle: strings.ToLower(feedTitle),
		pubdate:     item.Pubdate,
		duration:    item.Duration.Seconds(),
		unplayed:    unplayed,
		feedID:      feedID,
		timestamp:   timestamp,
		id:          id,
	}
}

func filterItemSortInfo(items []*itemSortInfo, test func(*itemSortInfo) bool) []*itemSortInfo {

	if len(items) == 0 {
		return items
	}

	results := make([]*itemSortInfo, 0, len(items))
	for _, it := range items {
		if test(it) {
			results = append(results, it)
		}
	}

	return results
}

func filterItemsUnplayed() func([]*itemSortInfo) []*itemSortInfo {
	return func(items []*itemSortInfo) []*itemSortInfo {
		return filterItemSortInfo(items, func(it *itemSortInfo) bool {
			return it.unplayed
		})
	}
}

func filterItemsStartDate(start time.Time) func([]*itemSortInfo) []*itemSortInfo {
	return func(items []*itemSortInfo) []*itemSortInfo {
		return filterItemSortInfo(items, func(it *itemSortInfo) bool {
			return !it.pubdate.Before(start)
		})
	}
}

func filterItemsTitle(title string) func([]*itemSortInfo) []*itemSortInfo {
	return func(items []*itemSortInfo) []*itemSortInfo {
		return filterItemSortInfo(items, func(it *itemSortInfo) bool {
			return strings.Contains(it.lcTitle, title)
		})
	}
}

func filterItemsFeedID(feedID int64) func([]*itemSortInfo) []*itemSortInfo {
	return func(items []*itemSortInfo) []*itemSortInfo {
		return filterItemSortInfo(items, func(it *itemSortInfo) bool {
			return it.feedID == feedID
		})
	}
}

func limitItemSortInfo(items []*itemSortInfo, n int) []*itemSortInfo {

	if len(items) <= n {
		return items
	}

	return items[:n]
}

func limitItems(n int) func([]*itemSortInfo) []*itemSortInfo {
	return func(items []*itemSortInfo) []*itemSortInfo {
		return limitItemSortInfo(items, n)
	}
}

func lessItemID(items []*itemSortInfo, i, j int) bool {
	return !(items[i].id < items[j].id)
}

func lessItemSecondary(items []*itemSortInfo, i, j int) bool {
	if items[i].pubdate.Equal(items[j].pubdate) {
		return lessItemID(items, i, j)
	}
	return !items[i].pubdate.Before(items[j].pubdate)
}

func sortItemSortInfo(items []*itemSortInfo, less func(i, j int) bool) []*itemSortInfo {

	if len(items) < 2 {
		return items
	}

	sort.Slice(items, less)
	return items
}

func sortItemsPubdate(ascending bool) func([]*itemSortInfo) []*itemSortInfo {
	return func(items []*itemSortInfo) []*itemSortInfo {
		return sortItemSortInfo(items, func(i, j int) bool {
			if items[i].pubdate.Equal(items[j].pubdate) {
				return lessItemID(items, i, j)
			}
			return ascending == (items[i].pubdate.Before(items[j].pubdate))
		})
	}
}

func sortItemsTitle(ascending bool) func([]*itemSortInfo) []*itemSortInfo {
	return func(items []*itemSortInfo) []*itemSortInfo {
		return sortItemSortInfo(items, func(i, j int) bool {
			if items[i].lcTitle == items[j].lcTitle {
				return lessItemSecondary(items, i, j)
			}
			return ascending == (items[i].lcTitle < items[j].lcTitle)
		})
	}
}

func sortItemsFeedTitle(ascending bool) func([]*itemSortInfo) []*itemSortInfo {
	return func(items []*itemSortInfo) []*itemSortInfo {
		return sortItemSortInfo(items, func(i, j int) bool {
			if items[i].lcFeedTitle == items[j].lcFeedTitle {
				return lessItemSecondary(items, i, j)
			}
			return ascending == (items[i].lcFeedTitle < items[j].lcFeedTitle)
		})
	}
}

func sortItemsDuration(ascending bool) func([]*itemSortInfo) []*itemSortInfo {
	return func(items []*itemSortInfo) []*itemSortInfo {
		return sortItemSortInfo(items, func(i, j int) bool {
			if items[i].duration == items[j].duration {
				return lessItemSecondary(items, i, j)
			}
			return ascending == (items[i].duration < items[j].duration)
		})
	}
}

func sortItemsTimestamp(ascending bool) func([]*itemSortInfo) []*itemSortInfo {
	return func(items []*itemSortInfo) []*itemSortInfo {
		return sortItemSortInfo(items, func(i, j int) bool {
			if items[i].timestamp.Equal(items[j].timestamp) {
				return lessItemSecondary(items, i, j)
			}
			return ascending == (items[i].timestamp.Before(items[j].timestamp))
		})
	}
}

func TestListItems(t *testing.T) {
	testWithInitDB(t, testListItems)
}

func testListItems(t *testing.T, db *sql.DB) {

	id := int64(1)
	timestamp := time.Now()
	allItems := make([]*itemSortInfo, 0, len(allTestCases))

	for name, test := range allTestCases {

		feed := test.NewFeed()

		feedID, err := saveFeed(db, feed, timestamp)
		if err != nil {
			t.Fatalf("%s: saveFeed returned error %q", name, err)
		}

		for i, item := range feed.Items {

			unplayed := i == 0
			if unplayed {
				markUnplayedItems(t, db, feedID, 1)
			}

			allItems = append(allItems, newItemSortInfo(item, feedID, feed.Title, unplayed, timestamp, id))
			id++
		}

		timestamp = timestamp.AddDate(0, 0, -1)
	}

	data := []struct {
		Name    string
		Opts    []listItemOptions
		Actions []func([]*itemSortInfo) []*itemSortInfo
	}{
		{
			Name: "Items By Pubdate DESC",
			Opts: []listItemOptions{
				{},
				{
					SortOrder: sortOrderDefault,
				},
				{
					SortOrder: sortOrderDesc,
				},
				{
					SortBy: sortItemsByPubdate,
				},
				{
					SortBy:    sortItemsByPubdate,
					SortOrder: sortOrderDefault,
				},
				{
					SortBy:    sortItemsByPubdate,
					SortOrder: sortOrderDesc,
				},
			},
			Actions: []func([]*itemSortInfo) []*itemSortInfo{
				sortItemsPubdate(false),
			},
		},
		{
			Name: "Items By Pubdate ASC",
			Opts: []listItemOptions{
				{
					SortOrder: sortOrderAsc,
				},
				{
					SortBy:    sortItemsByPubdate,
					SortOrder: sortOrderAsc,
				},
			},
			Actions: []func([]*itemSortInfo) []*itemSortInfo{
				sortItemsPubdate(true),
			},
		},
		{
			Name: "Items By Title ASC",
			Opts: []listItemOptions{
				{
					SortBy: sortItemsByTitle,
				},
				{
					SortBy:    sortItemsByTitle,
					SortOrder: sortOrderDefault,
				},
				{
					SortBy:    sortItemsByTitle,
					SortOrder: sortOrderAsc,
				},
			},
			Actions: []func([]*itemSortInfo) []*itemSortInfo{
				sortItemsTitle(true),
			},
		},
		{
			Name: "Items By Title DESC",
			Opts: []listItemOptions{
				{
					SortBy:    sortItemsByTitle,
					SortOrder: sortOrderDesc,
				},
			},
			Actions: []func([]*itemSortInfo) []*itemSortInfo{
				sortItemsTitle(false),
			},
		},
		{
			Name: "Items By Feed ASC",
			Opts: []listItemOptions{
				{
					SortBy: sortItemsByFeed,
				},
				{
					SortBy:    sortItemsByFeed,
					SortOrder: sortOrderDefault,
				},
				{
					SortBy:    sortItemsByFeed,
					SortOrder: sortOrderAsc,
				},
			},
			Actions: []func([]*itemSortInfo) []*itemSortInfo{
				sortItemsFeedTitle(true),
			},
		},
		{
			Name: "Items By Feed DESC",
			Opts: []listItemOptions{
				{
					SortBy:    sortItemsByFeed,
					SortOrder: sortOrderDesc,
				},
			},
			Actions: []func([]*itemSortInfo) []*itemSortInfo{
				sortItemsFeedTitle(false),
			},
		},
		{
			Name: "Items By Duration DESC",
			Opts: []listItemOptions{
				{
					SortBy: sortItemsByDuration,
				},
				{
					SortBy:    sortItemsByDuration,
					SortOrder: sortOrderDefault,
				},
				{
					SortBy:    sortItemsByDuration,
					SortOrder: sortOrderDesc,
				},
			},
			Actions: []func([]*itemSortInfo) []*itemSortInfo{
				sortItemsDuration(false),
			},
		},
		{
			Name: "Items By Duration ASC",
			Opts: []listItemOptions{
				{
					SortBy:    sortItemsByDuration,
					SortOrder: sortOrderAsc,
				},
			},
			Actions: []func([]*itemSortInfo) []*itemSortInfo{
				sortItemsDuration(true),
			},
		},
		{
			Name: "Items By Timestamp DESC",
			Opts: []listItemOptions{
				{
					SortBy: sortItemsByTimestamp,
				},
				{
					SortBy:    sortItemsByTimestamp,
					SortOrder: sortOrderDefault,
				},
				{
					SortBy:    sortItemsByTimestamp,
					SortOrder: sortOrderDesc,
				},
			},
			Actions: []func([]*itemSortInfo) []*itemSortInfo{
				sortItemsTimestamp(false),
			},
		},
		{
			Name: "Items By Timestamp ASC",
			Opts: []listItemOptions{
				{
					SortBy:    sortItemsByTimestamp,
					SortOrder: sortOrderAsc,
				},
			},
			Actions: []func([]*itemSortInfo) []*itemSortInfo{
				sortItemsTimestamp(true),
			},
		},
		{
			Name: "Top 10 Items By Pubdate",
			Opts: []listItemOptions{
				{
					Limit: 10,
				},
			},
			Actions: []func([]*itemSortInfo) []*itemSortInfo{
				sortItemsPubdate(false),
				limitItems(10),
			},
		},
		{
			Name: "Bottom 10 Items By Pubdate",
			Opts: []listItemOptions{
				{
					Limit:     10,
					SortOrder: sortOrderAsc,
				},
			},
			Actions: []func([]*itemSortInfo) []*itemSortInfo{
				sortItemsPubdate(true),
				limitItems(10),
			},
		},
		{
			Name: "Top 10 Unplayed Items By Pubdate",
			Opts: []listItemOptions{
				{
					Limit:    10,
					Unplayed: true,
				},
			},
			Actions: []func([]*itemSortInfo) []*itemSortInfo{
				filterItemsUnplayed(),
				sortItemsPubdate(false),
				limitItems(10),
			},
		},
		{
			Name: "Items Since October 2017",
			Opts: []listItemOptions{
				{
					StartDate: time.Date(2017, time.October, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			Actions: []func([]*itemSortInfo) []*itemSortInfo{
				filterItemsStartDate(time.Date(2017, time.October, 1, 0, 0, 0, 0, time.UTC)),
				sortItemsPubdate(false),
			},
		},
		{
			Name: "Items With Title 'The'",
			Opts: []listItemOptions{
				{
					Title: "the",
				},
			},
			Actions: []func([]*itemSortInfo) []*itemSortInfo{
				filterItemsTitle("the"),
				sortItemsPubdate(false),
			},
		},
		{
			Name: "Items With Title 'Pineapple'",
			Opts: []listItemOptions{
				{
					Title: "Pineapple",
				},
			},
			Actions: []func([]*itemSortInfo) []*itemSortInfo{
				filterItemsTitle("pineapple"),
				sortItemsPubdate(false),
			},
		},
		{
			Name: "Items With Nonsense Title",
			Opts: []listItemOptions{
				{
					Title: "xxxxxxxxxxxx",
				},
			},
			Actions: []func([]*itemSortInfo) []*itemSortInfo{
				filterItemsTitle("xxxxxxxxxxxx"),
				sortItemsPubdate(false),
			},
		},
		{
			Name: "Items With FeedID 1",
			Opts: []listItemOptions{
				{
					FeedID: 1,
				},
			},
			Actions: []func([]*itemSortInfo) []*itemSortInfo{
				filterItemsFeedID(1),
				sortItemsPubdate(false),
			},
		},
	}

	tmpl := template.Must(
		template.New("test").Parse(`{{range $i, $item := .Items}}{{with $i}}#{{end}}{{$item.FeedTitle}}: {{$item.Title}}{{end}}`),
	)

	for _, test := range data {

		items := allItems
		for _, apply := range test.Actions {
			items = apply(items)
		}

		titles := make([]string, len(items))
		for i, item := range items {
			titles[i] = item.FeedTitle + ": " + item.Title
		}

		exp := strings.Join(titles, "#")

		for i := range test.Opts {

			w := &bytes.Buffer{}

			err := listItems(db, w, tmpl, test.Opts[i])
			if err != nil {
				t.Fatalf("%s: listItems returned error %q", test.Name, err)
			}

			if result := w.String(); result != exp {
				t.Errorf("%s (%d/%d): Expected titles %s, got %s", test.Name, i+1, len(test.Opts), jsonify(titles), jsonify(strings.Split(result, "#")))
			}
		}
	}
}

func TestDefaultListTemplate(t *testing.T) {
	testWithInitDB(t, testDefaultListTemplate)
}

func testDefaultListTemplate(t *testing.T, db *sql.DB) {

	feeds := map[string]int{
		"Columbo": 3,
	}

	for name, unplayed := range feeds {

		id, err := saveFeed(db, allTestCases[name].NewFeed(), time.Now())
		if err != nil {
			t.Fatalf("%s: saveFeed returned error %q", name, err)
		}

		if unplayed > 0 {
			markUnplayedItems(t, db, id, unplayed)
		}
	}

	data := []struct {
		Opts   listItemOptions
		Now    time.Time
		Output string
	}{
		{
			Opts: listItemOptions{
				Title: "xxxxxxx",
			},
			Output: `No items found`,
		},
		{
			Opts: listItemOptions{
				// TODO: Sqlite LIKE is only case-sensitive for ASCII
				// characters. We should fix this.
				Title: "étude",
			},
			Output: `No items found`,
		},
		{
			Opts: listItemOptions{
				Title: "Étude",
			},
			Now: time.Date(1972, time.September, 18, 10, 30, 0, 0, time.UTC),
			Output: `Showing 1 item:

    * 1. Étude in Black
         From Columbo, Today
         Duration: 1h38m`,
		},
		{
			Opts: listItemOptions{
				Title:     "Dead",
				SortOrder: sortOrderAsc,
				Limit:     2,
			},
			Now: time.Date(1971, time.March, 2, 10, 30, 0, 0, time.UTC),
			Output: `Showing 2 items:

      1. Dead Weight
         From Columbo, Date unknown
         Duration: 1h13m

      2. Ransom for a Dead Man
         From Columbo, Today
         Duration: Unknown`,
		},
		{
			Opts: listItemOptions{
				Title:     "Murder",
				SortOrder: sortOrderAsc,
				Limit:     3,
			},
			Now: time.Date(1972, time.February, 14, 10, 30, 0, 0, time.UTC),
			Output: `Showing 3 items:

      1. Prescription: Murder
         From Columbo, in 1968
         Duration: Unknown

      2. Murder by the Book
         From Columbo, over 4 months ago
         Duration: 1h13m

    * 3. Blueprint for Murder
         From Columbo, 4 days ago
         Duration: 1h13m`,
		},
		{
			Opts: listItemOptions{
				FeedID:    1,
				Title:     "Murder",
				SortOrder: sortOrderAsc,
				Limit:     3,
			},
			Now: time.Date(1972, time.February, 14, 10, 30, 0, 0, time.UTC),
			Output: `Showing 3 items:

      1. Prescription: Murder
         Released in 1968
         Duration: Unknown

      2. Murder by the Book
         Released over 4 months ago
         Duration: 1h13m

    * 3. Blueprint for Murder
         Released 4 days ago
         Duration: 1h13m`,
		},
		{
			Opts: listItemOptions{
				FeedID:    1,
				Title:     "Murder",
				SortOrder: sortOrderAsc,
				Limit:     3,
				ShowDesc:  true,
			},
			Now: time.Date(1972, time.February, 14, 10, 30, 0, 0, time.UTC),
			Output: `Showing 3 items:

      1. Prescription: Murder
         Released in 1968
         Dr. Ray Fleming (Gene Barry), a psychiatrist, murders his wife (Nina
         Foch) and persuades his mistress Joan Hudson (Katherine Justice), who
         is an actress and one of his patients, to support his alibi by
         impersonating her.
         Duration: Unknown

      2. Murder by the Book
         Released over 4 months ago
         Ken Franklin (Jack Cassidy) is one-half of a mystery writing team, but
         partner Jim Ferris (Martin Milner) wants to go solo. This would expose
         the fact that Ferris did all the actual writing, and leave the
         high-living Franklin without his cash cow. Franklin tricks Ferris into
         taking a trip to his remote cabin two hours away. At the cabin, he
         convinces Ferris to call home and say he's working late at the office.
         During the call, Franklin shoots Ferris, then takes his body back
         north and dumps it on his lawn.
         Duration: 1h13m

    * 3. Blueprint for Murder
         Released 4 days ago
         Duration: 1h13m`,
		},
	}

	for _, test := range data {

		wexp := &bytes.Buffer{}
		fmt.Fprintln(wexp, test.Output)

		tmpl := defaultItemTemplate(test.Now)

		wgot := &bytes.Buffer{}
		err := listItems(db, wgot, tmpl, test.Opts)
		if err != nil {
			t.Fatalf("listItems returned error %q", err)
		}

		if exp, got := wexp.String(), wgot.String(); got != exp {
			t.Errorf("Expected listItems to output %q, got %q", exp, got)
		}
	}
}

func TestUpdateFeed(t *testing.T) {
	testWithInitDB(t, testUpdateFeed)
}

func testUpdateFeed(t *testing.T, db *sql.DB) {

	feed := allTestCases["Serial"].NewFeed()

	id, err := saveFeed(db, feed, time.Now())
	if err != nil {
		t.Fatalf("saveFeed returned error %q", err)
	}

	verifyFeed(t, db, id, feed)

	obj := reflect.ValueOf(feed).Elem()

	data := []struct {
		StructField reflect.Value
		DBField     string
		Chars       int
	}{
		{
			StructField: obj.FieldByName("Title"),
			DBField:     "title",
			Chars:       30,
		},
		{
			StructField: obj.FieldByName("Author"),
			DBField:     "author",
			Chars:       20,
		},
		{
			StructField: obj.FieldByName("Desc"),
			DBField:     "desc",
			Chars:       200,
		},
		{
			StructField: obj.FieldByName("Link"),
			DBField:     "link",
			Chars:       50,
		},
		{
			StructField: obj.FieldByName("Image"),
			DBField:     "image",
			Chars:       65,
		},
	}

	for i := range data {

		values := map[string]interface{}{}

		for j := 0; j <= i; j++ {

			test := data[j]
			value := randomString(test.Chars)

			values[test.DBField] = value
			test.StructField.Set(reflect.ValueOf(value))
		}

		if err := updateFeed(db, id, values); err != nil {
			t.Fatalf("updateFeed %v returned error %q", values, err)
		}

		verifyFeed(t, db, id, feed)
	}
}

func TestUpdateFeedNoValues(t *testing.T) {
	testWithDB(t, testUpdateFeedNoValues)
}

func testUpdateFeedNoValues(t *testing.T, db *sql.DB) {

	data := []map[string]interface{}{
		nil,
		{},
	}

	exp := errors.New("no values provided")

	for _, values := range data {
		if err := updateFeed(db, 1, values); !equalErrors(exp, err) {
			t.Errorf("Expected updateFeed %#v to return error %q, got %v", values, exp, err)
		}

	}
}

func TestUpdateFeedNoFeed(t *testing.T) {
	testWithInitDB(t, testUpdateFeedNoFeed)
}

func testUpdateFeedNoFeed(t *testing.T, db *sql.DB) {

	exp := errors.New("feed not found")

	values := map[string]interface{}{
		"title": "dummy",
	}

	if err := updateFeed(db, 1, values); !equalErrors(exp, err) {
		t.Errorf("Expected updateFeed to return error %q, got %v", exp, err)
	}
}

func TestLoadFeedURL(t *testing.T) {
	testWithInitDB(t, testLoadFeedURL)
}

func testLoadFeedURL(t *testing.T, db *sql.DB) {

	now := time.Now()
	for name, test := range allTestCases {

		feed := test.NewFeed()

		id, err := saveFeed(db, feed, now)
		if err != nil {
			t.Fatalf("%s: saveFeed returned error %q", name, err)
		}

		m := map[target]string{
			targetFeed:  feed.URL,
			targetImage: feed.Image,
			targetLink:  feed.Link,
		}

		for target, exp := range m {

			url, err := loadFeedURL(db, id, target)
			if err != nil {
				t.Fatalf("%s: loadFeedURL %v returned error %q", name, target, err)
			}

			if url != exp {
				t.Errorf("%s: expected loadFeedURL to return %q, got %q", name, exp, url)
			}
		}
	}
}

func TestLoadFeedURLNoFeeds(t *testing.T) {
	testWithInitDB(t, testLoadFeedURLNoFeeds)
}

func testLoadFeedURLNoFeeds(t *testing.T, db *sql.DB) {
	if _, err := loadFeedURL(db, 1, targetFeed); err != errNoFeedFound {
		t.Errorf("Expected loadFeedURL to return error %v, got %v", errNoFeedFound, err)
	}
}

func TestExportList(t *testing.T) {
	testWithInitDB(t, testExportList)
}

func testExportList(t *testing.T, db *sql.DB) {

	now := time.Now()
	urls := make([]string, 0, len(allTestCases))

	for name, test := range allTestCases {

		feed := test.NewFeed()
		urls = append(urls, feed.URL)

		_, err := saveFeed(db, feed, now)
		if err != nil {
			t.Fatalf("%s: saveFeed returned error %s", name, err)
		}
	}

	sort.Slice(urls, func(i, j int) bool {
		return strings.ToLower(urls[i]) < strings.ToLower(urls[j])
	})

	wexp := &bytes.Buffer{}
	for _, s := range urls {
		fmt.Fprintln(wexp, s)
	}

	wgot := &bytes.Buffer{}
	if err := exportList(db, wgot); err != nil {
		t.Fatalf("exportList returned error %q", err)
	}

	if got, exp := wgot.String(), wexp.String(); got != exp {
		t.Errorf("Expected exportList to write %q, got %q", exp, got)
	}
}

func TestExportListNoFeeds(t *testing.T) {
	testWithInitDB(t, testExportListNoFeeds)
}

func testExportListNoFeeds(t *testing.T, db *sql.DB) {

	exp := errors.New("no feeds to export")

	if got := exportList(db, ioutil.Discard); !equalErrors(exp, got) {
		t.Fatalf("expected exportList to return error %q, got %v", exp, got)
	}
}

func TestExportOPML(t *testing.T) {
	testWithInitDB(t, testExportOPML)
}

func testExportOPML(t *testing.T, db *sql.DB) {

	type entry struct {
		Type  string
		Title string
		Desc  string
		URL   string
	}

	now := time.Now()
	entries := make([]entry, 0, len(allTestCases))

	for name, test := range allTestCases {

		feed := test.NewFeed()
		entries = append(entries, entry{
			Type:  feed.Type,
			Title: feed.Title,
			Desc:  feed.Desc,
			URL:   feed.URL,
		})

		_, err := saveFeed(db, feed, now)
		if err != nil {
			t.Fatalf("%s: saveFeed returned error %s", name, err)
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		x, y := entries[i], entries[j]
		if strings.ToLower(x.Title) == strings.ToLower(y.Title) {
			return strings.ToLower(x.URL) < strings.ToLower(y.URL)
		}
		return strings.ToLower(x.Title) < strings.ToLower(y.Title)
	})

	wexp := &bytes.Buffer{}
	writeOPMLHeader(t, wexp, now)
	for _, e := range entries {
		writeOPMLEntry(t, wexp, e.Type, e.Title, e.Desc, e.URL)
	}
	writeOPMLFooter(t, wexp)

	wgot := &bytes.Buffer{}
	if err := exportOPML(db, wgot, now); err != nil {
		t.Fatalf("exportOPML returned error %q", err)
	}

	if got, exp := wgot.String(), wexp.String(); got != exp {
		t.Errorf("Expected exportOPML to write %q, got %q", exp, got)
	}
}

func TestExportOPMLNoFeeds(t *testing.T) {
	testWithInitDB(t, testExportOPMLNoFeeds)
}

func testExportOPMLNoFeeds(t *testing.T, db *sql.DB) {

	exp := errors.New("no feeds to export")

	if got := exportOPML(db, ioutil.Discard, time.Now()); !equalErrors(exp, got) {
		t.Fatalf("expected exportOPML to return error %q, got %v", exp, got)
	}
}

func TestExtractURLsList(t *testing.T) {

	w := &bytes.Buffer{}
	urls := make([]string, 0, len(allTestCases))

	fmt.Fprintln(w, "#", len(allTestCases), "Feed URLs")
	fmt.Fprintln(w)

	for _, test := range allTestCases {

		url := test.NewFeed().URL
		urls = append(urls, url)

		for i := 0; i <= rand.Intn(2); i++ {
			fmt.Fprint(w, strings.Repeat(" ", 1+rand.Intn(5)))
			fmt.Fprintln(w, url)
		}

		switch {
		case len(urls)%3 == 0:
			fmt.Fprintln(w)
		case len(urls)%5 == 0:
			fmt.Fprintln(w)
			fmt.Fprintln(w, "# Comment")
			fmt.Fprintln(w)
		}
	}

	fmt.Fprintln(w, "  ://invalid-url.com")

	got, err := extractURLsList(w)
	if err != nil {
		t.Fatalf("extractURLsList returned error %q", err)
	}

	if !equalStringSlices(got, urls) {
		t.Errorf("expected extractURLsList to return urls %s, got %s", jsonify(urls), jsonify(got))
	}
}

func TestExtractURLsOPML(t *testing.T) {

	w := &bytes.Buffer{}
	urls := make([]string, 0, len(allTestCases))

	writeOPMLHeader(t, w, time.Now())

	for _, test := range allTestCases {

		feed := test.NewFeed()
		urls = append(urls, feed.URL)

		for i := 0; i <= rand.Intn(2); i++ {
			writeOPMLEntry(t, w, feed.Type, feed.Title, feed.Desc, feed.URL)
		}
	}

	writeOPMLEntry(t, w, "rss", "Invalid URL", "This URL is invalid", "://invalid-url.com")
	writeOPMLFooter(t, w)

	got, err := extractURLsOPML(w)
	if err != nil {
		t.Fatalf("extractURLsOPML returned error %q", err)
	}

	if !equalStringSlices(got, urls) {
		t.Errorf("expected extractURLsOPML to return urls %s, got %s", jsonify(urls), jsonify(got))
	}
}

func TestAddCallbackTemplate(t *testing.T) {

	N := 100
	tmpl := template.Must(
		template.New("test").Parse(`{{range .Values}}{{template "double" .}}{{end}}`),
	)

	values := make([]int, N)
	for i := range values {
		values[i] = rand.Intn(1000)
	}

	var results []int
	double := func(i int) error {
		results = append(results, i*2)
		return nil
	}

	tmpl2, err := addCallbackTemplate(tmpl, "double", double)
	if err != nil {
		t.Fatalf("addCallbackTemplate returned error %q", err)
	}

	w := &bytes.Buffer{}
	err = tmpl2.Execute(w, map[string]interface{}{
		"Values": values,
	})
	if err != nil {
		t.Fatalf("template.Execute returned error %q", err)
	}

	if w.String() != "" {
		t.Errorf("Expected no output from template.Execute, got %q", w.String())
	}

	if len(results) != N {
		t.Fatalf("Expected %d results, got %d", N, len(results))
	}

	for i := 0; i < N; i++ {
		if results[i] != values[i]*2 {
			t.Errorf("Expected results[%d] to be %d, got %d", i, values[i]*2, results[i])
		}
	}
}

func verifyAddResult(t *testing.T, res *syncResult, feed *kibner.Feed) {

	if res.ID <= 0 {
		t.Errorf("%s: Expected addFeed to return ID > 0, got %d", feed.Title, res.ID)
	}

	if got, exp := res.URL, feed.URL; got != exp {
		t.Errorf("%s: Expected addFeed to return URL %q, got %q", feed.Title, exp, got)
	}

	if got, exp := res.Title, feed.Title; got != exp {
		t.Errorf("%s: Expected addFeed to return Title %q, got %q", feed.Title, exp, got)
	}

	if got, exp := res.Items, len(feed.Items); got != exp {
		t.Errorf("%s: Expected addFeed to return %d items, got %d items", feed.Title, exp, got)
	}
}

func verifySequenceNumbers(t *testing.T, db *sql.DB, seqnos map[string]int64) {

	got := getSequenceNumbers(t, db)

	if len(got) != len(seqnos) {
		t.Errorf("Expected %d sequence number(s), got %d", len(seqnos), len(got))
		return
	}

	for k, v := range seqnos {

		seqno, ok := got[k]

		if !ok {
			t.Errorf("Expected sequence number for %q", k)
			continue
		}

		if seqno != v {
			t.Errorf("Expected %q sequence number to be %d, got %d", k, v, seqno)
		}
	}
}

func verifyTables(t *testing.T, db *sql.DB, tables map[string]int) {

	got := getTableNames(t, db)

	if len(got) != len(tables) {
		t.Errorf("Expected tables %v, got %v", mapKeys(tables), got)
		return
	}

	checked := map[string]bool{}
	for _, tbl := range got {

		if checked[tbl] {
			t.Errorf("Expected tables %v, got %v", mapKeys(tables), got)
			return
		}
		checked[tbl] = true

		e, ok := tables[tbl]
		if !ok {
			t.Errorf("Expected tables %v, got %v", mapKeys(tables), got)
			return
		}

		if g := getRowCount(t, db, tbl); g != e {
			t.Errorf("Expected table %q to have %d row(s), got %d rows", tbl, e, g)
		}
	}
}

func verifyColumns(t *testing.T, db *sql.DB, tableName string, exp []sqlitemeta.Column) {

	got := getColumns(t, db, tableName)

	if len(got) != len(exp) {
		t.Errorf("Expected %d column(s) in table %q, got %d column(s)", len(exp), tableName, len(got))
		return
	}

	for i := range got {

		g := got[i]
		e := exp[i]

		if g.ID != e.ID {
			t.Errorf("Expected column %s.%s to have ID %d, got %d", tableName, e.Name, e.ID, g.ID)
		}

		if g.Name != e.Name {
			t.Errorf("Expected column %s.%s to have Name %s, got %s", tableName, e.Name, e.Name, g.Name)
		}

		if g.Type != e.Type {
			t.Errorf("Expected column %s.%s to have Type %s, got %s", tableName, e.Name, e.Type, g.Type)
		}

		if !bytes.Equal(g.Default, e.Default) {
			t.Errorf("Expected column %s.%s to have Default %v, got %v", tableName, e.Name, e.Default, g.Default)
		}

		if g.NotNull != e.NotNull {
			t.Errorf("Expected column %s.%s to have NotNull %v, got %v", tableName, e.Name, e.NotNull, g.NotNull)
		}

		if g.PrimaryKey != e.PrimaryKey {
			t.Errorf("Expected column %s.%s to have PrimaryKey %d, got %d", tableName, e.Name, e.PrimaryKey, g.PrimaryKey)
		}
	}
}

func verifyIndexes(t *testing.T, db *sql.DB, tableName string, exp []sqlitemeta.Index) {

	got := getIndexes(t, db, tableName)

	if len(got) != len(exp) {
		t.Errorf("Expected %d index(es) in table %s, got %d", len(exp), tableName, len(got))
		return
	}

	for i := range got {

		g := got[i]
		e := exp[i]

		if g.Name != e.Name {
			t.Errorf("Expected index %d in table %s to have Name %q, got %q", i, tableName, e.Name, g.Name)
		}

		if !reflect.DeepEqual(g.ColumnNames, e.ColumnNames) {
			t.Errorf("Expected index %d in table %s to have ColumnNames %v, got %v", i, tableName, e.ColumnNames, g.ColumnNames)
		}

		if g.Type != e.Type {
			t.Errorf("Expected index %d in table %s to have Type %q, got %q", i, tableName, e.Type, g.Type)
		}

		if g.IsUnique != e.IsUnique {
			t.Errorf("Expected index %d in table %s to have IsUnique %v, got %v", i, tableName, e.IsUnique, g.IsUnique)
		}

		if g.IsPartial != e.IsPartial {
			t.Errorf("Expected index %d in table %s to have IsPartial %v, got %v", i, tableName, e.IsPartial, g.IsPartial)
		}
	}
}

func verifyForeignKeys(t *testing.T, db *sql.DB, tableName string, exp []sqlitemeta.ForeignKey) {

	got := getForeignKeys(t, db, tableName)

	if len(got) != len(exp) {
		t.Errorf("Expected %d foreign key(s) in table %q, got %d", len(exp), tableName, len(got))
		return
	}

	for i := range got {

		g := got[i]
		e := exp[i]

		if g.ID != e.ID {
			t.Errorf("Expected foreign key %d in table %q to have ID %d, got %d", i, tableName, e.ID, g.ID)
		}

		if !equalStringSlices(g.ChildKey, e.ChildKey) {
			t.Errorf("Expected foreign key %d in table %q to have Child Key %q, got %q", i, tableName, e.ChildKey, g.ChildKey)
		}

		if !reflect.DeepEqual(g.ParentKey, e.ParentKey) {
			t.Errorf("Expected foreign key %d in table %q to have Parent Key %v, got %v", i, tableName, e.ParentKey, g.ParentKey)
		}

		if g.ParentTable != e.ParentTable {
			t.Errorf("Expected foreign key %d in table %q to have Parent Table %q, got %q", i, tableName, e.ParentTable, g.ParentTable)
		}

		if g.OnUpdate != e.OnUpdate {
			t.Errorf("Expected foreign key %d in table %q to have OnUpdate %q, got %q", i, tableName, e.OnUpdate, g.OnUpdate)
		}

		if g.OnDelete != e.OnDelete {
			t.Errorf("Expected foreign key %d in table %q to have OnDelete %q, got %q", i, tableName, e.OnDelete, g.OnDelete)
		}
	}
}

func verifyForeignKeysEnabled(t *testing.T, db *sql.DB, enabled bool) {
	if got := getForeignKeysEnabled(t, db); got != enabled {
		t.Errorf("Expected foreign keys enabled to be %v, got %v", enabled, got)
	}
}

func verifyFeed(t *testing.T, db *sql.DB, id int64, feed *kibner.Feed) {
	compareFeeds(t, feed, getFeed(t, db, id))
}

func compareFeeds(t *testing.T, feed, dbFeed *kibner.Feed) {

	if got, exp := dbFeed.Title, feed.Title; got != exp {
		t.Errorf("Expected feed %q to have Title %q, got %q", feed.Title, exp, got)
	}

	if got, exp := dbFeed.Author, feed.Author; got != exp {
		t.Errorf("Expected feed %q to have Author %q, got %q", feed.Title, exp, got)
	}

	if got, exp := dbFeed.Desc, feed.Desc; got != exp {
		t.Errorf("Expected feed %q to have Desc %q, got %q", feed.Title, exp, got)
	}

	if got, exp := dbFeed.Type, feed.Type; got != exp {
		t.Errorf("Expected feed %q to have Type %q, got %q", feed.Title, exp, got)
	}

	if got, exp := dbFeed.URL, feed.URL; got != exp {
		t.Errorf("Expected feed %q to have URL %q, got %q", feed.Title, exp, got)
	}

	if got, exp := dbFeed.Image, feed.Image; got != exp {
		t.Errorf("Expected feed %q to have Image %q, got %q", feed.Title, exp, got)
	}

	if got, exp := dbFeed.Link, feed.Link; got != exp {
		t.Errorf("Expected feed %q to have Link %q, got %q", feed.Title, exp, got)
	}

	if got, exp := len(dbFeed.Items), len(feed.Items); got != exp {
		t.Errorf("Expected feed %q to have %d item(s), got %d", feed.Title, exp, got)
		return
	}

	for i := range feed.Items {
		compareItems(t, feed.Title, i+1, feed.Items[i], dbFeed.Items[i])
	}
}

func verifyFeedView(t *testing.T, data *feedView, feed *kibner.Feed) {
	compareFeedView(t, feedViewFromFeed(feed), data)
}

func compareFeedView(t *testing.T, exp, got *feedView) {

	if got.Title != exp.Title {
		t.Errorf("Expected feedView for %q to have Title %q, got %q", exp.Title, exp.Title, got.Title)
	}

	if got.Author != exp.Author {
		t.Errorf("Expected feedView for %q to have Author %q, got %q", exp.Title, exp.Author, got.Author)
	}

	if got.Desc != exp.Desc {
		t.Errorf("Expected feedView for %q to have Desc %q, got %q", exp.Title, exp.Desc, got.Desc)
	}

	if got.Items != exp.Items {
		t.Errorf("Expected feedView for %q to have %d Items, got %d", exp.Title, exp.Items, got.Items)
	}

	if got.UnplayedItems != exp.UnplayedItems {
		t.Errorf("Expected feedView for %q to have UnplayedItems %d, got %d", exp.Title, exp.UnplayedItems, got.UnplayedItems)
	}

	if !got.LastPubdate.Equal(exp.LastPubdate) {
		t.Errorf("Expected feedView for %q to have LastPubdate %s, got %s", exp.Title, exp.LastPubdate.Format(time.RFC1123Z), got.LastPubdate.Format(time.RFC1123Z))
	}
}

func compareItems(t *testing.T, feedTitle string, i int, item, dbItem *kibner.Item) {

	if got, exp := dbItem.Title, item.Title; got != exp {
		t.Errorf("Expected %s item %d to have Title %q, got %q", feedTitle, i, exp, got)
	}

	if got, exp := dbItem.Pubdate, item.Pubdate; got != exp {
		t.Errorf("Expected %s item %d to have Pubdate %v, got %v", feedTitle, i, exp, got)
	}

	if got, exp := dbItem.Desc, item.Desc; got != exp {
		t.Errorf("Expected %s item %d to have Desc %q, got %q", feedTitle, i, exp, got)
	}

	if got, exp := dbItem.URL, item.URL; got != exp {
		t.Errorf("Expected %s item %d to have URL %q, got %q", feedTitle, i, exp, got)
	}

	if got, exp := dbItem.Duration, item.Duration; got != exp {
		t.Errorf("Expected %s item %d to have Duration %s, got %s", feedTitle, i, exp, got)
	}

	if got, exp := dbItem.Filesize, item.Filesize; got != exp {
		t.Errorf("Expected %s item %d to have Filesize %d, got %d", feedTitle, i, exp, got)
	}

	if got, exp := dbItem.GUID, item.GUID; got != exp {
		t.Errorf("Expected %s item %d to have GUID %q, got %q", feedTitle, i, exp, got)
	}
}

func verifyUnplayedItemCount(t *testing.T, db *sql.DB, feedID int64, count int) {

	if got := getUnplayedItemCount(t, db, feedID); got != count {
		t.Errorf("Expected %d unplayed items, got %d", count, got)
	}
}

func feedViewFromFeed(feed *kibner.Feed) *feedView {

	var pubdate time.Time
	if len(feed.Items) > 0 {
		pubdate = feed.Items[0].Pubdate
	}

	return &feedView{
		Title:         feed.Title,
		Author:        feed.Author,
		Desc:          feed.Desc,
		Items:         int64(len(feed.Items)),
		UnplayedItems: 0,
		LastPubdate:   pubdate,
	}
}

func getTableNames(t *testing.T, db *sql.DB) []string {
	return getMasterNames(t, db, "table")
}

func getMasterNames(t *testing.T, db *sql.DB, typ string) []string {

	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type = ? ORDER BY name", typ)
	if err != nil {
		t.Fatalf("Error querying sqlite_master for %q names: %s", typ, err)
	}
	defer rows.Close()

	var names []string

	for rows.Next() {

		var name string

		if err = rows.Scan(&name); err != nil {
			t.Fatalf("Error scanning %q rows from sqlite_master: %s", typ, err)
		}

		names = append(names, name)
	}

	if err = rows.Err(); err != nil {
		t.Fatalf("Error iterating over %q rows from sqlite_master: %s", typ, err)
	}

	return names
}

func getSequenceNumbers(t *testing.T, db *sql.DB) map[string]int64 {

	rows, err := db.Query("SELECT name, seq FROM sqlite_sequence")
	if err != nil {
		t.Fatalf("Error querying sqlite_sequence: %s", err)
	}
	defer rows.Close()

	var seqs map[string]int64

	for rows.Next() {

		var name string
		var seq int64

		if err = rows.Scan(&name, &seq); err != nil {
			t.Fatalf("Error scanning sqlite_sequence: %s", err)
		}

		if _, ok := seqs[name]; ok {
			t.Fatalf("Error scanning sqlite_sequence: duplicate name %q", name)
		}

		if seqs == nil {
			seqs = make(map[string]int64)
		}

		seqs[name] = seq
	}

	if err = rows.Err(); err != nil {
		t.Fatalf("Error iterating over rows from sqlite_sequence: %s", err)
	}

	return seqs
}

func getRowCount(t *testing.T, db *sql.DB, tableName string) int {

	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM " + tableName).Scan(&count)
	if err != nil {
		t.Fatalf("Error querying row count for table %q: %s", tableName, err)
	}

	return count
}

func getColumns(t *testing.T, db *sql.DB, tableName string) []sqlitemeta.Column {

	columns, err := sqlitemeta.Columns(db, tableName)
	if err != nil {
		t.Fatalf("Error querying %s columns: %s", tableName, err)
	}

	return columns
}

func getIndexes(t *testing.T, db *sql.DB, tableName string) []sqlitemeta.Index {

	indexes, err := sqlitemeta.Indexes(db, tableName)
	if err != nil {
		t.Fatalf("Error querying indexes: %s", err)
	}

	return indexes
}

func getForeignKeys(t *testing.T, db *sql.DB, tableName string) []sqlitemeta.ForeignKey {

	foreignKeys, err := sqlitemeta.ForeignKeys(db, tableName)
	if err != nil {
		t.Fatalf("Error querying foreign keys for %s: %s", tableName, err)
	}

	return foreignKeys
}

func getForeignKeysEnabled(t *testing.T, db *sql.DB) bool {

	var enabled bool
	err := db.QueryRow("SELECT foreign_keys FROM pragma_foreign_keys").Scan(&enabled)
	if err != nil {
		t.Fatalf("Error querying foreign key status: %s", err)
	}

	return enabled
}

func getFeed(t *testing.T, db *sql.DB, id int64) *kibner.Feed {

	q :=
		`SELECT
			title,
			author,
			desc,
			type,
			url,
            image,
            link,
			timestamp
        FROM
            feeds
        WHERE
            id = ?
        ORDER BY
            title`

	var title string
	var author string
	var desc string
	var typ string
	var url string
	var image string
	var link string
	var timestamp time.Time

	err := db.QueryRow(q, id).Scan(&title, &author, &desc, &typ, &url, &image, &link, &timestamp)
	if err != nil {
		t.Fatalf("Error querying table %q: %s", "feeds", err)
	}

	if timestamp.IsZero() {
		t.Errorf("Feed %q has zero timestamp", title)
	}

	return &kibner.Feed{
		Title:  title,
		Author: author,
		Desc:   desc,
		Type:   typ,
		URL:    url,
		Image:  image,
		Link:   link,
		Items:  getItems(t, db, id),
	}
}

func getItems(t *testing.T, db *sql.DB, feedID int64) []*kibner.Item {

	q :=
		`SELECT
			title,
			desc,
			pubdate,
			url,
			filesize,
			duration,
			guid,
			unplayed,
            timestamp
        FROM
            items
        WHERE
            feedid = ?
        ORDER BY
            pubdate DESC, GUID`

	rows, err := db.Query(q, feedID)
	if err != nil {
		t.Fatalf("Error querying table %q: %s", "items", err)
	}
	defer rows.Close()

	var items []*kibner.Item

	for rows.Next() {

		var title string
		var desc string
		var pubdate time.Time
		var url string
		var filesize int64
		var duration int
		var guid string
		var unplayed bool
		var timestamp time.Time

		if err := rows.Scan(&title, &desc, &pubdate, &url, &filesize, &duration, &guid, &unplayed, &timestamp); err != nil {
			t.Fatalf("Error scanning rows from %q: %s", "items", err)
		}

		if timestamp.IsZero() {
			t.Errorf("Expected item %q to have a non-zero timestamp, got zero", title)
		}

		items = append(items, &kibner.Item{
			Title:    title,
			Desc:     desc,
			Pubdate:  pubdate,
			URL:      url,
			Filesize: filesize,
			Duration: time.Duration(duration) * time.Second,
			GUID:     guid,
		})
	}

	if err = rows.Err(); err != nil {
		t.Fatalf("Error iterating over rows from %q: %s", "items", err)
	}

	return items
}

func getUnplayedItemCount(t *testing.T, db *sql.DB, feedID int64) int {

	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM items WHERE feedid = ? AND unplayed = 1", feedID).Scan(&count)
	if err != nil {
		t.Fatalf("Error querying unplayed items: %s", err)
	}

	return count
}

func markUnplayedItems(t *testing.T, db *sql.DB, feedID int64, count int) {

	q :=
		`UPDATE
			items
		SET
			unplayed = 1
		WHERE
			ROWID in
				(SELECT
					ROWID
				FROM
					items
				WHERE
					feedid = ?
				ORDER BY
					pubdate DESC
				LIMIT ?)`

	res, err := db.Exec(q, feedID, count)
	if err != nil {
		t.Fatalf("markUnplayedItems returned error %q", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		t.Fatalf("markUnplayedItems.RowsAffected returned error %q", err)
	}

	if n != int64(count) {
		t.Fatalf("Expected markUnplayedItems.RowsAffected to return %d rows, got %d", count, n)
	}
}

func writeOPMLHeader(t *testing.T, w io.Writer, timestamp time.Time) {
	fmt.Fprintln(w, `<?xml version="1.0" encoding="UTF-8"?>`)
	fmt.Fprintln(w, `<opml version="2.0">`)
	fmt.Fprintln(w, `    <head>`)
	fmt.Fprintln(w, `        <title>Kibner subscriptions</title>`)
	fmt.Fprintf(w, `        <dateCreated>%s</dateCreated>`, timestamp.Format(time.RFC1123Z))
	fmt.Fprintln(w)
	fmt.Fprintln(w, `        <docs>http://dev.opml.org/spec2.html</docs>`)
	fmt.Fprintln(w, `    </head>`)
	fmt.Fprintln(w, `    <body>`)
}

func writeOPMLFooter(t *testing.T, w io.Writer) {
	fmt.Fprintln(w, `    </body>`)
	fmt.Fprint(w, `</opml>`)
}

func writeOPMLEntry(t *testing.T, w io.Writer, typ, title, desc, url string) {
	fmt.Fprint(w, `        <outline type="`)
	xml.EscapeText(w, []byte(typ))
	fmt.Fprint(w, `" text="`)
	xml.EscapeText(w, []byte(title))
	fmt.Fprint(w, `" title="`)
	xml.EscapeText(w, []byte(title))
	fmt.Fprint(w, `" description="`)
	xml.EscapeText(w, []byte(desc))
	fmt.Fprint(w, `" xmlUrl="`)
	xml.EscapeText(w, []byte(url))
	fmt.Fprintln(w, `"></outline>`)
}

func equalErrors(exp, got error) bool {

	switch {
	case got == exp:
		return true
	case got == nil, exp == nil:
		return false
	case got.Error() == exp.Error():
		return true
	default:
		return false
	}
}

func equalStringSlices(got, exp []string) bool {

	if len(got) != len(exp) {
		return false
	}

	for i := range got {
		if got[i] != exp[i] {
			return false
		}
	}

	return true
}

func jsonify(v interface{}) string {

	w := &bytes.Buffer{}

	e := json.NewEncoder(w)
	e.SetEscapeHTML(false)
	e.SetIndent("", "  ")

	if err := e.Encode(v); err != nil {
		panic(err)
	}

	return w.String()
}

func randomString(size int) string {

	chars := " 0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	n := len(chars)
	s := make([]byte, size)

	for i := 0; i < len(s); i++ {
		s[i] = chars[rand.Intn(n)]
	}

	return string(s)
}

func mapKeys(m interface{}) []string {

	keys := reflect.ValueOf(m).MapKeys()
	if len(keys) == 0 {
		return nil
	}

	vals := make([]string, len(keys))
	for i := range keys {
		vals[i] = keys[i].String()
	}

	return vals
}

func serverURL(ts *httptest.Server, path string) string {
	return ts.URL + "/" + strings.TrimLeft(path, "/")
}

func nullStrings(vals ...string) []sql.NullString {

	if len(vals) == 0 {
		return nil
	}

	ns := make([]sql.NullString, len(vals))

	for i := range vals {
		ns[i].Valid = true
		ns[i].String = vals[i]
	}

	return ns
}
