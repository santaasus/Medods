package security

import "golang.org/x/crypto/bcrypt"

func GeneratePasswordHash(password string) (hash []byte, err error) {
	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	return
}
