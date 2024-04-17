package main

import "fmt"

func main() {
	fmt.Println("Basics")
	Switcharoo(true)
	Switcharoo(1)
	Switcharoo("hello world")
	fmt.Println("------")
	p := Person{
		Name: "Ashley",
		Age:  30,
	}
	p.NameValue()
	fmt.Println(p.Name)
	p.NamePtr()
	fmt.Println(p.Name)
}

func Switcharoo(v interface{}) {
	switch v.(type) {
	case bool:
		fmt.Println("Bool")
	case int:
		fmt.Println("Int")
	case string:
		fmt.Println("String")
	}
}

type Person struct {
	Name string
	Age  int
}

func (p Person) NameValue() {
	p.Name = "Changed my value"
}

func (p *Person) NamePtr() {
	p.Name = "Changed my ptr"
}
