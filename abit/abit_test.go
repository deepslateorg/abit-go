package abit

import (
	"fmt"
	"testing"
)

func TestNull(t *testing.T) {
	tree, _ := NewABITObject(&[]byte{})
	err := tree.Put("thingy", "segsss uwu åäö")
	if err != nil {
		t.Fatal(err.Error())
	}
	err = tree.Put("thingy2", "bruh")
	if err != nil {
		t.Fatal(err.Error())
	}
	s1, _ := tree.GetString("thingy")
	_, _ = tree.GetString("thingy2")
	fmt.Println(*s1)
	b, _ := tree.ToByteArray()
	for _, value := range b {
		fmt.Printf("%02X ", value)
	}
	fmt.Println()
}
