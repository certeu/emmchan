package emm

import (
	"encoding/xml"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/ics/emm/pkg/rss"
)

const (
	cd = `<directory>
		  <channel id="P_malekalssite">
			<dc:format>rss</dc:format>
			<dc:type>webnews</dc:type>
			<dc:subject>eucert</dc:subject>
			<dc:description>malekals site</dc:description>
			<dc:identifier>http://www.malekal.com/</dc:identifier>
			<iso:country>US</iso:country>
			<region>Global</region>
			<category>Specialist</category>
			<ranking>1</ranking>
			<iso:language>en</iso:language>
			<ocs:schedule>
			  <ocs:updatePeriod>daily</ocs:updatePeriod>
			  <ocs:updateFrequency>2</ocs:updateFrequency>
			</ocs:schedule>
			<feed title="malekals site" url="http://www.malekal.com/feed/"/>
		  </channel>
		</directory>`
)

var rssFeed *rss.Feed

func init() {
	rssFeed = &rss.Feed{
		XMLName: xml.Name{Space: "", Local: "rss"},
		Channel: &rss.Channel{
			Title:         "Research Blog",
			Link:          "https://www.zscalaer.com/",
			Links:         []string{"https://www.zscaler.com/", "", ""},
			Description:   "",
			Language:      "en",
			PubDate:       "Wed, 04 Oct 2017 03:54:41 -0700",
			LastBuildDate: "Fri, 06 Oct 2017 01:00:11 -0700",
			Items:         []rss.Item(nil),
		},
	}
}

func newDirectory(xmlstr string) *Directory {
	d := NewDirectory(xmlstr)
	return d
}

func TestLoadDump(t *testing.T) {
	in := strings.NewReader(cd)
	d := &Directory{}
	if err := d.Load(in); err != nil {
		t.Errorf("Could not load channel directory")
	}
	if err := d.Dump(ioutil.Discard); err != nil {
		t.Errorf("Could not dump channel directory")
	}
}

func TestIndex(t *testing.T) {
	d := newDirectory(cd)
	tests := []struct {
		in   string
		want int
	}{
		{"http://www.malekal.com/", 0},
		{"http://cert.europa.eu/", -1},
	}
	for _, test := range tests {
		idx := d.Channels.Index(test.in)
		if idx != test.want {
			t.Errorf("Error search for %s in channel directory", test.in)
		}
	}
}

func TestNewChannel(t *testing.T) {
	c := NewChannel(rssFeed)
	if c.ID != "ResearchBlog" && c.Identifier != "https://www.zscalaer.com/" {
		t.Errorf("Could not create new channel")
	}
}

func TestAdd(t *testing.T) {
	d := newDirectory(cd)
	c := NewChannel(rssFeed)
	d.Add(c)
	d.Add(c)
	if len(d.Channels) != 2 {
		t.Errorf("Channel wasn't properly added.")
	}
}
