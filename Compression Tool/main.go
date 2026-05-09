package main

import (
	"compression_tool/compressor"
)

func main() {
	err := compressor.Encode("./test.txt", "./test.enc")
	if err != nil {
		panic(err)
	}
	err = compressor.Decode("./test.enc", "decoded.txt")
	if err != nil {
		panic(err)
	}
}
