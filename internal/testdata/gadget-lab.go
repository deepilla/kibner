package testdata

import (
	"time"

	kibner "github.com/deepilla/kibner/internal/types"
)

func GadgetLab() *kibner.Feed {
	return &kibner.Feed{
		Title:  "The Gadget Lab Podcast",
		Author: "WIRED",
		Desc:   "The latest tech news from WIRED.",
		Type:   "rss",
		URL:    "https://www.wired.com/feed/podcast/gadget-lab",
		Link:   "https://www.wired.com/",
		Image:  "https://www.wired.com/wp-content/uploads/2016/09/Gadget-Lab-Podcast-3000.jpg",
		Items: []*kibner.Item{
			{
				Title:    "Google Recap",
				Desc:     "Google is now truly a hardware player, but it has a powerful AI software engine running across all of its gadgets.",
				Pubdate:  time.Date(2017, time.October, 5, 21, 29, 51, 0, time.UTC),
				URL:      "https://www.podtrac.com/pts/redirect.mp3/https://www.wired.com/podcast-download/2264829/gadget-lab-podcast-332.mp3",
				Filesize: 37036980,
				Duration: 38*time.Minute + 30*time.Second,
				GUID:     "https://www.wired.com/?p=2264829",
			},
			{
				Title:    "So Many Alexas",
				Desc:     "This week, new hardware and software from Amazon, and Twitter goes long.",
				Pubdate:  time.Date(2017, time.September, 29, 20, 17, 46, 0, time.UTC),
				URL:      "https://www.podtrac.com/pts/redirect.mp3/https://www.wired.com/podcast-download/2263889/gadget-lab-podcast-331.mp3",
				Filesize: 50364000,
				Duration: 52*time.Minute + 23*time.Second,
				GUID:     "https://www.wired.com/?p=2263889",
			},
			{
				Title:    "Scenes From Cupertino",
				Desc:     "We recap all of the news from Apple, offer analysis of the iPhone X, and give you advice about which iPhone to buy.",
				Pubdate:  time.Date(2017, time.September, 13, 1, 58, 49, 0, time.UTC),
				URL:      "https://www.podtrac.com/pts/redirect.mp3/https://www.wired.com/podcast-download/2262822/gadget-lab-podcast-330.mp3",
				Filesize: 57140096,
				Duration: 1*time.Hour + 7*time.Minute + 57*time.Second,
				GUID:     "https://www.wired.com/?p=2262822",
			},
			{
				Title:    "Apple Preview",
				Desc:     "The hosts open their iPhones and ask Siri to predict whats coming from Apple on Tuesday.",
				Pubdate:  time.Date(2017, time.September, 8, 17, 50, 39, 0, time.UTC),
				URL:      "https://www.podtrac.com/pts/redirect.mp3/https://www.wired.com/podcast-download/2261964/gadget-lab-podcast-329.mp3",
				Filesize: 34688558,
				Duration: 41*time.Minute + 13*time.Second,
				GUID:     "https://www.wired.com/?p=2261964",
			},
			{
				Title:    "Toy Wars",
				Desc:     "Its a special Force Friday edition of the show this week, with guest Brendan Nystedt.",
				Pubdate:  time.Date(2017, time.September, 1, 19, 17, 15, 0, time.UTC),
				URL:      "https://www.podtrac.com/pts/redirect.mp3/https://www.wired.com/podcast-download/2261645/gadget-lab-podcast-328.mp3",
				Filesize: 27362248,
				Duration: 32*time.Minute + 30*time.Second,
				GUID:     "https://www.wired.com/?p=2261645",
			},
			{
				Title:    "Note to Self",
				Desc:     "On this weeks Gadget Lab podcast, we run down the news on the Samsung Galaxy Note 8.",
				Pubdate:  time.Date(2017, time.August, 25, 16, 56, 26, 0, time.UTC),
				URL:      "https://www.podtrac.com/pts/redirect.mp3/https://www.wired.com/podcast-download/2260956/gadget-lab-podcast-327.mp3",
				Filesize: 42980090,
				Duration: 51*time.Minute + 6*time.Second,
				GUID:     "https://www.wired.com/?p=2260956",
			},
			{
				Title:    "New Phone, Who Dis?",
				Desc:     "The Essential Android phone hits the market in a couple of weeks. Heres what you need to know about it.",
				Pubdate:  time.Date(2017, time.August, 18, 23, 20, 43, 0, time.UTC),
				URL:      "https://www.podtrac.com/pts/redirect.mp3/https://www.wired.com/podcast-download/2260382/gadget-lab-podcast-326.mp3",
				Filesize: 48635494,
				Duration: 57*time.Minute + 50*time.Second,
				GUID:     "https://www.wired.com/?p=2260382",
			},
			{
				Title:    "Mickey's Out",
				Desc:     "This week: Disney ditches Netflix, bundles await, and Facebook makes TV now.",
				Pubdate:  time.Date(2017, time.August, 11, 18, 40, 29, 0, time.UTC),
				URL:      "https://www.podtrac.com/pts/redirect.mp3/https://www.wired.com/podcast-download/2259647/gadget-lab-podcast-325.mp3",
				Filesize: 29522916,
				Duration: 35*time.Minute + 5*time.Second,
				GUID:     "https://www.wired.com/?p=2259647",
			},
			{
				Title:    "The iPhone of Cars",
				Desc:     "We talk about the Tesla Model 3, the forthcoming iPhone, and how to spend that money burning a hole in your pocket.",
				Pubdate:  time.Date(2017, time.August, 4, 19, 20, 21, 0, time.UTC),
				URL:      "https://www.podtrac.com/pts/redirect.mp3/https://www.wired.com/podcast-download/2259230/gadget-lab-podcast-324.mp3",
				Filesize: 29379215,
				Duration: 40*time.Minute + 43*time.Second,
				GUID:     "https://www.wired.com/?p=2259230",
			},
			{
				Title:    "Muzzle Your Phone Before It Eats You Alive",
				Desc:     "Notifications and how to kill them all, this week on the Gadget Lab Podcast.",
				Pubdate:  time.Date(2017, time.July, 28, 19, 7, 20, 0, time.UTC),
				URL:      "https://www.podtrac.com/pts/redirect.mp3/https://www.wired.com/podcast-download/2258718/gadget-lab-podcast-323.mp3",
				Filesize: 38883851,
				Duration: 53*time.Minute + 56*time.Second,
				GUID:     "https://www.wired.com/?p=2258718",
			},
		},
	}
}
