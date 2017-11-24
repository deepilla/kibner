package testdata

import kibner "github.com/deepilla/kibner/internal/types"
import "time"

func MinimalFeed() *kibner.Feed {
	return &kibner.Feed{
		Title:  "Minimal Feed",
		Author: "Unknown Author",
		Type:   "rss",
		URL:    "http://deepilla.com/feeds/test-feed-minimal.xml",
		Items: []*kibner.Item{
			{
				Title:   "Untitled: Sep 21, 2015",
				Pubdate: time.Date(2015, time.September, 21, 1, 5, 38, 0, time.UTC),
				URL:     "http://deepilla.com/assets/mp3/minimal-2.mp3",
				GUID:    "http://deepilla.com/assets/mp3/minimal-2.mp3",
			},
			{
				Title: "Untitled",
				URL:   "http://deepilla.com/assets/mp3/minimal-1.mp3",
				GUID:  "http://deepilla.com/assets/mp3/minimal-1.mp3",
			},
		},
	}
}
