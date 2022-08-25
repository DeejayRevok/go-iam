package security

import "golang.org/x/crypto/bcrypt"

type BcryptHashComparator struct{}

func (comparator *BcryptHashComparator) Compare(source string, sourceHash string) error {
	return bcrypt.CompareHashAndPassword([]byte(sourceHash), []byte(source))
}

func NewBcryptHashComparator() *BcryptHashComparator {
	comparator := BcryptHashComparator{}
	return &comparator
}
