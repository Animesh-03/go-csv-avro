package avro

import "fmt"

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

type Schema struct {
	Type      SchemaType `json:"type"`
	Name      string     `json:"name"`
	Namespace string     `json:"namespace"`
	Doc       string     `json:"doc"`
	Aliases   []string   `json:"aliases"`
	Fields    []Field    `json:"fields"`
}

type Field struct {
	Name      string     `json:"name"`
	Namespace string     `json:"namespace"`
	Doc       string     `json:"doc"`
	Type      SchemaType `json:"type"`
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
		// a := record[f.Name]
		// fmt.Printf("Decoding %s: %s\n", f.Name, a.GetRecordValString())
	}

	// record.PrintRecord()
	// fmt.Println((*bytes)[:10])

	return record
}

type Record map[string]RecordVal

type RecordVal struct {
	Bytes []byte
	Type  SchemaType
}

func (r *Record) PrintRecord() {
	i := 0
	mapLen := len(map[string]RecordVal(*r))
	for k, v := range map[string]RecordVal(*r) {
		if i < mapLen-1 {
			fmt.Printf("%s: %s, ", k, v.GetRecordValString())
		} else {
			fmt.Printf("%s: %s\n", k, v.GetRecordValString())
		}
		i++
	}
}

func (rv *RecordVal) GetRecordValString() string {
	var s string

	switch rv.Type {
	case Null:
		s = "null"
	case Int, Long:
		s = fmt.Sprintf("%d", DecodeVInt(&rv.Bytes))
	case Float:
		s = fmt.Sprintf("%f", DecodeFloat32(&rv.Bytes))
	case Double:
		s = fmt.Sprintf("%f", DecodeFloat64(&rv.Bytes))
	case Bytes:
		s = fmt.Sprintf("%v", rv.Bytes)
	case String:
		s = DecodeString(&rv.Bytes)
	case Bool:
		s = fmt.Sprintf("%t", DecodeBool(&rv.Bytes))
	default:
		panic(fmt.Sprintf("%s not implemented yet", rv.Type))
	}

	return s
}
