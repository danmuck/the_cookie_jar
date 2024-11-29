package models

type Media struct {
	ID       string `bson:"_id" json:"id"`
	Path     string `json:"Path" bson:"Path"`
	AuthorID string `json:"AuthorID" bson:"AuthorID"`
}
