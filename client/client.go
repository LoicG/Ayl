package main

import (
	"bytes"
	"campaigns"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Result struct {
	Campaign string
	Content  data.Content
}

type BodyRequest struct {
	Country string
	Device  string
}

type Client struct {
	host   string
	client *http.Client
}

func (c *Client) post(url string, input []byte, output interface{}) error {
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(input))
	if err != nil {
		return err
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	response, err := c.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("server request failed, error code: %d",
			response.StatusCode)
	}
	return json.NewDecoder(response.Body).Decode(output)
}

func newClient(port uint32, maxIdleConnsPerHost int) *Client {
	transport := &http.Transport{
		MaxIdleConnsPerHost: maxIdleConnsPerHost,
	}
	return &Client{
		host: fmt.Sprintf("localhost:%d",
			port),
		client: &http.Client{Transport: transport},
	}
}

func (c *Client) getCampaign(placement, country, device string) (
	interface{}, error) {

	url, err := url.Parse(fmt.Sprintf("http://%s/ad", c.host))
	if err != nil {
		return nil, err
	}
	body := &BodyRequest{
		Country: country,
		Device:  device,
	}
	buf, err := json.MarshalIndent(body, "", "  ")
	if err != nil {
		return nil, err
	}
	values := url.Query()
	values.Add("placement", placement)
	url.RawQuery = values.Encode()
	result := &Result{}
	return result, c.post(url.String(), buf, result)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, strings.TrimSpace(`
adserver client [OPTIONS]
`)+"\n")
		flag.PrintDefaults()
	}
	jobs := flag.Int("jobs", 1, "concurrent jobs")
	number := flag.Int("number", 1, "number request per jobs")
	port := flag.Uint("port", 8080, "server port")
	placement := flag.String("placement", "f64551fcd6f07823cb87971cfb914464", "placement identifier")
	country := flag.String("country", "FRA", "country")
	device := flag.String("device", "DESKTOP", "device")
	flag.Parse()

	count := int64(0)
	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			log.Printf("%d requests\n", atomic.LoadInt64(&count))
		}
	}()

	start := time.Now()
	wg := sync.WaitGroup{}
	for i := 0; i != *jobs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := newClient(uint32(*port), *jobs)
			for j := 0; j != *number; j++ {
				result, err := client.getCampaign(*placement, *country, *device)
				fmt.Printf("%v\n", result)
				if err != nil {
					log.Fatal(err)
				}
				atomic.AddInt64(&count, 1)
			}
		}()
	}
	wg.Wait()
	ticker.Stop()
	end := time.Since(start)
	log.Printf("%d Requests in: %v\n", count, end)
	seconds := int64(end.Seconds())
	if seconds != 0 {
		log.Printf("%d Requests per Second", count/seconds)
	}
}
