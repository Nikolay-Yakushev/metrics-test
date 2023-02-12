package ports

type Iterator interface {
    getNext() *any
    hasNext() bool
}

type Collection interface {
    createIterator() Iterator
}