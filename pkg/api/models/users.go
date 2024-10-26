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
