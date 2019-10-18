package domain

import (
	"golang.org/x/crypto/bcrypt"
)

//GenerateByteHashedPassword Generates a hashed byte[] password
//from the user.passowrd data
func GenerateByteHashedPassword(password string) ([]byte, error) {
	res, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	return res, nil
}

//GenerateStringHashedPassword Generates a hashed string password
//from the user.passowrd data
func GenerateStringHashedPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	// Convert []byte to string
	// bytesRead := bytes.IndexByte(hashed, 0)
	res := string(hashed)

	return res, nil
}
