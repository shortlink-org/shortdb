package index

// Index - interface of index
type Index[T any] interface {
	Find(key T) []byte
	Insert(key T) error
	Delete(key T) error

	Marshal() ([]byte, error)
	UnMarshal(bytes []byte, item any) error
}
