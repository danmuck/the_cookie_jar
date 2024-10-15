/*
	Most of these functions are just wrappers right now. Using them should be pretty straightforward

*/

package middleware

import (
    "crypto/rand"
    "golang.org/x/crypto/bcrypt"
)

// This is just a wrapper for the bcrypt function. 
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

// Returns true if passwords match, false if not.
func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
//This generates a high entropy token string and returns it as a string.
func GenToken() (string, error) {
    // Generate random bytes
    token := make([]byte, 32)
    _, err := rand.Read(token)
    if err != nil {
        return "", err
    }
    
    return string(token), nil
}
//These area, again, just wrappers for bcrypt functions.
func HashToken(token string) (string) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
    if err != nil{
		return string(bytes)
	}
	//**TODO** Change the signature and allow for errors
	return ""
}
//Once again, a bcrypt wrapper.
func ValidateToken(token, hashToken string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashToken), []byte(token))
    return err == nil
}