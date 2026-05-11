package utils

import "github.com/google/uuid"

func CreateToken() (plainUUID string, hashedToken string) {
	// 1. Create the secret
	u := uuid.New().String()

	// 2. Hash it for storage
	hashed := HashPlainText(u)

	return u, hashed
}
