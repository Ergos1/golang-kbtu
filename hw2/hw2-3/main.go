package main

import (
	"golang.org/x/tour/reader"
)

type MyReader struct{}

// TODO: Add a Read([]byte) (int, error) method to MyReader.
func (r MyReader) Read(b []byte) (int, error) {
	for i:=0; i<len(b); i++{
		b[i] = 'A'
	}
	return len(b), nil
}

func main() {
	reader.Validate(MyReader{})
	// r := MyReader{}
	// b := make([]byte, 2)
	// for{
	// 	n, err := r.Read(b)
	// 	fmt.Printf("n = %v err = %v b = %v\n", n, err, b)
	// 	fmt.Printf("b[:n] = %q\n", b[:n])
	// }
}
