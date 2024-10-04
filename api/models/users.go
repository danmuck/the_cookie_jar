package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
	Hash string `bson:"password" json:"password" form:"password"`
	JWT  string `bson:"jwt" json:"jwt"`
	Auth string `bson:"auth_token" json:"auth_token"`
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
	pw_bytes := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(pw_bytes, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	err = bcrypt.CompareHashAndPassword(hash, pw_bytes)
	if err != nil {
		panic(err)
	}

	u.Auth.Hash = string(hash)
}

func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Auth.Hash), []byte(password))
	return err == nil
}

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
func NewUser(name string, password string) *User {
	id := uuid.New()
	s := NewStatus("I'm new here.")
	u := &User{
		ID:       id.String(),
		Username: name,
		Status:   &s,
		Auth:     Credentials{Hash: password, Auth: "nil_auth", JWT: "nil_jwt"},
		Org:      "no organization",
	}

	status := fmt.Sprintf("I am new so have my -- username: %v password: %v",
		u.Username, password)

	u.UpdateStatus(status)

	return u
}
