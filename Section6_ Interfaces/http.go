package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Exploring http package")

	resp, err := http.Get("http://google.com")

	if err != nil {
		os.Exit(1)
	}

	defer resp.Body.Close()

	fmt.Println(resp)

	// The response body implements a readCloser interface which is a very common interface implemented by many readCloser structs.

	// WAY 1 to read the response.
	// To read the function body, we have to create a slice with an educated guess of the size. But here if we make a very big slice, we could waste a lot of space here
	// buf := make([]byte, 999999)
	// resp.Body.Read(buf)
	// fmt.Println(string(buf))

	// WAY 2 to read the response. This way is good because there would be less wastage of space.
	buf := make([]byte, 4096) // Common size allocated is 4096 bytes
	respBody := make([]byte, 0)
	for {
		// n is basically the number of characters read in the response
		n, err := resp.Body.Read(buf)

		if err != nil && err != io.EOF {
			// Handle error
			break
		}

		if n == 0 {
			// We have read all the characters
			break
		}

		respBody = append(respBody, buf[:n]...)
	}

	//fmt.Println(string(respBody))

	// WAY 3
	// io.Copy(os.Stdout, resp.Body) // Here os.Stdout implements the writer interface. The implementation of this is similar to the approach that we have in WAY 2
}
