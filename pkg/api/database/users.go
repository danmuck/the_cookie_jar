package database

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/danmuck/the_cookie_jar/pkg/api/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

/*
Adds a user to the database.

NOTE: By default a user's authentication token hash is an empty string.
*/
func AddUser(username string, password string) error {
	// Hashing the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Creating the new user
	user := &models.User{
		ID:       uuid.New().String(),
		Username: username,
		Auth:     models.Credentials{PasswordHash: string(hashedPassword), AuthTokenHash: ""},
	}

	// Trying to add user to the database assuming they don't already exist
	userCollection := GetCollection("users")
	err = userCollection.FindOne(context.TODO(), gin.H{"username": username}).Decode(&user)
	if err != nil {
		return err
	}
	_, err = userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}

	return nil
}

/*
Gets user model from the database.
*/
func GetUser(username string) (*models.User, error) {
	var user *models.User
	err := GetCollection("users").FindOne(context.TODO(), gin.H{"username": username}).Decode(&user)
	return user, err
}

/*
Will search for the user in the database based on username of given user model
and then update it to the given model.
*/
func UpdateUser(user *models.User) error {
	// Does the user exist
	_, err := GetUser(user.GetUsername())
	if err != nil {
		return err
	}

	err = GetCollection("users").FindOneAndReplace(context.TODO(), gin.H{"username": user.GetUsername()}, user).Err()
	return err
}

/*
Returns nil if given password matches the password associated with the given
username in the database.
*/
func VerifyPassword(username string, password string) error {
	// Grabbing user from database
	user, err := GetUser(username)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Auth.PasswordHash), []byte(password))
	return err
}

/*
Generates a JWT (authentication) token and returns it, but also updates JWT
token for the user in the database.

Tokens expire an hour after being handed out.
*/
func GenerateAuthToken(username string) (string, error) {
	// Obtaining the user
	user, err := GetUser(username)
	if err != nil {
		return "", err
	}

	// Generating the JWT auth token
	claims := jwt.MapClaims{"username": username, "exp": time.Now().Add(time.Hour * 1)}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("JWT_AUTH_TOKEN_SECRET")))
	if err != nil {
		return "", err
	}

	// Hashing the new JWT auth token, associating it with user.
	//
	// NOTE: This hashed token loses its metadata (like expiration), so we will
	//       need to grab that from the user cookies. This is still safe though
	//       because we don't trust user cookie auth token unless it matches
	//       the hashed auth token in the database :)
	//
	hasher := sha256.New()
	hasher.Write([]byte(token))
	user.Auth.AuthTokenHash = hex.EncodeToString(hasher.Sum(nil))

	// Updating user in database with new JWT auth token hash
	err = UpdateUser(user)
	if err != nil {
		return "", err
	}

	return token, err
}

/*
Returns nil if given authentication token, when hashed, matches the
authentication token hash associated with the given username in the database.
*/
func VerifyAuthToken(username string, token string) error {
	// Obtaining the user
	user, err := GetUser(username)
	if err != nil {
		return err
	}

	// Hashing the given token then comparing it to the user hashed token
	hasher := sha256.New()
	hasher.Write([]byte(token))
	if hex.EncodeToString(hasher.Sum(nil)) == user.GetAuthTokenHash() {
		tokenObj, _, err := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
		if err != nil {
			return err
		}

		claims, properToken := tokenObj.Claims.(jwt.MapClaims)
		if properToken {
			if time.Unix(int64(claims["exp"].(float64)), 0).Before(time.Now()) {
				// To-Do: Remove the auth token because it is expired
			}
		} else {
			return err
		}

		return nil
	}

	return fmt.Errorf("given authentication token does not match user token")
}
