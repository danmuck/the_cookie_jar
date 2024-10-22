package models

type User struct {
	ID       string      `bson:"_id" json:"id"`
	Username string      `bson:"username" json:"username" form:"username"`
	Auth     Credentials `bson:"auth" json:"auth"`
}

type Credentials struct {
	PasswordHash  string `bson:"password" json:"password" form:"password"`
	AuthTokenHash string `bson:"hashed_token" json:"hashed_token"`
}

func (u *User) GetId() string {
	return u.ID
}

func (u *User) GetUsername() string {
	return u.Username
}

func (u *User) GetPasswordHash() string {
	return u.Auth.PasswordHash
}

func (u *User) GetAuthTokenHash() string {
	return u.Auth.AuthTokenHash
}

// constructor -->
/*
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
*/
