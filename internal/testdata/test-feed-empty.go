package testdata

import kibner "github.com/deepilla/kibner/internal/types"

func EmptyFeed() *kibner.Feed {
	return &kibner.Feed{
		Title:  "Empty Feed",
		Author: "deepilla",
		Desc:   "This is a test RSS feed with no items and no iTunes fields.",
		Type:   "rss",
		URL:    "http://deepilla.com/feeds/test-feed-empty.xml",
		Link:   "http://deepilla.com",
		Image:  "http://deepilla.com/assets/img/artwork.jpg",
	}
}
