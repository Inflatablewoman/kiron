package server

import "golang.org/x/crypto/bcrypt"

// createHashedPassword hashes a password and returns a hashed version of the password
func createHashedPassword(password string) (string, error) {
	// Strength based on...
	// http://chargen.matasano.com/chargen/2015/3/26/enough-with-the-salts-updates-on-secure-password-schemes.html
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 11)
	if err != nil {
		return "", err
	}


	return string(hash), err
}


// MatchPassword checks a raw password with a hashed password and returns true if they match
func MatchPassword(password string, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, err
	}

	return err == nil, err
}