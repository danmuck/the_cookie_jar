package models

import "github.com/danmuck/the_cookie_jar/pkg/utils"

type Classroom struct {
	ID            string   `bson:"_id" json:"id"`
	Name          string   `bson:"name" json:"name"`
	ProfessorID   string   `bson:"professor" json:"professor"`
	InstructorIDs []string `bson:"instructors" json:"instructors"`
	StudentIDs    []string `bson:"students" json:"students"`
	BoardIDs      []string `bson:"boards" json:"boards"`
}

func (c *Classroom) ContainsUserID(id string) bool {
	return c.ProfessorID == id || utils.Contains(c.InstructorIDs, id) || utils.Contains(c.StudentIDs, id)
}

func (c *Classroom) IsUserIDPrivileged(id string) bool {
	return c.ProfessorID == id || utils.Contains(c.InstructorIDs, id)
}
