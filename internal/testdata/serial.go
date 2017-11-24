package testdata

import (
	"time"

	kibner "github.com/deepilla/kibner/internal/types"
)

func Serial() *kibner.Feed {
	return &kibner.Feed{
		Title:  "Serial",
		Author: "This American Life",
		Desc:   "A podcast from the creators of This American Life",
		Type:   "rss",
		URL:    "http://feeds.serialpodcast.org/serialpodcast",
		Link:   "https://serialpodcast.org",
		Image:  "https://serialpodcast.org/sites/all/modules/custom/serial/img/serial-itunes-logo.png",
		Items: []*kibner.Item{
			{
				Title:    "S01 Episode 12: What We Know",
				Desc:     "On January 13, 1999, Adnan Syed was a hurt and vengeful ex-boyfriend who carried out a premeditated murder. Or he was a bewildered bystander, framed for a crime he could never have committed. After 15 months of reporting, we take out everything we’ve...",
				Pubdate:  time.Date(2014, time.December, 18, 10, 30, 0, 0, time.UTC),
				URL:      "https://dts.podtrac.com/redirect.mp3/dovetail.prxu.org/serial/ea810ac6-a771-492f-b56a-e1b988903631/serial-s01-e12.mp3",
				Filesize: 0,
				Duration: 56 * time.Minute,
				GUID:     "57 at http://serialpodcast.org",
			},
			{
				Title:    "S01 Episode 11: Rumors",
				Desc:     "Almost everyone describes the 17-year-old Adnan the same way: good kid, helpful at the mosque, respectful to his elders. But a couple of months ago, Sarah started getting phone calls from people who knew Adnan back then, and told her stories of a...",
				Pubdate:  time.Date(2014, time.December, 11, 10, 30, 0, 0, time.UTC),
				URL:      "https://dts.podtrac.com/redirect.mp3/dovetail.prxu.org/serial/08ee345b-1b30-4bab-9160-0506fea3b123/serial-s01-e11.mp3",
				Filesize: 0,
				Duration: 41 * time.Minute,
				GUID:     "56 at http://serialpodcast.org",
			},
			{
				Title:    "S01 Episode 10: The Best Defense is a Good Defense",
				Desc:     "Adnan’s trial lawyer was M. Cristina Gutierrez, a renowned defense attorney in Maryland – tough and savvy and smart. Other lawyers said she was exactly the kind of person you’d want defending you on a first-degree murder charge. But Adnan was convicted...",
				Pubdate:  time.Date(2014, time.December, 4, 10, 30, 0, 0, time.UTC),
				URL:      "https://dts.podtrac.com/redirect.mp3/dovetail.prxu.org/serial/f4626e87-319f-4af7-9fd8-f1b92afc067a/serial-s01-e10.mp3",
				Filesize: 0,
				Duration: 54 * time.Minute,
				GUID:     "54 at http://serialpodcast.org",
			},
			{
				Title:    "S01 Episode 09: To Be Suspected",
				Desc:     "New information is coming in about what maybe didn’t happen on January 13, 1999. And while Adnan’s memory of that day is foggy at best, he does remember what happened next: being questioned, being arrested and, a little more than a year later, being...",
				Pubdate:  time.Date(2014, time.November, 20, 10, 30, 0, 0, time.UTC),
				URL:      "https://dts.podtrac.com/redirect.mp3/dovetail.prxu.org/serial/46d17784-a8dc-4704-a30c-a5a34c91fed6/serial-s01-e09.mp3",
				Filesize: 0,
				Duration: 45 * time.Minute,
				GUID:     "49 at http://serialpodcast.org",
			},
			{
				Title:    "S01 Episode 08: The Deal with Jay",
				Desc:     "The state’s case against Adnan Syed hinged on Jay’s credibility; he was their star witness and also, because of his changing statements to police, their chief liability. Naturally, Adnan’s lawyer tried hard to make Jay look untrustworthy at trial. So,...",
				Pubdate:  time.Date(2014, time.November, 13, 10, 30, 0, 0, time.UTC),
				URL:      "https://dts.podtrac.com/redirect.mp3/dovetail.prxu.org/serial/c51c3825-2937-4734-a093-685ae7b35204/serial-s01-e08.mp3",
				Filesize: 0,
				Duration: 44 * time.Minute,
				GUID:     "45 at http://serialpodcast.org",
			},
			{
				Title:    "S01 Episode 07: The Opposite of the Prosecution",
				Desc:     "Adnan told Sarah about a case in Virginia that had striking similarities to his own: one key witness, incriminating cell phone records, young people, drugs - and a defendant who has always maintained his innocence. Sarah called up one of the defense...",
				Pubdate:  time.Date(2014, time.November, 6, 10, 30, 0, 0, time.UTC),
				URL:      "https://dts.podtrac.com/redirect.mp3/dovetail.prxu.org/serial/109d1f22-c8bc-468d-bea3-49b26bd18e81/serial-s01-e07.mp3",
				Filesize: 0,
				Duration: 33 * time.Minute,
				GUID:     "42 at http://serialpodcast.org",
			},
			{
				Title:    "S01 Episode 06: The Case Against Adnan Syed",
				Desc:     "The physical evidence against Adnan Syed was scant - a few underwhelming fingerprints. So aside from cell records, what did the prosecutors bring to the jury, to shore up Jay's testimony? Sarah weighs all the other circumstantial evidence they had...",
				Pubdate:  time.Date(2014, time.October, 30, 9, 30, 0, 0, time.UTC),
				URL:      "https://dts.podtrac.com/redirect.mp3/dovetail.prxu.org/serial/96a58cd4-22f9-4ec7-a24e-9f9e913d3b81/serial-s01-e06.mp3",
				Filesize: 0,
				Duration: 44 * time.Minute,
				GUID:     "34 at http://serialpodcast.org",
			},
			{
				Title:    "S01 Episode 05: Route Talk",
				Desc:     "Adnan once issued a challenge to Sarah. He told her to test the state’s timeline of the murder by driving from Woodlawn High School to Best Buy in 21 minutes. It can’t be done, he said. So Sarah and Dana take up the challenge, and raise him one: They...",
				Pubdate:  time.Date(2014, time.October, 23, 9, 30, 0, 0, time.UTC),
				URL:      "https://dts.podtrac.com/redirect.mp3/dovetail.prxu.org/serial/44446ba8-a848-44e3-a632-8d100a47249c/serial-s01-e05.mp3",
				Filesize: 0,
				Duration: 43 * time.Minute,
				GUID:     "31 at http://serialpodcast.org",
			},
			{
				Title:    "S01 Episode 04: Inconsistencies",
				Desc:     "A few days after Hae’s body is found, the detectives get a lead that opens the case up for them. They find Jay at work late one night and bring him down to Homicide. At first, he insists he doesn’t know anything about the murder. But eventually he...",
				Pubdate:  time.Date(2014, time.October, 16, 9, 30, 0, 0, time.UTC),
				URL:      "https://dts.podtrac.com/redirect.mp3/dovetail.prxu.org/serial/b0195641-b3fa-4019-9fc6-27692c3ec5d9/serial-s01-e04.mp3",
				Filesize: 0,
				Duration: 34 * time.Minute,
				GUID:     "28 at http://serialpodcast.org",
			},
			{
				Title:    "S01 Episode 03: Leakin Park",
				Desc:     "It’s February 9, 1999. Hae has been missing for three weeks. A man on his lunch break pulls off a road to pee, and stumbles on her body in a city forest. His odd recounting of the discovery makes Detectives Ritz and MacGillivary suspicious. For...",
				Pubdate:  time.Date(2014, time.October, 9, 10, 0, 0, 0, time.UTC),
				URL:      "https://dts.podtrac.com/redirect.mp3/dovetail.prxu.org/serial/034a15be-c970-4927-8a55-b6784045b030/serial-s01-e03.mp3",
				Filesize: 0,
				Duration: 28 * time.Minute,
				GUID:     "21 at http://serialpodcast.org",
			},
			{
				Title:    "S01 Episode 02: The Breakup",
				Desc:     "Their relationship began like a storybook high-school romance: a prom date, love notes, sneaking off to be alone. But unlike other kids at school, they had to keep their dating secret, because their parents disapproved. Both of them, but especially...",
				Pubdate:  time.Date(2014, time.October, 3, 14, 0, 0, 0, time.UTC),
				URL:      "https://dts.podtrac.com/redirect.mp3/dovetail.prxu.org/serial/75ea7f5c-e80f-46c5-9d5e-759b398d1d32/serial-s01-e02.mp3",
				Filesize: 0,
				Duration: 37 * time.Minute,
				GUID:     "17 at http://serialpodcast.org",
			},
			{
				Title:    "S01 Episode 01: The Alibi",
				Desc:     "It's Baltimore, 1999. Hae Min Lee, a popular high-school senior, disappears after school one day. Six weeks later detectives arrest her classmate and ex-boyfriend, Adnan Syed, for her murder. He says he's innocent - though he can't exactly remember...",
				Pubdate:  time.Date(2014, time.October, 3, 13, 45, 0, 0, time.UTC),
				URL:      "https://dts.podtrac.com/redirect.mp3/dovetail.prxu.org/serial/d7f03a15-be26-4634-8884-5fadd404ad75/serial-s01-e01.mp3",
				Filesize: 0,
				Duration: 54 * time.Minute,
				GUID:     "2 at http://serialpodcast.org",
			},
		},
	}
}
