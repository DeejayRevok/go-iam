package hash

type HashComparator interface {
	Compare(source string, sourceHash string) error
}
