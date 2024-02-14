package avro

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type AvroReader struct {
	path       string
	MetaData   MetaData
	buf        []byte
	syncMarker []byte
	Records    []Record
}

func NewReader(path string) *AvroReader {
	file, err := os.Open(path)
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

	ar := AvroReader{
		path:     path,
		MetaData: MetaData{},
		buf:      buf,
	}

	// Get MetaData
	buf = buf[4:]

	for k, v := range DecodeMap(&buf) {
		switch k {
		case "avro.schema":
			schema := Schema{}
			json.Unmarshal([]byte(v), &schema)
			ar.MetaData.schema = schema

		case "avro.codec":
			ar.MetaData.codec = v
		}
	}

	ar.syncMarker = buf[:16]
	buf = buf[16:]

	numRecords := DecodeVInt(&buf)

	ar.Records = make([]Record, 0)
	// Block Size
	DecodeVInt(&buf)

	for range numRecords {
		ar.Records = append(ar.Records, DecodeFields(&buf, ar.GetSchema().Fields))
	}

	return &ar
}

func (ar *AvroReader) GetMetadata() MetaData {
	return ar.MetaData
}

func (ar *AvroReader) GetSchema() Schema {
	return ar.MetaData.schema
}

func (ar *AvroReader) GetCodec() string {
	return ar.MetaData.codec
}

func (ar *AvroReader) PrintRecords() {
	for i, record := range ar.Records {
		fmt.Printf("Record %d: ", i)
		record.PrintRecord()
	}
}
