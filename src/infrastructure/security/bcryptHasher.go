package security

import "golang.org/x/crypto/bcrypt"

type BcryptHasher struct {
}

func (hasher *BcryptHasher) Hash(password string) (*string, error) {
	passwordBytes := []byte(password)

	hash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)
	if err != nil {
		return nil, err
	}
	strHash := string(hash)
	return &strHash, nil
}

func NewBcryptHasher() *BcryptHasher {
	hasher := BcryptHasher{}
	return &hasher
}
