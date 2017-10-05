package emm

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"regexp"

	"github.com/ics/emm/pkg/rss"
)

type Channels []*EMMChannel

// Index returns the index of first found channel with the given identifier(id).
func (c *Channels) Index(id string) int {
	for idx, ch := range *c {
		if ch.Identifier == id {
			return idx
		}
	}
	return -1
}

// Directory represents a channel directory tree.
type Directory struct {
	XMLName  xml.Name `xml:"directory"`
	Channels Channels `xml:"channel"`
}

// Add appends an EMMChannel to directory channel slice.
func (d *Directory) Add(ec *EMMChannel) {
	idx := d.Channels.Index(ec.Identifier)
	if idx != -1 {
		feeds := d.Channels[idx].Feeds
		if feeds == nil {
			feeds = ec.Feeds
		} else {
			for _, f := range *ec.Feeds {
				feeds.Add(f)
			}
		}
		d.Channels[idx].Feeds = feeds
		return
	}
	d.Channels = append(d.Channels, ec)
}

// Load will load a channel directory tree from an io.Reader.
func (d *Directory) Load(ch io.Reader) error {
	f, err := ioutil.ReadAll(ch)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(f, &d)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return err
	}

	return nil
}

// Dump writes the channel directory to an io.Writer.
func (d *Directory) Dump(ch io.Writer) error {
	out, err := xml.MarshalIndent(d, "", "  ")
	if err != nil {
		return err
	}
	ch.Write(out)
	return nil
}

// NewEMMChannel creates a new EMM channel from a RSS feed.
func NewEMMChannel(r *rss.RSSFeed) *EMMChannel {
	rc := r.Channel
	u, _ := url.Parse(rc.URL)
	if rc.Link == "" {
		rc.Link = fmt.Sprintf("%s://%s/", u.Scheme, u.Host)
	}
	feeds := Feeds{
		Feed{rc.Title, FeedURL(*u)},
	}
	e := &EMMChannel{
		RSSFeed:         r,
		Format:          "rss",
		Type:            "webnews",
		Subject:         "eucert",
		Description:     rc.Description,
		Identifier:      rc.Link,
		Encoding:        r.Encoding,
		CountryCode:     "US",
		Region:          "Global",
		Category:        "Specialist",
		Ranking:         1,
		Language:        rc.Language,
		UpdatePeriod:    "daily",
		UpdateFrequency: 4,
		Feeds:           &feeds,
	}
	e.genID()
	e.setEncoding()
	return e

}

// EMMChannel represents a channel entry.
type EMMChannel struct {
	// the RSS feed from which this channel was generared
	RSSFeed         *rss.RSSFeed `xml:"-"`
	ID              string       `xml:"id,attr"`
	Format          string       `xml:"format"`
	Type            string       `xml:"type"`
	Subject         string       `xml:"subject"`
	Description     string       `xml:"description"`
	Identifier      string       `xml:"identifier"`
	Encoding        string       `xml:"encoding"`
	CountryCode     string       `xml:"country"`
	Region          string       `xml:"region"`
	Category        string       `xml:"category"`
	Ranking         int          `xml:"ranking"`
	Language        string       `xml:"language"`
	UpdatePeriod    string       `xml:"schedule>updatePeriod"`
	UpdateFrequency int          `xml:"schedule>updateFrequency"`
	Feeds           *Feeds       `xml:"feed"`
}

func (e *EMMChannel) genID() {
	t := e.RSSFeed.Channel.Title
	if t == "" {
		t = e.RSSFeed.Channel.URL
	}
	alph, _ := regexp.Compile("[^a-zA-Z0-9]*")
	title := alph.ReplaceAllString(t, "")
	e.ID = title
}

func (e *EMMChannel) setEncoding() {
	if e.Encoding == "" {
		e.Encoding = "UTF-8"
	}
}

type Feeds []Feed

func (f *Feeds) Add(other Feed) {
	for _, feed := range *f {
		if feed.URL != other.URL {
			*f = append(*f, other)
		}
	}
}

// Feed represents a channel feed.
type Feed struct {
	Title string  `xml:"title,attr"`
	URL   FeedURL `xml:"url,attr"`
}

type FeedURL url.URL

// UnmarshalXMLAttr unmarshals the URL string into FeedURL.
// FeedURL is a url.URL.
func (f *FeedURL) UnmarshalXMLAttr(attr xml.Attr) error {
	u, err := url.Parse(attr.Value)
	if err != nil {
		return err
	}
	*f = FeedURL(*u)
	return nil
}

func (f *FeedURL) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	u := url.URL(*f)
	attr := xml.Attr{Name: name, Value: u.String()}
	return attr, nil
}

// FromFile loads a channel directory from a XML file.
func FromFile(path, name string) (*Directory, error) {
	d := &Directory{}
	f, err := os.Open(path)
	if err != nil {
		return d, err
	}
	defer f.Close()

	err = d.Load(f)
	if err != nil {
		return d, err
	}
	return d, nil
}
