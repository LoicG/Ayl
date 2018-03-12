package main

import (
	"campaigns"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"

	gouuid "github.com/pborman/uuid"
)

var defaultContent = data.Content{
	Title:       "title",
	Description: "description",
	Landing:     "landing",
}
var defaultDevices = []string{"DESKTOP", "MOBILE"}
var defaultCountries = []string{"FRA", "ALL", "BEL", "ESP", "EN"}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, strings.TrimSpace(`
campaigns.json generator [OPTIONS]
`)+"\n")
		flag.PrintDefaults()
	}
	campaigns := flag.Uint("campaigns", 1, "campaign count")
	placements := flag.Uint("placements", 1, "placement count")
	path := flag.String("path", "", "output path")
	flag.Parse()

	if *campaigns < 1 || *placements < 1 {
		log.Fatal(fmt.Errorf("invalid argument"))
	}

	pIds := []string{}
	for c := uint(0); c < *placements; c++ {
		pIds = append(pIds, gouuid.New())
	}
	output := &data.Campaigns{}
	output.Elements = map[string]*data.Campaign{}
	for c := uint(0); c < *campaigns; c++ {
		output.Elements[gouuid.New()] = &data.Campaign{
			Price:      rand.Float32() * 4,
			Devices:    data.MakeList(defaultDevices[:rand.Intn(len(defaultDevices)+1)]),
			Content:    &defaultContent,
			Countries:  data.MakeList(defaultCountries[:rand.Intn(len(defaultCountries)+1)]),
			Placements: data.MakeList(pIds[:rand.Intn(len(pIds)+1)]),
		}
	}
	err := data.SaveJSONFile(*path, output)
	if err != nil {
		log.Fatal(err)
	}
}
