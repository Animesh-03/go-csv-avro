package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/animesh-03/go-csv-avro/avro"
)

func main() {
	// Encode and Decode Int
	fmt.Println(avro.EncodeVInt(64))
	fmt.Println(avro.DecodeVInt(&avro.EncodeVInt(372)[:]))
	fmt.Println(avro.DecodeVInt(&[]byte{0x46}))

	// Encode and Decode Strings
	fmt.Println(avro.EncodeString("Hello World"))
	fmt.Println(avro.DecodeString(avro.EncodeString("Hello World")))

	// Encode and Decode booleans
	fmt.Println(avro.EncodeBool(true))
	fmt.Println(avro.DecodeBool(avro.EncodeBool(false)))

	// Encode and Decode Floats
	fmt.Println(avro.EncodeFloat32(float32(3.456)))
	fmt.Println(avro.DecodeFloat32(avro.EncodeFloat32(float32(3.456789123456))))
	fmt.Println(avro.EncodeFloat64(float64(3.456789123456)))
	fmt.Println(avro.DecodeFloat64(avro.EncodeFloat64(float64(3.456789123456))))

	fmt.Println("\n----------------------\n")

	file, err := os.Open("twitter.avro")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fileStats, err := file.Stat()
	if err != nil {
		panic(err)
	}

	buf := make([]byte, fileStats.Size())

	_, err = bufio.NewReader(file).Read(buf)
	if err != nil {
		panic(err)
	}

	avro.GetMetaData(buf[4:])

}
