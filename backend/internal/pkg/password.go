package pkg

import "golang.org/x/crypto/bcrypt"

// DefaultPassword is the default password for new members.
const DefaultPassword = "home123"

// HashPassword hashes a password using bcrypt.
func HashPassword(pwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword compares a hashed password with a plain text one.
func CheckPassword(hash, pwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd)) == nil
}
