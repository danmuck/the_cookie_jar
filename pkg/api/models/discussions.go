package models

import (
	"time"
)

type Board struct {
	ID        string   `bson:"_id" json:"id"`
	Name      string   `bson:"name" json:"name"`
	ThreadIDs []string `bson:"threads" json:"threads"`
}

type Thread struct {
	ID         string   `bson:"_id" json:"id"`
	Name       string   `bson:"name" json:"name"`
	CommentIDs []string `bson:"comments" json:"comments"`
}

type Comment struct {
	ID           string    `bson:"_id" json:"id"`
	Username     string    `bson:"username" json:"username"`
	Text         string    `bson:"text" json:"text"`
	LikedUserIDs []string  `bson:"likes" json:"likes"`
	Date         time.Time `bson:"date" json:"date"`
}
