package main

import (
	"fmt"
	"net/http"
	"time"
)

type linkStatus struct {
	link   string
	status bool
}

func main() {
	links := []string{"http://google.com", "http://stackoverflow.com", "http://youtube.com"}

	// Make a channel using make command. Channels are used to share data between two goroutines
	linkStatusChannel := make(chan linkStatus)

	for _, link := range links {
		// Start a new go routine
		go getLinkStatus(link, linkStatusChannel)
	}

	for linkStatus := range linkStatusChannel {
		if linkStatus.status {
			fmt.Println(linkStatus.link, "is up!")
		} else {
			fmt.Println(linkStatus.link, "is down!")
		}

		go func(l string) {
			// Wait for 10 seconds for this particular link to again get the status
			time.Sleep(10 * time.Second)
			getLinkStatus(l, linkStatusChannel)
		}(linkStatus.link)
	}
}

func getLinkStatus(link string, c chan linkStatus) {
	_, err := http.Get(link)

	if err != nil {
		c <- linkStatus{
			link:   link,
			status: false,
		}
		return
	}

	c <- linkStatus{
		link:   link,
		status: true,
	}
}
