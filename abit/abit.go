package abit

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sort"
)

type ABITObject struct {
	dataType uint8
	blob     *([]byte)
	text     *string
	tree     map[string](*ABITObject)
	boolean  bool
	integer  int64
	array    *ABITArray
}

type ABITArray struct {
	array [](*ABITObject)
}

type Null struct{}

func NewABITObject(document *[]byte) (*ABITObject, error) {
	if len(*document) > 0 {
		tree, _, err := decodeTree(document, 0, false)
		if err != nil {
			return nil, err
		}
		return &([]ABITObject{tree}[0]), nil
	} else {
		tree := ABITObject{
			dataType: 0b0110,
			tree:     map[string]*ABITObject{},
		}
		return &([]ABITObject{tree}[0]), nil
	}
}

func NewABITArray() *ABITArray {
	arr := ABITArray{}
	return &([]ABITArray{arr}[0])
}

func (t *ABITObject) Put(key string, value interface{}) error {
	// Must be tree type to put an object
	if t.dataType != 0b0110 {
		return fmt.Errorf("ABITObject is invalid type")
	}
	switch b := value.(type) {
	case Null:
		o := &ABITObject{
			dataType: 0b0000,
		}
		t.tree[key] = o
	case bool:
		o := &ABITObject{
			dataType: 0b0001,
			boolean:  b,
		}
		t.tree[key] = o
	case int64:
		o := &ABITObject{
			dataType: 0b0010,
			integer:  b,
		}
		t.tree[key] = o
	case []byte:
		o := &ABITObject{
			dataType: 0b0011,
			blob:     &b,
		}
		t.tree[key] = o
	case string:
		o := &ABITObject{
			dataType: 0b0100,
			text:     &b,
		}
		t.tree[key] = o
	case ABITArray:
		o := &ABITObject{
			dataType: 0b0101,
			array:    &b,
		}
		t.tree[key] = o
	case ABITObject:
		if b.dataType == 0b0110 {
			t.tree[key] = &b
		} else {
			return fmt.Errorf("ABITObject is invalid type")
		}
	default:
		return fmt.Errorf("unknown type")
	}
	return nil
}

func (a *ABITArray) Add(value interface{}) error {
	o := &ABITObject{}
	switch b := value.(type) {
	case Null:
		o.dataType = 0b0000
	case bool:
		o.dataType = 0b0001
		o.boolean = b
	case int64:
		o.dataType = 0b0010
		o.integer = b
	case []byte:
		o.dataType = 0b0011
		o.blob = &b
	case string:
		o.dataType = 0b0100
		o.text = &b
	case ABITArray:
		o.dataType = 0b0101
		o.array = &b
	case ABITObject:
		if b.dataType == 0b0110 {
			o = &b
		} else {
			return fmt.Errorf("ABITObject is not a valid type")
		}
	default:
		return fmt.Errorf("unsupported type")
	}
	a.array = append(a.array, o)
	return nil
}

func (t *ABITObject) Remove(key string) {
	delete(t.tree, key)
}

func (a *ABITArray) Remove(index int64) {
	ret := make([](*ABITObject), 0)
	ret = append(ret, (a.array)[:index]...)
	a.array = append(ret, (a.array)[index+1:]...)
}

func (t *ABITObject) get(key string) (interface{}, error) {
	// Must be tree type to get an object
	if t.dataType != 0b0110 {
		return 0, fmt.Errorf("ABITObject is not of type tree")
	}
	o := t.tree[key]
	switch o.dataType {
	case 0b0000:
		return Null{}, nil
	case 0b0001:
		return o.boolean, nil
	case 0b0010:
		return o.integer, nil
	case 0b0011:
		return o.blob, nil
	case 0b0100:
		return o.text, nil
	case 0b0101:
		return o.array, nil
	case 0b0110:
		return o, nil
	default:
		return 0, fmt.Errorf("object trying to be fetched is invalid")
	}
}

func (t *ABITObject) GetNull(key string) (Null, error) {
	obj, err := t.get(key)
	if err != nil {
		return Null{}, err
	}
	switch o := obj.(type) {
	case Null:
		return o, nil
	}
	return Null{}, fmt.Errorf("object trying to be fetched is not a null")
}

func (t *ABITObject) GetBool(key string) (bool, error) {
	obj, err := t.get(key)
	if err != nil {
		return false, err
	}
	switch o := obj.(type) {
	case bool:
		return o, nil
	}
	return false, fmt.Errorf("object trying to be fetched is not a boolean")
}

func (t *ABITObject) GetInteger(key string) (int64, error) {
	obj, err := t.get(key)
	if err != nil {
		return 0, err
	}
	switch o := obj.(type) {
	case int64:
		return o, nil
	}
	return 0, fmt.Errorf("object trying to be fetched is not an integer")
}

func (t *ABITObject) GetBlob(key string) (*[]byte, error) {
	obj, err := t.get(key)
	if err != nil {
		return nil, err
	}
	switch o := obj.(type) {
	case *[]byte:
		return o, nil
	}
	return nil, fmt.Errorf("object trying to be fetched is not a blob")
}

func (t *ABITObject) GetString(key string) (*string, error) {
	obj, err := t.get(key)
	if err != nil {
		return nil, err
	}
	switch o := obj.(type) {
	case *string:
		return o, nil
	}
	return nil, fmt.Errorf("object trying to be fetched is not a string")
}

func (t *ABITObject) GetArray(key string) (*ABITArray, error) {
	obj, err := t.get(key)
	if err != nil {
		return nil, err
	}
	switch o := obj.(type) {
	case *ABITArray:
		return o, nil
	}
	return nil, fmt.Errorf("object trying to be fetched is not an array")
}

func (t *ABITObject) GetTree(key string) (*ABITObject, error) {
	obj, err := t.get(key)
	if err != nil {
		return nil, err
	}
	switch o := obj.(type) {
	case *ABITObject:
		return o, nil
	}
	return nil, fmt.Errorf("object trying to be fetched is not an array")
}

func encodeKey(value string) (*[]byte, error) {
	keyBytes := []byte(value)
	if len(keyBytes) > 256 {
		return nil, fmt.Errorf("key too long")
	} else if len(keyBytes) < 1 {
		return nil, fmt.Errorf("key too short")
	}
	buf := make([]byte, 1+len(keyBytes))

	buf[0] = uint8(len(keyBytes) - 1)
	copy(buf[1:], keyBytes)

	return &buf, nil
}

func encodeNull() *[]byte {
	return &[]byte{0}
}

func encodeBoolean(value bool) *[]byte {
	if value {
		return &[]byte{0x11}
	} else {
		return &[]byte{0x01}
	}
}

func encodeInteger(value int64, type_n uint8) *[]byte {
	var byteCount uint8 = 0
	switch {
	case value >= -128 && value <= 127:
		byteCount = 1
	case value >= -32768 && value <= 32767:
		byteCount = 2
	case value >= -8388608 && value <= 8388607:
		byteCount = 3
	case value >= -2147483648 && value <= 2147483647:
		byteCount = 4
	case value >= -549755813888 && value <= 549755813887:
		byteCount = 5
	case value >= -140737488355328 && value <= 140737488355327:
		byteCount = 6
	case value >= -36028797018963968 && value <= 36028797018963967:
		byteCount = 7
	default:
		byteCount = 8
	}

	buf := make([]byte, 8)

	binary.LittleEndian.PutUint64(buf, uint64(value))

	out := append([]byte{((byteCount - 1) << 4) | (type_n & 0x0f)}, buf[:byteCount]...)

	return &out
}

func encodeBlob(value *[]byte, type_n uint8) *[]byte {
	var buffer bytes.Buffer
	buffer.Write(*encodeInteger(int64(len(*value)), type_n))
	buffer.Write(*value)
	out := buffer.Bytes()
	return &out
}

func encodeString(value *string) *[]byte {
	stringBytes := []byte(*value)
	return encodeBlob(&stringBytes, 0b0100)
}

func encodeArray(value *ABITArray) (*[]byte, error) {
	var buffer bytes.Buffer
	for _, obj := range value.array {
		switch obj.dataType {
		case 0b0000:
			buffer.Write(*encodeNull())
		case 0b0001:
			buffer.Write(*encodeBoolean(obj.boolean))
		case 0b0010:
			buffer.Write(*encodeInteger(obj.integer, 0b0010))
		case 0b0011:
			buffer.Write(*encodeBlob(obj.blob, 0b0011))
		case 0b0100:
			buffer.Write(*encodeString(obj.text))
		case 0b0101:
			p, err := encodeArray(obj.array)
			if err != nil {
				return nil, err
			}
			buffer.Write(*p)
		case 0b0110:
			p, err := encodeTree(obj, true)
			if err != nil {
				return nil, err
			}
			buffer.Write(*p)
		default:
			return nil, fmt.Errorf("object in array is of invalid type")
		}
	}
	p := buffer.Bytes()
	return encodeBlob(&p, 0b0101), nil
}

func encodeTree(value *ABITObject, nested bool) (*[]byte, error) {
	keys := make([]string, 0, len(value.tree))
	for k := range value.tree {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		if len(keys[i]) == len(keys[j]) {
			// If lengths are equal, sort lexicographically
			return keys[i] < keys[j]
		}
		// Otherwise, sort by length
		return len(keys[i]) < len(keys[j])
	})

	var buffer bytes.Buffer
	for _, key := range keys {
		p, err := encodeKey(key)
		if err != nil {
			return nil, err
		}
		buffer.Write(*p)
		obj := value.tree[key]
		switch obj.dataType {
		case 0b0000:
			buffer.Write(*encodeNull())
		case 0b0001:
			buffer.Write(*encodeBoolean(obj.boolean))
		case 0b0010:
			buffer.Write(*encodeInteger(obj.integer, 0b0010))
		case 0b0011:
			buffer.Write(*encodeBlob(obj.blob, 0b0011))
		case 0b0100:
			buffer.Write(*encodeString(obj.text))
		case 0b0101:
			p, err := encodeArray(obj.array)
			if err != nil {
				return nil, err
			}
			buffer.Write(*p)
		case 0b0110:
			p, err := encodeTree(obj, true)
			if err != nil {
				return nil, err
			}
			buffer.Write(*p)
		default:
			return nil, fmt.Errorf("object in array is of invalid type")
		}
	}
	p := buffer.Bytes()
	if nested {
		return encodeBlob(&p, 0b0110), nil
	} else {
		return &p, nil
	}
}

func (t *ABITObject) ToByteArray() ([]byte, error) {
	p, err := encodeTree(t, false)
	return *p, err
}

func decodeKey(blob *[]byte, offset int64) (string, int64, error) {
	keyLength := int64((*blob)[offset] + 1)
	return string((*blob)[offset : offset+int64(keyLength)]), offset + 1 + keyLength, nil
}

func decodeType(blob *[]byte, offset int64) uint8 {
	return (*blob)[offset] & 0x0f
}

func decodeNull(blob *[]byte, offset int64) (int64, error) {
	if (*blob)[offset] != 0x00 {
		return 0, fmt.Errorf("byte is not null")
	}
	return offset + 1, nil
}

func decodeBoolean(blob *[]byte, offset int64) (bool, int64, error) {
	switch (*blob)[offset] {
	case 0b00010001:
		return true, offset + 1, nil
	case 0b00000001:
		return false, offset + 1, nil
	}
	return false, 0, fmt.Errorf("byte is not boolean")
}

func decodeInteger(blob *[]byte, offset int64, maxSize int) (int64, int64, error) {
	intSize := ((*blob)[offset] >> 4) + 1
	if maxSize < int(intSize) {
		return 0, 0, fmt.Errorf("integer is too big")
	}

	extended := make([]byte, 8)
	copy(extended, (*blob)[offset+1:])

	// If the sign bit (most significant bit of the original byte slice) is set, perform sign-extension
	if (*blob)[offset+1+int64(intSize)-1]&0x80 != 0 {
		for i := intSize; i < 8; i++ {
			extended[i] = 0xFF
		}
	}

	// Convert to int64 by interpreting the extended slice as a little-endian 8-byte integer
	result := int64(extended[0]) |
		int64(extended[1])<<8 |
		int64(extended[2])<<16 |
		int64(extended[3])<<24 |
		int64(extended[4])<<32 |
		int64(extended[5])<<40 |
		int64(extended[6])<<48 |
		int64(extended[7])<<56

	return result, offset + 1 + int64(intSize), nil
}

func decodeBlob(blob *[]byte, offset int64) ([]byte, int64, error) {
	blobLength, offset, err := decodeInteger(blob, offset, 4)
	if err != nil {
		return nil, 0, err
	}
	if blobLength < 0 {
		return nil, 0, fmt.Errorf("negative length for blob")
	}
	var buf []byte = (*blob)[offset : offset+blobLength]
	return buf, offset + blobLength, nil
}

func decodeString(blob *[]byte, offset int64) (string, int64, error) {
	rawString, offset, err := decodeBlob(blob, offset)
	return string(rawString), offset, err
}

func decodeArray(blob *[]byte, offset int64) (ABITArray, int64, error) {
	arr := ABITArray{}
	arrBlob, offset, err := decodeBlob(blob, offset)
	if err != nil {
		return arr, 0, err
	}
	var index int64 = 0
	for int(index) < len(arrBlob) {
		switch decodeType(&arrBlob, index) {
		case 0b0000:
			index, err = decodeNull(&arrBlob, index)
			if err != nil {
				return arr, 0, err
			}
			arr.array = append(arr.array, &ABITObject{
				dataType: 0,
			})
		case 0b0001:
			var b bool
			b, index, err = decodeBoolean(&arrBlob, index)
			if err != nil {
				return arr, 0, err
			}
			arr.array = append(arr.array, &ABITObject{
				dataType: 1,
				boolean:  b,
			})
		case 0b0010:
			var b int64
			b, index, err = decodeInteger(&arrBlob, index, 8)
			if err != nil {
				return arr, 0, err
			}
			arr.array = append(arr.array, &ABITObject{
				dataType: 2,
				integer:  b,
			})
		case 0b0011:
			var b []byte
			b, index, err = decodeBlob(&arrBlob, index)
			if err != nil {
				return arr, 0, err
			}
			arr.array = append(arr.array, &ABITObject{
				dataType: 3,
				blob:     &([]([]byte){b}[0]),
			})
		case 0b0100:
			var b string
			b, index, err = decodeString(&arrBlob, index)
			if err != nil {
				return arr, 0, err
			}
			arr.array = append(arr.array, &ABITObject{
				dataType: 4,
				text:     &([]string{b}[0]),
			})
		case 0b0101:
			var b ABITArray
			b, index, err = decodeArray(&arrBlob, index)
			if err != nil {
				return arr, 0, err
			}
			arr.array = append(arr.array, &ABITObject{
				dataType: 5,
				array:    &([]ABITArray{b}[0]),
			})
		case 0b0110:
			var b ABITObject
			b, index, err = decodeTree(&arrBlob, index, true)
			if err != nil {
				return arr, 0, err
			}
			arr.array = append(arr.array, &([]ABITObject{b}[0]))
		default:
			return arr, 0, fmt.Errorf("invalid type")
		}
	}
	if int(index) > len(arrBlob) {
		return arr, 0, fmt.Errorf("corrupt array")
	}
	return arr, offset, nil
}

func keyCompare(a, b string) bool {
	if len(a) == len(b) {
		// If lengths are equal, sort lexicographically
		return a < b
	}
	// Otherwise, sort by length
	return len(a) < len(b)
}

func decodeTree(blob *[]byte, offset int64, nested bool) (ABITObject, int64, error) {
	tree := ABITObject{
		dataType: 6,
		tree:     map[string]*ABITObject{},
	}

	var treeBlob []byte
	var err error
	if nested {
		treeBlob, offset, err = decodeBlob(blob, offset)
		if err != nil {
			return tree, 0, err
		}
	} else {
		treeBlob = *blob
		offset = int64(len(*blob) - 1)
	}

	var key, lastKey string
	var index int64 = 0
	for int(index) < len(treeBlob) {
		key, index, err = decodeKey(&treeBlob, index)
		if err != nil {
			return tree, 0, err
		}

		if !keyCompare(lastKey, key) {
			return tree, 0, fmt.Errorf("invalid key order")
		}
		lastKey = key

		switch decodeType(&treeBlob, index) {
		case 0b0000:
			index, err = decodeNull(&treeBlob, index)
			if err != nil {
				return tree, 0, err
			}
			tree.tree[key] = &ABITObject{
				dataType: 0,
			}
		case 0b0001:
			var b bool
			b, index, err = decodeBoolean(&treeBlob, index)
			if err != nil {
				return tree, 0, err
			}
			tree.tree[key] = &ABITObject{
				dataType: 1,
				boolean:  b,
			}
		case 0b0010:
			var b int64
			b, index, err = decodeInteger(&treeBlob, index, 8)
			if err != nil {
				return tree, 0, err
			}
			tree.tree[key] = &ABITObject{
				dataType: 2,
				integer:  b,
			}
		case 0b0011:
			var b []byte
			b, index, err = decodeBlob(&treeBlob, index)
			if err != nil {
				return tree, 0, err
			}
			tree.tree[key] = &ABITObject{
				dataType: 3,
				blob:     &([]([]byte){b}[0]),
			}
		case 0b0100:
			var b string
			b, index, err = decodeString(&treeBlob, index)
			if err != nil {
				return tree, 0, err
			}
			tree.tree[key] = &ABITObject{
				dataType: 4,
				text:     &([]string{b}[0]),
			}
		case 0b0101:
			var b ABITArray
			b, index, err = decodeArray(&treeBlob, index)
			if err != nil {
				return tree, 0, err
			}
			tree.tree[key] = &ABITObject{
				dataType: 5,
				array:    &([]ABITArray{b}[0]),
			}
		case 0b0110:
			var b ABITObject
			b, index, err = decodeTree(&treeBlob, index, true)
			if err != nil {
				return tree, 0, err
			}
			tree.tree[key] = &([]ABITObject{b}[0])
		default:
			return tree, 0, fmt.Errorf("invalid type")
		}
	}
	if int(index) > len(treeBlob) {
		return tree, 0, fmt.Errorf("corrupt array")
	}
	return tree, offset, nil
}
