package db_types

import "time"
import "github.com/google/uuid"

// go naming conventions make anything starting with lowercase letter --> private
type User struct {
	ID       string  `bson:"_id"`
	Username string  `bson:"username"`
	Org      string  `bson:"org"`
	Status   *status `bson:"status,omitempty"`
}

// since this starts with lowercase letter it is private and cannot be accessed outside of this package
type status struct {
	ID        string    `bson:"_id"`
	Status    string    `bson:"status"`
	Timestamp time.Time `bson:"timestamp"`
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
		Org:      "Not Verified",
	}
	u.updateStatus("I am new here.")

	return u
}
