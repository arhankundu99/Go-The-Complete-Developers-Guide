package main

import "fmt"

func printAddr(m map[string]string) {
	fmt.Printf("Address of map: %p\n", &m)
}

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

	fmt.Printf("Address of map: %p\n", &colors2)
	printAddr(colors2)
}
