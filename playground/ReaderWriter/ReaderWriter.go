package readerwriter

import "fmt"

type MyReader interface {
	Read(p []byte) (n int, err error)
}

type MyWriter interface {
	Write(p []byte) (n int, err error)
}

type MyReaderWriter interface {
	MyReader
	MyWriter
}

type ReaderWriter struct{}

func (rw ReaderWriter) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (rw ReaderWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}

func ReadAndWrite(rw MyReaderWriter) {
	rw.Read([]byte{0})
	rw.Write([]byte{0})
	fmt.Println("Done!")
}
