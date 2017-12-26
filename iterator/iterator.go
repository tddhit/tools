package iterator

type Iterator interface {
	First()
	Next()
	End() bool
	HasNext() bool
}

type KVIterator interface {
	First()
	Next()
	End() bool
	Key() []byte
	Value() []byte
}
