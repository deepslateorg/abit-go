# How to Use
```go
// create an array
arr := NewABITArray()

// add string to array
arr.Add("päror")

// add integer to array
arr.Add(int64(46410))

// create an empty tree
tree, _ := NewABITObject(&[]byte{})

// add array to tree
tree.Put("landet", *arr)

// add blob to tree
blob := []byte{0, 4, 1, 0}
tree.Put("riktnummer", blob)

// convert ABITObject to an abit binary / document
doc, err := tree.ToByteArray()
if err != nil {
    // Handle error TT~~TT
    // Something went seriously wrong
}

// create an ABITObject from an abit document
tree2, err := NewABITObject(&doc)
if err != nil {
    // Handle error TT~~TT
    // invalid document
}

// get array from tree
arr2, err := tree2.GetArray("landet")
if err != nil {
    // Handle error TT~~TT
}

// get value from array
vegetable, err := arr2.GetString(0)
if err != nil {
    // Handle error TT~~TT
    // item at index wasn't string / index out of bounds
}
```
# Spec
### key:
```
    [NNNNNNNN]   ╠ unsigned byte
    [BBBBBBBB]   ╗
       ....      ╠ number of bytes used for the key string is N+1 (UTF-8 encoded string)
    [BBBBBBBB]   ╝
```

### null:
```
    [0000|0000]
```

### boolean:
```
    true:  [0001|0001]
    false: [0000|0001]
```

### integer:
```
    [0XXX|0010] 
    [NNNNNNNN]   ╗
       ....      ╠ number of bytes used for the integer is X+1 (2s compliment & little-endian)    ║
    [NNNNNNNN]   ╝
```

### blob:
```
    [00XX|0011] 
    [NNNNNNNN]   ╗
       ....      ╠ number of bytes used for the integer is X+1 (2s compliment & little-endian)
    [NNNNNNNN]   ╝
    [BBBBBBBB]   ╗
       ....      ╠ number of bytes used for the blob is N
    [BBBBBBBB]   ╝
```

### string:
```
    [00XX|0100] 
    [NNNNNNNN]   ╗
       ....      ╠ number of bytes used for the integer is X+1 (2s compliment & little-endian)
    [NNNNNNNN]   ╝
    [SSSSSSSS]   ╗
       ....      ╠ number of bytes used for the string is N (UTF-8 encoded string)
    [SSSSSSSS]   ╝
```

### array:
```
    [00XX|0101] 
    [NNNNNNNN]   ╗
       ....      ╠ number of bytes used for the integer is X+1 (2s compliment & little-endian)
    [NNNNNNNN]   ╝
    [AAAAAAAA]   ╗
       ....      ╠ number of bytes used for the array is N
    [AAAAAAAA]   ╝
```

### tree:
```
    [00XX|0110] 
    [NNNNNNNN]   ╗
       ....      ╠ number of bytes used for the integer is X+1 (2s compliment & little-endian)
    [NNNNNNNN]   ╝
    [TTTTTTTT]   ╗
       ....      ╠ number of bytes used for the tree is N
    [TTTTTTTT]   ╝
```

### tree syntax:
```
    [  key   ] [ object ] ... [  key   ] [ object ]
```

### array syntax:
```
    [ object ] ... [ object ]
```

### other syntax:
* An integer must be the minimum amount of bytes required to represent it.
* While an array can be any order, trees need to be ordered such that smaller keys are first, if keys are of equal length, treat it as a big-endian integer and put the smaller integer first.