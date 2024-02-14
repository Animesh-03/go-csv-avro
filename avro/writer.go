package avro

import (
	"bufio"
	"crypto/rand"
	"os"
)

type AvroWriter struct {
	path           string
	Metadata       MetaData
	Records        []Record
	syncFileMarker []byte
}

func NewWriter(path string, schemaJson string, codec string) *AvroWriter {
	schema := NewSchemaFromJSON(schemaJson)
	metadata := NewMetadata(*schema, codec)

	return &AvroWriter{
		path:     path,
		Metadata: metadata,
	}
}

func (aw *AvroWriter) Write() {
	file, err := os.OpenFile(aw.path, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	writer.WriteString("Obj")
	writer.WriteByte(1)
	writer.Flush()

	writer.Write(aw.Metadata.Encode())
	writer.Write([]byte{0})
	writer.Flush()

	aw.syncFileMarker = make([]byte, 16)
	rand.Read(aw.syncFileMarker)
	// aw.syncFileMarker = []byte{103, 199, 53, 41, 115, 239, 223, 148, 173, 211, 0, 126, 158, 235, 255, 174}

	writer.Write(aw.syncFileMarker)
	writer.Flush()

	writer.Write(*EncodeVInt(int64(len(aw.Records))))

	serializedRecords := make([]byte, 0)

	for _, record := range aw.Records {
		recordBytes := record.EncodeRecord(aw.Metadata.schema.Fields)
		serializedRecords = append(serializedRecords, recordBytes...)
	}

	writer.Write(*EncodeVInt(int64(len(serializedRecords))))
	writer.Write(serializedRecords)
	writer.Flush()

	writer.Write(aw.syncFileMarker)
	writer.Flush()
}

func (aw *AvroWriter) GetSchema() Schema {
	return aw.Metadata.schema
}

func (aw *AvroWriter) GetCodec() string {
	return aw.Metadata.codec
}
