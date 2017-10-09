// The emm command reads URLS from STDIN and adds them to a channel directory.
// On EOF the channel directory is written to STDOUT.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"sync"

	"github.com/ics/emm/pkg/emm"
	"github.com/ics/emm/pkg/rss"
)

var buildInfo string

var (
	chDir   = flag.String("d", "", "Channel directory file path")
	version = flag.Bool("v", false, "Display version and exit")
)

func getFeed(feedURL string) (*rss.Feed, error) {
	client := emm.NewClient(nil)

	resp, err := client.Get(feedURL)
	if err != nil {
		return nil, err
	}
	defer func() {
		// Drain up to 512 bytes and close the body to let
		// the Transport reuse the connection
		io.CopyN(ioutil.Discard, resp.Body, 512)
		resp.Body.Close()
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	rssFeed, err := rss.NewFeed(body)
	if err != nil {
		return nil, err
	}
	rssFeed.Channel.URL = feedURL
	return rssFeed, nil
}

func processChannel(inCh chan string, done <-chan bool, d *emm.Directory, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case u := <-inCh:
			rssFeed, err := getFeed(u)
			if err != nil {
				log.Printf("Error in %s: %s", u, err)
			} else {
				emmCh := emm.NewChannel(rssFeed)
				d.Add(emmCh)
			}
		case <-done:
			return
		}
	}
}

func validInput(in string) error {
	u, err := url.Parse(in)
	if err != nil {
		return fmt.Errorf("Could not parse URL: %s", err)
	}
	if u.Scheme == "" || u.Host == "" || u.Path == "" {
		return fmt.Errorf("Invalid URL %s", u)
	}
	return nil
}

func main() {
	flag.Parse()
	if *version {
		fmt.Printf("Version: %s\n", buildInfo)
		return
	}
	if *chDir == "" {
		fmt.Printf("Could not load channel directory\n")
		flag.Usage()
		os.Exit(1)
	}
	d, err := emm.FromFile(*chDir, "Public")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Loaded channel directory with %d channels", len(d.Channels))

	var wg sync.WaitGroup
	doneCh := make(chan bool)
	inCh := make(chan string)
	numWorkers := 100

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go processChannel(inCh, doneCh, d, &wg)
	}

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		if s.Text() == "" {
			continue
		}
		if err := validInput(s.Text()); err == nil {
			inCh <- s.Text()
		}
	}

	if err := s.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	close(doneCh)
	wg.Wait()

	if err := d.Dump(os.Stdout); err != nil {
		log.Fatal(err)
	}
}
