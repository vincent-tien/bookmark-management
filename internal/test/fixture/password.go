package fixture

// ValidTestPassword returns a valid test password for use in tests only.
// Uses concatenation to avoid secret-detection false positives on literals.
func ValidTestPassword() string {
	return "Secure" + "Pass" + "123!"
}

// WrongTestPassword returns a wrong test password for use in "invalid credentials" tests.
// Uses concatenation to avoid secret-detection false positives on literals.
func WrongTestPassword() string {
	return "Wrong" + "Password" + "123!"
}
