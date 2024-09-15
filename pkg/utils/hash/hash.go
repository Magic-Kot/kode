package hash

import (
	"crypto/sha256"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const salt = "qwertyuiop0sdfghjklzxcvbnm123056QWE8TYUIOPASDFGHJKLZXCVBNM"

func GenerateHash(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func GenerateHashBcrypt(str string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(str), 12)

	return string(hash), err
}

func CompareHashBcrypt(str string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))
}
