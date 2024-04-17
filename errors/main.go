package main

import (
	"errors"
	"fmt"
)

func main() {
	err := GetSentinalError()
	_, ok := err.(CustomError)
	fmt.Printf("Is CustomError via Type assertation: %v\n", ok)
	_, ok = err.(InnerError)
	fmt.Printf("Is InnerError via Type assertation: %v\n", ok)
	_, ok = err.(SomeOtherError)
	fmt.Printf("Is SomeOtherError via Type assertation: %v\n", ok)
	fmt.Println()

	fmt.Printf("Is Sentinal CustomError via Is: %v\n", errors.Is(err, SentinalCustomErr))
	fmt.Printf("Is Sentinal InnerError via Is: %v\n", errors.Is(err, SentinalInnerErr))
	fmt.Printf("Is Sentinal SomeOtherError via Is: %v\n", errors.Is(err, SentinalSomeOtherErr))
	fmt.Println()

	err = GetError()
	fmt.Printf("Is CustomError via Is %v\n", errors.Is(err, CustomError{}))
	fmt.Printf("Is InnerError via Is %v\n", errors.Is(err, InnerError{}))
	fmt.Printf("Is SomeOtherError via Is %v\n", errors.Is(err, SomeOtherError{}))
	fmt.Println()

	t := errors.As(err, &CustomError{})
	fmt.Printf("Is CustomError via As: %v\n", t)
	t = errors.As(err, &InnerError{})
	fmt.Printf("Is InnerError via As: %v\n", t)
	t = errors.As(err, &SomeOtherError{})
	fmt.Printf("Is SomeOtherError via As: %v\n", t)
	fmt.Println()
}

var SentinalInnerErr = InnerError{msg: "InnerError"}
var SentinalCustomErr = CustomError{err: SentinalInnerErr}
var SentinalSomeOtherErr = SomeOtherError{}

type SomeOtherError struct{}

func (e SomeOtherError) Error() string {
	return "SomeOtherError"
}

type CustomError struct {
	err error
}

func (e CustomError) Error() string {
	return fmt.Sprintf("CustomError: %v", e.err)
}

func (e CustomError) Unwrap() error {
	return e.err
}

type InnerError struct {
	msg string
}

func (e InnerError) Error() string {
	return e.msg
}

func GetSentinalError() error {
	return SentinalCustomErr
}

func GetError() error {
	return CustomError{InnerError{msg: "InnerError"}}
}
