package abit

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"strings"
	"testing"

	"github.com/multiformats/go-multibase"
)

// Helper function to check if a function panics
func shouldPanic(t *testing.T, f func()) {
	t.Helper() // Marks this function as a helper in the test logs

	// Use defer and recover to catch a panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but function completed successfully.")
		}
	}()

	// Run the function that we expect to panic
	f()
}

func TestNull(t *testing.T) {
	tree, _ := NewABITObject(&[]byte{})
	// legal keys
	tree.Put("null obj", Null{})
	tree.Put(strings.Repeat(" ", 128), Null{})
	tree.Put(strings.Repeat(" ", 129), Null{})
	tree.Put(strings.Repeat(" ", 255), Null{})
	tree.Put(strings.Repeat(" ", 256), Null{})

	// illegal keys
	shouldPanic(t, func() { tree.Put(strings.Repeat(" ", 257), Null{}) })
	shouldPanic(t, func() { tree.Put(strings.Repeat(" ", 6969), Null{}) })
	shouldPanic(t, func() { tree.Put("", Null{}) })

	// fetch values from keys
	tree.GetNull("null obj")
	tree.GetNull(strings.Repeat(" ", 128))
	tree.GetNull(strings.Repeat(" ", 129))
	tree.GetNull(strings.Repeat(" ", 255))
	tree.GetNull(strings.Repeat(" ", 256))
}

func TestBoolean(t *testing.T) {
	tree, _ := NewABITObject(&[]byte{})
	obj := true

	for i := 0; i < 2; i++ {
		// legal keys
		tree.Put("bool obj", obj)
		tree.Put(strings.Repeat(" ", 128), obj)
		tree.Put(strings.Repeat(" ", 129), obj)
		tree.Put(strings.Repeat(" ", 255), obj)
		tree.Put(strings.Repeat(" ", 256), obj)

		// illegal keys
		shouldPanic(t, func() { tree.Put(strings.Repeat(" ", 257), obj) })
		shouldPanic(t, func() { tree.Put(strings.Repeat(" ", 6969), obj) })
		shouldPanic(t, func() { tree.Put("", obj) })

		// fetch values from keys
		var out bool
		out = tree.GetBool("bool obj")
		if out != obj {
			t.Fatal("incorrect value")
		}
		out = tree.GetBool(strings.Repeat(" ", 128))
		if out != obj {
			t.Fatal("incorrect value")
		}
		out = tree.GetBool(strings.Repeat(" ", 129))
		if out != obj {
			t.Fatal("incorrect value")
		}
		out = tree.GetBool(strings.Repeat(" ", 255))
		if out != obj {
			t.Fatal("incorrect value")
		}
		out = tree.GetBool(strings.Repeat(" ", 256))
		if out != obj {
			t.Fatal("incorrect value")
		}
		obj = false
	}
}

func TestInteger(t *testing.T) {
	tree, _ := NewABITObject(&[]byte{})
	obj := int64(6969696969420)

	// legal keys
	tree.Put("int obj", obj)
	tree.Put(strings.Repeat(" ", 128), obj)
	tree.Put(strings.Repeat(" ", 129), obj)
	tree.Put(strings.Repeat(" ", 255), obj)
	tree.Put(strings.Repeat(" ", 256), obj)

	// illegal keys
	shouldPanic(t, func() { tree.Put(strings.Repeat(" ", 257), obj) })
	shouldPanic(t, func() { tree.Put(strings.Repeat(" ", 6969), obj) })
	shouldPanic(t, func() { tree.Put("", obj) })

	// fetch values from keys
	var out int64
	out = tree.GetInteger("int obj")
	if out != obj {
		t.Fatal("incorrect value")
	}

	for i := 0; i < 10000; i++ {
		obj = int64(rand.Uint64())
		tree.Put("int obj", obj)
		tree.Put("meow", obj+5)
		tree.Put("meowmeow", -obj)
		out = tree.GetInteger("int obj")
		if out != obj {
			t.Fatal("didn't fetch same number as inputted")
		}
	}
}

func randBytes(length int64) []byte {
	blob := make([]byte, length)
	for i := range blob {
		blob[i] = byte(rand.Intn(256)) // Generate a random byte (0-255)
	}
	return blob
}

func TestBlob(t *testing.T) {
	tree, _ := NewABITObject(&[]byte{})
	var obj []byte
	var out *[]byte

	for i := int64(0); i < 5000; i++ {
		// legal keys
		obj = randBytes(rand.Int63n(2) * i)
		tree.Put("blob obj", obj)
		out = tree.GetBlob("blob obj")
		if !bytes.Equal(obj, *out) {
			t.Fatal("Input not same as output")
		}
		obj = randBytes(rand.Int63n(2) * i)
		tree.Put(strings.Repeat(" ", 128), obj)
		out = tree.GetBlob(strings.Repeat(" ", 128))
		if !bytes.Equal(obj, *out) {
			t.Fatal("Input not same as output")
		}
		obj = randBytes(rand.Int63n(2) * i)
		tree.Put(strings.Repeat(" ", 129), obj)
		out = tree.GetBlob(strings.Repeat(" ", 129))
		if !bytes.Equal(obj, *out) {
			t.Fatal("Input not same as output")
		}
		obj = randBytes(rand.Int63n(2) * i)
		tree.Put(strings.Repeat(" ", 255), obj)
		out = tree.GetBlob(strings.Repeat(" ", 255))
		if !bytes.Equal(obj, *out) {
			t.Fatal("Input not same as output")
		}
		obj = randBytes(rand.Int63n(2) * i)
		tree.Put(strings.Repeat(" ", 256), obj)
		out = tree.GetBlob(strings.Repeat(" ", 256))
		if !bytes.Equal(obj, *out) {
			t.Fatal("Input not same as output")
		}

		// illegal keys
		shouldPanic(t, func() { tree.Put(strings.Repeat(" ", 257), obj) })
		shouldPanic(t, func() { tree.Put(strings.Repeat(" ", 6969), obj) })
		shouldPanic(t, func() { tree.Put("", obj) })
	}
}

func TestString(t *testing.T) {
	tree, _ := NewABITObject(&[]byte{})
	var obj string
	var out *string

	for i := int64(0); i < 5000; i++ {
		// legal keys
		obj = string(randBytes(rand.Int63n(2) * i))
		tree.Put("string obj", string(obj))
		out = tree.GetString("string obj")
		if (*out) != obj {
			t.Fatal("Input not same as output")
		}
		obj = string(randBytes(rand.Int63n(2) * i))
		tree.Put(strings.Repeat(" ", 128), obj)
		out = tree.GetString(strings.Repeat(" ", 128))
		if (*out) != obj {
			t.Fatal("Input not same as output")
		}
		obj = string(randBytes(rand.Int63n(2) * i))
		tree.Put(strings.Repeat(" ", 129), obj)
		out = tree.GetString(strings.Repeat(" ", 129))
		if (*out) != obj {
			t.Fatal("Input not same as output")
		}
		obj = string(randBytes(rand.Int63n(2) * i))
		tree.Put(strings.Repeat(" ", 255), obj)
		out = tree.GetString(strings.Repeat(" ", 255))
		if (*out) != obj {
			t.Fatal("Input not same as output")
		}
		obj = string(randBytes(rand.Int63n(2) * i))
		tree.Put(strings.Repeat(" ", 256), obj)
		out = tree.GetString(strings.Repeat(" ", 256))
		if (*out) != obj {
			t.Fatal("Input not same as output")
		}

		// illegal keys
		shouldPanic(t, func() { tree.Put(strings.Repeat(" ", 257), obj) })
		shouldPanic(t, func() { tree.Put(strings.Repeat(" ", 6969), obj) })
		shouldPanic(t, func() { tree.Put("", obj) })
	}
}

func TestGenericTree(t *testing.T) {
	tree, _ := NewABITObject(&[]byte{})

	tree.Put("null obj", Null{})
	tree.Put("boolean t obj", true)
	tree.Put("boolean f obj", false)
	tree.Put("integer p big", int64(69696969420))
	tree.Put("integer n big", int64(-69696969420))
	tree.Put("integer p small", int64(69))
	tree.Put("integer n small", int64(-69))

	blobs := []byte{0, 1, 2, 3, 4, 5, 6, 7}
	blobb := randBytes(1024)

	tree.Put("small blob obj", blobs)
	tree.Put("big blob obj", blobb)
	tree.Put("string small", "Hello 💀")
	tree.Put("string big", "Lorem ipsum dolor sit 💺🧘‍♂️ amet, consectetur adipiscing elit, sed do 👌 eiusmod tempor incididunt ut 🅱️🤫 labore et 🎻📯🎺 dolore magna aliqua. Ut 🅱️🤫 enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut 🅱️🤫 aliquip ex 🎻🎻🎻🎻🎻🎻 ea commodo consequat. Duis aute irure dolor in 😩 reprehenderit in 🍝🥫 voluptate velit esse cillum dolore eu 🤲 fugiat nulla pariatur. Excepteur sint occaecat cupidatat non ❌ proident, sunt in 😩😂 culpa qui officia deserunt mollit anim id 😗 est ⏩💆👷🏿 laborum.")
	arr := NewABITArray()
	arr.Add("1")
	arr.Add(int64(2))
	arr.Add(true)
	arr.Add(int64(4))
	arr.Add(int64(3))
	arr.Add("2")
	arr.Add(int64(1))
	tree.Put("array obj", *arr)
	nestedTree, _ := NewABITObject(&[]byte{})
	nestedTree.Put("thing", "AMOGUS")
	tree.Put("nesty", *nestedTree)
	tree.Put("a very very very very very very very very very very very very very  very very very very very very very very very very very very very  very very very very very very very very very very very very very  very very very very very very very very very ve long key", "meow")

	treeBlob1 := tree.ToByteArray()
	/*for i := 0; i < len(treeBlob1); i++ {
		fmt.Printf("%02X ", treeBlob1[i])
	}
	fmt.Println()*/
	tree2, err := NewABITObject(&treeBlob1)
	if err != nil {
		t.Fatal(err.Error())
	}
	treeBlob2 := tree2.ToByteArray()

	if !bytes.Equal(treeBlob1, treeBlob2) {
		t.Fatal("abit not equal")
	}
}

func TestGenericArray(t *testing.T) {
	tree, _ := NewABITObject(&[]byte{})
	arr := NewABITArray()
	arr.Add("1")
	arr.Add(int64(2))
	arr.Add(true)
	arr.Add(int64(4))
	arr.Add(int64(3))
	arr.Add("2")
	arr.Add(int64(1))
	tree.Put("array obj", *arr)

	treeBlob1 := tree.ToByteArray()

	tree2, err := NewABITObject(&treeBlob1)
	if err != nil {
		t.Fatal(err.Error())
	}
	treeBlob2 := tree2.ToByteArray()

	if !bytes.Equal(treeBlob1, treeBlob2) {
		t.Fatal("abit not equal")
	}

	arr2 := tree2.GetArray("array obj")

	s1 := arr.GetString(0)
	s2 := arr2.GetString(0)
	if *s1 != *s2 {
		t.Fatalf("strings not equal: %s <-> %s", *s1, *s2)
	}

	i1 := arr.GetInteger(1)
	i2 := arr2.GetInteger(1)
	if i1 != i2 {
		t.Fatalf("integers not equal: %d <-> %d", i1, i2)
	}

	if arr2.Length() != 7 {
		t.Fatalf("array is not expected length")
	}

	if tree.Keys()[0] != "array obj" {
		t.Fatalf("tree not containing expected key")
	}
}

func TestInvalidTree(t *testing.T) {
	for i := 0; i < 200000; i++ {
		obj := randBytes(rand.Int63n(256) + 512)
		_, err := NewABITObject(&obj)
		if err == nil {
			t.Fatal("this tree should be invalid, try rerunning test if this happens")
		}
	}
}

var replacementValues = []string{"null", "boolean", "integer", "blob", "string"}

func modifyValues(data interface{}) {
	hasModified := false
	for !hasModified {
		switch v := data.(type) {
		case map[string]interface{}:
			for key, val := range v {
				switch val := val.(type) {
				case string:
					// Check if val is one of the target values and replace it
					if val == "null" || val == "boolean" || val == "integer" || val == "blob" || val == "string" {
						v[key] = getRandomReplacement(&hasModified, val)
					}
				case map[string]interface{}, []interface{}:
					// Recur for nested maps and slices
					modifyValues(val)
				}
			}
		case []interface{}:
			for i, val := range v {
				switch val := val.(type) {
				case string:
					// Check if val is one of the target values and replace it
					if val == "null" || val == "boolean" || val == "integer" || val == "blob" || val == "string" {
						v[i] = getRandomReplacement(&hasModified, val)
					}
				case map[string]interface{}, []interface{}:
					// Recur for nested maps and slices
					modifyValues(val)
				}
			}
		}
	}
}

func getRandomReplacement(modified *bool, currentVal string) string {
	if !(*modified) && (rand.Intn(30) == 0) {
		val := replacementValues[rand.Intn(len(replacementValues))]
		for val == currentVal {
			val = replacementValues[rand.Intn(len(replacementValues))]
		}
		*modified = true
		return val
	}
	return currentVal
}

func TestLexicon(t *testing.T) {
	jsonData := `{
		"key1": "null",
		"key2": "boolean",
		"key3": "integer",
		"key4": "blob",
		"key5": "string",
		"key6": [
			"null",
			"boolean",
			"integer",
			"blob",
			"string",
			{
			"key1": "boolean",
			"key2": "integer"
			},
			[
				"boolean",
				"integer"
			]
		],
		"key7": {
			"key1": "boolean",
			"key2": "integer"
		}
	}`
	lex := InitLexicon(jsonData)

	tree, _ := NewABITObject(&[]byte{})
	tree.Put("key1", Null{})
	tree.Put("key2", true)
	tree.Put("key3", int64(69))
	tree.Put("key4", []byte{0, 1, 2, 3, 128})
	tree.Put("key5", "mrrowp :3")

	arr := NewABITArray()
	arr.Add(Null{})
	arr.Add(false)
	arr.Add(int64(7331))
	arr.Add([]byte{0, 1, 2, 3, 33})
	arr.Add("shlirp shlorp")

	arrtree, _ := NewABITObject(&[]byte{})
	arrtree.Put("key1", true)
	arrtree.Put("key2", int64(123456789))

	arr.Add(*arrtree)

	arrarr := NewABITArray()
	arrarr.Add(true)
	arrarr.Add(int64(410))

	arr.Add(*arrarr)

	tree.Put("key6", *arr)

	treetree, _ := NewABITObject(&[]byte{})
	treetree.Put("key1", true)
	treetree.Put("key2", int64(3897))

	tree.Put("key7", *treetree)

	if !(&lex).Matches(tree) {
		t.Fatalf("Doesn't match when should")
	}

	for i := 0; i < 10000; i++ {

		var data interface{}
		err := json.Unmarshal([]byte(jsonData), &data)
		if err != nil {
			t.Fatal(err.Error())
		}

		modifyValues(data)

		modifiedJson, err := json.MarshalIndent(data, "", "    ")
		if err != nil {
			t.Fatal(err.Error())
			return
		}

		lex := InitLexicon(string(modifiedJson))

		if (&lex).Matches(tree) {
			t.Fatalf("match when shouldn't")
		}
	}
}

func TestConvertionToJson(t *testing.T) {
	tree, _ := NewABITObject(&[]byte{})

	tree.Put("null obj", Null{})
	tree.Put("boolean t obj", true)
	tree.Put("boolean f obj", false)
	tree.Put("integer p big", int64(69696969420))
	tree.Put("integer n big", int64(-69696969420))
	tree.Put("integer p small", int64(69))
	tree.Put("integer n small", int64(-69))

	blobs := []byte{0, 1, 2, 3, 4, 5, 6, 7}
	_, blobb, _ := multibase.Decode("z4YJbaPkw5uDwGtEQ3NMvau2EngLQFjmuR8ewchb6xN7LxqBNennqsyE6kR1Xes3azGBZ8HV6KL9A1mMcnP64XvmJCd5JXaXEUaaMGimjaccSSUvLEKCAWrGPqagiSihrqKE6mYyUDgKLfv8yqbNwHoXqi1AFxyCxX6AfaNiRDYAL9z7jcb965rurAFxBHoyWksFL7bQGm1zPuw43Y7JouByY7Y7pFtm374g5gPPGetVbfgsiqvvL7oqJ68vSwwFXDqeJ5Z1uNdu5hVs9qBQMEuaLQ8qjyhjjBnfcQRWvtgqcQnsunfheDpaqmRi9Cx39wBMgJQ44R8R8mZz3UGWvHqmgknHuBWMkXEkHBmJD754gfSqHuGTvNtWtPbgm212QvTeLiiCCn28EH7d2WSGkqQJHc1fhc6ebUwQWw1xTNtjdvSX7shrB4ErNTBDvuAfx2Buaicw5ngUNpWVYqzcdrDND6MkXL1vrcW67tAa1rqQZnacrfydAgBJVcPPKBLW6YF2LUdijAyvZA9rDaHmWXTMwwygx68UFpjZED23noHdVBAcuHts4AwGt4WenVQYC8ffcJSrXrqsB32gyTRj7mvQ14rDjEAYcfYnSmmqFJ114dVV51quZxLuMnUwy4nEZFq7K8UYQkAxN9XS2H8KSHdpX6BLLHRHvoL4X77vcovSEWV4qPn3fJsg8Jvn8kPUaE4LrpGpaE6oTgem5q6mvZ7zsMbP1NA1bKzV4RSKQNrimLtJPzWREJWymk9e54Jk3qTWdJjXrLRP3oWYtyQ4X9nr8moGMvySAsJZccAN3n6L3242LScHBBAEYzwHbjCVCxJwYMVTppxZ4ZxP6JUzupPwwFv4vDRotz8xpswmJgQzCUA1VNmDYQQE5tHQQTPQac5X1ndvJKrijSsdEk7E3v7RV9NiWDeCSdQLKhbUNnhTJNXsqABEPfGaZZ28bGzSCATZSz9mKsp75WN1QNu43MGVmnnjsVYmNnFZ71cam7S1BdeRsogDi2ThYTPumz5wpa5RQXXjaWMFJ3zKmdeRY19on5LuZHnrW4zVeyNznJZ1pWTzfXzbCwPAWs71RDv1BgAzgjgVdvFTH6RsnA1gP1xXHVSRWupYZeGPieDCsDLegNMm2459rA4k6xKoo4mk9z3bXH7sbEbFWFpJkKwYMLTG6fzX2bncJitpgSTyHkEa1dEoQPDDujcKFUCJ7fDbPRrde9NPuwQpcdeWxv4xFQWbeWBPDDoGHqbpyLzhDKVSr4MB9mbd34ZyQZUnh7XH2KuGVdS93kyo7o7SnDVLC3SnXvvX7Lq1uhMJQcRugMPbfBdg9z9MDfYFP933binE8YgV1GeA9TZPRWRFg9CUVRCKXhg2VqhcGSFCUyELukU8rNc3raooQjMA")

	tree.Put("small blob obj", blobs)
	tree.Put("big blob obj", blobb)
	tree.Put("string small", "Hello 💀")
	tree.Put("string big", "Lorem ipsum dolor sit 💺🧘‍♂️ amet, consectetur adipiscing elit, sed do 👌 eiusmod tempor incididunt ut 🅱️🤫 labore et 🎻📯🎺 dolore magna aliqua. Ut 🅱️🤫 enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut 🅱️🤫 aliquip ex 🎻🎻🎻🎻🎻🎻 ea commodo consequat. Duis aute irure dolor in 😩 reprehenderit in 🍝🥫 voluptate velit esse cillum dolore eu 🤲 fugiat nulla pariatur. Excepteur sint occaecat cupidatat non ❌ proident, sunt in 😩😂 culpa qui officia deserunt mollit anim id 😗 est ⏩💆👷🏿 laborum.")
	arr := NewABITArray()
	arr.Add("1")
	arr.Add(int64(2))
	arr.Add(true)
	arr.Add(int64(4))
	arr.Add(int64(3))
	arr.Add("2")
	arr.Add(int64(1))
	tree.Put("array obj", *arr)
	nestedTree, _ := NewABITObject(&[]byte{})
	nestedTree.Put("thing", "AMOGUS")
	tree.Put("nesty", *nestedTree)
	tree.Put("a very very very very very very very very very very very very very  very very very very very very very very very very very very very  very very very very very very very very very very very very very  very very very very very very very very very ve long key", "meow")

	c := `{"nesty":{"thing":"AMOGUS"},"null obj":null,"array obj":["1",2,true,4,3,"2",1],"string big":"Lorem ipsum dolor sit 💺🧘‍♂️ amet, consectetur adipiscing elit, sed do 👌 eiusmod tempor incididunt ut 🅱️🤫 labore et 🎻📯🎺 dolore magna aliqua. Ut 🅱️🤫 enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut 🅱️🤫 aliquip ex 🎻🎻🎻🎻🎻🎻 ea commodo consequat. Duis aute irure dolor in 😩 reprehenderit in 🍝🥫 voluptate velit esse cillum dolore eu 🤲 fugiat nulla pariatur. Excepteur sint occaecat cupidatat non ❌ proident, sunt in 😩😂 culpa qui officia deserunt mollit anim id 😗 est ⏩💆👷🏿 laborum.","big blob obj_b":"z4YJbaPkw5uDwGtEQ3NMvau2EngLQFjmuR8ewchb6xN7LxqBNennqsyE6kR1Xes3azGBZ8HV6KL9A1mMcnP64XvmJCd5JXaXEUaaMGimjaccSSUvLEKCAWrGPqagiSihrqKE6mYyUDgKLfv8yqbNwHoXqi1AFxyCxX6AfaNiRDYAL9z7jcb965rurAFxBHoyWksFL7bQGm1zPuw43Y7JouByY7Y7pFtm374g5gPPGetVbfgsiqvvL7oqJ68vSwwFXDqeJ5Z1uNdu5hVs9qBQMEuaLQ8qjyhjjBnfcQRWvtgqcQnsunfheDpaqmRi9Cx39wBMgJQ44R8R8mZz3UGWvHqmgknHuBWMkXEkHBmJD754gfSqHuGTvNtWtPbgm212QvTeLiiCCn28EH7d2WSGkqQJHc1fhc6ebUwQWw1xTNtjdvSX7shrB4ErNTBDvuAfx2Buaicw5ngUNpWVYqzcdrDND6MkXL1vrcW67tAa1rqQZnacrfydAgBJVcPPKBLW6YF2LUdijAyvZA9rDaHmWXTMwwygx68UFpjZED23noHdVBAcuHts4AwGt4WenVQYC8ffcJSrXrqsB32gyTRj7mvQ14rDjEAYcfYnSmmqFJ114dVV51quZxLuMnUwy4nEZFq7K8UYQkAxN9XS2H8KSHdpX6BLLHRHvoL4X77vcovSEWV4qPn3fJsg8Jvn8kPUaE4LrpGpaE6oTgem5q6mvZ7zsMbP1NA1bKzV4RSKQNrimLtJPzWREJWymk9e54Jk3qTWdJjXrLRP3oWYtyQ4X9nr8moGMvySAsJZccAN3n6L3242LScHBBAEYzwHbjCVCxJwYMVTppxZ4ZxP6JUzupPwwFv4vDRotz8xpswmJgQzCUA1VNmDYQQE5tHQQTPQac5X1ndvJKrijSsdEk7E3v7RV9NiWDeCSdQLKhbUNnhTJNXsqABEPfGaZZ28bGzSCATZSz9mKsp75WN1QNu43MGVmnnjsVYmNnFZ71cam7S1BdeRsogDi2ThYTPumz5wpa5RQXXjaWMFJ3zKmdeRY19on5LuZHnrW4zVeyNznJZ1pWTzfXzbCwPAWs71RDv1BgAzgjgVdvFTH6RsnA1gP1xXHVSRWupYZeGPieDCsDLegNMm2459rA4k6xKoo4mk9z3bXH7sbEbFWFpJkKwYMLTG6fzX2bncJitpgSTyHkEa1dEoQPDDujcKFUCJ7fDbPRrde9NPuwQpcdeWxv4xFQWbeWBPDDoGHqbpyLzhDKVSr4MB9mbd34ZyQZUnh7XH2KuGVdS93kyo7o7SnDVLC3SnXvvX7Lq1uhMJQcRugMPbfBdg9z9MDfYFP933binE8YgV1GeA9TZPRWRFg9CUVRCKXhg2VqhcGSFCUyELukU8rNc3raooQjMA","string small":"Hello 💀","boolean f obj":false,"boolean t obj":true,"integer n big":-69696969420,"integer p big":69696969420,"small blob obj_b":"z13DUyZY2dc","integer n small":-69,"integer p small":69,"a very very very very very very very very very very very very very  very very very very very very very very very very very very very  very very very very very very very very very very very very very  very very very very very very very very very ve long key":"meow"}`

	if c != tree.ToJson() {
		t.Fatal("incorrectly converted abit to json")
	}
}
