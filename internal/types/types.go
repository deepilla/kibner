// Package types contains type definitions shared between
// the main package and the "testdata" testing package.
// It is usually imported with the alias "kibner".
package types

import "time"

// Feed represents a podcast feed.
type Feed struct {
	Title  string
	Author string
	Desc   string
	Type   string
	URL    string
	Link   string
	Image  string
	Items  []*Item
}

// Item represents an individual podcast episode.
type Item struct {
	Title    string
	Desc     string
	Pubdate  time.Time
	URL      string
	Filesize int64
	Duration time.Duration
	GUID     string
}
