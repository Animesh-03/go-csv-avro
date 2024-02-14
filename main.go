package main

import (
	"fmt"

	"github.com/animesh-03/go-csv-avro/avro"
)

func main() {
	// Encode and Decode Int
	fmt.Println(avro.EncodeVInt(64))
	fmt.Println(avro.DecodeVInt(avro.EncodeVInt(372)))
	fmt.Println(avro.DecodeVInt(&[]byte{0xc8, 0x01}))

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

	fmt.Println("\n----------------------")

	ar := avro.NewReader("twitter.avro")

	fmt.Printf("Schema:\n%+v\n", ar.GetSchema())
	fmt.Printf("Codec: %s\n", ar.GetCodec())
	ar.PrintRecords()
}
