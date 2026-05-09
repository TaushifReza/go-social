package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(bytes), err
}

func CheckPassword(inputPassword, hashPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(inputPassword))

	return err
}
