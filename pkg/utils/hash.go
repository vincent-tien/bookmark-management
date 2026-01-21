package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(s string) string {
	hashBytes, _ := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)

	return string(hashBytes)
}

func VerifyPassword(pw, hashPw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPw), []byte(pw))

	return err == nil
}
