package testdata

import (
	"time"

	kibner "github.com/deepilla/kibner/internal/types"
)

func WhyWeEat2() *kibner.Feed {
	return &kibner.Feed{
		Title:  "Why We Eat What We Eat",
		Author: "Blue Apron / Gimlet Creative",
		Desc:   "A podcast from Blue Apron and Gimlet Creative for anyone who has ever eaten.",
		Type:   "rss",
		URL:    "http://feeds.gimletcreative.com/whyweeatshow",
		Link:   "http://whyweeat.show/",
		Image:  "http://static.megaphone.fm/podcasts/d56755fc-a848-11e7-9635-8bf60d5c6344/image/uploads_2F1507041637303-x7r5xio7tuh-36a66f626df9c478608850a675889158_2FBA_Art_Final.png",
		Items: []*kibner.Item{
			{
				Title:    "The Search for Big Kale",
				Desc:     "How in the world did this bitter, leafy green find its way on to 1 out of every 5 menus in the US?",
				Pubdate:  time.Date(2017, time.October, 11, 19, 32, 21, 0, time.UTC),
				URL:      "https://traffic.megaphone.fm/GLT5421811798.mp3?updated=1507748492",
				Filesize: 32409600,
				Duration: 1350 * time.Second,
				GUID:     "f366df52-ae09-11e7-bcda-835224aa380b",
			},
			{
				Title:    "Introducing: Why We Eat What We Eat",
				Desc:     "Welcome to Why Eat What We Eat, a podcast about the not-obvious-answers to our strange eating habits.",
				Pubdate:  time.Date(2017, time.October, 3, 16, 2, 0, 0, time.UTC),
				URL:      "https://traffic.megaphone.fm/GLT5623364696.mp3",
				Filesize: 3636244,
				Duration: 151 * time.Second,
				GUID:     "8b4e2dd4-a84c-11e7-8147-2b03377712c4",
			},
		},
	}
}
