// The emm command reads URLS from STDIN and adds them to a channel directory.
// It is unfortunate that EMM channels share the name with the Go primitive.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ics/emm/pkg/emm"
	"github.com/ics/emm/pkg/rss"
)

var (
	outFile = flag.String("o", "out.xml", "New channel directory file")
)

func getRSSFeed(feedURL string) (*rss.RSSFeed, error) {
	client := emm.NewClient(nil)

	resp, err := client.Get(feedURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	rssFeed, err := rss.NewRSSFeed(body)
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
			log.Printf("Downloading %v\n", u)
			rssFeed, _ := getRSSFeed(u)
			emmCh := emm.NewEMMChannel(rssFeed)

			log.Printf("Adding to channel directory %v\n", u)
			fmt.Printf("%#v\n", emmCh)
			d.Add(emmCh)
			//if err := processURL(url, d); err != nil {
			//	log.Printf("Error processing %s\n", url)
			//}
		case <-done:
			// fmt.Printf("[%d] Processing done!\n", workerId)
			return
		}
	}
}

func main() {
	startTime := time.Now()

	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s [options] <channelsfile>\n", os.Args[0])
		os.Exit(1)
	}
	d, err := emm.FromFile(os.Args[1], "Public")
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
	log.Printf("Processing took %s\n", elapsed)

}
