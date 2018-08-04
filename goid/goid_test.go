package goid

import "testing"

func TestGoid(t *testing.T) {
	println(Get() == goid())
}

func BenchmarkGoidFromASM(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get()
	}
}

func BenchmarkGoidFromStack(b *testing.B) {
	for i := 0; i < b.N; i++ {
		goid()
	}
}
