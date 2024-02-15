# Installation and Running Instructions

The project requires `go1.22.0` or higher to run. It can be installed from [here](https://go.dev/doc/install).

To run:
```sh
go run .
```

# AVRO Encoding

## Zig-Zag Encoding

**All The Integers in AVRO are encoded using zig-zag encoding and then encoded using [Variable Length Integer Encoding](#variable-length-integer-encoding).**

The integers are encoded in the following format:

```
Value       Encoded Value
0               0
-1              1
1               2
-2              3
2               4
...
-64             127
64              128
...
```

The above can formulated as follows and is implemented [here](avro/encoding.go#43):

```
encoded_value = val*2 if val < 0 else val*2 -  1
```

The decoding formula is below and is implemented [here](avro/encoding.go#34)

```
value = encoded_value/2 if encoded_value%2 == 0 else -1*(encoded_value+1)/2
```

## Variable Length Integer Encoding

The [zig-zag encoded integers](#zig-zag-encoding) are then encoded using Variable Length Encoding.

This is a variable-length format for positive integers where the high-order bit of each byte indicates whether more bytes remain to be read. The low-order seven bits are appended as increasingly more significant bits in the resulting integer value. Thus values from zero to 127 may be stored in a single byte, values from 128 to 16,383 may be stored in two bytes, and so on.

```
Value	Byte 1	    Byte 2	    Byte 3
0	    00000000		
1	    00000001		
2	    00000010		
...			
127	    01111111		
128	    10000000	00000001	
129	    10000001	00000001	
130	    10000010	00000001	
...			
16,383	11111111	01111111	
16,384	10000000	10000000	00000001
16,385	10000001	10000000	00000001
...
```

This encoding is useful as it does not use excess bytes to represent a relatively small number. For example the int `372` would be represented by 4 bytes without any encoding but this can be represented in 2 bytes using variable length encoding i.e, `0xe8 0x05`.

The encoding and decoding functions are implemented [here](avro/encoding.go#43) and [here](avro/encoding.go#9) respectively.

## Primitive Types Encoding

### Strings And Bytes

Strings and bytes are fundamentally stored in the same way. A string is just an array of bytes. A string is encoded by providing the [VLE](#variable-length-integer-encoding) length of the string and then proceeded by that many bytes of the string.

For Example, The string `"Hello World"` is of length `11` which is encoded using [VLE](#variable-length-integer-encoding) as `0x16`. This then proceeded by the byte values of the individual [ASCII](https://upload.wikimedia.org/wikipedia/commons/1/1b/ASCII-Table-wide.svg) characters of the string i.e, `0x48 0x65 0x6c 0x6c 0x6f 0x20 0x57 0x6f 0x72 0x6c 0x64`.

```
(string_length) (string_bytes)
```

### Floats and Doubles

Floats and Doubles are stored in a similar way with the difference being the number of bytes needed for each. The bits for a float/double are encoded in the little endian format. A float takes up 4 bytes while a double takes up 8 bytes of data.

The float encoding and decoding are implemented [here](avro/encoding.go#116) and [here](avro/encoding.go#123)

The double encoding and decoding are implemented [here](avro/encoding.go#131) and [here](avro/encoding.go#138)

```
(little_endian_encoded_float)
```

### Records

The format of a record is defined in the schema.

A record is defined as follows:
```json
{
    "type": "record", 
    "name": "test",
    "fields" : [
        {"name": "a", "type": "long"},
        {"name": "b", "type": "string"}
    ]
}
```

A record is encoded by concatenating its encoded fields as per the schema defined.

The following record

```json
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
}
```

is encoded as `0x36 0x06 0x66 0x6f 0x6f` with `0x36` being the VLE encoded value of `27`, `0x06` being the VLE encoded length of `"foo"` and `0x66 0x6f 0x6f` being the ascii bytes of the string.

The encoding for a record is implemented [here](avro/record.go#75)s

### Blocks

Blocks store data in chunks which enables skipping reading data that is not needed.

A Block is encoded in the following format:

```
(VLE Number of Objects in the block) (VLE Size of all the objects in block) (Encoded Objects) (null byte representing end of block) (16 byte File Sync Marker)
```

The sync marker is randomly generated 16 bytes.

### Maps

Maps are encoded as a series of blocks. Each block consists of a long count value, followed by that many key/value pairs. A block with count zero indicates the end of the map. Each item is encoded per the map's value schema.

If a block's count is negative, then the count is followed immediately by a long block size, indicating the number of bytes in the block. The actual count in this case is the absolute value of the count written.

Decoding a map is implemented [here](avro/encoding.go#145)

## AVRO File Format

An avro file contains of 2 parts
- [Header](#header)
- [Data](#data)

### Header

The header is of the following format:

```
(magic_number) (metadata_map) (null_byte) (16 byte file sync marker)

magic_number -> "Obj" 0x01
metadata_map -> (num_keys) (key_value_pairs)
key_value_pair -> (encoded metadata_key) (encoded value_string)
metadata_key -> "avro.schema" | "avro.codec"
null_byte -> 0x00
```

### Data

The data block consists of the encoded records.

```
(num_records) (objects_size) (encoded_records) (16 byte file sync marker)

num_records -> number of records in block
objects_size -> VLE size in bytes of all the records in the block
```

# Example

Consider the following file:

```
od -t x1z -v test.avro 

0000000 4f 62 6a 01 04 16 61 76 72 6f 2e 73 63 68 65 6d  >Obj...avro.schem<
0000020 61 a8 02 7b 22 6e 61 6d 65 22 3a 22 73 6f 6d 65  >a..{"name":"some<
0000040 5f 73 63 68 65 6d 61 22 2c 22 74 79 70 65 22 3a  >_schema","type":<
0000060 22 72 65 63 6f 72 64 22 2c 22 6e 61 6d 65 73 70  >"record","namesp<
0000100 61 63 65 22 3a 22 63 6f 6d 2e 73 6f 6d 65 74 68  >ace":"com.someth<
0000120 69 6e 67 2e 61 76 72 6f 22 2c 22 66 69 65 6c 64  >ing.avro","field<
0000140 73 22 3a 5b 7b 22 6e 61 6d 65 22 3a 22 66 69 65  >s":[{"name":"fie<
0000160 6c 64 31 22 2c 22 74 79 70 65 22 3a 22 6c 6f 6e  >ld1","type":"lon<
0000200 67 22 7d 2c 7b 22 6e 61 6d 65 22 3a 22 66 69 65  >g"},{"name":"fie<
0000220 6c 64 32 22 2c 22 74 79 70 65 22 3a 22 73 74 72  >ld2","type":"str<
0000240 69 6e 67 22 7d 5d 7d 14 61 76 72 6f 2e 63 6f 64  >ing"}]}.avro.cod<
0000260 65 63 08 6e 75 6c 6c 00 01 75 8d 2d fa b7 8e 5f  >ec.null..u.-..._<
0000300 af 1f c4 8e d2 60 e5 b6 04 50 e2 f3 ee 96 0a 16  >.....`...P......<
0000320 48 65 6c 6c 6f 20 57 6f 72 6c 64 e4 f3 ee 96 0a  >Hello World.....<
0000340 22 48 65 6c 6c 6f 20 57 6f 72 6c 64 20 41 67 61  >"Hello World Aga<
0000360 69 6e 01 75 8d 2d fa b7 8e 5f af 1f c4 8e d2 60  >in.u.-..._.....`<
0000400 e5 b6 f5 d0 2b e4 09 6c 73 e5                    >....+..ls.<
```

## Header

The first 4 bytes of the file `4f 62 6a 01` is the magic number for avro i.e, "Obj1".

The next few bytes represent the metadata map of the file which contain the schema definition and the codec used.

The 5th byte i.e, `0x04` represents the VLE number of records in the map which is `2`.

The next byte `0x16` represents the VLE length of the key which is `11`. The next 11 bytes make up the key which is "avro.schema".

The next 2 bytes `0xa8 0x02` is VLE length of the value of `"avro.schema"` which is `148`. The next 148 bytes makes up the json schema of the records.

The schema json is:
```json
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
}
```

The next byte `0x14` represents the VLE length of next key which is `10`. The next 10 bytes make up the string `"avro.codec"`.

The next by `0x08` represents the VLE length of value of `"avro.codec"` which `4`. The next 4 bytes make up the string `"null"`.

The next byte is `0x00` which is the `null byte` and represents the end of the map.

The next 16 bytes are the randomly generated file sync marker
```
0x7b 0xe1 0x69 0x1e 0x4d 0xfa a80x 0x00 0x8b 0x12 0x44 0x61 0x6c 0xbc 0xa7 0xc7 
```
## Data

The next byte `0x04` is the VLE number of records in the block which `2`.

The next byte `ox50` is the VLE length of records in the block which `80`.

The next few bytes `0xe2 0xf3 0xee 0x96 0x0a` is the VLE long for `"field1"` as per the schema which is `1366154481`.

The next field is "field2" as per the schema which is a string. The next byte `0x16` represents the length of the string i.e, `11`. The next `11` bytes are the string `"Hello World"`. This makes the first record
```json
{"field1": "1366154481", "field2": "Hello World"}
```

Similarly the next record is 
```json
{"field1": "1366154482", "field2": "Hello World Again"}
```

The last `16` bytes is the file sync marker.