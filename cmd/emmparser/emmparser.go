// The emm command reads URLS from STDIN and adds them to a channel directory.
// It is unfortunate that EMM channels share the name with the Go primitive.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ics/emm/pkg/emm"
	"github.com/ics/emm/pkg/rss"
)

var (
	chDir   = flag.String("d", "", "Load channel directory file")
	outFile = flag.String("o", "out.xml", "Write channel directory to file")
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
			//log.Printf("Downloading %v\n", u)
			rssFeed, err := getFeed(u)
			fmt.Printf("%#v", rssFeed.Channel)
			if err != nil {
				log.Printf("Error in %s: %s", u, err)
				return
			}
			emmCh := emm.NewChannel(rssFeed)
			//log.Printf("Adding to channel directory %v\n", u)
			d.Add(emmCh)
		case <-done:
			// fmt.Printf("[%d] Processing done!\n", workerId)
			return
		}
	}
}

func main() {
	startTime := time.Now()
	flag.Parse()
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
		inCh <- s.Text()
	}
	if err := s.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	close(doneCh)
	wg.Wait()
	log.Printf("All done! Writing new channel directory to %s...\n", *outFile)

	f, err := os.OpenFile(*outFile, os.O_RDWR|os.O_CREATE, 0640)
	if err != nil {
		log.Fatal(err)
	}

	if err := d.Dump(f); err != nil {
		log.Fatal(err)
	}

	elapsed := time.Since(startTime)
	log.Printf("Elapsed time: %s\n", elapsed)

}
