package main

import (
	"errors"
	"fmt"

	dostuff "github.com/aes421/GoPlayground/playground/dostuff"
	"github.com/aes421/GoPlayground/playground/imports"
	readerwriter "github.com/aes421/GoPlayground/playground/readerwriter"
)

func main() {
	fmt.Println("Hello, World!")
	imports.Import1()
	dostuff.Import3()
	err := fmt.Errorf("Error: %v", CustomError{})
	fmt.Println(errors.Is(err, CustomError{}))

	greeter := PrintGreeting
	greeter("Ashley")

	p := &Person{
		Name: "Ashley",
		Age:  30,
	}
	fmt.Println(p.String())

	i := MyInt(5)
	fmt.Println(i.Double())

	readerwriter.ReadAndWrite(readerwriter.ReaderWriter{})
}

type CustomError struct{}

func (e CustomError) Error() string {
	return "Custom Error"
}

func PrintGreeting(name string) {
	fmt.Printf("Hello, %s!\n", name)
}

type Person struct {
	Name string
	Age  int
}

func (p *Person) String() string {
	return fmt.Sprintf("%s is %d years old", p.Name, p.Age)
}

type MyInt int

func (i MyInt) Double() MyInt {
	return i * 2
}
