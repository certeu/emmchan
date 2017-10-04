package emm

import (
	"io/ioutil"
)

const (
	Pending = iota
	Polling
	Error
	Done
)

type State struct {
	url    string
	status int
}

type Resource struct {
	URL  string
	Body []byte
}

func (r *Resource) Fetch() error {
	client := NewClient(nil)
	resp, err := client.Get(r.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	r.Body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func Fetcher(in <-chan string, out chan<- *Resource) {
	for url := range in {
		r := &Resource{URL: url}
		r.Fetch()
		out <- r
	}
	close(out)
}
