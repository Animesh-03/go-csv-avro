package avro

import (
	"encoding/json"
	"fmt"
)

type SchemaType string

const (
	Null   SchemaType = "null"
	Bool   SchemaType = "boolean"
	Int    SchemaType = "int"
	Long   SchemaType = "long"
	Float  SchemaType = "float"
	Double SchemaType = "double"
	Bytes  SchemaType = "bytes"
	String SchemaType = "string"

	RecordSchemaType SchemaType = "record"
	Enum             SchemaType = "enum"
	Array            SchemaType = "array"
	Map              SchemaType = "map"
	Fixed            SchemaType = "fixed"
)

type MetaData struct {
	schema Schema
	codec  string
}

func NewMetadata(s Schema, codec string) MetaData {
	if codec == "" {
		codec = "null"
	}

	return MetaData{
		schema: s,
		codec:  codec,
	}
}

func (m *MetaData) Encode() []byte {
	bytes := make([]byte, 0)

	// Number of keys in metadata
	bytes = append(bytes, *EncodeVInt(2)...)
	bytes = append(bytes, m.schema.Encode()...)

	bytes = append(bytes, *EncodeString("avro.codec")...)
	bytes = append(bytes, *EncodeString(m.codec)...)

	return bytes
}

type Schema struct {
	Name      string     `json:"name,omitempty"`
	Type      SchemaType `json:"type,omitempty"`
	Doc       string     `json:"doc,omitempty"`
	Namespace string     `json:"namespace,omitempty"`
	Aliases   []string   `json:"aliases,omitempty"`
	Fields    []Field    `json:"fields,omitempty"`
}

func NewSchemaFromJSON(v string) *Schema {
	schema := Schema{}
	json.Unmarshal([]byte(v), &schema)

	return &schema
}

func (s *Schema) Encode() []byte {
	bytes := make([]byte, 0)

	bytes = append(bytes, *EncodeVInt(int64(len("avro.schema")))...)
	bytes = append(bytes, []byte("avro.schema")...)

	schemaString, err := json.Marshal(*s)
	if err != nil {
		panic(err)
	}

	schemaBytes := *EncodeString(string(schemaString))
	bytes = append(bytes, schemaBytes...)

	return bytes
}

type Field struct {
	Name      string     `json:"name,omitempty"`
	Namespace string     `json:"namespace,omitempty"`
	Doc       string     `json:"doc,omitempty"`
	Type      SchemaType `json:"type,omitempty"`
}

func DecodeFields(bytes *[]byte, fields []Field) Record {
	var record Record = make(Record)

	for _, f := range fields {
		switch f.Type {
		case Null:
			record[f.Name] = RecordVal{
				Bytes: *EncodeVInt(DecodeVInt(bytes)),
				Type:  f.Type,
			}
		case Bool:
			record[f.Name] = RecordVal{
				Bytes: *EncodeBool(DecodeBool(bytes)),
				Type:  f.Type,
			}
		case Int, Long:
			record[f.Name] = RecordVal{
				Bytes: *EncodeVInt(DecodeVInt(bytes)),
				Type:  f.Type,
			}
		case Float:
			record[f.Name] = RecordVal{
				Bytes: *EncodeFloat32(DecodeFloat32(bytes)),
				Type:  f.Type,
			}
		case Double:
			record[f.Name] = RecordVal{
				Bytes: *EncodeFloat64(DecodeFloat64(bytes)),
				Type:  f.Type,
			}
		case Bytes:
			numBytes := DecodeVInt(bytes)
			record[f.Name] = RecordVal{
				Bytes: (*bytes)[:numBytes],
				Type:  f.Type,
			}
			*bytes = (*bytes)[numBytes:]
		case String:
			record[f.Name] = RecordVal{
				Bytes: *EncodeString(DecodeString(bytes)),
				Type:  f.Type,
			}
		default:
			panic(fmt.Sprintf("%s not implemented yet", f.Type))
		}
	}

	return record
}
