package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// go naming conventions make anything starting with lowercase letter --> private
type User struct {
	ID       string `bson:"_id" json:"id"`
	Username string `bson:"username" json:"username" form:"username"`
	Org      string `bson:"org" json:"org" form:"org"`
	role     role
	Status   *status `bson:"status,omitempty" json:"status" form:"status"`
	Hash     string  `bson:"password" json:"password" form:"password"`
}

type role struct {
	s    string
	auth string
}

// since this starts with lowercase letter it is private and cannot be accessed outside of this package
type status struct {
	ID        string    `bson:"_id" json:"id" form:"id"`
	Status    string    `bson:"status" json:"status" form:"status"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp" form:"timestamp"`
}

func (u *User) updateStatus(s string) {
	u.Status = &status{
		ID:        uuid.New().String(),
		Status:    s,
		Timestamp: time.Now(),
	}
}

func NewStatus(s string) status {
	status := status{
		ID:        uuid.New().String(),
		Status:    s,
		Timestamp: time.Now(),
	}
	return status
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
		role:     role{s: "default", auth: "nil"},
		Org:      "no organization",
		Hash:     password,
	}
	// pw_bytes := []byte(password)
	// hash, err := bcrypt.GenerateFromPassword(pw_bytes, bcrypt.DefaultCost)
	// if err != nil {
	// 	panic(err)
	// }

	// Comparing the password with the hash
	// err = bcrypt.CompareHashAndPassword(hash, pw_bytes)
	// if err != nil {
	// 	panic(err)
	// }

	status := fmt.Sprintf("I am new so have my -- username: %v password: %v",
		u.Username, password)

	u.updateStatus(status)

	return u
}
