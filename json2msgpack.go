package json2msgpack

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
)

// EncodeJSON ...
// generic MessagePack JSON serialiser
func EncodeJSON(bin []byte) []byte {
	var obj interface{}
	if err := json.Unmarshal(bin, &obj); err != nil {
		panic(fmt.Sprintf("Error unmarshalling json: '%v'", err))
	}
	fmt.Printf("input: %v", obj)
	return encodeObj(obj.(map[string]interface{}))
}

func isFloat(n float64) bool {
	return n != math.Floor(n)
}

func encode(v interface{}) (buf []byte) {
	switch vv := v.(type) {
	case string:
		len := len(v.(string))
		if len < 32 {
			buf = make([]byte, 1+len)
			buf[0] = 0xa0 | byte(len)
			for i := 0; i < len; i++ {
				buf[i+1] = v.(string)[i]
			}
		} else if len <= 0xff {
			// str8, but only when not in compatibility mode
			buf = make([]byte, 2+len)
			buf[0] = 0xd9
			buf[1] = byte(len)
			for i := 0; i < len; i++ {
				buf[i+2] = v.(string)[i]
			}
		} else if len <= 0xffff {
			buf = make([]byte, 3+len)
			buf[0] = 0xda
			binary.LittleEndian.PutUint16(buf[1:], uint16(len))
			for i := 0; i < len; i++ {
				buf[i+3] = v.(string)[i]
			}
		} else {
			buf = make([]byte, 5+len)
			buf[0] = 0xdb
			binary.LittleEndian.PutUint32(buf[1:], uint32(len))
			for i := 0; i < len; i++ {
				buf[i+5] = v.(string)[i]
			}
		}
	case float64:
		if isFloat(v.(float64)) {
			buf = make([]byte, 9)
			buf[0] = 0xcb
			n := math.Float64bits(v.(float64))
			binary.LittleEndian.PutUint64(buf[1:], uint64(n))
		} else {
			n := v.(float64)
			if n >= 0 {
				if n < 128 {
					buf = make([]byte, 1)
					buf[0] = byte(n)
				} else if n < 256 {
					buf = make([]byte, 2)
					buf[0] = 0xcc
					buf[1] = byte(n)
				} else if n < 65536 {
					buf = make([]byte, 3)
					buf[0] = 0xcd
					binary.LittleEndian.PutUint16(buf[1:], uint16(n))
				} else if n <= 0xffffffff {
					buf = make([]byte, 5)
					buf[0] = 0xce
					binary.LittleEndian.PutUint32(buf[1:], uint32(n))
				} else if n <= 9007199254740991 {
					buf = make([]byte, 9)
					buf[0] = 0xcf
					binary.LittleEndian.PutUint64(buf[1:], uint64(n))
				}
			} else {
				if n >= -32 {
					buf = make([]byte, 1)
					buf[0] = byte(0x100 + n)
				} else if n >= -128 {
					buf = make([]byte, 2)
					buf[0] = 0xd0
					buf[1] = byte(n)
				} else if n >= -32768 {
					buf = make([]byte, 3)
					buf[0] = 0xd1
					binary.LittleEndian.PutUint16(buf[1:], uint16(n))
				} else if n > -214748365 {
					buf = make([]byte, 5)
					buf[0] = 0xd2
					binary.LittleEndian.PutUint32(buf[1:], uint32(n))
				} else if n >= -9007199254740991 {
					buf = make([]byte, 9)
					buf[0] = 0xd3
					binary.LittleEndian.PutUint64(buf[1:], uint64(n))
				}
			}
		}
	case nil:
		buf = make([]byte, 1)
		buf[0] = byte(0xc0)
	case bool:
		codeBool := 0xc2
		if v.(bool) {
			codeBool = 0xc3
		}
		buf = make([]byte, 1)
		buf[0] = byte(codeBool)
	case map[string]interface{}:
		buf = encodeObj(vv)
	case []interface{}:
		length := len(vv)
		if length < 16 {
			buf = make([]byte, 1)
			buf[0] = byte(0x90 | length)
		} else if length < 65536 {
			buf = make([]byte, 3)
			buf[0] = 0xdc
			binary.LittleEndian.PutUint16(buf[1:], uint16(length))
		} else {
			buf = make([]byte, 5)
			buf[0] = 0xdd
			binary.LittleEndian.PutUint32(buf[1:], uint32(length))
		}

		var acc []byte
		for i, u := range vv {
			if i == 0 {
				acc = encode(u)
			} else {
				acc = append(acc, encode(u)...)
			}
		}
		buf = append(buf, acc...)
	default:
		panic(fmt.Sprintf("%v (%T) : Parser unknown type", v, vv))
	}
	return
}

func encodeObj(obj map[string]interface{}) (buf []byte) {
	length := len(obj)

	if length < 16 {
		buf = make([]byte, 1)
		buf[0] = byte(0x80 | length)
	} else {
		buf = make([]byte, 3)
		buf[0] = byte(0xde)
		binary.LittleEndian.PutUint16(buf[1:], uint16(length))
	}

	for k, v := range obj {
		pair := append(encode(k), encode(v)...)
		buf = append(buf, pair...)
	}
	return
}
