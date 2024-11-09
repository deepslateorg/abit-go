// Package abit implements abit documents into go.
package abit

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"sort"
)

// ABITObject is used to store ABIT trees.
//
//	// Initialize an empty ABITObject
//	tree, _ := abit.NewABITObject([]byte{})
type ABITObject struct {
	dataType uint8
	blob     *([]byte)
	text     *string
	tree     map[string](*ABITObject)
	boolean  bool
	integer  int64
	array    *ABITArray
}

// ABITArray is used to store ABIT arrays.
//
//	// Initialize an empty ABITArray
//	arr := abit.NewABITArray()
type ABITArray struct {
	array [](*ABITObject)
}

// Null is a helper object to represent null values in abit.
type Null struct{}

// ABITLexicon stores a schema to see if a given ABITObject matches the schema.
type ABITLexicon struct {
	lexicon ABITObject
}

// NewABITObject Creates an ABIT object from a binary ABIT document.
//
// error is nil on success or an error if the document is invalid.
//
// # Example 1
//
//	// Initialize an empty ABITObject
//	tree, _ := abit.NewABITObject([]byte{})
//
// # Example 2
//
//	// Initialize an ABITObject from a byte slice of unknown validity
//	tree, err := abit.NewABITObject(doc)
//	if err == nil {
//		// Code to handle a valid ABIT Document here
//	}
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

// NewABITArray Initializes and returns an empty ABIT array.
//
// # Example
//
//	// Initialize ABITArray
//	arr := abit.NewABITArray()
func NewABITArray() *ABITArray {
	arr := ABITArray{}
	return &([]ABITArray{arr}[0])
}

// Put adds a value to the corresponding key in the ABIT object.
//
// # Requirements
//   - key must be less than or equal to 256 bytes when encoded with UTF-8, but also more than or equal to 1 byte.
//   - value can be of types: abit.Null, bool, int64, []byte, string, ABITArray, ABITObject
func (t *ABITObject) Put(key string, value interface{}) {
	// Must be tree type to put an object
	if t.dataType != 0b0110 {
		panic("ABITObject is invalid type")
	}
	if len([]byte(key)) > 256 || 0 >= len([]byte(key)) {
		panic("key too long")
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
			panic("ABITObject is invalid type")
		}
	default:
		panic("unknown type")
	}
}

// Add adds a value to the ABITArray.
//
// # Requirements
//   - Value can be of types: abit.Null, bool, int64, []byte, string, ABITArray, ABITObject
func (a *ABITArray) Add(value interface{}) {
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
			panic("ABITObject is not a valid type")
		}
	default:
		panic("unsupported type")
	}
	a.array = append(a.array, o)
}

// Keys gets all the keys in a tree.
//
// # Example
//
//	keys:= tree.Keys()
//	for i, key := range keys {
//		// Code iterating through the tree here.
//	}
func (t *ABITObject) Keys() []string {
	if t.dataType != 0b0110 {
		panic("the ABITObject is not of correct type")
	}
	keys := make([]string, 0, len(t.tree))
	for k := range t.tree {
		keys = append(keys, k)
	}
	return keys
}

// Length gets the number of objects in the array.
//
// # Example
//
//	for i := 0; i < arr.Length(); i++ {
//		// Code iterating through the array here.
//	}
func (a *ABITArray) Length() int {
	return len(a.array)
}

// Remove the key and its associated value from the ABITObject.
func (t *ABITObject) Remove(key string) {
	delete(t.tree, key)
}

// Remove the value at index from the ABITArray.
//
// # Requirements
//   - 0 <= index < length of array
func (a *ABITArray) Remove(index int64) {
	if index < 0 || int(index) >= len(a.array) {
		panic("index out of bounds")
	}
	ret := make([](*ABITObject), 0)
	ret = append(ret, (a.array)[:index]...)
	a.array = append(ret, (a.array)[index+1:]...)
}

func (a *ABITArray) get(index int64) interface{} {
	o := a.array[index]
	switch o.dataType {
	case 0b0000:
		return Null{}
	case 0b0001:
		return o.boolean
	case 0b0010:
		return o.integer
	case 0b0011:
		return o.blob
	case 0b0100:
		return o.text
	case 0b0101:
		return o.array
	case 0b0110:
		return o
	default:
		panic("object trying to be fetched is invalid")
	}
}

func (t *ABITObject) get(key string) interface{} {
	// Must be tree type to get an object
	if t.dataType != 0b0110 {
		panic("ABITObject is not of type tree")
	}
	o := t.tree[key]
	switch o.dataType {
	case 0b0000:
		return Null{}
	case 0b0001:
		return o.boolean
	case 0b0010:
		return o.integer
	case 0b0011:
		return o.blob
	case 0b0100:
		return o.text
	case 0b0101:
		return o.array
	case 0b0110:
		return o
	default:
		panic("object trying to be fetched is invalid")
	}
}

// GetNull fetches abit.Null assosiated with key.
//
// # Requirements
//   - value associated with key is null
func (t *ABITObject) GetNull(key string) Null {
	obj := t.get(key)
	switch o := obj.(type) {
	case Null:
		return o
	}
	panic("object trying to be fetched is not a null")
}

// GetBool fetches bool assosiated with key.
//
// # Requirements
//   - value associated with key is a boolean
func (t *ABITObject) GetBool(key string) bool {
	obj := t.get(key)
	switch o := obj.(type) {
	case bool:
		return o
	}
	panic("object trying to be fetched is not a boolean")
}

// GetInteger fetches integer assosiated with key.
//
// # Requirements
//   - value associated with key is an integer
func (t *ABITObject) GetInteger(key string) int64 {
	obj := t.get(key)
	switch o := obj.(type) {
	case int64:
		return o
	}
	panic("object trying to be fetched is not an integer")
}

// GetBlob fetches blob assosiated with key.
//
// # Requirements
//   - value associated with key is a blob
func (t *ABITObject) GetBlob(key string) *[]byte {
	obj := t.get(key)
	switch o := obj.(type) {
	case *[]byte:
		return o
	}
	panic("object trying to be fetched is not a blob")
}

// GetString fetches string assosiated with key.
//
// # Requirements
//   - value associated with key is a string
func (t *ABITObject) GetString(key string) *string {
	obj := t.get(key)
	switch o := obj.(type) {
	case *string:
		return o
	}
	panic("object trying to be fetched is not a string")
}

// GetArray fetches array assosiated with key.
//
// # Requirements
//   - value associated with key is an array
func (t *ABITObject) GetArray(key string) *ABITArray {
	obj := t.get(key)
	switch o := obj.(type) {
	case *ABITArray:
		return o
	}
	panic("object trying to be fetched is not an array")
}

// GetTree fetches tree assosiated with key.
//
// # Requirements
//   - value associated with key is a tree
func (t *ABITObject) GetTree(key string) *ABITObject {
	obj := t.get(key)
	switch o := obj.(type) {
	case *ABITObject:
		return o
	}
	panic("object trying to be fetched is not a tree")
}

// GetNull fetches abit.Null at index.
//
// # Requirements
//   - value at index is null
func (a *ABITArray) GetNull(index int64) Null {
	obj := a.get(index)
	switch o := obj.(type) {
	case Null:
		return o
	}
	panic("object trying to be fetched is not a null")
}

// GetBool fetches bool at index.
//
// # Requirements
//   - value at index is a boolean
func (a *ABITArray) GetBool(index int64) bool {
	obj := a.get(index)
	switch o := obj.(type) {
	case bool:
		return o
	}
	panic("object trying to be fetched is not a boolean")
}

// GetInteger fetches integer at index.
//
// # Requirements
//   - value at index is an integer
func (a *ABITArray) GetInteger(index int64) int64 {
	obj := a.get(index)
	switch o := obj.(type) {
	case int64:
		return o
	}
	panic("object trying to be fetched is not an integer")
}

// GetBlob fetches blob at index.
//
// # Requirements
//   - value at index is a blob
func (a *ABITArray) GetBlob(index int64) *[]byte {
	obj := a.get(index)
	switch o := obj.(type) {
	case *[]byte:
		return o
	}
	panic("object trying to be fetched is not a blob")
}

// GetString fetches string at index.
//
// # Requirements
//   - value at index is a string
func (a *ABITArray) GetString(index int64) *string {
	obj := a.get(index)
	switch o := obj.(type) {
	case *string:
		return o
	}
	panic("object trying to be fetched is not a string")
}

// GetArray fetches array at index.
//
// # Requirements
//   - value at index is an array
func (a *ABITArray) GetArray(index int64) *ABITArray {
	obj := a.get(index)
	switch o := obj.(type) {
	case *ABITArray:
		return o
	}
	panic("object trying to be fetched is not an array")
}

// GetTree fetches tree at index.
//
// # Requirements
//   - value at index is a tree
func (a *ABITArray) GetTree(index int64) *ABITObject {
	obj := a.get(index)
	switch o := obj.(type) {
	case *ABITObject:
		return o
	}
	panic("object trying to be fetched is not a tree")
}

func encodeKey(value string) *[]byte {
	keyBytes := []byte(value)
	if len(keyBytes) > 256 {
		panic("key too long")
	} else if len(keyBytes) < 1 {
		panic("key too short")
	}
	buf := make([]byte, 1+len(keyBytes))

	buf[0] = uint8(len(keyBytes) - 1)
	copy(buf[1:], keyBytes)

	return &buf
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

func encodeArray(value *ABITArray) *[]byte {
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
			p := encodeArray(obj.array)
			buffer.Write(*p)
		case 0b0110:
			p := encodeTree(obj, true)
			buffer.Write(*p)
		default:
			panic("object in array is of invalid type")
		}
	}
	p := buffer.Bytes()
	return encodeBlob(&p, 0b0101)
}

func encodeTree(value *ABITObject, nested bool) *[]byte {
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
		p := encodeKey(key)
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
			p := encodeArray(obj.array)
			buffer.Write(*p)
		case 0b0110:
			p := encodeTree(obj, true)
			buffer.Write(*p)
		default:
			panic("object in array is of invalid type")
		}
	}
	p := buffer.Bytes()
	if nested {
		return encodeBlob(&p, 0b0110)
	} else {
		return &p
	}
}

// ToByteArray converts the ABITObject to a binary document in abit format.
func (t *ABITObject) ToByteArray() []byte {
	p := encodeTree(t, false)
	return *p
}

func decodeKey(blob *[]byte, offset int64) (string, int64, error) {
	if offset < 0 || int(offset) >= len(*blob) {
		return "", 0, fmt.Errorf("key out of bounds")
	}
	keyLength := int64((*blob)[offset]) + 1
	if int(offset+1+keyLength) > len(*blob) {
		return "", 0, fmt.Errorf("key out of bounds")
	}
	return string((*blob)[offset+1 : offset+1+keyLength]), offset + 1 + keyLength, nil
}

func decodeType(blob *[]byte, offset int64) (uint8, error) {
	if offset < 0 || int(offset) >= len(*blob) {
		return 0, fmt.Errorf("type exceeds blob")
	}
	return (*blob)[offset] & 0x0f, nil
}

func decodeNull(blob *[]byte, offset int64) (int64, error) {
	if offset < 0 || int(offset) >= len(*blob) {
		return 0, fmt.Errorf("null exceeds blob")
	}
	if (*blob)[offset] != 0x00 {
		return 0, fmt.Errorf("byte is not null")
	}
	return offset + 1, nil
}

func decodeBoolean(blob *[]byte, offset int64) (bool, int64, error) {
	if offset < 0 || int(offset) >= len(*blob) {
		return false, 0, fmt.Errorf("bool exceeds blob")
	}
	switch (*blob)[offset] {
	case 0b00010001:
		return true, offset + 1, nil
	case 0b00000001:
		return false, offset + 1, nil
	}
	return false, 0, fmt.Errorf("byte is not boolean")
}

func decodeInteger(blob *[]byte, offset int64, maxSize int) (int64, int64, error) {
	if offset < 0 || int(offset) >= len(*blob) {
		return 0, 0, fmt.Errorf("integer exceeds blob")
	}
	intSize := ((*blob)[offset] >> 4) + 1
	if maxSize < int(intSize) {
		return 0, 0, fmt.Errorf("integer is too big at %d", offset)
	}
	if int(offset+1+int64(intSize)) > len(*blob) {
		return 0, 0, fmt.Errorf("integer is out of bounds")
	}

	extended := make([]byte, 8)
	copy(extended, (*blob)[offset+1:offset+1+int64(intSize)])

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
	if len(*blob) < int(offset+blobLength) {
		return nil, 0, fmt.Errorf("length for blob exceeds the blob")
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
		typ, err := decodeType(&arrBlob, index)
		if err != nil {
			return arr, 0, err
		}
		switch typ {
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

	var err error
	var index int64 = 0
	if nested {
		index = offset
		var treeSize int64
		treeSize, index, err = decodeInteger(blob, offset, 4)
		offset = index + treeSize
		if err != nil {
			return tree, 0, err
		}
	} else {
		offset = int64(len(*blob))
	}

	var key, lastKey string = "", ""
	for index < offset {
		key, index, err = decodeKey(blob, index)
		if err != nil {
			return tree, 0, err
		}

		if !keyCompare(lastKey, key) {
			return tree, 0, fmt.Errorf("invalid key order: (%d)->(%d), %s -> %s", len(lastKey), len(key), lastKey, key)
		}
		lastKey = key

		typ, err := decodeType(blob, index)
		if err != nil {
			return tree, 0, err
		}
		switch typ {
		case 0b0000:
			index, err = decodeNull(blob, index)
			if err != nil {
				return tree, 0, err
			}
			tree.tree[key] = &ABITObject{
				dataType: 0,
			}
		case 0b0001:
			var b bool
			b, index, err = decodeBoolean(blob, index)
			if err != nil {
				return tree, 0, err
			}
			tree.tree[key] = &ABITObject{
				dataType: 1,
				boolean:  b,
			}
		case 0b0010:
			var b int64
			b, index, err = decodeInteger(blob, index, 8)
			if err != nil {
				return tree, 0, err
			}
			tree.tree[key] = &ABITObject{
				dataType: 2,
				integer:  b,
			}
		case 0b0011:
			var b []byte
			b, index, err = decodeBlob(blob, index)
			if err != nil {
				return tree, 0, err
			}
			tree.tree[key] = &ABITObject{
				dataType: 3,
				blob:     &([]([]byte){b}[0]),
			}
		case 0b0100:
			var b string
			b, index, err = decodeString(blob, index)
			if err != nil {
				return tree, 0, err
			}
			tree.tree[key] = &ABITObject{
				dataType: 4,
				text:     &([]string{b}[0]),
			}
		case 0b0101:
			var b ABITArray
			b, index, err = decodeArray(blob, index)
			if err != nil {
				return tree, 0, err
			}
			tree.tree[key] = &ABITObject{
				dataType: 5,
				array:    &([]ABITArray{b}[0]),
			}
		case 0b0110:
			var b ABITObject
			b, index, err = decodeTree(blob, index, true)
			if err != nil {
				return tree, 0, err
			}
			tree.tree[key] = &([]ABITObject{b}[0])
		default:
			return tree, 0, fmt.Errorf("invalid type")
		}
	}
	if int(index) > len(*blob) {
		return tree, 0, fmt.Errorf("corrupt array")
	}
	return tree, offset, nil
}

// InitLexicon creates an ABITLexicon used for seeing if any ABITObject follows the given schema or not.
//
//   - lexicon is a json with the specified scema.
//   - Returns ABITLexicon
//
// # Allowed types:
//   - "null"
//   - "boolean"
//   - "integer"
//   - "blob"
//   - "string"
//
// # Example:
//
//	abit.InitLexicon(`{
//		"key1":"boolean",
//		"key2":{
//			"nested":"key"
//		},
//		"key3":[
//			"string"
//			"string"
//			"null"
//		]
//	}`)
func InitLexicon(lexicon string) ABITLexicon {
	// Unmarshal JSON into a map
	var lexiconMap map[string]interface{}
	err := json.Unmarshal([]byte(lexicon), &lexiconMap)
	if err != nil {
		panic(err.Error())
	}

	return ABITLexicon{
		lexicon: jsonTypeTreeToABIT(lexiconMap),
	}
}

func jsonTypeArrayToABIT(lexicon []interface{}) ABITArray {
	arr := NewABITArray()

	for i := range lexicon {
		switch t := lexicon[i].(type) {
		case string:
			switch t {
			case "null":
				arr.Add(Null{})
			case "boolean":
				arr.Add(false)
			case "integer":
				arr.Add(int64(0))
			case "blob":
				arr.Add([]byte{})
			case "string":
				arr.Add("")
			default:
				panic("value must be any of: \"null\", \"boolean\", \"integer\", \"blob\", \"string\"")
			}
		case []interface{}: // Array
			arr.Add(jsonTypeArrayToABIT(t))
		case map[string]interface{}: // Tree
			arr.Add(jsonTypeTreeToABIT(t))
		default:
			panic("value to every key in lexicon must be either a string, array or tree")
		}
	}

	return *arr
}

func jsonTypeTreeToABIT(lexicon map[string]interface{}) ABITObject {
	// Create ABITObject
	tree, err := NewABITObject(&[]byte{})
	if err != nil {
		panic(err.Error())
	}
	keys := make([]string, 0, len(lexicon))
	for k := range lexicon {
		keys = append(keys, k)
	}
	for i := range keys {
		switch t := lexicon[keys[i]].(type) {
		case string:
			switch t {
			case "null":
				tree.Put(keys[i], Null{})
			case "boolean":
				tree.Put(keys[i], false)
			case "integer":
				tree.Put(keys[i], int64(0))
			case "blob":
				tree.Put(keys[i], []byte{})
			case "string":
				tree.Put(keys[i], "")
			default:
				panic("Value must be any of: \"null\", \"boolean\", \"integer\", \"blob\", \"string\"")
			}
		case []interface{}: // Array
			tree.Put(keys[i], jsonTypeArrayToABIT(t))
		case map[string]interface{}: // Tree
			tree.Put(keys[i], jsonTypeTreeToABIT(t))
		default:
			panic("Value to every key in lexicon must be either a string, array or tree")
		}
	}
	return *tree
}

// Matches checks if the given ABITObject matches the schema in the lexicon.
//
//   - Returns true if matches, returns false otherwise.
//
// # Example:
//
//	if lex.Matches(doc) {
//		// functions to handle the valid abit document here.
//	}
func (l *ABITLexicon) Matches(doc *ABITObject) bool {
	return matchTree(&l.lexicon, doc)
}

func matchTree(a *ABITObject, b *ABITObject) bool {
	keys1 := make([]string, 0, len(a.tree))
	for k := range a.tree {
		keys1 = append(keys1, k)
	}

	keys2 := make([]string, 0, len(b.tree))
	for k := range b.tree {
		keys2 = append(keys2, k)
	}

	// Same number of items?
	if len(keys1) != len(keys2) {
		return false
	}

	sort.Slice(keys1, func(i, j int) bool {
		if len(keys1[i]) == len(keys1[j]) {
			// If lengths are equal, sort lexicographically
			return keys1[i] < keys1[j]
		}
		// Otherwise, sort by length
		return len(keys1[i]) < len(keys1[j])
	})

	sort.Slice(keys2, func(i, j int) bool {
		if len(keys2[i]) == len(keys2[j]) {
			// If lengths are equal, sort lexicographically
			return keys2[i] < keys2[j]
		}
		// Otherwise, sort by length
		return len(keys2[i]) < len(keys2[j])
	})

	// Are keys identical?
	for i := int64(0); int(i) < len(keys1); i++ {
		if keys1[i] != keys2[i] {
			return false
		}
	}

	for i := range keys1 {
		if a.tree[keys1[i]].dataType != b.tree[keys1[i]].dataType {
			return false
		}

		switch a.tree[keys1[i]].dataType {
		case 0b0101: // Array
			if !matchArray(a.tree[keys1[i]], b.tree[keys1[i]]) {
				return false
			}
		case 0b0110: // Tree
			if !matchTree(a.tree[keys1[i]], b.tree[keys1[i]]) {
				return false
			}
		}
	}

	return true
}

func matchArray(a *ABITObject, b *ABITObject) bool {
	if a.dataType != 0b0101 || b.dataType != 0b0101 {
		return false
	}

	if len(a.array.array) != len(b.array.array) {
		return false
	}

	for i := range a.array.array {
		if a.array.array[i].dataType != b.array.array[i].dataType {
			return false
		}

		switch a.array.array[i].dataType {
		case 0b0101: // Array
			if !matchArray(a.array.array[i], b.array.array[i]) {
				return false
			}
		case 0b0110: // Tree
			if !matchTree(a.array.array[i], b.array.array[i]) {
				return false
			}
		}
	}

	return true
}
