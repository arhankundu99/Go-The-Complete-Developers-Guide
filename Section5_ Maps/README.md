# Section 5: Maps

```go
package main

import "fmt"

func main() {
	// Nil map declaration
	var colors map[string]string
	fmt.Println(colors == nil) // True

	colors2 := map[string]string{
		"black": "#000000",
		"white": "#ffffff",
	}
	fmt.Println(colors2)
	colors2["green"] = "#b012f3"
	fmt.Println(colors2)

	// Map declaration using make
	var colors3 = make(map[string]string)
	fmt.Println(colors3 == nil) // false. (Just like the case we saw with slices)

	// Delete a key from map
	delete(colors2, "green")
	fmt.Println(colors2)

	// Iterating over a map
	for key, value := range colors2 {
		fmt.Println(key, value)
	}
}
```

Note: Maps are passed as reference types in go. Which means if we change the map (For eg., add new keys), the change would be reflected in the original map also. This was not the case with slice and that is why we return a slice when we append some elements.

Maybe The underlying structure of maps contain a pointer which then contains a pointer to the underlying array. So that is why when array is reallocated in case of adding new keys inside a function, the changes are reflected in the original map also.

Why are maps designed this way? This design is chosen for efficiency and convinience, as maps are often used for quick data access and modifications.

So why are not slices designed this way? 
Consistency with array behaviour: This design maintains some consistency with how array operations behave. When we pass an array to a function, changes to the array elements are reflected back, but modifying array size is not permitted once the array is defined. So that is why appending to the slice returns a new array.

