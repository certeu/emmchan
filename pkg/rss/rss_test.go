package rss

import (
	"encoding/xml"
	"testing"
)

var tests = []struct {
	in   string
	want *Feed
}{
	{`<?xml version="1.0" encoding="UTF-8"?>
	<rss xmlns:atom="http://www.w3.org/2005/Atom" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:feedburner="http://rssnamespace.org/feedburner/ext/1.0" version="2.0" xml:base="https://www.zscaler.com/">
	   <channel>
		  <title>Research Blog</title>
		  <link>https://www.zscaler.com/</link>
		  <language>en</language>
		  <pubDate>Wed, 04 Oct 2017 03:54:41 -0700</pubDate>
		  <lastBuildDate>Fri, 06 Oct 2017 01:00:11 -0700</lastBuildDate>
		  <atom10:link xmlns:atom10="http://www.w3.org/2005/Atom" rel="self" type="application/rss+xml" href="http://feeds.feedburner.com/zscaler/research" />
		  <feedburner:info uri="zscaler/research" />
		  <atom10:link xmlns:atom10="http://www.w3.org/2005/Atom" rel="hub" href="http://pubsubhubbub.appspot.com/" />
		  <item>
			 <title>Infostealer spreading through a compromised website</title>
			 <link>...</link>
			 <description>The Zscaler ThreatLabZ ...</description>
			 <author>tdewan@zscaler.com</author>
			 <pubDate>October 04, 2017</pubDate>
			 <source url="https://www.zscaler.com/...">Research Blog</source>
			 <feedburner:origLink>http://www.zscaler.com/...</feedburner:origLink>
		  </item>
		</channel>
		</rss>`,
		&Feed{
			XMLName: xml.Name{Space: "", Local: "rss"},
			Channel: &Channel{
				Title:         "Research Blog",
				Link:          "https://www.zscalaer.com/",
				Links:         []string{"https://www.zscaler.com/", "", ""},
				Description:   "",
				Language:      "en",
				PubDate:       "Wed, 04 Oct 2017 03:54:41 -0700",
				LastBuildDate: "Fri, 06 Oct 2017 01:00:11 -0700",
				Items:         []Item(nil),
			},
		},
	},
}

var testsBad = []struct {
	in   string
	want *Feed
}{
	{`<html></html>`, &Feed{}},
}

func TestNewFeed(t *testing.T) {

	for _, test := range tests {
		buf := []byte(test.in)
		actual, err := NewFeed(buf)
		if err != nil {
			t.Fatalf("In %v: unexpected error: %s", test.in, err)
		}
		if actual.Channel.Title != test.want.Channel.Title {
			t.Errorf("Have %#v, want %#v", actual, test.want)
		}
	}

	for _, test := range testsBad {
		buf := []byte(test.in)
		_, err := NewFeed(buf)
		if err != nil {
			if err.Error() != "expected element type <rss> but have <html>" {
				t.Errorf("%s", err)
			}
		}
	}
}
