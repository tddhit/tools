package skiplist

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/tddhit/tools/log"
)

func TestSK(t *testing.T) {
	sk := New()
	for i := 500; i <= 1200; i++ {
		sk.Put([]byte("hello"+strconv.Itoa(i)), []byte("world"+strconv.Itoa(i)))
	}
	for i := 200; i <= 300; i++ {
		sk.Put([]byte("hello"+strconv.Itoa(i)), []byte("world"+strconv.Itoa(i)))
	}
	for i := 1; i <= 100; i++ {
		sk.Put([]byte("hello"+strconv.Itoa(i)), []byte("world"+strconv.Itoa(i)))
	}
	for i := 100; i <= 200; i++ {
		sk.Put([]byte("hello"+strconv.Itoa(i)), []byte("world"+strconv.Itoa(i)))
	}
	for i := 201; i <= 500; i++ {
		sk.Put([]byte("hello"+strconv.Itoa(i)), []byte("world"+strconv.Itoa(i)))
	}
	for i := 237; i <= 337; i++ {
		sk.Delete([]byte("hello" + strconv.Itoa(i)))
	}
	for i := 537; i <= 637; i++ {
		sk.Delete([]byte("hello" + strconv.Itoa(i)))
	}
	for i := 1; i <= 1300; i++ {
		value1 := sk.Get([]byte("hello" + strconv.Itoa(i)))
		value2 := []byte("world" + strconv.Itoa(i))
		if i >= 1 && i < 237 {
			assert(bytes.Compare(value1, value2) == 0)
		} else if i >= 237 && i <= 337 {
			assert(value1 == nil)
		} else if i >= 537 && i <= 637 {
			assert(value1 == nil)
		} else if i > 337 && i <= 1200 {
			assert(bytes.Compare(value1, value2) == 0)
		} else if i > 1200 {
			assert(value1 == nil)
		}
	}
}

func assert(equal bool) {
	if !equal {
		log.Panic()
	}
}
