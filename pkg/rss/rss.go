package rss

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"

	"golang.org/x/text/encoding/charmap"
)

// Feed represents the RSS feed.
type Feed struct {
	XMLName  xml.Name `xml:"rss"`
	Encoding string   `xml:"encoding,attr"`
	Channel  *Channel `xml:"channel"`
}

// Channel represents a RSS channel within a feed.
type Channel struct {
	XMLName       xml.Name `xml:"channel"`
	URL           string   `xml:"-"`
	Title         string   `xml:"title"`
	Link          string   `xml:"-"`
	LinkSlice     []string `xml:"link"`
	Description   string   `xml:"description"`
	Language      string   `xml:"language"`
	PubDate       string   `xml:"pubDate"`
	LastBuildDate string   `xml:"lastBuildDate"`
	Items         []Item
}

// Item represents a channel item
type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Creator     string `xml:"creator"`
	GUID        string `xml:"guid"`
	Description string `xml:"description"`
	Content     string `xml:"content"`
}

// NewFeed creates a new Feed from a given byte slice and returns a
// pointer to it.
func NewFeed(buf []byte) (*Feed, error) {
	f := RSSFeed{}
	d := xml.NewDecoder(bytes.NewReader(buf))
	d.CharsetReader = makeCharsetReader
	if err := d.Decode(&f); err != nil {
		return nil, err
	}
	for _, l := range f.Channel.LinkSlice {
		if l != "" {
			f.Channel.Link = l
			break
		}
	}
	return &f, nil
}

func makeCharsetReader(charset string, input io.Reader) (io.Reader, error) {
	if charset == "ISO-8859-1" || charset == "Windows-1252" {
		// Windows-1252 is a superset of ISO-8859-1, so should do here
		return charmap.Windows1252.NewDecoder().Reader(input), nil
	}
	if charset == "Windows-1255" {
		return charmap.Windows1255.NewDecoder().Reader(input), nil
	}
	return nil, fmt.Errorf("Unknown charset: %s", charset)
}
