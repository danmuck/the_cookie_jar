package models

import (
	"github.com/danmuck/the_cookie_jar/pkg/utils"
)

type Classroom struct {
	ID          string    `json:"id" bson:"_id"`
	Name        string    `json:"Name" bson:"Name"`
	ProfessorID string    `json:"ProfessorID" bson:"ProfessorID"`
	StudentIDs  []string  `json:"StudentIDs" bson:"StudentIDs"`
	ThreadIDs   []string  `json:"ThreadIDs" bson:"ThreadIDs"`
	Game        ClassGame `json:"Game" bson:"Game"`
}
type Thread struct {
	ID         string   `json:"id" bson:"_id"`
	Title      string   `json:"Title" bson:"Title"`
	Date       string   `json:"Date" bson:"Date"`
	AuthorID   string   `json:"AuthorID" bson:"AuthorID"`
	CommentIDs []string `json:"CommentIDs" bson:"CommentIDs"`
}
type Comment struct {
	ID           string   `json:"id" bson:"_id"`
	Text         string   `json:"Text" bson:"Text"`
	AuthorID     string   `json:"AuthorID" bson:"AuthorID"`
	LikedUserIDs []string `json:"LikedUserIDs" bson:"LikedUserIDs"`
}
type ClassGame struct {
	// Dunno yet
}

func (c *Classroom) ContainsUserID(id string) bool {
	return c.ProfessorID == id || utils.Contains(c.StudentIDs, id)
}

func (c *Classroom) IsUserIDPrivileged(id string) bool {
	return c.ProfessorID == id
}
