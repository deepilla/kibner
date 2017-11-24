package testdata

import (
	"time"

	kibner "github.com/deepilla/kibner/internal/types"
)

func EarHustle() *kibner.Feed {
	return &kibner.Feed{
		Title:  "Ear Hustle",
		Author: "Ear Hustle & Radiotopia",
		Desc:   "Ear Hustle brings you stories of life inside prison, shared and produced by those living it.",
		Type:   "rss",
		URL:    "http://feeds.earhustlesq.com/earhustlesq",
		Link:   "https://www.earhustlesq.com/",
		Image:  "https://f.prxu.org/59/images/d452ebe9-7814-4f3c-bd2e-be28c2468002/EH_Logo_zag_1500x1500.jpg",
		Items: []*kibner.Item{
			{
				Title:    "Left Behind",
				Desc:     "How do inmates with profoundly long sentences cope with their realities, and maintain a sense of hope and well-being as the years pass?",
				Pubdate:  time.Date(2017, time.September, 27, 13, 5, 17, 0, time.UTC),
				URL:      "https://dts.podtrac.com/redirect.mp3/dovetail.prxu.org/59/d8ec9498-302d-41be-818f-1ece35512d6f/Left_Behind_final_A.mp3",
				Filesize: 33904569,
				Duration: 35*time.Minute + 19*time.Second,
				GUID:     "prx_59_d8ec9498-302d-41be-818f-1ece35512d6f",
			},
			{
				Title:    "Unwritten",
				Desc:     "The color of your skin influences your life on the inside, from sharing food to celebrating birthdays.",
				Pubdate:  time.Date(2017, time.September, 13, 13, 6, 23, 0, time.UTC),
				URL:      "https://dts.podtrac.com/redirect.mp3/dovetail.prxu.org/59/f2e78445-e2aa-496d-b20e-5bc8bf364f67/Unwritten_A.mp3",
				Filesize: 28811331,
				Duration: 30 * time.Minute,
				GUID:     "prx_59_f2e78445-e2aa-496d-b20e-5bc8bf364f67",
			},
			{
				Title:    "The Boom Boom Room",
				Desc:     "Being married in prison is common. Opportunities to get intimate with your spouse are not.",
				Pubdate:  time.Date(2017, time.August, 30, 12, 39, 49, 0, time.UTC),
				URL:      "https://dts.podtrac.com/redirect.mp3/dovetail.prxu.org/59/8dd8a4e9-4e1a-42c9-99d5-f65c84fcc9ea/Boom_Boom_Room_billboard.mp3",
				Filesize: 27192954,
				Duration: 28*time.Minute + 19*time.Second,
				GUID:     "prx_59_8dd8a4e9-4e1a-42c9-99d5-f65c84fcc9ea",
			},
			{
				Title:    "Catch a Kite",
				Desc:     `Earlonne and Nigel answer questions mailed in via postcard, aka "kites," from listeners around the world.`,
				Pubdate:  time.Date(2017, time.August, 9, 12, 41, 37, 0, time.UTC),
				URL:      "https://dts.podtrac.com/redirect.mp3/dovetail.prxu.org/59/8bf421ae-eac5-4964-ac82-230c91caee19/Q_A_mix4_billboard.mp3",
				Filesize: 26942337,
				Duration: 28*time.Minute + 3*time.Second,
				GUID:     "prx_59_8bf421ae-eac5-4964-ac82-230c91caee19",
			},
			{
				Title:    "The SHU",
				Desc:     "The hole, the box, solitary confinement… stories from four men who spent years completely isolated in prison.",
				Pubdate:  time.Date(2017, time.July, 26, 12, 43, 43, 0, time.UTC),
				URL:      "https://dts.podtrac.com/redirect.mp3/dovetail.prxu.org/59/1e0592e7-4a24-4627-a73f-a4e95ca8007b/EH_ep4_billboard.mp3",
				Filesize: 27298872,
				Duration: 28*time.Minute + 26*time.Second,
				GUID:     "prx_59_1e0592e7-4a24-4627-a73f-a4e95ca8007b",
			},
			{
				Title:    "Looking Out",
				Desc:     "Meet Rauch, a man incarcerated at San Quentin who has figured out how to “look out” for dozens of critters around the place.",
				Pubdate:  time.Date(2017, time.July, 12, 12, 52, 25, 0, time.UTC),
				URL:      "https://dts.podtrac.com/redirect.mp3/dovetail.prxu.org/59/81291ee0-a9f4-4508-a71f-3a0f1d07bb23/EH_EP3_billboard.mp3",
				Filesize: 22841976,
				Duration: 23*time.Minute + 47*time.Second,
				GUID:     "prx_59_81291ee0-a9f4-4508-a71f-3a0f1d07bb23",
			},
			{
				Title:    "Misguided Loyalty",
				Desc:     "While the perceived thrill and glamour of gang life are often too strong to resist, the hard consequences of this seduction can last a lifetime.",
				Pubdate:  time.Date(2017, time.June, 28, 11, 36, 6, 0, time.UTC),
				URL:      "https://dts.podtrac.com/redirect.mp3/dovetail.prxu.org/59/3cf87185-1149-4963-b065-ba6d84cf2c39/EH_EP2_billboard_2.mp3",
				Filesize: 26500317,
				Duration: 27*time.Minute + 36*time.Second,
				GUID:     "prx_59_3cf87185-1149-4963-b065-ba6d84cf2c39",
			},
			{
				Title:    "Cellies",
				Desc:     "Finding a roommate can be tough. Finding someone to share a 4' x 9' space with is a whole 'nother story.",
				Pubdate:  time.Date(2017, time.June, 14, 13, 56, 33, 0, time.UTC),
				URL:      "https://dts.podtrac.com/redirect.mp3/dovetail.prxu.org/59/fac58065-7b8d-4772-ba0c-b07bbec9564e/Cellies_ep1_billboard.mp3",
				Filesize: 22616796,
				Duration: 23*time.Minute + 33*time.Second,
				GUID:     "prx_59_fac58065-7b8d-4772-ba0c-b07bbec9564e",
			},
			{
				Title:    "Meet Ear Hustle",
				Desc:     "Ear Hustle brings you the stories of life inside prison, shared and produced by those living it.",
				Pubdate:  time.Date(2017, time.May, 25, 16, 15, 57, 0, time.UTC),
				URL:      "https://dts.podtrac.com/redirect.mp3/dovetail.prxu.org/59/439f17a3-7907-425f-9fa6-54f3d9aa1132/Meet_Ear_Hustle.mp3",
				Filesize: 3867057,
				Duration: 4*time.Minute + 1*time.Second,
				GUID:     "prx_59_439f17a3-7907-425f-9fa6-54f3d9aa1132",
			},
		},
	}
}
