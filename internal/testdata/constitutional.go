package testdata

import (
	"time"

	kibner "github.com/deepilla/kibner/internal/types"
)

func Constitutional() *kibner.Feed {
	return &kibner.Feed{
		Title:  "Constitutional",
		Author: "The Washington Post",
		Desc:   "With the writing of the Constitution in 1787, the framers set out a young nation’s highest ideals. And ever since, we’ve been fighting over it — what is in it and what was left out. At the heart of these arguments is the story of America. As a follow-up to the popular Washington Post podcast “Presidential,” reporter Lillian Cunningham returns with this series exploring the Constitution and the people who framed and reframed it — revolutionaries, abolitionists, suffragists, teetotalers, protesters, justices, presidents – in the ongoing struggle to form a more perfect union across a vast and diverse land.",
		Type:   "rss",
		URL:    "https://podcast.posttv.com/itunes-5951e18ce4b056ae8fa8e4b2.xml",
		Link:   "https://www.washingtonpost.com",
		Image:  "https://podcast.posttv.com/series/20170724/t_1500924687921_name_FINAL_WP_actual_square.jpg",
		Items: []*kibner.Item{
			{
				Title:    "Episode 07: Congress and citizens",
				Desc:     "Is it a feature or a bug of the amendment process that an idea of James Madison's, more than 200 years ago, could be recently resurrected and etched into the U.S. Constitution?",
				Pubdate:  time.Date(2017, time.September, 25, 7, 0, 0, 0, time.UTC),
				URL:      "http://www.podtrac.com/pts/redirect.mp3/podcast.posttv.com/washpost-production/5951e18ce4b056ae8fa8e4b2/20170922/59c57779e4b08d2bc76a9cd2/59c5777ae4b0dc47945913c3_1351620000001-300030_t_1506113406759_44100_160_2.mp3",
				Filesize: 48390383,
				Duration: 2420 * time.Second,
				GUID:     "59c57779e4b08d2bc76a9cd2",
			},
			{
				Title:    "Episode 06: Senate and states",
				Desc:     `When the United States changed its process for electing senators, did that lead to a decline in state power? Or did it instead bring us closer to a "more perfect union"?`,
				Pubdate:  time.Date(2017, time.September, 11, 7, 0, 0, 0, time.UTC),
				URL:      "http://www.podtrac.com/pts/redirect.mp3/podcast.posttv.com/washpost-production/5951e18ce4b056ae8fa8e4b2/20170908/59b308d4e4b08d2bc76a9b0b/59b308d4e4b0dc4794590ffb_1351620000001-300030_t_1504905432307_44100_160_2.mp3",
				Filesize: 56137775,
				Duration: 2807 * time.Second,
				GUID:     "59b308d4e4b08d2bc76a9b0b",
			},
			{
				Title:    "Episode 05: Gender",
				Desc:     "From the American Revolution through today, women have been leading a long-burning rebellion to gain rights not originally guaranteed under the Constitution.",
				Pubdate:  time.Date(2017, time.August, 28, 7, 0, 0, 0, time.UTC),
				URL:      "http://www.podtrac.com/pts/redirect.mp3/podcast.posttv.com/washpost-production/5951e18ce4b056ae8fa8e4b2/20170825/59a08392e4b07a5a3360bd82/59a08392e4b0b07aa66102e1_1351620000001-300040_t_1503691680514_44100_128_2.mp3",
				Filesize: 48036996,
				Duration: 3003 * time.Second,
				GUID:     "59a08392e4b07a5a3360bd82",
			},
			{
				Title:    "Episode 04: Race",
				Desc:     "As powerful as it was to change the Constitution after the Civil War, and enshrine racial equality into our governing document, that wasn’t enough to change the reality of life in America.",
				Pubdate:  time.Date(2017, time.August, 21, 7, 0, 0, 0, time.UTC),
				URL:      "http://www.podtrac.com/pts/redirect.mp3/podcast.posttv.com/washpost-production/5951e18ce4b056ae8fa8e4b2/20170818/5997666ee4b09529fad4f97f/5997666fe4b0b07aa660ed17_1351620000001-300030_t_1503094392481_44100_160_2.mp3",
				Filesize: 63593128,
				Duration: 3180 * time.Second,
				GUID:     "5997666ee4b09529fad4f97f",
			},
			{
				Title:    "Episode 03: Nationality",
				Desc:     "What makes someone American? A landmark Supreme Court case in 1898, involving a child born in San Francisco to Chinese immigrant parents, would help answer that question.",
				Pubdate:  time.Date(2017, time.August, 14, 7, 0, 0, 0, time.UTC),
				URL:      "http://www.podtrac.com/pts/redirect.mp3/podcast.posttv.com/washpost-production/5951e18ce4b056ae8fa8e4b2/20170811/598e092ee4b02ea76027a81d/598e092fe4b0b07aa660d825_1351620000001-300040_t_1502480697030_44100_128_2.mp3",
				Filesize: 46807364,
				Duration: 2926 * time.Second,
				GUID:     "598e092ee4b02ea76027a81d",
			},
			{
				Title:    "Episode 02: Ancestry",
				Desc:     `This "we the people" episode explores indigenous rights. In 1879, a case involving Chief Standing Bear came before a Nebraska courtroom and demanded an answer to the question: Are Native Americans considered human beings under the U.S. Constitution?`,
				Pubdate:  time.Date(2017, time.August, 7, 7, 0, 0, 0, time.UTC),
				URL:      "http://www.podtrac.com/pts/redirect.mp3/podcast.posttv.com/washpost-production/5951e18ce4b056ae8fa8e4b2/20170804/5984cdf2e4b08ca13090aca6/5984cdf3e4b0b07aa660c44a_1351620000001-300040_t_1501875708969_44100_128_2.mp3",
				Filesize: 40430518,
				Duration: 2527 * time.Second,
				GUID:     "5984cdf2e4b08ca13090aca6",
			},
			{
				Title:    "Episode 01: Framed",
				Desc:     "In the premier episode of “Constitutional,” we go back in time to that hot Philadelphia summer in 1787 when a group of revolutionary Americans debated, drank and together drafted the U.S. Constitution.",
				Pubdate:  time.Date(2017, time.July, 24, 7, 0, 0, 0, time.UTC),
				URL:      "http://www.podtrac.com/pts/redirect.mp3/podcast.posttv.com/washpost-production/5951e18ce4b056ae8fa8e4b2/20170721/59725fbae4b05ba1300f5b04/59725fbde4b0b07aa6609a10_1351620000001-300040_t_1500667850244_44100_128_2.mp3",
				Filesize: 62025215,
				Duration: 3877 * time.Second,
				GUID:     "59725fbae4b05ba1300f5b04",
			},
			{
				Title:    "Introducing 'Constitutional'",
				Desc:     "Preview The Washington Post's newest podcast, a narrative series about the revolutionary figures who shaped America's story. Subscribe now to get the first episode when it launches July 24.",
				Pubdate:  time.Date(2017, time.June, 29, 14, 0, 0, 0, time.UTC),
				URL:      "http://www.podtrac.com/pts/redirect.mp3/podcast.posttv.com/washpost-production/5951e18ce4b056ae8fa8e4b2/20170629/59547941e4b056ae8fa8e4b7/59547941e4b0b07aa6606517_1351620000001-300030_t_1498708302801_44100_160_2.mp3",
				Filesize: 5788808,
				Duration: 290 * time.Second,
				GUID:     "59547941e4b056ae8fa8e4b7",
			},
		},
	}
}
