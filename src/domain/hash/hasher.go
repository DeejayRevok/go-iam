package hash

type Hasher interface {
	Hash(source string) (*string, error)
}
