package rss

import (
	"bytes"
	"encoding/xml"
)

type RSSFeed struct {
	XMLName  xml.Name    `xml:"rss"`
	Encoding string      `xml:"encoding,attr"`
	Channel  *RSSChannel `xml:"channel"`
}
type RSSChannel struct {
	XMLName       xml.Name `xml:"channel"`
	URL           string   `xml:"-"`
	Title         string   `xml:"title"`
	Link          string   `xml:"link"`
	Description   string   `xml:"description"`
	Language      string   `xml:"language"`
	PubDate       string   `xml:"pubDate"`
	LastBuildDate string   `xml:"LastBuilDate"`
	Items         []Item
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Creator     string `xml:"creator"`
	Guid        string `xml:"guid"`
	Description string `xml:"description"`
	Content     string `xml:"content"`
}

func NewRSSFeed(buf []byte) (*RSSFeed, error) {
	f := RSSFeed{}
	d := xml.NewDecoder(bytes.NewReader(buf))
	if err := d.Decode(&f); err != nil {
		return nil, err
	}
	return &f, nil
}
