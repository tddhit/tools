package sort

import (
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/tddhit/tools/log"
)

const MAX = 100

func assert(b bool, err string) {
	if !b {
		log.Fatal(err)
	}
}

func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	if (a == nil) != (b == nil) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func shuffle(data []int) []int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ret := make([]int, len(data))
	perm := r.Perm(len(data))
	for i, randIndex := range perm {
		ret[i] = data[randIndex]
	}
	return ret
}

func TestSort(t *testing.T) {
	var orig []int
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < MAX; i++ {
		orig = append(orig, rand.Intn(MAX))
	}
	sort.Ints(orig)

	data := shuffle(orig)
	assert(!equal(orig, data), "shuffle equal")
	QuickSort(data)
	assert(equal(orig, data), "quick sort not equal")

	data = shuffle(orig)
	assert(!equal(orig, data), "shuffle equal")
	InsertSort(data)
	assert(equal(orig, data), "insert sort not equal")

	data = shuffle(orig)
	assert(!equal(orig, data), "shuffle equal")
	SelectSort(data)
	assert(equal(orig, data), "select sort not equal")

	data = shuffle(orig)
	assert(!equal(orig, data), "shuffle equal")
	BubbleSort(data)
	assert(equal(orig, data), "bubble sort not equal")

	assert(!equal([]int{}, nil), "nil/empty equal")
}
