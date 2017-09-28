package directory

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"unicode"
)

// Directory represents a channel directory tree.
type Directory struct {
	XMLName  xml.Name `xml:"directory"`
	Channels Channels `xml:"channel"`
}

// Channels represents the channels slice.
type Channels []*Channel

func (c Channels) Len() int {
	return len(c)
}

func (c Channels) Swap(i, j in) {
	d[i], d[j] = d[j], d[i]
}

func (c Channels) Less(i, j int) bool {
	iRunes := []rune(c[i])
	jRunes := []rune(c[j])

	max := len(iRunes)
	if max > len(jRunes) {
		max = len(jRunes)
	}

	for idx := 0; idx < max; idx++ {
		ir := iRunes[idx]
		jr := jRunes[idx]

		lir := unicode.ToLower(ir)
		ljr := unicode.ToLower(jr)

		if lir != ljr {
			return lir < ljr
		}

		// the lowercase runes are the same, so compare the original
		if ir != jr {
			return ir < jr
		}
	}

	return false
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
	Feed            []Feed `xml:"feed"`
}

// Feed represents a channel feed.
type Feed struct {
	Title string `xml:"title,attr"`
	URL   string `xml:"url,attr"`
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
