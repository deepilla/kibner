package main

import (
	"fmt"
	"regexp"

	"github.com/mmcdole/gofeed"
	"github.com/mmcdole/gofeed/rss"
)

// RSSTranslator is a stripped down gofeed translator
// for RSS feeds. It only handles fields that kibner
// is interested in.
type RSSTranslator struct{}

// NewRSSTranslator creates a new RSSTranslator.
func NewRSSTranslator() gofeed.Translator {
	return &RSSTranslator{}
}

// Translate converts an RSS Feed into a gofeed Feed.
// Note that not all of the gofeed fields will be set
// -- just the ones that kibner needs.
func (t *RSSTranslator) Translate(feed interface{}) (*gofeed.Feed, error) {

	rss, ok := feed.(*rss.Feed)
	if !ok {
		return nil, fmt.Errorf("Feed did not match expected type of *rss.Feed")
	}

	return &gofeed.Feed{
		Title:       t.translateFeedTitle(rss),
		Author:      t.translateFeedAuthor(rss),
		Link:        t.translateFeedLink(rss),
		Image:       t.translateFeedImage(rss),
		Description: t.translateFeedDescription(rss),
		Items:       t.translateFeedItems(rss),
		ITunesExt:   rss.ITunesExt,
		FeedVersion: rss.Version,
		FeedType:    "rss",
	}, nil
}

func (t *RSSTranslator) translateFeedItem(item *rss.Item) *gofeed.Item {

	return &gofeed.Item{
		Title:           t.translateItemTitle(item),
		GUID:            t.translateItemGUID(item),
		Description:     t.translateItemDescription(item),
		Enclosures:      t.translateItemEnclosures(item),
		PublishedParsed: item.PubDateParsed,
		ITunesExt:       item.ITunesExt,
	}
}

func (t *RSSTranslator) translateFeedTitle(rss *rss.Feed) string {
	return rss.Title
}

func (t *RSSTranslator) translateFeedAuthor(rss *rss.Feed) *gofeed.Person {

	var author string

	// Prefer the iTunes author as it tends to be more accurate.
	// Managing Editor and Webmaster may contain contact info
	// for the feed maintainer or some other arbitrary value (e.g.
	// "SoundCloud Feeds" for podcasts hosted on SoundCloud).
	switch {
	case rss.ITunesExt != nil && rss.ITunesExt.Author != "":
		author = rss.ITunesExt.Author
	case rss.ManagingEditor != "":
		author = rss.ManagingEditor
	case rss.WebMaster != "":
		author = rss.WebMaster
	default:
		return nil
	}

	name, email := parseNameEmail(author)
	if name == "" {
		return nil
	}

	return &gofeed.Person{
		Name:  name,
		Email: email,
	}
}

func (t *RSSTranslator) translateFeedLink(rss *rss.Feed) string {

	switch {
	case rss.Link != "":
		return rss.Link
	case rss.Image != nil && rss.Image.Link != "":
		return rss.Image.Link
	default:
		return ""
	}
}

func (t *RSSTranslator) translateFeedImage(rss *rss.Feed) *gofeed.Image {

	var url, title string

	// Prefer the iTunes image as it tends to be larger
	// (per Apple's requirements). Fall back to the RSS
	// image.
	switch {
	case rss.ITunesExt != nil && rss.ITunesExt.Image != "":
		url = rss.ITunesExt.Image
	case rss.Image != nil && rss.Image.URL != "":
		url = rss.Image.URL
		title = rss.Image.Title
	default:
		return nil
	}

	return &gofeed.Image{
		URL:   url,
		Title: title,
	}
}

func (t *RSSTranslator) translateFeedDescription(rss *rss.Feed) string {

	title := t.translateFeedTitle(rss)

	// Prefer the iTunes subtitle as publishers tend to
	// keep it short and sweet.
	if rss.ITunesExt != nil && rss.ITunesExt.Subtitle != "" && rss.ITunesExt.Subtitle != title {
		return rss.ITunesExt.Subtitle
	}

	// Fall back to the shorter of the itunes summary and
	// the channel description.
	descs := []string{
		rss.Description,
	}
	if rss.ITunesExt != nil {
		descs = append(descs, rss.ITunesExt.Summary)
	}

	return shortestDescription(descs, title)
}

func (t *RSSTranslator) translateFeedItems(rss *rss.Feed) []*gofeed.Item {

	results := make([]*gofeed.Item, len(rss.Items))
	for i := range rss.Items {
		results[i] = t.translateFeedItem(rss.Items[i])
	}

	return results
}

func (t *RSSTranslator) translateItemTitle(item *rss.Item) string {
	return item.Title
}

func (t *RSSTranslator) translateItemGUID(item *rss.Item) string {

	if item.GUID == nil {
		return ""
	}

	return item.GUID.Value
}

func (t *RSSTranslator) translateItemEnclosures(item *rss.Item) []*gofeed.Enclosure {

	if item.Enclosure == nil {
		return nil
	}

	return []*gofeed.Enclosure{
		{
			URL:    item.Enclosure.URL,
			Type:   item.Enclosure.Type,
			Length: item.Enclosure.Length,
		},
	}
}

func (t *RSSTranslator) translateItemDescription(item *rss.Item) string {

	var descs []string

	// Use the shortest of the iTunes subtitle, iTunes summary
	// and RSS item description.
	if itunes := item.ITunesExt; itunes != nil {
		descs = append(descs, itunes.Subtitle, itunes.Summary)
	}
	descs = append(descs, item.Description)

	return shortestDescription(descs, t.translateItemTitle(item))
}

func shortestDescription(descs []string, title string) string {

	var desc string

	for _, s := range descs {
		if s == "" || s == desc || s == title {
			continue
		}
		if len(desc) > 0 && len(desc) <= len(s) {
			continue
		}
		desc = s
	}

	return desc
}

// Name and email parsing adapted from the internal package
// gofeed/internals/shared.

var (
	rxNameEmail = regexp.MustCompile(`^([^@]+)\s+\(([^@]+@[^)]+)\)$`)
	rxEmailName = regexp.MustCompile(`^([^@]+@[^\s]+)\s+\(([^@]+)\)$`)
	rxNameOnly  = regexp.MustCompile(`^([^@()]+)$`)
	rxEmailOnly = regexp.MustCompile(`^([^@()]+@[^@()]+)$`)
)

func parseNameEmail(text string) (string, string) {

	if text == "" {
		return "", ""
	}

	var name, email string

	switch {
	case rxNameEmail.MatchString(text):
		matches := rxNameEmail.FindStringSubmatch(text)
		name = matches[1]
		email = matches[2]
	case rxEmailName.MatchString(text):
		matches := rxEmailName.FindStringSubmatch(text)
		email = matches[1]
		name = matches[2]
	case rxNameOnly.MatchString(text):
		matches := rxNameOnly.FindStringSubmatch(text)
		name = matches[1]
	case rxEmailOnly.MatchString(text):
		matches := rxEmailOnly.FindStringSubmatch(text)
		email = matches[1]
	}

	return name, email
}
