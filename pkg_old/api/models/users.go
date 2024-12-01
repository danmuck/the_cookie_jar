package models

type User struct {
	ID           string      `bson:"_id" json:"id"`
	Username     string      `bson:"username" json:"username" form:"username"`
	ClassroomIDs []string    `bson:"classroom_ids" json:"classroom_ids"`
	Auth         Credentials `bson:"auth" json:"auth"`
}

type Credentials struct {
	PasswordHash  string `bson:"password" json:"password" form:"password"`
	AuthTokenHash string `bson:"hashed_token" json:"hashed_token"`
}

type UserCommentStats struct {
    ID            string    `bson:"_id" json:"id"`
    Username      string    `bson:"username" json:"username"`
    TotalComments int       `bson:"total_comments" json:"total_comments"`
    TotalLikes    int       `bson:"total_likes" json:"total_likes"`
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
