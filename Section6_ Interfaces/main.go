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
