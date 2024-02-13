package avro

import "fmt"

func GetMetaData(bytes []byte) (map[string]string, uint) {
	numRecords, _ := DecodeVInt(&bytes)
	// bytes = bytes[offset:]

	fmt.Println(numRecords)

	metadata := make(map[string]string)

	for range numRecords {
		s, offset := DecodeString(&bytes)
		bytes = bytes[offset:]

		v, offset := DecodeString(&bytes)
		bytes = bytes[offset:]

		// fmt.Printf("%v	: %s\n", []byte(s), v)

		metadata[s] = v
	}

	fmt.Println(metadata)

	return metadata, 1
}
