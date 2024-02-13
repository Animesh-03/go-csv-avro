package avro

import (
	"encoding/binary"
	"math"
)

// Return the parsed Variable Length Int and the number of bytes parsed
func DecodeVInt(bytes *[]byte) (int64, uint) {
	moreBit := false
	value := uint64(0)
	last := uint(0)

	for i, b := range *bytes {
		// Check if more bytes are to be read
		moreBit = (b & 128) != 0
		// Get the last 7 bits of current byte
		valBytes := b & 0x7f

		value |= uint64(valBytes) << (i * 7)

		if !moreBit {
			last = uint(i)
			break
		}
	}

	*bytes = (*bytes)[last+1:]

	return DecodeZigZag(uint64(value)), last + 1
}

// Return decoded int
func DecodeZigZag(val uint64) int64 {
	if uint64(val)>>63 == 1 {
		return -1 * int64((uint64(val)+1)/2)
	} else {
		return int64((uint64(val) + 1) / 2)
	}
}

// Return zig-zag encoded int
func EncodeZigZag(val int64) uint64 {
	if val < 0 {
		return uint64(math.Abs(float64(val))*2 - 1)
	} else {
		return uint64(val * 2)
	}
}

// Return byte array of a Variable Length Zig-Zag encoded int
func EncodeVInt(val int64) []byte {
	vIntBytes := make([]byte, 0)

	// Get the zig-zag encoded value
	valBytes := uint64(EncodeZigZag(val))

	for valBytes > 0 {
		// Get the last 7 bits
		lowerBits := valBytes & 0x7f
		valBytes = valBytes >> 7

		// Set the more bit if needed
		moreBit := uint64(0)
		if valBytes > 0 {
			moreBit = 1
		}

		// Append a byte with moreBit as MSB and the val as the remaining bits
		vIntBytes = append(vIntBytes, byte((moreBit<<7)|lowerBits))
	}

	return vIntBytes
}

// Returns an array of bytes for the encoded string
// The first few bytes represent the VInt encoded length of the string
func EncodeString(s string) []byte {
	lenBytes := EncodeVInt(int64(len(s)))
	strBytes := lenBytes
	strBytes = append(strBytes, []byte(s)...)

	return strBytes
}

// Return the parsed string and the number of bytes parsed
func DecodeString(bytes *[]byte) (string, uint) {
	length, offset := DecodeVInt(bytes)
	*bytes = (*bytes)[offset:]
	return string((*bytes)[offset : offset+uint(length)]), uint(length) + offset
}

// Return a byte with value 1 if true else 0
func EncodeBool(b bool) byte {
	if b {
		return byte(1)
	} else {
		return byte(0)
	}
}

// Return bool from byte
func DecodeBool(b byte) (bool, uint) {
	if b == byte(1) {
		return true, 1
	} else {
		return false, 1
	}
}

// Return a byte array from float32
func EncodeFloat32(val float32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, math.Float32bits(val))
	return buf
}

// Return float32 from byte array
func DecodeFloat32(bytes []byte) (float32, uint) {
	return math.Float32frombits(binary.LittleEndian.Uint32(bytes)), 4
}

// Return a byte array from float32
func EncodeFloat64(val float64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, math.Float64bits(val))
	return buf
}

// Return float64 from byte array
func DecodeFloat64(bytes []byte) (float64, uint) {
	return math.Float64frombits(binary.LittleEndian.Uint64(bytes)), 8
}
