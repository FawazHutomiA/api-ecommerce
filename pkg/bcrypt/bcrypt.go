package bcrypt

import "golang.org/x/crypto/bcrypt"

func HashPassword(cost int, plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), cost)
	return string(bytes), err
}

func ComparePasswordHash(plain, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	return err == nil
}
