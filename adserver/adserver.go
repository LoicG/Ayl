package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, strings.TrimSpace(`
Adserver [OPTIONS]
`)+"\n")
		flag.PrintDefaults()
	}
	port := flag.Int("p", 8080, "server port")
	campaigns := flag.String("f", "", "campaigns file")
	flag.Parse()

	server, err := newServer(*port, *campaigns)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Port: ", *port)
	log.Println("Campaigns file: ", *campaigns)
	server.Run()
}
