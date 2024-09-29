package sandbox

import "time"

// go naming conventions make anything starting with lowercase letter --> private
type User struct {
	id       []byte
	username string
	status   *status
}

// since this starts with lowercase letter it is private and cannot be accessed outside of this package
type status struct {
	s string
	d time.Time
}

func (u *User) updateStatus(s string) {
	u.status = &status{
		s: s,
		d: time.Now(),
	}
}

// use public methods starting with capital letter to interface with private attributes
func (u *User) GetId() []byte {
	return u.id
}

func (u *User) GetUsername() string {
	return u.username
}

func (u *User) GetStatus_String() string {
	return u.status.s
}

// constructor -->
func NewUser() *User {
	u := &User{
		id:       []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09},
		username: "Big",
		status:   nil,
	}
	u.updateStatus("Chillin")

	return u
}
