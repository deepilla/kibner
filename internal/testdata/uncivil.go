package testdata

import (
	"time"

	kibner "github.com/deepilla/kibner/internal/types"
)

func Uncivil() *kibner.Feed {
	return &kibner.Feed{
		Title:  "Uncivil",
		Author: "Gimlet",
		Desc:   "A new history podcast from Gimlet Media, where we go back to the time our divisions led to war. We ransack America's past, and bring you stories left out of the official history. Hosted by Jack Hitt and Chenjerai Kumanyika.",
		Type:   "rss",
		URL:    "https://feeds.megaphone.fm/uncivil",
		Link:   "http://uncivil.show",
		Image:  "http://static.megaphone.fm/podcasts/60cff344-5135-11e7-a8cc-e3b7831c52a1/image/uploads_2F1506371489244-7tntk0du7bx-60a7278c8262d8209c89a0de0fca6fef_2Funcivil_show_art_season1.png",
		Items: []*kibner.Item{
			{
				Title:    "The Raid",
				Desc:     "A group of ex-farmers, a terrorist from Kansas, and a schoolteacher attempt the greatest covert operation of the Civil War.",
				Pubdate:  time.Date(2017, time.October, 4, 10, 0, 0, 0, time.UTC),
				URL:      "https://traffic.megaphone.fm/GLT6754684783.mp3",
				Filesize: 35388186,
				Duration: 1474 * time.Second,
				GUID:     "6d592f74-9409-11e7-b285-f76a850c5be4",
			},
			{
				Title:    "Give Us a Call",
				Desc:     "We want to hear from you! Do you have a family story from the Civil War? A myth you've seen on social media you want busted? Let us know! Call 347-395-5078 and leave us a voicemail.",
				Pubdate:  time.Date(2017, time.September, 20, 15, 46, 39, 0, time.UTC),
				URL:      "https://traffic.megaphone.fm/GLT8145162015.mp3",
				Filesize: 2251337,
				Duration: 93 * time.Second,
				GUID:     "9ec7ea3a-9e1a-11e7-9005-7b27fae7db0f",
			},
			{
				Title:    "Coming Soon",
				Desc:     "A new history podcast from Gimlet Media, where we go back to the time our divisions turned into a war, and bring you stories left out of the official history.",
				Pubdate:  time.Date(2017, time.August, 24, 10, 44, 0, 0, time.UTC),
				URL:      "https://traffic.megaphone.fm/GLT7465151929.mp3",
				Filesize: 3179206,
				Duration: 132 * time.Second,
				GUID:     "95ddd8b2-88da-11e7-bf3d-3ff6bef2d7c4",
			},
		},
	}
}
