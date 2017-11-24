/*
Open a series of urls.

Check status code for each url and store urls I could not
open in a dedicated array.
Fetch urls concurrently using goroutines.
*/

package main

import (
	"fmt"
	"net/http"
)

// -------------------------------------

// Custom user agent.
const (
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/53.0.2785.143 " +
		"Safari/537.36"
)

// -------------------------------------

// fetchUrl opens a url with GET method and sets a custom user agent.
// If url cannot be opened, then log it to a dedicated channel.
func fetchUrl(url string, chFailedUrls chan string, chIsFinished chan bool) {

	// Open url.
	// Need to use http.Client in order to set a custom user agent:
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)

	// Inform the channel chIsFinished that url fetching is done (no
	// matter whether successful or not). Defer triggers only once
	// we leave fetchUrl():
	defer func() {
		chIsFinished <- true
	}()

	// If url could not be opened, we inform the channel chFailedUrls:
	if err != nil || resp.StatusCode != 200 {
		chFailedUrls <- url
		return
	}

}

func main() {

	// Create a random urls list just as an example:
	urlsList := [10]string{
		"http://example1.com",
		"http://example2.com",
		"http://example3.com",
		"http://example4.com",
		"http://example5.com",
		"http://example10.com",
		"http://example20.com",
		"http://example30.com",
		"http://example40.com",
		"http://example50.com",
	}

	// Create 2 channels, 1 to track urls we could not open
	// and 1 to inform url fetching is done:
	chFailedUrls := make(chan string)
	chIsFinished := make(chan bool)

	// Open all urls concurrently using the 'go' keyword:
	for _, url := range urlsList {
		go fetchUrl(url, chFailedUrls, chIsFinished)
	}

	// Receive messages from every concurrent goroutine. If
	// an url fails, we log it to failedUrls array:
	failedUrls := make([]string, 0)
	for i := 0; i < len(urlsList); {
		select {
		case url := <-chFailedUrls:
			failedUrls = append(failedUrls, url)
		case <-chIsFinished:
			i++
		}
	}

	// Print all urls we could not open:
	fmt.Println("Could not fetch these urls: ", failedUrls)

}
