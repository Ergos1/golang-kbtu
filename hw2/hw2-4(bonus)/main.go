package main

import (
	"io"
	"os"
	"strings"
)

type rot13Reader struct {
	r io.Reader
}

func rot13(x byte) byte{
	if x >= 'A' && x <= 'Z'{
		x = byte((x-'A'+13)%26 + 'A')
	} else if x >= 'a' && x <= 'z'{
		x = byte((x-'a'+13)%26 + 'a')
	}
	return x
}

func (rot rot13Reader) Read(b []byte) (int, error){
	n, err := rot.r.Read(b)
	if err!=nil{
		return n, err
	}
	for i:=0; i < len(b); i++{
		b[i] = rot13(b[i])
	}
	return n, err
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}
