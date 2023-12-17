# Section 6: Interfaces
## Interfaces
Interfaces in go is represented by two components internally: A type and a value. The type of the interface holds information about the type and value part holds the actual value. 
```go
var a interface = 2
```

Interface variables are dynamic. Which means they can store values of any type. ```interface{}``` is an interface which has no method declarations. So all the datatypes implicitly implement this interface. And this is why the below function can be called with any datatype.
```go
func anyValue(a interface{}) {
	fmt.Println(a)
}
```

How is it different from generics? 
1. In generics, we know at the compile time only that what types are going to implement the generic / template class. And so corresponding classes are generated (At semantic analysis phase. So if there are templates, then AST is again updated before handing over the AST to the compiling phase). And so there is type safety. 
2. However, In case of interface{}, there would be no type safety as all the types implement this interface. (By type safety, we mean that any type can be passed in this function)
3. And in case of interface{}, we also have to do type assertions like below. This could cause unsignificant performance decrease since we have to use runtime assertions like this.
```go
if val, ok := a.(int); ok {
    fmt.Println("a is an int:", val)
}
```

And using interfaces, we can have maps, slices etc with any datatype as shown below
```go
package main

import "fmt"

func main() {
	slice := make([]interface{}, 2)
	slice[0] = 2
	slice[1] = "Hi"

	fmt.Println(slice)

	m := make(map[interface{}]interface{})
	m[2] = "Hi"
	m["Hi"] = 2

	fmt.Println(m)
}
```

Nil Interface
```go
var c interface{}
fmt.Println(c == nil) // true
```

NOTE:
```go
// For the below interface, type and value does not hold any relavance.
type i interface{
	someMethod()
}
```

## Usage of interfaces
```go
package main

import "fmt"

type englishBot struct{}

type spanishBot struct{}

func (eb englishBot) getGreeting() string {
	return "Hello!"
}

// If the receiver member variable name is not used, we can omit the name
func (spanishBot) getGreeting() string {
	return "Hola!"
}

// Now if we want a common method that prints the greeting of a bot, we have to write duplicated logic

// The below two functions would result in an error because there is NO compile time polymorphism in go.
func printGreeting(eb englishBot) {
	fmt.Println(eb.getGreeting())
}

func printGreeting(sb spanishBot) {
	fmt.Println(sb.getGreeting())
}

func main() {
	var eb englishBot
	var sb spanishBot

	fmt.Println(eb.getGreeting())
	fmt.Println(sb.getGreeting())
}
```

Instead we can either write two methods with different name or 2 receiver methods with the same name but there is a duplication of logic here, And we can solve this using interfaces.

```go
package main

import "fmt"

type bot interface {
	getGreeting() string
}

type englishBot struct{}

type spanishBot struct{}

// Implicitly implements the interface bot by providing definition to getGreeting() method
func (eb englishBot) getGreeting() string {
	return "Hello!"
}

// Implicitly implements the interface bot by providing definition to getGreeting() method
// If the receiver member variable name is not used, we can omit the name
func (spanishBot) getGreeting() string {
	return "Hola!"
}

func printGreeting(b bot) {
	fmt.Println(b.getGreeting())
}

func main() {
	var eb englishBot
	var sb spanishBot

	printGreeting(eb)
	printGreeting(sb)
}
```
Notes
1) Interfaces are not generic types.
2) Interfaces are implicit. It saves us from some broilerplate, but the flip side is that Checking whether a struct has implemented an interface may take some time because we have to check whether it has implemented all the methods of that interface or not.

We can also embed interfaces in other interfaces like the below example
```go
package main

import "fmt"

type interface1 interface {
	method1()
}

type interface2 interface {
	method1()
	method2()
}

type combinedInterface interface {
	interface1
	interface2
}

type struct1 struct{}

func (s struct1) method1() { fmt.Println("method1") }
func (s struct1) method2() { fmt.Println("method2") }

func someMethod(ci combinedInterface) {
	ci.method1()
	ci.method2()
}

func main() {
	someMethod(struct1{})
}
```

## http package
```go
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

	fmt.Println(string(respBody))

    // WAY 3
	io.Copy(os.Stdout, resp.Body) // Here os.Stdout implements the writer interface. The implementation of this is similar to the approach that we have in WAY 2
}
```

See io.Copy function: https://pkg.go.dev/io

This function implements the above function in a similar way

## NOTE

Reader and writer interfaces are super improtant and they are used in many places in documentation of go. Reader takes in a source of information and reads the info into a byte slice that we provide. Whereas the writer takes in a byte slice that is already filled with information and writes to destinations like console output, file, etc, 

READ MORE ABOUT THEM IN THE DOCUMENTATION.