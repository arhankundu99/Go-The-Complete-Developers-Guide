# Section 7: Channels and Go Routines
## Parallelism vs Concurrency
A lot of people use this terms interchangebly but this two are different concepts.

### Concurrency
Concurrency is about handling multiple threads and context switching between the threads to process them. Let's say we have 2 threads and one core. That means only one thread at a time can be executed at that core. Now when the thread is being processed and the runtime encounters a blocking call (like http request), then the scheduler (usually integrated with the runtime) may move this thread and process the other thread in that core. And when the previous thread finishes its blocking call, then the scheduler may again process this thread in that core.

It is more about context switching between the threads and efficiently processing them. It is not about handling threads in parallel. In the above example, we only had one thread. But due to context switching, these threads got processed efficiently.

### Parallelism
Let's say we have multiple cores, that means now if there are multiple threads, each thread can be assigned to each core and now the threads can be processed in parallel.

## Detour for javascript
Now javascript is a single threaded and synchronous language. Does it support parallelism and concurrency? Let's find out.

Javascript does not have threads, but when we have things like network io, file io, setTimeout etc, these things are handled by node runtime.

For operations like fileio, dns lookup, crypto and zlib, libuv handles it using thread pool. Network IO is handled by OS (The operations are offloaded to libuv and then libuv offloads to os) and the timers are handled by event loop (Before the timer phase of the event loop, the event loop checks whether any timer has passed and executes the callbacks). But this is the part of the javascript runtime (Eg Node) and not javascript itself. So we can say that javascript does not support parallelism but the runtime supports parallelism.

And regarding concurrency, javascript achieves concurrency with the help of its event loop (See, Although its not concurrency using multiple threads, it acheives a similar effect using event loop. For eg., When a promise is encountered, the context is switched and when the promise is resolved or rejected, the context may again be switched to execute the callback).

## Goroutines (Theory)
1. Goroutines are basic units of concurrency in go. They are blocks of code that are executed concurrently. 
2. Each go routine DOES NOT correspond to a os thread. The Go Scheduler (Which is integrated with the go runtime and manages the goroutines) maps goroutines to OS threads in M:N fashion, where M >> N generally and M is the number of go routines and N is the number of OS threads. 
3. When a goroutine is getting executed by a OS thread, and if the go routine has a blocking call (For eg., Http request get), then go scheduler would replace that go routine with another go routine in that thread. And when the blocking call finishes, the go routine is again assigned to a thread.
4. They are light weight as compared to threads.
5. Because of this efficient scheduling behaviour, Go is more suitable when we need to handle millions of concurrent requests at scale. 
6. Goroutines does not give control over OS threads. If we want a thread to be fully dedicated in a task (Even if it has a blocking call. In go, if the go routine has a blocking call, then that goroutine is switched with another goroutine by the scheduler), then we cannot do so. (We can do this in C++, because each thread in C++ corresponds to each OS thread).

## Difference with C++ Threads
1. Each C++ thread corresponds to an OS thread.
2. C++ is more suitable if we want control over threads. (For eg., if we want to dedicate a complete thread for a particular function).
3. If we want to handle concurrent requests at scale, then Go is more suitable.

## Difference between node and go
1. Node is more suitable if we have more IO intensive tasks whereas Go is more suitable, if we have combination of CPU Intensive + IO Intensive tasks.
2. Well, Node can also executed CPU intensive tasks using worker threads. But again each worker thread corresponds to each os thread unlike go routines, which makes scheduling and processing tasks concurrently very efficient in case of go routines.
3. So if go is better than Node in handling both CPU Intensive tasks + IO Intensive tasks, why don't be use go then? Well, Node has a better ecosystem, vast npm package for dependencies and good support because it can be used in both frontend and backend. But performance wise, Go is better than Node.

## Goroutines (Practice)
Consider the below code which fetches whether the status of each of the links is up or down
```go
package main

import (
	"fmt"
	"net/http"
)

func main() {
	links := []string{"http://google.com", "http://stackoverflow.com", "http://youtube.com"}

	for _, link := range links {
		linkStatus := getLinkStatus(link)

		if linkStatus {
			fmt.Println(link, "is up!")
		} else {
			fmt.Println(link, "is down!")
		}
	}
}

func getLinkStatus(link string) bool {
	_, err := http.Get(link)

	if err != nil {
		return false
	}

	return true
}
```

In the above code, the link status are checked one by one and not concurrently, We can make it concurrent using goroutines. Note Go channels are used to share data between two goroutines and the data in the channels follow FIFO structure like queues.

```go
package main

import (
	"fmt"
	"net/http"
)

func main() {
	links := []string{"http://google.com", "http://stackoverflow.com", "http://youtube.com"}

	// Make a channel using make command. Channels are used to share data between two goroutines
	linkStatusChannel := make(chan string)

	for _, link := range links {
		// Start a new go routine
		go getLinkStatus(link, linkStatusChannel)
	}

	// Now print the status of the link
	for idx := 0; idx < len(links); idx++ {
		// <- linkStatusChannel: This is a blocking call, similar to future.get() in cpp
		linkStatus := <-linkStatusChannel
		fmt.Println(linkStatus)
	}
}

func getLinkStatus(link string, c chan string) {
	_, err := http.Get(link)

	if err != nil {
        // add the data in the channel. 
		c <- link + " is down!"
		return
	}

	c <- link + " is up!"
}
```

Now what if we want to check the status of the links every 10 seconds
```go
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
```

### Checking if a channel is closed.
```go
// The below code returns 2 values
val, ok := <-channel
if ok {
    // Channel is open
} else {
    // Channel is closed
}
```
Note: The channel should be closed using ```close(channel)``` command in goroutine. 

## Wait groups
We can use wait groups in go, to wait for the goroutines to be executed completely. The below example is taken from https://www.freecodecamp.org/news/concurrent-programming-in-go. This is similar to thread.join() in cpp.
```go
package main

import (
	"fmt"
	"sync"
)

func main() {
    // Create a wait group
	var wg sync.WaitGroup

    // Add the number of go routines
	wg.Add(2)
	go helloworld(&wg)
	go goodbye(&wg)
	
    // Wait for the go routines to complete
    wg.Wait()
}

func helloworld(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Hello World!")
}

func goodbye(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Good Bye!")
}
```

Notes:
1. ```wg.Add(int)```: This method indicates the number of goroutines to wait. In the above code, We had provided 2 for 2 different goroutines. Hence the internal counter wait becomes 2.
2. ```wg.Wait()```: This method blocks the execution of code until the internal counter becomes 0.
3. ```wg.Done()```: This will reduce the internal counter value by 1.
NOTE: If a WaitGroup is explicitly passed into functions, it should be added by a pointer.
