package abit

import (
	"bytes"
	"math/rand"
	"strings"
	"testing"
)

func TestNull(t *testing.T) {
	tree, _ := NewABITObject(&[]byte{})
	// legal keys
	err := tree.Put("null obj", Null{})
	if err != nil {
		t.Fatal(err.Error())
	}
	err = tree.Put(strings.Repeat(" ", 128), Null{})
	if err != nil {
		t.Fatal(err.Error())
	}
	err = tree.Put(strings.Repeat(" ", 129), Null{})
	if err != nil {
		t.Fatal(err.Error())
	}
	err = tree.Put(strings.Repeat(" ", 255), Null{})
	if err != nil {
		t.Fatal(err.Error())
	}
	err = tree.Put(strings.Repeat(" ", 256), Null{})
	if err != nil {
		t.Fatal(err.Error())
	}

	// illegal keys
	err = tree.Put(strings.Repeat(" ", 257), Null{})
	if err == nil {
		t.Fatal("was able to add a too long key")
	}
	err = tree.Put(strings.Repeat(" ", 6969), Null{})
	if err == nil {
		t.Fatal("was able to add a too long key")
	}
	err = tree.Put("", Null{})
	if err == nil {
		t.Fatal("was able to add a too short key")
	}

	// fetch values from keys
	_, err = tree.GetNull("null obj")
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = tree.GetNull(strings.Repeat(" ", 128))
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = tree.GetNull(strings.Repeat(" ", 129))
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = tree.GetNull(strings.Repeat(" ", 255))
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = tree.GetNull(strings.Repeat(" ", 256))
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestBoolean(t *testing.T) {
	tree, _ := NewABITObject(&[]byte{})
	obj := true

	for i := 0; i < 2; i++ {
		// legal keys
		err := tree.Put("bool obj", obj)
		if err != nil {
			t.Fatal(err.Error())
		}
		err = tree.Put(strings.Repeat(" ", 128), obj)
		if err != nil {
			t.Fatal(err.Error())
		}
		err = tree.Put(strings.Repeat(" ", 129), obj)
		if err != nil {
			t.Fatal(err.Error())
		}
		err = tree.Put(strings.Repeat(" ", 255), obj)
		if err != nil {
			t.Fatal(err.Error())
		}
		err = tree.Put(strings.Repeat(" ", 256), obj)
		if err != nil {
			t.Fatal(err.Error())
		}

		// illegal keys
		err = tree.Put(strings.Repeat(" ", 257), obj)
		if err == nil {
			t.Fatal("was able to add a too long key")
		}
		err = tree.Put(strings.Repeat(" ", 6969), obj)
		if err == nil {
			t.Fatal("was able to add a too long key")
		}
		err = tree.Put("", obj)
		if err == nil {
			t.Fatal("was able to add a too short key")
		}

		// fetch values from keys
		var out bool
		out, err = tree.GetBool("bool obj")
		if err != nil {
			t.Fatal(err.Error())
		}
		if out != obj {
			t.Fatal("incorrect value")
		}
		out, err = tree.GetBool(strings.Repeat(" ", 128))
		if err != nil {
			t.Fatal(err.Error())
		}
		if out != obj {
			t.Fatal("incorrect value")
		}
		out, err = tree.GetBool(strings.Repeat(" ", 129))
		if err != nil {
			t.Fatal(err.Error())
		}
		if out != obj {
			t.Fatal("incorrect value")
		}
		out, err = tree.GetBool(strings.Repeat(" ", 255))
		if err != nil {
			t.Fatal(err.Error())
		}
		if out != obj {
			t.Fatal("incorrect value")
		}
		out, err = tree.GetBool(strings.Repeat(" ", 256))
		if err != nil {
			t.Fatal(err.Error())
		}
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
	err := tree.Put("int obj", obj)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = tree.Put(strings.Repeat(" ", 128), obj)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = tree.Put(strings.Repeat(" ", 129), obj)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = tree.Put(strings.Repeat(" ", 255), obj)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = tree.Put(strings.Repeat(" ", 256), obj)
	if err != nil {
		t.Fatal(err.Error())
	}

	// illegal keys
	err = tree.Put(strings.Repeat(" ", 257), obj)
	if err == nil {
		t.Fatal("was able to add a too long key")
	}
	err = tree.Put(strings.Repeat(" ", 6969), obj)
	if err == nil {
		t.Fatal("was able to add a too long key")
	}
	err = tree.Put("", obj)
	if err == nil {
		t.Fatal("was able to add a too short key")
	}

	// fetch values from keys
	var out int64
	out, err = tree.GetInteger("int obj")
	if err != nil {
		t.Fatal(err.Error())
	}
	if out != obj {
		t.Fatal("incorrect value")
	}

	for i := 0; i < 10000; i++ {
		obj = int64(rand.Uint64())
		tree.Put("int obj", obj)
		tree.Put("meow", obj+5)
		tree.Put("meowmeow", -obj)
		out, err = tree.GetInteger("int obj")
		if err != nil {
			t.Fatal(err.Error())
		}
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
		err := tree.Put("blob obj", obj)
		if err != nil {
			t.Fatal(err.Error())
		}
		out, err = tree.GetBlob("blob obj")
		if err != nil {
			t.Fatal(err.Error())
		}
		if !bytes.Equal(obj, *out) {
			t.Fatal("Input not same as output")
		}
		obj = randBytes(rand.Int63n(2) * i)
		err = tree.Put(strings.Repeat(" ", 128), obj)
		if err != nil {
			t.Fatal(err.Error())
		}
		out, err = tree.GetBlob(strings.Repeat(" ", 128))
		if err != nil {
			t.Fatal(err.Error())
		}
		if !bytes.Equal(obj, *out) {
			t.Fatal("Input not same as output")
		}
		obj = randBytes(rand.Int63n(2) * i)
		err = tree.Put(strings.Repeat(" ", 129), obj)
		if err != nil {
			t.Fatal(err.Error())
		}
		out, err = tree.GetBlob(strings.Repeat(" ", 129))
		if err != nil {
			t.Fatal(err.Error())
		}
		if !bytes.Equal(obj, *out) {
			t.Fatal("Input not same as output")
		}
		obj = randBytes(rand.Int63n(2) * i)
		err = tree.Put(strings.Repeat(" ", 255), obj)
		if err != nil {
			t.Fatal(err.Error())
		}
		out, err = tree.GetBlob(strings.Repeat(" ", 255))
		if err != nil {
			t.Fatal(err.Error())
		}
		if !bytes.Equal(obj, *out) {
			t.Fatal("Input not same as output")
		}
		obj = randBytes(rand.Int63n(2) * i)
		err = tree.Put(strings.Repeat(" ", 256), obj)
		if err != nil {
			t.Fatal(err.Error())
		}
		out, err = tree.GetBlob(strings.Repeat(" ", 256))
		if err != nil {
			t.Fatal(err.Error())
		}
		if !bytes.Equal(obj, *out) {
			t.Fatal("Input not same as output")
		}

		// illegal keys
		err = tree.Put(strings.Repeat(" ", 257), obj)
		if err == nil {
			t.Fatal("was able to add a too long key")
		}
		err = tree.Put(strings.Repeat(" ", 6969), obj)
		if err == nil {
			t.Fatal("was able to add a too long key")
		}
		err = tree.Put("", obj)
		if err == nil {
			t.Fatal("was able to add a too short key")
		}
	}
}

func TestString(t *testing.T) {
	tree, _ := NewABITObject(&[]byte{})
	var obj string
	var out *string

	for i := int64(0); i < 5000; i++ {
		// legal keys
		obj = string(randBytes(rand.Int63n(2) * i))
		err := tree.Put("string obj", string(obj))
		if err != nil {
			t.Fatal(err.Error())
		}
		out, err = tree.GetString("string obj")
		if err != nil {
			t.Fatal(err.Error())
		}
		if (*out) != obj {
			t.Fatal("Input not same as output")
		}
		obj = string(randBytes(rand.Int63n(2) * i))
		err = tree.Put(strings.Repeat(" ", 128), obj)
		if err != nil {
			t.Fatal(err.Error())
		}
		out, err = tree.GetString(strings.Repeat(" ", 128))
		if err != nil {
			t.Fatal(err.Error())
		}
		if (*out) != obj {
			t.Fatal("Input not same as output")
		}
		obj = string(randBytes(rand.Int63n(2) * i))
		err = tree.Put(strings.Repeat(" ", 129), obj)
		if err != nil {
			t.Fatal(err.Error())
		}
		out, err = tree.GetString(strings.Repeat(" ", 129))
		if err != nil {
			t.Fatal(err.Error())
		}
		if (*out) != obj {
			t.Fatal("Input not same as output")
		}
		obj = string(randBytes(rand.Int63n(2) * i))
		err = tree.Put(strings.Repeat(" ", 255), obj)
		if err != nil {
			t.Fatal(err.Error())
		}
		out, err = tree.GetString(strings.Repeat(" ", 255))
		if err != nil {
			t.Fatal(err.Error())
		}
		if (*out) != obj {
			t.Fatal("Input not same as output")
		}
		obj = string(randBytes(rand.Int63n(2) * i))
		err = tree.Put(strings.Repeat(" ", 256), obj)
		if err != nil {
			t.Fatal(err.Error())
		}
		out, err = tree.GetString(strings.Repeat(" ", 256))
		if err != nil {
			t.Fatal(err.Error())
		}
		if (*out) != obj {
			t.Fatal("Input not same as output")
		}

		// illegal keys
		err = tree.Put(strings.Repeat(" ", 257), obj)
		if err == nil {
			t.Fatal("was able to add a too long key")
		}
		err = tree.Put(strings.Repeat(" ", 6969), obj)
		if err == nil {
			t.Fatal("was able to add a too long key")
		}
		err = tree.Put("", obj)
		if err == nil {
			t.Fatal("was able to add a too short key")
		}
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

	treeBlob1, _ := tree.ToByteArray()
	/*for i := 0; i < len(treeBlob1); i++ {
		fmt.Printf("%02X ", treeBlob1[i])
	}
	fmt.Println()*/
	tree2, err := NewABITObject(&treeBlob1)
	if err != nil {
		t.Fatal(err.Error())
	}
	treeBlob2, _ := tree2.ToByteArray()

	if !bytes.Equal(treeBlob1, treeBlob2) {
		t.Fatal("abit not equal")
	}
}

func TestInvalidTree(t *testing.T) {
	for i := 0; i < 500000; i++ {
		obj := randBytes(rand.Int63n(256) + 512)
		_, err := NewABITObject(&obj)
		if err == nil {
			t.Fatal("this tree should be invalid, try rerunning test if this happens")
		}
	}
}
