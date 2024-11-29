package models

type User struct {
	Username     string      `json:"Username" bson:"Username"`
	ClassroomIDs []string    `json:"ClassroomIDs" bson:"ClassroomIDs"`
	Auth         Credentials `json:"Auth" bson:"Auth"`
}
type Credentials struct {
	PasswordHash  string `json:"PasswordHash" bson:"PasswordHash"`
	AuthTokenHash string `json:"AuthTokenHash" bson:"AuthTokenHash"`
}
