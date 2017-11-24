package testdata

import (
	"time"

	kibner "github.com/deepilla/kibner/internal/types"
)

func Rework() *kibner.Feed {
	return &kibner.Feed{
		Title:  "REWORK",
		Author: "Basecamp",
		Desc:   "A podcast by Basecamp about a better way to work and run your business. We bring you stories and unconventional wisdom from Basecamp’s co-founders and other business owners.",
		Type:   "rss",
		URL:    "http://feeds.feedburner.com/basecamp/rework",
		Link:   "https://rework.fm",
		Image:  "https://dfkfj8j276wwv.cloudfront.net/images/1d/c9/97/91/1dc99791-61f8-4e5b-b6a3-69bb542fb652/7a36227f38a6bea4d8ec172a982500cb5dde7967886a27b2dc8d87e4cc19bb3c8a508bd9cafd4fa9214d7f564b3f76706b912b54a1a7c8bc92825e2db21e5dd7.jpeg",
		Items: []*kibner.Item{
			{
				Title:    "04 - Say No",
				Desc:     "It's easy to say yes, whether it's to a customer request or a deadline from your boss. But saying yes too many times can result in an unmanageable workload or distract you from the stuff you really want to be doing. It's good to practice saying no and se",
				Pubdate:  time.Date(2017, time.September, 26, 12, 0, 0, 0, time.UTC),
				URL:      "http://feedproxy.google.com/~r/basecamp/rework/~5/pmV2EJ-hBIM/35a5b45e-acf3-4805-b1f0-a7e81ed5bcd0.mp3",
				Filesize: 24444342,
				Duration: 25*time.Minute + 27*time.Second,
				GUID:     "gid://art19-episode-locator/V0/blkpmlLBlMvbwAWdGWJ195y5ssuMHermE1vqahTz4pk",
			},
			{
				Title:    "03 - Pick A Fight (on Twitter)",
				Desc:     "Basecamp CTO David Heinemeier Hansson is known for many things, including creating Ruby on Rails and writing business books. He also has a knack for arguing with people on the Internet. This cheerfully profane conversation explores how Twitter is like a",
				Pubdate:  time.Date(2017, time.September, 12, 12, 0, 0, 0, time.UTC),
				URL:      "http://feedproxy.google.com/~r/basecamp/rework/~5/4dRxfaiKp7o/77508be1-639d-4d21-a838-6dc8bc2cb35e.mp3",
				Filesize: 27760431,
				Duration: 28*time.Minute + 55*time.Second,
				GUID:     "gid://art19-episode-locator/V0/4r8rZpBSOkQMYpB0e-K71awM6AbNgGwWB1yYeMqohYE",
			},
			{
				Title:    "02 - Workaholics Aren't Heroes",
				Desc:     "Being tired isn't a badge of honor. There, we said it. We've been saying this for a while now, because our culture loves to glorify toiling long hours for its own sake and we think that leads to subpar work and general misery. In this episode, we talk to",
				Pubdate:  time.Date(2017, time.August, 29, 12, 0, 0, 0, time.UTC),
				URL:      "http://feedproxy.google.com/~r/basecamp/rework/~5/hYJu_xPRg7E/ca99d845-34be-49b6-88fb-d8dd5dee5517.mp3",
				Filesize: 31561351,
				Duration: 32*time.Minute + 52*time.Second,
				GUID:     "gid://art19-episode-locator/V0/swCEXakdVEDQbhdznJJ0f0-C1ojW7fAWbnkZQMG6yHU",
			},
			{
				Title:    "01 - Sell Your By-products",
				Desc:     "Welcome to the first episode of Rework! This podcast is based on Jason Fried and David Heinemeier Hansson's 2010 best-selling business book, which was itself based on years of blogging. So what better way to kick off this show than talking about byproduc",
				Pubdate:  time.Date(2017, time.August, 15, 12, 0, 0, 0, time.UTC),
				URL:      "http://feedproxy.google.com/~r/basecamp/rework/~5/6m3PFvpdw34/33f6278f-cc8d-4694-b599-b237eebad022.mp3",
				Filesize: 29878230,
				Duration: 31*time.Minute + 7*time.Second,
				GUID:     "gid://art19-episode-locator/V0/ElO45FaLvlp3ZI4hJcZpVpv_2wAYGMcMR4dWJrjkSgQ",
			},
			{
				Title:    "Rework Teaser",
				Desc:     "Rework is a podcast by the makers of Basecamp about a better way to work and run your business. While the prevailing narrative around successful entrepreneurship tells you to scale fast and raise money, we think there's a better way. We'll take you behin",
				Pubdate:  time.Date(2017, time.July, 26, 20, 42, 14, 0, time.UTC),
				URL:      "http://feedproxy.google.com/~r/basecamp/rework/~5/SUdV0TgEsh8/cfc2cc54-bb15-42f2-9e02-5d8708c7d6e4.mp3",
				Filesize: 1446556,
				Duration: 1*time.Minute + 30*time.Second,
				GUID:     "gid://art19-episode-locator/V0/FP72lI08fQa5uymp3rxmu_gTBY3uDNX2zCjMFebpzZc",
			},
		},
	}
}
