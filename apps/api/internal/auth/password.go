package auth

import "golang.org/x/crypto/bcrypt"

// bcryptCost balances resistance against latency for interactive logins.
const bcryptCost = 12

// HashPassword returns a bcrypt hash of the plaintext password. bcrypt rejects
// inputs longer than 72 bytes; callers validate the length before reaching here
// (see the auth service), so a too-long password surfaces as an error.
func HashPassword(plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword returns nil when the plaintext matches the stored hash, and a
// non-nil error otherwise. It runs in constant time relative to the hash.
func VerifyPassword(hash string, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}
