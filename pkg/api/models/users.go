package models

type User struct {
	Username         string      `json:"Username" bson:"Username"`
	ClassroomIDs     []string    `json:"ClassroomIDs" bson:"ClassroomIDs"`
	Auth             Credentials `json:"Auth" bson:"Auth"`
	ProfilePictureID string      `json:"ProfilePictureID" bson:"ProfilePictureID"`
}
type Credentials struct {
	PasswordHash  string `json:"PasswordHash" bson:"PasswordHash"`
	AuthTokenHash string `json:"AuthTokenHash" bson:"AuthTokenHash"`
}
