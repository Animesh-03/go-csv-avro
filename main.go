package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/animesh-03/go-csv-avro/avro"
)

func main() {
	fmt.Println("\n--------AVRO Encoding/Decoding Functions--------------")

	// Encode and Decode Int
	fmt.Println(avro.EncodeVInt(64))
	fmt.Println(avro.DecodeVInt(avro.EncodeVInt(372)))
	fmt.Println(avro.DecodeVInt(&[]byte{0xe2, 0x01}))

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

	fmt.Println("\n--------AVRO Reader Functions--------------")

	ar := avro.NewReader("twitter.avro")

	fmt.Printf("Schema:\n%+v\n", ar.GetSchema())
	fmt.Printf("Codec: %s\n", ar.GetCodec())
	ar.PrintRecords()

	fmt.Println("\n--------AVRO Writer Functions--------------")

	fields := []avro.Field{
		{
			Name: "field1",
			Type: avro.String,
		},
		{
			Name: "field2",
			Type: avro.Int,
		},
	}
	r := avro.NewRecord(fields, []string{"Hello World", "243"})
	fmt.Println(r.EncodeRecord(fields))

	schemaJson := `
	{
		"type": "record",
		"name": "some_schema",
		"namespace": "com.something.avro",
		"fields": [
			{
				"name":"field1",
				"type":"long"
			},
			{
				"name":"field2",
				"type":"string"
			}
		]
	}`

	aw := avro.NewWriter("test.avro", schemaJson, "")
	fmt.Printf("%+v\n", aw.GetSchema())
	fmt.Printf("Codec: %s\n", aw.GetCodec())

	aw.Records = []avro.Record{
		*avro.NewRecord(aw.GetSchema().Fields, []string{"1366154481", "Hello World"}),
		*avro.NewRecord(aw.GetSchema().Fields, []string{"1366154482", "Hello World Again"}),
	}
	aw.Write()

	fmt.Println("\n--------CSV To AVRO--------------")

	schemaJson = `
		{
			"type": "record",
			"name": "csv_schema",
			"fields": [
				{
					"name": "id",
					"type": "long"
				},
				{
					"name": "first_name",
					"type": "string"
				},
				{
					"name": "last_name",
					"type": "string"
				},
				{
					"name": "email",
					"type": "string"
				},
				{
					"name": "gender",
					"type": "string"
				},
				{
					"name": "ip_address",
					"type": "string"
				}
			]
		}
	`

	file, err := os.Open("data.csv")
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

	csvData := string(buf)
	csvRows := strings.Split(csvData, "\n")

	csvToAvroWriter := avro.NewWriter("data.avro", schemaJson, "null")
	fmt.Println(csvToAvroWriter.GetSchema())

	records := make([]avro.Record, 0)

	for i, row := range csvRows {
		if i == 0 {
			continue
		}
		attrs := strings.Split(row, ",")
		records = append(records, *avro.NewRecord(csvToAvroWriter.GetSchema().Fields, attrs))
	}

	csvToAvroWriter.Records = records
	csvToAvroWriter.Write()

	ar = avro.NewReader("data.avro")

	fmt.Printf("Schema:\n%+v\n", ar.GetSchema())
	fmt.Printf("Codec: %s\n", ar.GetCodec())
	ar.PrintRecords()

}
