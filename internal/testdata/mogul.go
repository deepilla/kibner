package testdata

import (
	"time"

	kibner "github.com/deepilla/kibner/internal/types"
)

func Mogul() *kibner.Feed {
	return &kibner.Feed{
		Title:  "Mogul: The Life and Death of Chris Lighty",
		Author: "Gimlet",
		Desc:   "Chris Lighty was a giant in hip-hop. He managed Foxy Brown, Fat Joe, Missy Elliott, Busta Rhymes, LL Cool J, 50 Cent—anyone who was anyone worked with Lighty. But in 2012 he was found dead at his home in the Bronx, a death that left the music world reelin",
		Type:   "rss",
		URL:    "http://feeds.gimletmedia.com/mogulshow",
		Link:   "https://gimletmedia.com/mogul/",
		Image:  "http://static.megaphone.fm/podcasts/44250a2c-2089-11e7-a1f3-bb6a31bc38f5/image/uploads_2F1497559518113-jto0imaspjo-698c8877f94c23bd4f7298fc7abe1d94_2Fcover_art_moguldarker.png",
		Items: []*kibner.Item{
			{
				Title:    "Part 6: August 30, 2012",
				Desc:     "August 30th, 2012. A day that shook hip hop.",
				Pubdate:  time.Date(2017, time.July, 28, 4, 0, 0, 0, time.UTC),
				URL:      "http://traffic.megaphone.fm/GLT6609848575.mp3",
				Filesize: 53570037,
				Duration: 2232 * time.Second,
				GUID:     "e777ff7c-208e-11e7-bea3-b76f85198655",
			},
			{
				Title:    "Part 5: How Heavy It Was",
				Desc:     "In this episode: cold hard cash.",
				Pubdate:  time.Date(2017, time.July, 21, 4, 0, 0, 0, time.UTC),
				URL:      "http://traffic.megaphone.fm/GLT5538340591.mp3",
				Filesize: 49766400,
				Duration: 2073 * time.Second,
				GUID:     "e7898378-208e-11e7-bea3-b33c6c8c60aa",
			},
			{
				Title:    "Part 4: Gucci Boots",
				Desc:     "Lighty is at the top of his game.",
				Pubdate:  time.Date(2017, time.July, 14, 4, 0, 0, 0, time.UTC),
				URL:      "http://traffic.megaphone.fm/GLT4758054050.mp3",
				Filesize: 45338331,
				Duration: 1889 * time.Second,
				GUID:     "e78091a0-208e-11e7-bea3-2b6946537b3f",
			},
			{
				Title:    "Part 3: Rice Pilaf",
				Desc:     "Chris Lighty meets Warren G.",
				Pubdate:  time.Date(2017, time.June, 30, 4, 0, 0, 0, time.UTC),
				URL:      "http://traffic.megaphone.fm/GLT2905886229.mp3",
				Filesize: 43274449,
				Duration: 1803 * time.Second,
				GUID:     "e76f3ce8-208e-11e7-bea3-df8ed23719ba",
			},
			{
				Title:    "Part 2: Not Just Me and Snakes",
				Desc:     "Chris is headed for the big time.",
				Pubdate:  time.Date(2017, time.June, 23, 4, 0, 0, 0, time.UTC),
				URL:      "http://traffic.megaphone.fm/GLT5760318992.mp3",
				Filesize: 40802429,
				Duration: 1700 * time.Second,
				GUID:     "e76661b8-208e-11e7-bea3-6b2aee15fcd0",
			},
			{
				Title:    "Part 1: That Beat, That Beat Right There",
				Desc:     "Let’s start at the end—at a funeral.",
				Pubdate:  time.Date(2017, time.June, 16, 0, 0, 0, 0, time.UTC),
				URL:      "http://traffic.megaphone.fm/GLT1554207124.mp3",
				Filesize: 48110027,
				Duration: 2004 * time.Second,
				GUID:     "e75cb49c-208e-11e7-bea3-5bed2eb5f7c7",
			},
			{
				Title:    "Trailer",
				Desc:     "One man’s story, from the first breakbeat to the last heartbeat. A hip-hop miniseries from Gimlet Media and Loud Speakers Network, hosted by Reggie Ossé.",
				Pubdate:  time.Date(2017, time.June, 10, 14, 35, 40, 0, time.UTC),
				URL:      "http://traffic.megaphone.fm/GLT5491724384.mp3",
				Filesize: 4994194,
				Duration: 208 * time.Second,
				GUID:     "162fd028-208b-11e7-8d13-b7ad5ca2b7cf",
			},
		},
	}
}
