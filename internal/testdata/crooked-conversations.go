package testdata

import (
	"time"

	kibner "github.com/deepilla/kibner/internal/types"
)

func Crooked() *kibner.Feed {
	return &kibner.Feed{
		Title:  "Crooked Conversations",
		Author: "Crooked Media",
		Desc:   "One side effect of our national addiction to Trump’s tweets and other news cycle garbage is that fascinating issues, brilliant books and important debates aren't getting the attention they deserve. With a rotating crew of your favorite Crooked Media hosts, contributors, and special guests, Crooked Conversations brings Pod Save America's no-b.s., conversational style to topics in politics, media, culture, sports, and technology that aren’t making headlines but still have a major impact on our world.",
		Type:   "rss",
		URL:    "http://feeds.feedburner.com/crooked-conversations",
		Link:   "https://art19.com/shows/crooked-conversations",
		Image:  "https://dfkfj8j276wwv.cloudfront.net/images/a9/b3/c8/d4/a9b3c8d4-8aa1-4859-8ca0-cd3e135c877a/c65b236279e39274b31c79b63c48e50dd32bb20d9f20b3da854ecb776a58be5113da5a81aba365d1a99b9ebdb241cb89f4eae81fb69e71304218c5bf2d9efb11.jpeg",
		Items: []*kibner.Item{
			{
				Title:    "Introduction to Crooked Conversations",
				Desc:     "With a rotating crew of your favorite Crooked Media hosts, contributors, and special guests, Crooked Conversations brings Pod Save America's no-b.s., conversational style to topics in politics, media, culture, sports, and technology that aren’t making headlines but still have a major impact on our world.",
				Pubdate:  time.Date(2017, time.October, 4, 13, 52, 37, 0, time.UTC),
				URL:      "http://rss.art19.com/episodes/bae66e2a-6e51-4c1e-bbf0-e19f358fa1f4.mp3",
				Filesize: 2014145,
				Duration: 2*time.Minute + 5*time.Second,
				GUID:     "gid://art19-episode-locator/V0/C7asogvo9pu1cq-ex_MehFnndlyjXj-C9YErvzdn22M",
			},
		},
	}
}
