package main

import (
	"encoding/xml"
	"errors"
	"io"
	"time"
)

type entry struct {
	Type  string `xml:"type,attr"`
	Text  string `xml:"text,attr"`
	Title string `xml:"title,attr"`
	Desc  string `xml:"description,attr"`
	URL   string `xml:"xmlUrl,attr"`
}

type opml struct {
	Version string   `xml:"version,attr"`
	Title   string   `xml:"head>title"`
	Pubdate opmltime `xml:"head>dateCreated"`
	Spec    string   `xml:"head>docs"`
	Entries []*entry `xml:"body>outline"`
}

func newOPML(title string) *opml {
	return &opml{
		Title:   title,
		Spec:    "http://dev.opml.org/spec2.html",
		Version: "2.0",
	}
}

func (o *opml) writeTo(w io.Writer) error {

	_, err := w.Write([]byte(xml.Header))
	if err != nil {
		return err
	}

	e := xml.NewEncoder(w)
	e.Indent("", "    ")

	return e.Encode(o)
}

func (o *opml) readFrom(r io.Reader) error {
	return xml.NewDecoder(r).Decode(o)
}

type opmltime time.Time

func (ot *opmltime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}

	layouts := []string{
		time.RFC1123Z, // "Mon, 02 Jan 2006 15:04:05 -0700"
		time.RFC1123,  // "Mon, 02 Jan 2006 15:04:05 MST"
		time.RFC822Z,  // "02 Jan 06 15:04 -0700"
		time.RFC822,   // "02 Jan 06 15:04 MST"
	}

	for _, l := range layouts {
		if t, err := time.Parse(l, s); err == nil {
			*ot = opmltime(t)
			return nil
		}
	}

	return errors.New("unsupported date format: " + s)
}

func (ot *opmltime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	s := time.Time(*ot).Format(time.RFC1123Z) // "Mon, 02 Jan 2006 15:04:05 -0700"
	return e.EncodeElement(s, start)
}
