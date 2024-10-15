package models

import (
	"fmt"
	"time"
	"github.com/danmuck/the_cookie_jar/pkg/api/middleware"
	"github.com/google/uuid"
)

// go naming conventions make anything starting with lowercase letter --> private
type User struct {
	ID       string      `bson:"_id" json:"id"`
	Username string      `bson:"username" json:"username" form:"username"`
	Org      string      `bson:"org" json:"org" form:"org"`
	Auth     Credentials `bson:"auth" json:"auth"`
	Status   *status     `bson:"status,omitempty" json:"status" form:"status"`
}

type Credentials struct {
	Hash           string    `bson:"password" json:"password" form:"password"`
	HashedToken    string    `bson:"hashed_token" json:"hashed_token"`
	TokenExpiration time.Time `bson:"token_expiration" json:"token_expiration"`
}
// since this starts with lowercase letter it is private and cannot be accessed outside of this package
type status struct {
	ID        string    `bson:"_id" json:"id" form:"id"`
	Status    string    `bson:"status" json:"status" form:"status"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp" form:"timestamp"`
}

func NewStatus(s string) status {
	status := status{
		ID:        uuid.New().String(),
		Status:    s,
		Timestamp: time.Now(),
	}
	return status
}

func (u *User) UpdateStatus(s string) {
	u.Status = &status{
		ID:        uuid.New().String(),
		Status:    s,
		Timestamp: time.Now(),
	}
}

func (u *User) UpdatePassword(password string) {
	if !middleware.CheckPasswordHash(password, u.Auth.Hash) {
		fmt.Errorf("current password is incorrect")
		return
	}
	
	hashedPassword, err := middleware.HashPassword(password)
	if err != nil{
		panic(err)
	}

	u.Auth.Hash = hashedPassword
}

func (u *User) VerifyPassword(password string) bool {
    result := middleware.CheckPasswordHash(password, u.Auth.Hash)
    fmt.Printf("Password verification result: %v\n", result)
    return result}

// use public methods starting with capital letter to interface with private attributes
func (u *User) GetId() string {
	return u.ID
}

func (u *User) GetUsername() string {
	return u.Username
}

func (u *User) GetStatus_String() string {
	return u.Status.Status
}

// constructor -->
func NewUser(name string, password string) (*User, error) {
	id := uuid.New()
	s := NewStatus("I'm new here.")
	
	hashedPassword, err := middleware.HashPassword(password)
	if err != nil {
		return nil, err
	}
	var placeholderTime time.Time
	u := &User{
		ID:       id.String(),
		Username: name,
		Status:   &s,
		Auth:     Credentials{Hash: hashedPassword, HashedToken: "nil_auth", TokenExpiration: placeholderTime},
		Org:      "no organization",
	}

	


	return u, nil
}