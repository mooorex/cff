package CFF

import (
	"crypto/rand"
	"fmt"
	"testing"
)

func TestFindBSet(t *testing.T) {
	cf, err := NewCFF(8, 31, 8, 3)
	if err != nil {
		t.Fatal(err)
	}

	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		t.Fatal(err)
	}

	str := ""
	for _, item := range b {
		str += fmt.Sprintf("%08b", item)
	}

	kk, err := cf.FindBSet(str)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(kk, len(kk))
}

func BenchmarkNewCFF(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewCFF(8, 31, 8, 3)
	}
}

func BenchmarkFindBSet(b *testing.B) {
	cf, err := NewCFF(8, 31, 8, 3)
	if err != nil {
		b.Fatal(err)
	}

	x := make([]byte, 32)
	_, err = rand.Read(x)
	if err != nil {
		b.Fatal(err)
	}

	str := ""
	for _, item := range x {
		str += fmt.Sprintf("%08b", item)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cf.FindBSet(str)
	}
}
