package avro

func GetMetaData(bytes *[]byte) map[string]string {
	numRecords := DecodeVInt(bytes)

	metadata := make(map[string]string)

	for range numRecords {
		s := DecodeString(bytes)
		v := DecodeString(bytes)

		metadata[s] = v
	}

	return metadata
}
