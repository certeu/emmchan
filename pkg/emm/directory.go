package emm

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/ics/emm/pkg/rss"
)

// Channels hold all EMM channels.
type Channels []*Channel

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
	sync.Mutex
	// Marks EMM instance
	Instance string
	XMLName  xml.Name `xml:"directory"`
	Channels Channels `xml:"channel"`
}

// Add appends an Channel to directory channel slice.
func (d *Directory) Add(ec *Channel) {
	d.Lock()
	defer d.Unlock()
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

// NewDirectory create a new channel directory from an XML string
func NewDirectory(xmlstr string) *Directory {
	d := &Directory{}
	d.Load(strings.NewReader(xmlstr))
	return d
}

// NewChannel creates a new EMM channel from a RSS feed.
func NewChannel(r *rss.Feed, inst string) *Channel {
	if inst == "" {
		inst = "Public"
	}
	rc := r.Channel
	u, _ := url.Parse(rc.URL)
	if rc.Link == "" {
		rc.Link = fmt.Sprintf("%s://%s/", u.Scheme, u.Host)
	}
	feeds := Feeds{
		Feed{rc.Title, FeedURL(*u)},
	}
	e := &Channel{
		Feed:            r,
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
	e.genID(inst)
	e.setEncoding()
	return e

}

// Channel represents a channel entry.
type Channel struct {
	// the RSS feed from which this channel was generared
	Feed            *rss.Feed `xml:"-"`
	ID              string    `xml:"id,attr"`
	Format          string    `xml:"format"`
	Type            string    `xml:"type"`
	Subject         string    `xml:"subject"`
	Description     string    `xml:"description"`
	Identifier      string    `xml:"identifier"`
	Encoding        string    `xml:"encoding"`
	CountryCode     string    `xml:"country"`
	Region          string    `xml:"region"`
	Category        string    `xml:"category"`
	Ranking         int       `xml:"ranking"`
	Language        string    `xml:"language"`
	UpdatePeriod    string    `xml:"schedule>updatePeriod"`
	UpdateFrequency int       `xml:"schedule>updateFrequency"`
	Feeds           *Feeds    `xml:"feed"`
}

func (e *Channel) genID(inst string) {
	t := e.Feed.Channel.Title
	if t == "" {
		t = e.Feed.Channel.URL
	}
	alph, _ := regexp.Compile("[^a-zA-Z0-9]*")
	title := alph.ReplaceAllString(t, "")
	if inst == "Private" {
		title = fmt.Sprintf("P_%s", title)
	}
	e.ID = title
}

func (e *Channel) setEncoding() {
	if e.Encoding == "" {
		e.Encoding = "UTF-8"
	}
}

// Feeds reprents a collection of feeds within an EMM channel.
type Feeds []Feed

// Add appends a new Feed to the feed collections.
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

// FeedURL is the custom type for a feed URL.
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

// MarshalXMLAttr serializes a FeedURL.
func (f *FeedURL) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	u := url.URL(*f)
	attr := xml.Attr{Name: name, Value: u.String()}
	return attr, nil
}

// FromFile loads a channel directory from a XML file.
func FromFile(path, inst string) (*Directory, error) {
	d := &Directory{Instance: inst}
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
