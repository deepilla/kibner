package testdata

import (
	"time"

	kibner "github.com/deepilla/kibner/internal/types"
)

func Rabbits() *kibner.Feed {
	return &kibner.Feed{
		Title:  "RABBITS",
		Author: "Public Radio Alliance",
		Desc:   "When Carly Parker’s friend Yumiko goes missing under very mysterious circumstances, Carly’s search for her friend leads her headfirst into a mysterious alternate reality game known only as Rabbits: an ancient dangerous game as old as time itself.",
		Type:   "rss",
		URL:    "http://rabbits.libsyn.com/rss",
		Link:   "http://rabbitspodcast.com",
		Image:  "http://static.libsyn.com/p/assets/9/8/c/3/98c3e4236dfe2b20/RABBITS-ICON-022717.png",
		Items: []*kibner.Item{
			{
				Title:    "Episode 110: The Future We Deserve",
				Desc:     "In the Season Finale of Rabbits, Carly encounters some old friends before leaving the country in search of Yumiko, and, perhaps, the final mystery at the heart of Rabbits.",
				Pubdate:  time.Date(2017, time.July, 4, 18, 16, 32, 0, time.UTC),
				URL:      "http://traffic.libsyn.com/rabbits/RABBITS_EPISODE_110_-_The_Future_We_Deserve.mp3?dest-id=479849",
				Filesize: 82644860,
				Duration: 1*time.Hour + 4*time.Minute + 2*time.Second,
				GUID:     "6f9cdd70b91c3093d03e502a6f91b116",
			},
			{
				Title:    "Episode 109: Hazel",
				Desc:     "In the penultimate episode of Season One, Carly and Jones dig further into The mysterious Gatewick Institute, and another conversation with Alan Scarpio reveals more about the game, and the enigmatic figure known as Hazel.",
				Pubdate:  time.Date(2017, time.June, 20, 23, 29, 43, 0, time.UTC),
				URL:      "http://traffic.libsyn.com/rabbits/RABBITS_EPISODE_109_-_Hazel.mp3?dest-id=479849",
				Filesize: 88544984,
				Duration: 1*time.Hour + 9*time.Minute + 9*time.Second,
				GUID:     "056424f8b33d5ea696c4a64da9f91eb3",
			},
			{
				Title:    "Episode 108: Elysian Drift",
				Desc:     "Carly and Jones visit a reclusive billionaire who appears to know a whole lot more than he’s letting on, and Carly begins to seriously consider the fact that Rabbits might be much more than just a game.",
				Pubdate:  time.Date(2017, time.June, 6, 3, 14, 35, 0, time.UTC),
				URL:      "http://traffic.libsyn.com/rabbits/RABBITS_EPISODE_108_-_Elysian_Drift.mp3?dest-id=479849",
				Filesize: 62413430,
				Duration: 47*time.Minute + 34*time.Second,
				GUID:     "88dbee1acad5e0850ea93bc1dfacb4ce",
			},
			{
				Title:    "Episode 107: Arcadia",
				Desc:     "Carly finds something while playing a game. Harper and Carly meet Batman, Carly and Jones visit Arcadia, and Marigold has another message.",
				Pubdate:  time.Date(2017, time.May, 23, 17, 22, 1, 0, time.UTC),
				URL:      "http://traffic.libsyn.com/rabbits/RABBITS_EPISODE_107_-_Arcadia.mp3?dest-id=479849",
				Filesize: 66620184,
				Duration: 51*time.Minute + 14*time.Second,
				GUID:     "4559ccafc9f65799700cdd7cddec1515",
			},
			{
				Title:    "Episode 106: Strange Attractors",
				Desc:     "In the sixth episode of Rabbits, the Marigold recordings appear to reveal the impossible, and we learn something remarkable about Carly’s past.",
				Pubdate:  time.Date(2017, time.May, 9, 15, 30, 55, 0, time.UTC),
				URL:      "http://traffic.libsyn.com/rabbits/RABBITS_EPISODE_106_-_Strange_Attractors.mp3?dest-id=479849",
				Filesize: 64951419,
				Duration: 50*time.Minute + 2*time.Second,
				GUID:     "b3b9dc6536fd2eb5a21c1dabeb374e23",
			},
			{
				Title:    "Episode 105: Priesthood One",
				Desc:     "In episode five, Carly digs into the mysterious Gatewick Institute, Jones finds something interesting in the Marigold Recordings, and the Magician provides some additional information about the game.",
				Pubdate:  time.Date(2017, time.April, 25, 18, 27, 28, 0, time.UTC),
				URL:      "http://traffic.libsyn.com/rabbits/RABBITS_EPISODE_105_-_Priesthood_One.mp3?dest-id=479849",
				Filesize: 70448523,
				Duration: 54*time.Minute + 45*time.Second,
				GUID:     "2d8ae974ebcecf7d87d830ec6c528a0a",
			},
			{
				Title:    "Episode 104: Doglover in Hell",
				Desc:     "In the fourth episode of Rabbits, Carly and Jones dig further into the history of the game, a new voice points them in a new direction, and the mysterious Hazel makes another appearance.",
				Pubdate:  time.Date(2017, time.April, 11, 15, 32, 16, 0, time.UTC),
				URL:      "http://traffic.libsyn.com/rabbits/RABBITS_EPISODE_104_-_Doglover_in_Hell.mp3?dest-id=479849",
				Filesize: 66406177,
				Duration: 51*time.Minute + 33*time.Second,
				GUID:     "503c9b96980598cdcbd93149a5312f39",
			},
			{
				Title:    "Episode 103: Marigold and Persephone",
				Desc:     "In our third episode, Carly’s search for Yumiko leads to an obscure underground radio station, a vintage t-shirt reveals an impossible photograph, and a work of classic modern art appears to contain a strange secret.",
				Pubdate:  time.Date(2017, time.March, 28, 14, 35, 45, 0, time.UTC),
				URL:      "http://traffic.libsyn.com/rabbits/RABBITS_EPISODE_103_-_Marigold_and_Persephone.mp3?dest-id=479849",
				Filesize: 69336915,
				Duration: 50*time.Minute + 35*time.Second,
				GUID:     "f6106928e5f37faccccfa744cdf7ffa7",
			},
			{
				Title:    "Episode 102: Concernicus Jones",
				Desc:     "In the second episode of Rabbits, Carly continues her search for Yumiko aided by a potential new ally, a stranger who appears to know a great deal about the mysterious game known as “Rabbits.”",
				Pubdate:  time.Date(2017, time.March, 14, 14, 31, 19, 0, time.UTC),
				URL:      "http://traffic.libsyn.com/rabbits/RABBITS_EPISODE_102_-_Concernicus_Jones.mp3?dest-id=479849",
				Filesize: 78294748,
				Duration: 1*time.Hour + 1*time.Minute + 52*time.Second,
				GUID:     "fec3a470ed3812f0eff5f72aa626993a",
			},
			{
				Title:    "Episode 101: Game On",
				Desc:     "In the series premiere of Rabbits, Carly Parker's search for her missing best friend leads her to a mysterious nameless ancient game the players refer to only as “Rabbits;” a secret, dangerous, and occasionally fatal underground game, where the prizes",
				Pubdate:  time.Date(2017, time.February, 28, 8, 30, 0, 0, time.UTC),
				URL:      "http://traffic.libsyn.com/rabbits/RABBITS_EPISODE_101_-_Game_On.mp3?dest-id=479849",
				Filesize: 62261240,
				Duration: 48*time.Minute + 44*time.Second,
				GUID:     "c1f3f359986b6bca595e0f19537dce18",
			},
			{
				Title:    "Episode 000: Introducing Rabbits",
				Desc:     "When Carly Parker’s friend Yumiko goes missing under very mysterious circumstances, Carly’s search for her friend leads her headfirst into a mysterious game known only as Rabbits.",
				Pubdate:  time.Date(2017, time.February, 21, 6, 33, 17, 0, time.UTC),
				URL:      "http://traffic.libsyn.com/rabbits/RABBITS_EPISODE_000.mp3?dest-id=479849",
				Filesize: 5359625,
				Duration: 2*time.Minute + 45*time.Second,
				GUID:     "62240d858321b26b6b18b3673f43c105",
			},
		},
	}
}
