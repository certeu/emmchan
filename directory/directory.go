package directory

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
)

// Directory represents a channel directory tree.
type Directory struct {
	XMLName  xml.Name   `xml:"directory"`
	Channels []*Channel `xml:"channel"`
}

// Channel represents a channel entry.
type Channel struct {
	ID              string `xml:"id,attr"`
	Format          string `xml:"format"`
	Type            string `xml:"type"`
	Subject         string `xml:"subject"`
	Description     string `xml:"description"`
	Identifier      string `xml:"identifier"`
	CountryCode     string `xml:"country"`
	Region          string `xml:"region"`
	Category        string `xml:"category"`
	Ranking         int    `xml:"ranking"`
	Language        string `xml:"language"`
	UpdatePeriod    string `xml:"schedule>updatePeriod"`
	UpdateFrequency int    `xml:"schedule>updateFrequency"`
	Feeds           []Feed `xml:"feed"`
}

// Feed represents a channel feed.
type Feed struct {
	Title string  `xml:"title,attr"`
	URL   FeedURL `xml:"url,attr"`
}

func (f *Feed) EqualTo(other *Feed) bool {
	this := f.URL
	that := other.URL
	return this.Host == that.Host && this.Path == that.Path
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

// Load will load a channel directory tree from a io.Reader.
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
