package testdata

import (
	"time"

	kibner "github.com/deepilla/kibner/internal/types"
)

func STown() *kibner.Feed {
	return &kibner.Feed{
		Title:  "S-Town",
		Author: "Serial & This American Life",
		Desc:   "S-Town is a new podcast from Serial and This American Life, hosted by Brian Reed, about a man named John who despises his Alabama town and decides to do something about it. He asks Brian to investigate the son of a wealthy family who's allegedly been brag",
		Type:   "rss",
		URL:    "http://feeds.stownpodcast.org/stownpodcast",
		Link:   "https://stownpodcast.org",
		Image:  "https://files.stownpodcast.org/img/s-town-itunes.jpg",
		Items: []*kibner.Item{
			{
				Title:    "Chapter I",
				Desc:     "“If you keep your mouth shut, you’ll be surprised what you can learn.”",
				Pubdate:  time.Date(2017, time.March, 28, 10, 0, 0, 0, time.UTC),
				URL:      "https://dts.podtrac.com/redirect.mp3/dovetail.prxu.org/stown/817df96e-1d5e-48b2-964c-1f31d8c2d7ff/s-town-ch01.mp3",
				Duration: 53 * time.Minute,
				GUID:     "e01",
			},
			{
				Title:    "Chapter II",
				Desc:     "“Has anybody called you?”",
				Pubdate:  time.Date(2017, time.March, 28, 9, 45, 0, 0, time.UTC),
				URL:      "https://dts.podtrac.com/redirect.mp3/dovetail.prxu.org/stown/ff212be1-e3e6-429f-be98-7e2cc743954d/s-town-ch02.mp3",
				Duration: 48 * time.Minute,
				GUID:     "e02",
			},
			{
				Title:    "Chapter III",
				Desc:     "“Tedious and brief.”",
				Pubdate:  time.Date(2017, time.March, 28, 9, 30, 0, 0, time.UTC),
				URL:      "https://dts.podtrac.com/redirect.mp3/dovetail.prxu.org/stown/e4f5a88b-b383-493a-b7f3-8b8ea52cbf35/s-town-ch03.mp3",
				Duration: 54 * time.Minute,
				GUID:     "e03",
			},
			{
				Title:    "Chapter IV",
				Desc:     "“If anybody could find it, it would be me.”",
				Pubdate:  time.Date(2017, time.March, 28, 9, 15, 0, 0, time.UTC),
				URL:      "https://dts.podtrac.com/redirect.mp3/dovetail.prxu.org/stown/bd934bb6-82d8-4114-ade6-3d9109ebc855/s-town-ch04.mp3",
				Duration: 1*time.Hour + 2*time.Minute,
				GUID:     "e04",
			},
			{
				Title:    "Chapter V",
				Desc:     "“Nobody’ll ever change my mind about it.”",
				Pubdate:  time.Date(2017, time.March, 28, 9, 0, 0, 0, time.UTC),
				URL:      "https://dts.podtrac.com/redirect.mp3/dovetail.prxu.org/stown/1d2076a7-739f-49bd-aa29-5f0501ee2f90/s-town-ch05.mp3",
				Duration: 1*time.Hour + 2*time.Minute,
				GUID:     "e05",
			},
			{
				Title:    "Chapter VI",
				Desc:     "“Since everyone around here thinks I’m a queer anyway.”",
				Pubdate:  time.Date(2017, time.March, 28, 8, 45, 0, 0, time.UTC),
				URL:      "https://dts.podtrac.com/redirect.mp3/dovetail.prxu.org/stown/1c79b084-4451-4f19-95f1-c62ab6578cc0/s-town-ch06.mp3",
				Duration: 47 * time.Minute,
				GUID:     "e06",
			},
			{
				Title:    "Chapter VII",
				Desc:     "“You’re beginning to figure it out now, aren’t you?”",
				Pubdate:  time.Date(2017, time.March, 28, 8, 30, 0, 0, time.UTC),
				URL:      "https://dts.podtrac.com/redirect.mp3/dovetail.prxu.org/stown/7acf74b8-7b0a-4e9e-90be-f69052064b77/s-town-ch07.mp3",
				Duration: 1*time.Hour + 3*time.Minute,
				GUID:     "e07",
			},
		},
	}
}
