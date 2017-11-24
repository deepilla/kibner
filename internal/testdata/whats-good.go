package testdata

import (
	"time"

	kibner "github.com/deepilla/kibner/internal/types"
)

func WhatsGood() *kibner.Feed {
	return &kibner.Feed{
		Title:  "What's Good with Stretch & Bobbito",
		Author: "NPR",
		Desc:   "What's Good with Stretch & Bobbito is your source for untold stories and uncovered truths. Hosts Stretch Armstrong and Bobbito Garcia interview cultural influencers, bringing their warmth, humor, and a fresh perspective. They're talking about art, music, politics, sports and what's good!",
		Type:   "rss",
		URL:    "https://www.npr.org/rss/podcast.php?id=510323",
		Image:  "https://media.npr.org/assets/img/2017/07/11/whatsgood_podcasttile_final_sq-91d7f2af7f80547f8ae2f34f832fa573861233a2.png?s=1400",
		Items: []*kibner.Item{
			{
				Title:    "Hill Harper",
				Desc:     `Hill Harper is an author, activist and actor, who recently became a regular on ABC's drama "The Good Doctor". The hosts and Harper talk about justice in America, how playing a cop on TV informs his views on policing, and about playing ball with former president Barack Obama.`,
				Pubdate:  time.Date(2017, time.October, 4, 4, 1, 16, 0, time.UTC),
				URL:      "https://play.podtrac.com/npr-510323/npr.mc.tritondigital.com/NPR_510323/media/anon.npr-mp3/npr/stretchbobbito/2017/10/20171004_stretchbobbito_stretchbobharperreworkfinal.mp3?orgId=1&d=1559&p=510323&story=553860280&t=podcast&e=553860280&ft=pod&f=510323",
				Duration: 1559 * time.Second,
				GUID:     "e6936b6e-90bc-4bb9-8275-49dc88d40faf",
			},
			{
				Title:    "Franchesca Ramsey",
				Desc:     "Stretch and Bobbito talk to Franchesca Ramsey about the early days of YouTube, her new Comedy Central pilot, and what she's learned while hosting her own podcast.",
				Pubdate:  time.Date(2017, time.September, 27, 4, 1, 29, 0, time.UTC),
				URL:      "https://play.podtrac.com/npr-510323/npr.mc.tritondigital.com/NPR_510323/media/anon.npr-mp3/npr/stretchbobbito/2017/09/20170927_stretchbobbito_stretchbobramseyfinal.mp3?orgId=1&d=1760&p=510323&story=553784105&t=podcast&e=553784105&ft=pod&f=510323",
				Duration: 1760 * time.Second,
				GUID:     "42dbfffd-d9d0-40ef-9eb2-85d46e8e86d0",
			},
			{
				Title:    "Jose Parla",
				Desc:     "Artist Jose Parla joins the guys to talk about Cuba, taking graffiti seriously as an artform, and the joy of digging up old music.",
				Pubdate:  time.Date(2017, time.September, 20, 4, 1, 11, 0, time.UTC),
				URL:      "https://play.podtrac.com/npr-510323/npr.mc.tritondigital.com/NPR_510323/media/anon.npr-mp3/npr/stretchbobbito/2017/09/20170920_stretchbobbito_stretchbobjoseparlafinal.mp3?orgId=1&d=1854&p=510323&story=552164427&t=podcast&e=552164427&ft=pod&f=510323",
				Duration: 1854 * time.Second,
				GUID:     "5dbca284-a5e0-4487-9bb9-1f29618a3a44",
			},
			{
				Title:    "Ana Navarro",
				Desc:     "The CNN commentator joins the guys to talk about life as a political refugee, winning shouting matches on cable news, and her patriotic love of peanut butter.",
				Pubdate:  time.Date(2017, time.September, 13, 4, 1, 31, 0, time.UTC),
				URL:      "https://play.podtrac.com/npr-510323/npr.mc.tritondigital.com/NPR_510323/media/anon.npr-mp3/npr/stretchbobbito/2017/09/20170913_stretchbobbito_stretchbobananewfinal.mp3?orgId=1&d=1656&p=510323&story=550542307&t=podcast&e=550542307&ft=pod&f=510323",
				Duration: 1656 * time.Second,
				GUID:     "4eb5a380-c9e2-458e-90d2-d905cc2c0f4c",
			},
			{
				Title:    "Run The Jewels",
				Desc:     "El-P and Killer Mike talk to Stretch and Bobbito about how they turned slang into a slogan, homophobia in hip-hop, and the importance of cherishing shared humanity over political differences.",
				Pubdate:  time.Date(2017, time.September, 6, 4, 1, 24, 0, time.UTC),
				URL:      "https://play.podtrac.com/npr-510323/npr.mc.tritondigital.com/NPR_510323/media/anon.npr-mp3/npr/stretchbobbito/2017/09/20170906_stretchbobbito_stretchbobrtjnewfinal.mp3?orgId=1&d=1755&p=510323&story=548751251&t=podcast&e=548751251&ft=pod&f=510323",
				Duration: 1755 * time.Second,
				GUID:     "7a7daa7a-0e32-4759-a558-acc294e89ccd",
			},
			{
				Title:    "Stevie Wonder",
				Desc:     "The one and only Stevie Wonder talks to Stretch, Bobbito and DJ Spinna about making music at Motown Records, pushing for the creation of Martin Luther King Jr. Day and missing his friend Prince.",
				Pubdate:  time.Date(2017, time.August, 30, 4, 1, 19, 0, time.UTC),
				URL:      "https://play.podtrac.com/npr-510323/npr.mc.tritondigital.com/NPR_510323/media/anon.npr-mp3/npr/stretchbobbito/2017/08/20170830_stretchbobbito_stretchbobstevienewfinal.mp3?orgId=1&d=2868&p=510323&story=547150392&t=podcast&e=547150392&ft=pod&f=510323",
				Duration: 2868 * time.Second,
				GUID:     "6de18329-92d5-4e2e-8ec6-949f9b7e9be2",
			},
			{
				Title:    "Regina King",
				Desc:     "Stretch & Bobbito sit down with the actor and director to talk about making the move from movies to TV, how growing up in LA informs her work, and partying in NYC.",
				Pubdate:  time.Date(2017, time.August, 23, 4, 1, 16, 0, time.UTC),
				URL:      "https://play.podtrac.com/npr-510323/npr.mc.tritondigital.com/NPR_510323/media/anon.npr-mp3/npr/codeswitch/2017/08/20170823_codeswitch_stretchbobreginakingfinal.mp3?orgId=1&d=1962&p=510323&story=545397519&t=podcast&e=545397519&ft=pod&f=510323",
				Duration: 1962 * time.Second,
				GUID:     "fbcd18f6-65f7-45fb-b15e-1555be9f45d3",
			},
			{
				Title:    "Linda Sarsour",
				Desc:     "The guys sit down with activist Linda Sarsour to talk protests, growing up in Brooklyn, and the radical act of speaking to your neighbors.",
				Pubdate:  time.Date(2017, time.August, 16, 4, 1, 0, 0, time.UTC),
				URL:      "https://play.podtrac.com/npr-510323/npr.mc.tritondigital.com/NPR_510323/media/anon.npr-mp3/npr/stretchbobbito/2017/08/20170816_stretchbobbito_stretchboblindafinal.mp3?orgId=1&d=1792&p=510323&story=533842536&t=podcast&e=533842536&ft=pod&f=510323",
				Duration: 1792 * time.Second,
				GUID:     "0649b8f6-67ed-4da9-80ed-455e08831590",
			},
			{
				Title:    "Chance The Rapper",
				Desc:     "Chance the Rapper joins the guys to talk about redefining mixtapes, using his influence for political good, and how having a daughter has changed his music.",
				Pubdate:  time.Date(2017, time.August, 9, 4, 1, 0, 0, time.UTC),
				URL:      "https://play.podtrac.com/npr-510323/npr.mc.tritondigital.com/NPR_510323/media/anon.npr-mp3/npr/stretchbobbito/2017/08/20170808_stretchbobbito_stretchbobchancefinal.mp3?orgId=1&d=2095&p=510323&story=540673885&t=podcast&e=540673885&ft=pod&f=510323",
				Duration: 2095 * time.Second,
				GUID:     "8e9f9d40-b6e8-4082-8565-b36af3df8d27",
			},
			{
				Title:    "Eddie Huang",
				Desc:     "The best-selling author of Fresh Off The Boat and host of Huang's World sits down with Stretch and Bobbito to talk sneakers, hip-hop and the American dream.",
				Pubdate:  time.Date(2017, time.August, 2, 4, 1, 0, 0, time.UTC),
				URL:      "https://play.podtrac.com/npr-510323/npr.mc.tritondigital.com/NPR_510323/media/anon.npr-mp3/npr/stretchbobbito/2017/08/20170802_stretchbobbito_stretchandbobeddiehuangfinal.mp3?orgId=1&d=1889&p=510323&story=540673272&t=podcast&e=540673272&ft=pod&f=510323",
				Duration: 1889 * time.Second,
				GUID:     "5db34407-3e00-40d9-84e5-cbae21e70bba",
			},
			{
				Title:    "Mahershala Ali",
				Desc:     "The Oscar winning actor joins Stretch and Bobbito to talk about his career as an MC, his process for getting into character, and how his faith informs his work.",
				Pubdate:  time.Date(2017, time.July, 26, 4, 1, 0, 0, time.UTC),
				URL:      "https://play.podtrac.com/npr-510323/npr.mc.tritondigital.com/NPR_510323/media/anon.npr-mp3/npr/stretchbobbito/2017/07/20170725_stretchbobbito_stretchbobbitomahershalaalifinal.mp3?orgId=1&d=1962&p=510323&story=539317543&t=podcast&e=539317543&ft=pod&f=510323",
				Duration: 1962 * time.Second,
				GUID:     "bef498cc-5d44-47f5-b8c5-68aaff49dc70",
			},
			{
				Title:    "Dave Chappelle & Donnell Rawlings",
				Desc:     "Comedians Dave Chappelle and Donnell Rawlings preview their month long residency at Radio City Music Hall. They talk to Stretch and Bobbito about their early days in stand-up, hosting Saturday Night Live after the election, and the changing face of their hometown Washington D.C.",
				Pubdate:  time.Date(2017, time.July, 19, 4, 1, 0, 0, time.UTC),
				URL:      "https://play.podtrac.com/npr-510323/npr.mc.tritondigital.com/NPR_510323/media/anon.npr-mp3/npr/stretchbobbito/2017/07/20170719_stretchbobbito_stretchbobbitodavechappellefinalmix.mp3?orgId=1&d=2034&p=510323&story=537942533&t=podcast&e=537942533&ft=pod&f=510323",
				Duration: 2034 * time.Second,
				GUID:     "971be916-cb55-4403-a6dc-d01e7a5a74b2",
			},
			{
				Title:    "It's What's Good with Stretch and Bobbito",
				Desc:     "The legendary duo Stretch and Bobbito are back! The first episode launches on Wednesday, July 19th, featuring an interview with Dave Chappelle.",
				Pubdate:  time.Date(2017, time.July, 11, 18, 27, 0, 0, time.UTC),
				URL:      "https://play.podtrac.com/npr-510323/npr.mc.tritondigital.com/NPR_510323/media/anon.npr-mp3/npr/stretchbobbito/2017/07/20170711_stretchbobbito_stretch_and_bob_database__-_apple_trailer_7_11_150_pm_cut.mp3?orgId=1&d=120&p=510323&story=536655610&t=podcast&e=536655610&ft=pod&f=510323",
				Duration: 120 * time.Second,
				GUID:     "012e4ea7-f042-4665-974f-17e8708b437c",
			},
		},
	}
}
