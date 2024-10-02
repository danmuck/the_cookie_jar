package models

import (
	"github.com/google/uuid"
	"time"
)

// go naming conventions make anything starting with lowercase letter --> private
type User struct {
	ID       string `bson:"_id" json:"id"`
	Username string `bson:"username" json:"username"`
	Org      string `bson:"org" json:"org"`
	role     role
	Status   *status `bson:"status,omitempty" json:"status"`
}

type role struct {
	s    string
	auth string
}

// since this starts with lowercase letter it is private and cannot be accessed outside of this package
type status struct {
	ID        string    `bson:"_id" json:"id"`
	Status    string    `bson:"status" json:"status"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
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
func NewUser(name string) *User {
	id := uuid.New()
	s := NewStatus("I'm new here.")
	u := &User{
		ID:       id.String(),
		Username: name,
		Status:   &s,
		role:     role{s: "default", auth: "nil"},
		Org:      "Not Verified",
	}
	u.updateStatus("I am new here.")

	return u
}
