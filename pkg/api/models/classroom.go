package models

type Classroom struct {
	ID            string   `bson:"_id" json:"id"`
	Name          string   `bson:"name" json:"name"`
	ProfessorID   string   `bson:"professor" json:"professor"`
	InstructorIDs []string `bson:"instructors" json:"instructors"`
	StudentIDs    []string `bson:"students" json:"students"`
	BoardIDs      []string `bson:"boards" json:"boards"`
}
