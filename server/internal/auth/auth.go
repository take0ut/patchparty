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

// CreateSession creates a session for a user and returns the cookie.
func CreateSession(
	users *models.UserModel,
	userID int,
) string {
	sessionToken := HashAndSalt([]byte(strconv.Itoa(userID) + time.Now().String()))
	users.SetUserSession(userID, sessionToken)
	return "patchparty_session=" + sessionToken + "; Path=/; Expires=" + getExpirationDate(7) + ";"
}

func getExpirationDate(days time.Duration) string {
	return time.Now().Add(days * (24 * time.Hour)).String()
}

func CheckIfUserExists(
	users *models.UserModel,
	email string,
	username string,
) bool {
	exists, err := users.CheckIfUserExists(email, username)
	if err != nil || exists {
		return false
	}
	return true
}

func LogoutUser(
	users *models.UserModel,
	userID int,
) bool {
	err := users.SetUserSession(userID, "")
	return err == nil
}
