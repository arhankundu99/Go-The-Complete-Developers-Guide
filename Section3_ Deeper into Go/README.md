# Section 3: Deeper into Go

## Variable declarations
Question: Why is short form syntax of variable declaration only allowed inside the functions?
```go
package main

import "fmt"

func main() {
	// One way of declaring and defining variables
	var card string = "Red five diamonds"

	// We can omit datatype part above (And It would work in the same way. Go compiler would infer the datatype from rhs)
	var card2 = "Black three spades"

	// Why Go uses var and not more specific datatypes like int, string, bool etc
	// Bacause in go, we don't have to explicitly specify the datatype and the datatype would be inferred automatically.

    // Short form syntax for the above declaration. (NOTE that this short form can only be used in the functions only and it cannot be used outside the functions (Why?))
	card3 := "Black four clubs"

	// If we do not assign any value, it would take default values (For eg., in case of bool, it is false, in case of int or float64, it is zero
	// and in case of string, it is empty string). And in this case, we have to specify the datatype
	var card4 string
    
    // Reassigning an already defined variable. Note that there is no colon here.
    card3 = "Black three clubs"

	fmt.Println(card)
	fmt.Println(card2)
	fmt.Println(card3)
	fmt.Println(card4)
}

```

## Functions

```go
// Here the return type is explicitly required.
func newCard() string {
    return "Black five diamonds"
}

```
Why does go require function return type to be explicitly declared? Why can't it infer the datatype from the return statement? Like in case of variable declarations, the return type is inferred from the rhs.

1. This is a design choice of Go. Specifying the return type explicitly in function signatures establishes a clear contract. Anyone using the function knows what to expect without having to look at the implementation.

2. Explicit return types ensure all implementations of a function adhere to the same contract. This is crucial for intefaces and polymorphism where different implementations return same type. 

Functions in go, can also return multiple data types. The function signature is shown below
```go
func someFunc() (val int, e error)
```
Go implements this in the same way of handling single return values. But why don't languages like C++ / Java implement this? Because if historical reasons. When C++ / Java was being developed, function having multiple return values was not a common pattern.

## Slices and arrays

Please Refer: 
1. https://go.dev/blog/slices-intro
2. https://go.dev/blog/slices

### Arrays

```go
// Declaration of an array
var a [4]int
a := [4]int

// Declaration and definition in the same line
var a [4]int = {0, 1, 2, 3}
a := [4]int{0, 1, 2, 3}

// If the array is not explicitly initialised, then the array elements would have zero values. For eg., 0 in case of int, false, in case of bool etc
```

<i>NOTE: Go arrays are <b>values</b>. The array variable denotes the entire array. It is not a pointer to the first array element (as would be the case in C). This means that when you assign or pass around an array value you will make a copy of its contents. (To avoid the copy, we could pass a pointer to the array, but then that’s a pointer to an array, not an array.) One way to think about arrays is as a sort of struct but with indexed rather than named fields: a fixed-size composite value.</i>

### Why are arrays in Go values (i.e, why are they not like C++ / C arrays where the name of the array is the address of the first element of the array)

1. When we pass an array as a function parameter, Go compiler copies the entire array, so if we make any change to that array inside the function, the original array won't be affected by that change. (i.e., we would have predictable behaviour)

2. But this also introduces performance overhead since go compiler would have to copy the array.

3. This design decision is taken since the design philosophy of go is to keep things simple. 

### Slices

Slices in go are basically dynamic arrays.

Slices have this structure.
```go
struct sliceHeader{
    length int
    capacity int
    zerothElement *byte
}
```
Note that the above struct is not available to programmer. This is for understanding how slices work under the hood

```go
// Declararion of slices
var slice int[] // Creates a nil slice
// This is equivalent to
slice := sliceHeader {
    length: 0,
    capacity: 0
    zerothElement: nil
}

// Another declaration. Here the length and capacity both are set to 3
var slice []int = []int{1, 2, 3}

// Short hand declaration
slice := []int{1, 2, 3}

// Or if we have an array defined. 
buffer := [5]int{1, 2, 3, 4, 5}

// Here the length is 3, But capacity is 4. Because the buffer can hold 5 elements and 0th element was cut off, so that capacity is 4.
var slice []int = buffer[1:4]
slice := buffer[1:4] // {2, 3, 4}

// SliceHeader of the above slice looks like this
slice := sliceHeader {
    length: 3,
    capacity: 4,
    zerothElement: &buffer[1]
}

// Now we can slice slices also
slice2 := slice[2:3] // {3}
// Slice2 looks like the below
slice2 := sliceHeader {
    length: 1,
    capacity: 3,
    zerothElement: &buffer[2]
}
```

<i>NOTE: The Capacity field is equal to the length of the underlying array, minus the index in the array of the first element of the slice.</i>

```go
slice2 := slice[a:b] // returns a slice with elements from index a to b - 1

slice2 = slice[a:] // returns a slice with elements from index a to len(slice) - 1

slice2 = slice[:b] // returns a slice with elements from index 0 to b - 1

slice2 = slice[:] // returns a slice with all the elements, and same capacity and length

slice2 = slice[0:0] // returns a slice with zero length, capacity same as the capacity of the original slice and a pointer to the first element of the underlying array.

// Short syntax for nil slice
slice3 := []int{} // returns a slice with zero length, zero capacity and the pointer to the array is nil.

var slice4 []int // This also creates a nil slice

```
<i>Note: The difference between zero slice and nil slice is that the pointer to the underlying array is not nil in case of former. And also the capacity may or may not be zero for the zero slice</i>

### Nil vs empty slice
```go
func main() {
	var slice []int
	// This creates a slice where it has no pointer to the underlying array.
	if slice == nil {
		fmt.Println("slice is nil")
	}

	// This creates an slice with zero capacity and zero length, but a pointer is assigned
	// to the underlying array. However this pointer does not point to a valid memory location
	slice2 := make([]int, 0, 0)
	if slice2 == nil {
		fmt.Println("slice2 is nil")
	} else {
		fmt.Println("slice2 is not nil")
	}
}
```

### Nil
<b>For a value to be assigned nil, type must be a pointer, channel, func, interface, map, or slice type.</b> Do not confuse it with null in other languages like java. See below content for more examples

### Passing slices in functions

Consider the below code
```go
func AddOneToEachElement(slice []byte) {
    for i := range slice {
        slice[i]++
    }
}
func main() {
    slice := buffer[10:20]
    for i := 0; i < len(slice); i++ {
        slice[i] = byte(i)
    }
    fmt.Println("before", slice)
    AddOneToEachElement(slice)
    fmt.Println("after", slice)
}
```
The output is
```
before [0 1 2 3 4 5 6 7 8 9]
after [1 2 3 4 5 6 7 8 9 10]
```

Go is pass by value, but since the slice has a pointer to the array, the same array gets modified.

### Make command 
Make command is used to create a new slice with the given length and capacity
```go
// Create a slice of type int, with length 10 and capacity 15.
slice := make([]int, 10, 15)

// if we do not specify the capacity, the capacity would be same as length
slice2 := make([]int, 10)
```

make() allocates memory on the heap (Since this function returns a pointer to the array and the array memory should not be deallocated when the make() function call ends) (Source: https://www.educative.io/answers/what-is-golang-function-maket-type-size-integertype-type)

```go
slice := make([]int, 10, 15)
fmt.Println(slice)

// fmt.Println(slice[10]) // Will throw an error since the idx >= length

slice = slice[:len(slice)+1]
fmt.Println(slice[10]) // This will not throw any error

// But we cannot increase the length of the slice beyond the capacity.
// slice = slice[:len(slice)+5] // Will throw an error
// fmt.Println(slice)

slice3 := make([]int, 2, 5)
fmt.Println(slice3) // [0, 0]

fmt.Println(slice3[2]) // Will give error since the idx >= len(slice3)
fmt.Println(slice3[2:4]) // [0, 0]. Prints a slice.
```


### Appending an element to our slice
```go
slice := []int{1, 2, 3} // Create a slice with length and capacity set to 3

slice = append(slice, 4, 5, 6) // Append integers to the slice

fmt.Println(slice) // [1, 2, 3, 4, 5, 6]
```

### Internal implementation of append
```go
package main

import (
	"fmt"
)

// elements is also a slice (Analogous to variadic arguments in C++)
func customAppend(slice []int, elements ...int) []int {
	sliceLength := len(slice)
	elementsLength := len(elements)
	sliceCapacity := cap(slice)

	if sliceLength + elementsLength > sliceCapacity {

		// That means we have to create a new slice
		newSliceCapacity := 2 * (sliceLength + elementsLength)
		newSliceLength := sliceLength + elementsLength
		newSlice := make([]int, newSliceLength, newSliceCapacity)

		copy(newSlice, slice)
		// The above copy statement expands to
		// for idx, element := range slice {
		// 	newSlice[idx] = element
		// }
		slice = newSlice
	}

	copy(slice[sliceLength : sliceLength + elementsLength], elements)
	return slice

}

func main() {
	slice := []int{1, 2, 3}        // Create a slice with length and capacity set to 3
	slice = append(slice, 4, 5, 6) // Append integers to the slice
	fmt.Println(slice)             // [1, 2, 3, 4, 5, 6]

	slice = customAppend(slice, 7, 8, 9)
	fmt.Println(slice) // [1, 2, 3, 4, 5, 6, 7, 8, 9]
}

```

Question: While appending to the slices, if the length exceeds the capacity, it returns a new slice with the copied and appended elements. Instead, why doesn't the append function just create a new array, and then change the pointer of the existing slice to the new array pointer?

Because <b>Go is a pass by value</b> language. If we change the pointer inside the append function, then if we are just changing the pointer internally. The original slice won't get modified.



## Strings

Strings are <b>READ ONLY SLICES</b> of runes. (Runes is an alias for int32). Why not use bytes as characters are 1 byte in size? Because by using runes, characters of other languages would also be taken care of.

```go
slash := "/usr/ken"[0] // yields the byte value '/'.
usr := "/usr/ken"[0:4] // yields the string "/usr"

// Convert a slice to string
str := string(slice)

// Convert a string to slice
slice := []byte(usr) // prints ascii values of the characters.
```

## For loop
```go
nums := [4]int{1, 2, 3, 4}

// range is a special keyword 
for idx, num := range nums {
    fmt.Println(idx, num)
}
```

## OO Approach vs Go Approach
<b>NOTE: Go is not an object oriented language!</b>. There is not idea of classes in go

### type in go
the type keyword in go is used to create a new datatype based on an underlying datatype and also we can define methods on the newly created datatype. 
```go
// If the first letter of the new type is Capital, then this type can be used (exported) in different packages
type Deck []int

// If the first letter of the new type is small, then this type can only be used in the current package.
type cards []int
```
<b><i>Note: Even if we have 2 types with same underlying type, WE CANNOT use them interchangebly.</i></b>

<b><i>Note: Difference between typedef in C++ and type in go is that typedef in C++ is used to create an alias, whereas type in go is used to create a new datatype and also we can create receiver methods (explained below) for this newly created datatype.</i></b>

If we want to use them interchangebly, we would have to typecast into other datatype
```go
newCards := cards{1, 2, 3}
newDeck := Deck{1, 2, 3}

// typecast cards to deck
newDeck2 := Deck(newCards)
```

And then we can attach methods to the datatype we created. These are called <b>Receiver Methods</b>. In the below method, d is the receiver for this method.
```go

// Common convention is that the parameter name is usually is first letter of the type and it should be lowercase.
func (d Deck) print() {
	for _, num := range d {
		fmt.Println(num)
	}
}

// In the above receiver method, the deck d is passed by value, that means if we update any value of deck d, it wont be reflected outside. If we want it to reflect outside also, then use pointers (See below syntax)

func (d *Deck) print() {
	for _, num := range d {
		fmt.Println(num)
	}
}

```

### Difference between Go's approach and normal object oriented language.
#### Polymorphism
Go supports runtime polymorphism and NOT compile time polymorphism. 

Eg of runtime polymorphism:
```go
package main

import (
	"fmt"
)

type shape interface {
	area() float64
}

type square struct {
	side float64
}

func (s square) area() float64 {
	return s.side * s.side
}

// Runtime polymorphism.
func getArea(s shape) float64 {
	return s.area()
}

func main() {
	sq := square{side: 2}
	fmt.Println(getArea(sq))
}
```

Why there is no compile time polymorphism in go? As go's design philosophy is to keep things simple and readable, go wants to give an unique name for each method. 

#### Inheritance
Go's struct embedding is somewhat analogous to inheritance, But it's more accurate to think of it has composition. Go focusses more on composition because the outer structure has the complete ownership of the embedded struct which makes the code simple and easy to read.

```go
package main

import (
	"fmt"
)

type shape struct {
	color string
	name  string
}

func (s shape) print() {
	fmt.Println("Color:", s.color)
	fmt.Println("Name:", s.name)
}

type square struct {
	shape // struct embedding. Now this struct will have access to the fields and the receiver methods of the embedded struct in the shape struct

	// the normal declartion of the embedded struct also would work
	// s shape or
	// shape shape (this is equivalent to just writing shape)
	side  int
}

func main() {
	sq := square{
		shape: shape{color: "red", name: "square"},
		side:  4,
	}
	fmt.Println(sq)       // {{red square} 4}
	fmt.Println(sq.shape) // {red square}

	sq.print() // Now we can access shape's methods from the square object
}

```

But if there are multiple embedded structs and some of them have same method names, then invoking the method on the struct would result in ambiguity error (See example below)

```go
package main

import (
	"fmt"
)

type shape struct {
	color string
	name  string
}

func (s shape) print() {
	fmt.Println("Color:", s.color)
	fmt.Println("Name:", s.name)
}

type dimension struct {
	size int
}

func (d dimension) print() {
	fmt.Println("Size:", d.size)
}

type square struct {
	shape // struct embedding. Now this struct will have access to the fields in the shape struct
	dimension
}

func main() {
	sq := square{
		shape:     shape{color: "red", name: "square"},
		dimension: dimension{size: 4},
	}
	fmt.Println(sq)       // {{red square} 4}
	fmt.Println(sq.shape) // {red square}

	sq.print() // But if there are ambiguos methods, (i.e, methods with the same name, there would be ambiguity).
	// In this case, this results in error
}
```

Overriding the methods in the embedded structs
```go
package main

import (
	"fmt"
)

type shape struct {
	color string
	name  string
}

func (s shape) print() {
	fmt.Println("Color:", s.color)
	fmt.Println("Name:", s.name)
}

type square struct {
	shape // struct embedding. Now this struct will have access to the fields in the shape struct
	size  int
}

// if we want to override a method of the embedded struct
func (s square) print() {
	fmt.Println("This has overriden the print method of shape.")
}

func main() {
	sq := square{
		shape: shape{color: "red", name: "square"},
		size:  4,
	}
	fmt.Println(sq)       // {{red square} 4}
	fmt.Println(sq.shape) // {red square}

	sq.print() // This has overriden the print method of shape.

	// But we can still print the shape's print method
	sq.shape.print()
	// Color: red
	// Name: square
}
```

#### Abstraction
In go, we have interfaces, which can contain only empty methods, just like interfaces in oo languages like C++, Java
```go
type shape interface {
	// Notice there is no func() declaration
	area() float64
}
```


#### Encapsulation
In go, there is visibility rules are applied in package level. It is somewhat different than the traditional encapsulation that we have in OO language.

```go
type square struct {
	size  int    // the first letter of size is small, so it cannot be accessed from different package
	Color string // the first letter is capital, so it can be accessed from different package.
}

// The first letter of the method is capital, so it can be accessed from different packages
func (s square) GetSize() int {
	return s.size
}

// The first letter is small, so it cannot be accessed from different packages
func (s *square) setSize(newSize int) {
	s.size = newSize
}
```

And also invoking methods on pointers is same as invoking methods on the types itself. The syntax is same (Refer below example).

```go
func (s *square) setSize(newSize int) {
	s.size = newSize // Same as (*s).size = newSize
}
```

Do not confuse with runtime polymorphism with the embedded structs
```go
package main

import (
	"fmt"
)

type square struct {
	size  int    // the first letter of size is small, so it cannot be accessed from different package
	Color string // the first letter is capital, so it can be accessed from different package.
}

// The first letter of the method is capital, so it can be accessed from different packages
func (s square) GetSize() int {
	return s.size * s.size
}

// The first letter is small, so it cannot be accessed from different packages
func (s *square) setSize(newSize int) {
	s.size = newSize // Same as (*s).size = newSize
}

func print(sq square) {
	fmt.Println(sq)
}

type shapes struct {
	square
}

func main() {
	sq := square{
		size:  4,
		Color: "red",
	}
	fmt.Println(sq)
	sq.setSize(3)
	fmt.Println(sq)

	sh := shapes{
		square: sq,
	}

	fmt.Println(sh.GetSize()) // This would work and invoke the embedded struct's method
 
	fmt.Println(print(sh)) // This results in ERROR because Runtime polymorphism is not supported for composition
}

```

### Why go follows this approach than regular OO approach?
1. As we discussed, Go's approach is to keep things simple.
2. In the interface implementation, We don't even have to mention what interface we are implementing for a struct. Go compiler will implicitly determine the interface the struct is implementing. But that also makes it harder to find which struct has implemented the interface because we have to check the implementation part of the structs.
3. Go focusses more on composition (Embedded struct). This leads to more simpler and readable code. 

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

## Error types in go
in go, error type is defined as an interface
```go
type interface error {
	Error() string
}
```

## Conventions in go
For filenames, follow snake case or kebab case and for variable names follow camel case. Also dont forget that if the first letter is capital in the variable name, then it becomes visible to other packages also.

Folder and module names are lowercase (Not camel case or snake case). Which means the module name can be ```mygreatproject``` and not ```my_great_project```.

## Variadic arguments
Variadic arguments in Go allows us to pass an arbitary number of arguments of the same type to a function. This is particularly useful when we dont know beforehand how many arguments a function might take. 

```go
func foo(va ...int) {
	fmt.Println(va)
	fmt.Println(va[1])
}

func main() {
	foo(1, 2, 3)

}
```
Internally, the variadic arguments are implemented using slices. 

Importantly, a variadic parameter must be the last parameter in a function parameter list.

And a slice can be passed to a variadic function using ```...``` syntax. 
```go
func foo(va ...int) {
	fmt.Println(va)
	fmt.Println(va[1])
}

func main() {
	slice := []int{1, 2, 3}
	foo(slice...)

}
```

NOTE: The below code does not work if we directly assign an interface{} slice with a slice of any other datatype.
```go
var i []interface{}
i = make([]int, 3) // This will throw an error

// If we have a slice with any other datatype, then we have the slice to []interface{}
buf := []byte{1, 2}

is := make([]interface{}, 2)

for idx, val := range buf {
	is[idx] = interface{}(val)
}

fmt.Println(is)
```


## Defer

```defer``` is a special keyword in go that is used for ensuring that a function call is performed later in a program's execution, typically for purposes of cleanup. `defer` is often used in situations involving file, network operations or mutex locks, where we need to ensure some operations occur even if a function call exits early due to an error.

When we defer a function call, that call is placed on a stack and executed when the surrounding function exits. If there are multiple defer calls, then those functions are executed in LIFO manner.

### Examples
```go
func example() {
	fmt.Println("Start")
	defer fmt.Println("This is deferred")
	fmt.Println("End")
}
// Output:
// Start
// End
// This is deferred
```

```go
func readFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Do something with the file

	return nil
}
```
In case of function arguments, deferred function arguments are evaluated when the ```defer``` statement is called, not when the function is executed.
```go
func argumentEvaluation() {
	i := 0
	defer fmt.Println(i) // This will print 0
	i++
	fmt.Println(i) // This will print 1
}
// Output:
// 1
// 0
```

### Why defer statement calls are executed in LIFO manner and not in FIFO manner
It's a common pattern in programming to clean up resources in reverse order to ensure proper release of dependecies or related resources

```go
func writeFile(filename string) {
	file, _ := os.Open(filename)
	defer file.Close()

	lock := acquireFileLock(file)
	defer releaseFileLock(file)
}
```

Now in the above code, first releaseFileLock would happen and then file.close() would be executed. If this was in FIFO approach, then it may result in undefined behaviour because we are releasing the lock after closing the file. And if FIFO approach was implemented, then we would have difficulties in writing the correct order of defer statements (Ofcourse we can write the statements but it would be more difficult because we would have to scroll up to note what defer statement has to be written and then again scroll down and write it...)

## go.mod and go.sum
go.mod file is analogous to package.json in node projects. It maintains the dependencies (Both direct and indirect). Here indirect means the dependencies of the direct dependencies. (In case of node, the indirect dependencies are maintained by package-lock.json).

mod stands for module. A module in go is a collection of packages.

To create a module, use the following command
```go mod init playingcards```

And in go.sum, sum is a shortcut for checksum. It maintains the checksums of all the dependencies (direct and indirect) so that we know the dependencies are authentic.

## Tests
In go, unit tests are typically written along side the implementation files (This is different from the project structure that we follow in node or c++ projects where we typically have a different folder for tests.). And for integration tests, the convention is to make a different package for the tests.

<b><i>NOTE: The test files must end with _test.go</i></b>

```
package_directory/
        deck.go
		deck_test.go
		card.go
		card_test.go
		...

integration/
		integration_test.go
```

Example of a test file:

```go
package main

import "testing"

func TestNewDeck(t *testing.T) {
	d := newDeck()

	if len(d) != 48 {
		t.Errorf("Expected 48, but got %v", len(d))
	}
}
```

To run the tests, use ```go test``` command.

Why is a pointer passed in the testing function?

Because as we know, Go is a pass by value language, by passing pointer to the testing function, we can share the state across other tests. (State for eg, can be logs, how many tests have failed etc). 

And one more reason is that Testing.t object can be quite large, So passing it as pointer would be more memory effecient.


