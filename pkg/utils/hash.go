package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes the provided password using bcrypt.
func HashPassword(s string) string {
	hashBytes, _ := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)

	return string(hashBytes)
}

// VerifyPassword verifies if the provided password matches the hashed password.
func VerifyPassword(pw, hashPw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPw), []byte(pw))

	return err == nil
}
