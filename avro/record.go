package avro

import (
	"fmt"
	"strconv"
)

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

func NewRecord(fields []Field, attrs []string) *Record {
	r := make(Record)
	for i, field := range fields {
		rv := RecordVal{
			Type: field.Type,
		}

		switch field.Type {
		case Null:
			rv.Bytes = *EncodeVInt(0)
		case Int, Long:
			val, err := strconv.ParseInt(attrs[i], 10, 64)
			if err != nil {
				panic(err)
			}
			rv.Bytes = *EncodeVInt(val)
		case Float:
			val, err := strconv.ParseFloat(attrs[i], 32)
			if err != nil {
				panic(err)
			}
			rv.Bytes = *EncodeFloat32(float32(val))
		case Double:
			val, err := strconv.ParseFloat(attrs[i], 64)
			if err != nil {
				panic(err)
			}
			rv.Bytes = *EncodeFloat64(val)
		case Bytes, String:
			byteLen := EncodeVInt(int64(len(attrs[i])))
			rv.Bytes = append(*byteLen, []byte(attrs[i])...)
		case Bool:
			val, err := strconv.ParseBool(attrs[i])
			if err != nil {
				panic(err)
			}
			rv.Bytes = *EncodeBool(val)
		default:
			panic(fmt.Sprintf("%s not implemented yet", rv.Type))
		}

		r[field.Name] = rv
	}

	return &r
}

func (r *Record) EncodeRecord(fields []Field) []byte {
	bytes := make([]byte, 0)

	for _, field := range fields {
		bytes = append(bytes, (*r)[field.Name].Bytes...)
	}

	return bytes
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
