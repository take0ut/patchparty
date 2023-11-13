package auth

import (
	"log"
	"strconv"
	"time"

	"github.com/take0ut/patch-party/server/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// User is a user.
func HashAndSalt(
	password []byte,
) string {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	return string(hash)
}

// VerifyPassword verifies a password against a hash.
func VerifyPassword(
	password []byte,
	hash string,
) bool {
	byteHash := []byte(hash)
	err := bcrypt.CompareHashAndPassword(byteHash, password)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

// CreateSession creates a session for a user.
func CreateSession(
	users models.UserModel,
	userID int,
) string {
	sessionToken := HashAndSalt([]byte(strconv.Itoa(userID) + time.Now().String()))
	users.SetUserSession(userID, sessionToken)
	return sessionToken
}

func CheckIfUserExists(
	users models.UserModel,
	email string,
	username string,
) bool {
	exists, err := users.CheckIfUserExists(email, username)
	if err != nil || exists {
		return false
	}
	return true
}
