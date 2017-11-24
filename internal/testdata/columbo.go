package testdata

import (
	"time"

	kibner "github.com/deepilla/kibner/internal/types"
)

func Columbo() *kibner.Feed {
	return &kibner.Feed{
		Title:  "Columbo",
		Author: "Richard Levinson & William Link",
		Desc:   "Columbo is an American television series starring Peter Falk as Columbo, a homicide detective with the Los Angeles Police Department.",
		Type:   "rss",
		URL:    "http://www.columbopodcast.com/feed/podcast/",
		Link:   "http://www.columbopodcast.com",
		Image:  "http://www.columbopodcast.com/wp-content/uploads/powerpress/big-columbo.png",
		Items: []*kibner.Item{
			{
				Title:    "Étude in Black",
				Desc:     "Alex Benedict (John Cassavetes), the married conductor of the Los Angeles Philharmonic Orchestra, murders his mistress, Jennifer Welles (Anjanette Comer), after she insists on going public with their affair.",
				Pubdate:  time.Date(1972, time.September, 18, 1, 30, 0, 0, time.UTC),
				URL:      "https://libsyn.com/columbo/s02e01.mp3",
				Filesize: 82644860,
				Duration: 98 * time.Minute,
				GUID:     "S02E01",
			},
			{
				// Missing Description
				Title:    "Blueprint for Murder",
				Pubdate:  time.Date(1972, time.February, 10, 1, 30, 0, 0, time.UTC),
				URL:      "https://libsyn.com/columbo/s01e07.mp3",
				Filesize: 82644860,
				Duration: 73 * time.Minute,
				GUID:     "S01E07",
			},
			{
				// Missing Title
				Title:    "Untitled: Jan 20, 1972",
				Desc:     "Roger Stanford (Roddy McDowall) is a chemist and photography buff whose uncle, David (James Gregory), has taken over a business that his parents built and his aunt (Ida Lupino) controls. David proposes selling the business to a conglomerate in return for a seat on the board of directors, then tries to blackmail Roger into resigning. Roger decides to murder his uncle with a box of exploding cigars.",
				Pubdate:  time.Date(1972, time.January, 20, 1, 30, 0, 0, time.UTC),
				URL:      "https://libsyn.com/columbo/s01e06.mp3",
				Filesize: 82644860,
				Duration: 73 * time.Minute,
				GUID:     "S01E06",
			},
			{
				// Missing GUID
				Title:    "Lady in Waiting",
				Desc:     "Beth Chadwick (Susan Clark) murders her domineering older brother, Bryce (Richard Anderson), after he attempts to break up her relationship with one of his executives, Peter (played by Leslie Nielsen). His mailing a letter threatening to terminate Peter's employment if he didn't break things off with Beth causes her to reach the tipping point and act to gain control of her own life and, it turns out, the family business.",
				Pubdate:  time.Date(1971, time.December, 16, 1, 30, 0, 0, time.UTC),
				URL:      "https://libsyn.com/columbo/s01e05.mp3",
				Filesize: 82644860,
				Duration: 73 * time.Minute,
				GUID:     "https://libsyn.com/columbo/s01e05.mp3",
			},
			{
				// Invalid filesize
				Title:    "Death Lends a Hand",
				Desc:     "Private investigator Carl Brimmer (Robert Culp) is hired by Arthur Kennicut (Ray Milland), a powerful publishing magnate who suspects his wife of infidelity. Although Brimmer indeed finds evidence of the wife being unfaithful, he attempts to blackmail her into revealing secrets about her husband. She refuses and tells him she will expose his plot, at which point Brimmer accidentally kills her in a fit of rage.",
				Pubdate:  time.Date(1971, time.October, 7, 1, 30, 0, 0, time.UTC),
				URL:      "https://libsyn.com/columbo/s01e02.mp3",
				Duration: 73 * time.Minute,
				GUID:     "S01E02",
			},
			{
				// Blank filesize
				Title:    "Murder by the Book",
				Desc:     "Ken Franklin (Jack Cassidy) is one-half of a mystery writing team, but partner Jim Ferris (Martin Milner) wants to go solo. This would expose the fact that Ferris did all the actual writing, and leave the high-living Franklin without his cash cow. Franklin tricks Ferris into taking a trip to his remote cabin two hours away. At the cabin, he convinces Ferris to call home and say he's working late at the office. During the call, Franklin shoots Ferris, then takes his body back north and dumps it on his lawn.",
				Pubdate:  time.Date(1971, time.September, 16, 1, 30, 0, 0, time.UTC),
				URL:      "https://libsyn.com/columbo/s01e01.mp3",
				Duration: 73 * time.Minute,
				GUID:     "S01E01",
			},
			{
				// Invalid duration
				Title:    "Ransom for a Dead Man",
				Desc:     "Leslie Williams (Lee Grant), a brilliant lawyer and pilot, murders her husband Paul (Harlan Warde) to get his money, arranging the act to look as if he had been kidnapped and killed by his captors.",
				Pubdate:  time.Date(1971, time.March, 2, 1, 30, 0, 0, time.UTC),
				URL:      "https://libsyn.com/columbo/s00e02.mp3",
				Filesize: 82644860,
				GUID:     "S00E02",
			},
			{
				// Missing duration
				Title:    "Prescription: Murder",
				Desc:     "Dr. Ray Fleming (Gene Barry), a psychiatrist, murders his wife (Nina Foch) and persuades his mistress Joan Hudson (Katherine Justice), who is an actress and one of his patients, to support his alibi by impersonating her.",
				Pubdate:  time.Date(1968, time.February, 21, 1, 30, 0, 0, time.UTC),
				URL:      "https://libsyn.com/columbo/s00e01.mp3",
				Filesize: 82644860,
				GUID:     "S00E01",
			},
			{
				// Missing pubdate
				Title:    "Dead Weight",
				Desc:     "Major General Martin Hollister (Eddie Albert), a retired Marine Corps war hero, learns he is being investigated for embezzling military funds, then shoots his skittish accomplice (John Kerr). The act is partially witnessed by Helen Stewart (Suzanne Pleshette) from a passing boat, but Hollister woos her into doubting her own story.",
				URL:      "https://libsyn.com/columbo/s01e03.mp3",
				Filesize: 82644860,
				Duration: 73 * time.Minute,
				GUID:     "S01E03",
			},
			{
				// Invalid pubdate
				Title:    "Suitable for Framing",
				Desc:     "Art critic Dale Kingston (Ross Martin) murders his uncle and tries to frame his aunt (Kim Hunter), to obtain what is considered to be one of the most valuable art collections in the world.",
				URL:      "https://libsyn.com/columbo/s01e04.mp3",
				Filesize: 82644860,
				Duration: 73 * time.Minute,
				GUID:     "S01E04",
			},
		},
	}
}
