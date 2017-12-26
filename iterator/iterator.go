package iterator

type Iterator interface {
	First()
	Next()
	HasNext() bool
}

type KVIterator interface {
	First()
	Next()
	HasNext() bool
	Key() []byte
	Value() []byte
}
