package main

import "fmt"

type card string

func (c card) print() {
	fmt.Println(c)
}

func (c card) toString() string {
	return string(c)
}
