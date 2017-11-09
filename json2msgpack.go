package json2msgpack

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"sort"
)

// EncodeJSON ...
// generic MessagePack JSON serialiser
func EncodeJSON(bin []byte) []byte {
	var obj interface{}
	if err := json.Unmarshal(bin, &obj); err != nil {
		panic(fmt.Sprintf("Error unmarshalling json: '%v'", err))
	}
	return Encode(obj)
}

func isFloat(n float64) bool {
	return n != math.Floor(n)
}

// Encode ...
// encodes memory data to MessagePack
func Encode(v interface{}) (buf []byte) {
	numbers := map[reflect.Kind]bool{
		reflect.Int:     true,
		reflect.Int8:    true,
		reflect.Int16:   true,
		reflect.Int32:   true,
		reflect.Int64:   true,
		reflect.Uint:    true,
		reflect.Uint8:   true,
		reflect.Uint16:  true,
		reflect.Uint32:  true,
		reflect.Uint64:  true,
		reflect.Float32: true,
		reflect.Float64: true}

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
			binary.BigEndian.PutUint16(buf[1:], uint16(len))
			for i := 0; i < len; i++ {
				buf[i+3] = v.(string)[i]
			}
		} else {
			buf = make([]byte, 5+len)
			buf[0] = 0xdb
			binary.BigEndian.PutUint32(buf[1:], uint32(len))
			for i := 0; i < len; i++ {
				buf[i+5] = v.(string)[i]
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
			binary.BigEndian.PutUint16(buf[1:], uint16(length))
		} else {
			buf = make([]byte, 5)
			buf[0] = 0xdd
			binary.BigEndian.PutUint32(buf[1:], uint32(length))
		}

		var acc []byte
		for i, u := range vv {
			if i == 0 {
				acc = Encode(u)
			} else {
				acc = append(acc, Encode(u)...)
			}
		}
		buf = append(buf, acc...)
	default:
		if numbers[reflect.TypeOf(vv).Kind()] {
			buf = encodeNumber(number(v))
		} else {
			panic(fmt.Sprintf("%v (%T) : Parser unknown type", v, vv))
		}
	}
	return
}

func number(v interface{}) (num float64) {
	switch vv := v.(type) {
	case int:
		num = float64(v.(int))
	case int8:
		num = float64(v.(int8))
	case int16:
		num = float64(v.(int16))
	case int32:
		num = float64(v.(int32))
	case int64:
		num = float64(v.(int64))
	case float32:
		num = float64(v.(float32))
	case float64:
		num = v.(float64)
	case uint:
		num = float64(v.(uint))
	case uint8:
		num = float64(v.(uint8))
	case uint16:
		num = float64(v.(uint16))
	case uint32:
		num = float64(v.(uint32))
	case uint64:
		num = float64(v.(uint64))
	default:
		panic(fmt.Sprintf("%v (%T): unknown type casted to Number", v, vv))
	}
	return
}

func encodeNumber(n float64) (buf []byte) {
	if isFloat(n) {
		buf = make([]byte, 9)
		buf[0] = 0xcb
		mem := math.Float64bits(n)
		binary.BigEndian.PutUint64(buf[1:], uint64(mem))
	} else {
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
				binary.BigEndian.PutUint16(buf[1:], uint16(n))
			} else if n <= 0xffffffff {
				buf = make([]byte, 5)
				buf[0] = 0xce
				binary.BigEndian.PutUint32(buf[1:], uint32(n))
			} else if n <= 9007199254740991 {
				buf = make([]byte, 9)
				buf[0] = 0xcf
				binary.BigEndian.PutUint64(buf[1:], uint64(n))
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
				binary.BigEndian.PutUint16(buf[1:], uint16(n))
			} else if n > -214748365 {
				buf = make([]byte, 5)
				buf[0] = 0xd2
				binary.BigEndian.PutUint32(buf[1:], uint32(n))
			} else if n >= -9007199254740991 {
				buf = make([]byte, 9)
				buf[0] = 0xd3
				binary.BigEndian.PutUint64(buf[1:], uint64(n))
			}
		}
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
		binary.BigEndian.PutUint16(buf[1:], uint16(length))
	}

	// GOlang WTF : If you require a stable iteration order in maps
	// you must maintain a separate data structure that specifies that order.
	var keys []string
	for k := range obj {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		// fmt.Println("Key:", k, "Value:", obj[k])
		pair := append(Encode(k), Encode(obj[k])...)
		buf = append(buf, pair...)
	}

	// undefined order map below is commented in favor to ordering trick , see above
	// for k, v := range obj {
	// 	pair := append(encode(k), encode(v)...)
	// 	buf = append(buf, pair...)
	// }
	return
}
