package main

import (
	"fmt"
	"log"

	"github.com/ics/emm/directory"
)

func main() {
	d, err := directory.FromFile("./channeldirectory.xml", "Public")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v", d.Channels[7])

}
